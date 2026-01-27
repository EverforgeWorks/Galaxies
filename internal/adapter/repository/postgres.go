package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"galaxies/internal/core/entity"
	"galaxies/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


// PostgresRepository holds the connection pool.
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository initializes the repo with a pool.
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// --- PLAYER & SHIP PERSISTENCE ---

// SavePlayer persists the Player and their linked Ship in a single Atomic Transaction.
func (r *PostgresRepository) SavePlayer(ctx context.Context, p *entity.Player) error {
	// 1. Thread Safety: Lock the player while we read/marshal their data.
	// We don't want the credits changing while we are writing to the DB.
	// (Assumes you added a Public RLock helper or just lock explicitly if allowed)
    // p.Mu.RLock() 
    // defer p.Mu.RUnlock() 
    // Note: Since 'mu' is private in your struct, you might need to add a `Lock/Unlock` method 
    // to Player, or just ensure this is called from the Engine which owns the lock.

	// 2. Start Transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	// Always Rollback. If Commit() is called later, Rollback is a no-op.
	defer tx.Rollback(ctx)

	// 3. Prepare Data for "ships" table (JSON Marshaling)
	// We marshal the sub-structs to match our JSONB columns
	statsJson, _ := json.Marshal(p.Ship.Stats)
	cargoJson, _ := json.Marshal(p.Ship.Cargo)
	crewJson, _ := json.Marshal(p.Ship.Crew)
	passJson, _ := json.Marshal(p.Ship.Passengers)

	// 4. Upsert Player
	// We use ON CONFLICT to handle both "New Character" and "Save Game" scenarios
	qPlayer := `
		INSERT INTO players (id, name, credits, current_system_id, last_login)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE 
		SET credits = $3, current_system_id = $4, last_login = NOW();
	`
	_, err = tx.Exec(ctx, qPlayer, 
		p.ID, 
		p.Name, 
		p.Credits, 
		p.CurrentSystem.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to save player: %w", err)
	}

	// 5. Upsert Ship
	qShip := `
		INSERT INTO ships (
			id, player_id, name, model_name, 
			current_hull, current_shield, current_fuel,
			stats, cargo, crew, passengers
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE 
		SET 
			current_hull = $5, current_shield = $6, current_fuel = $7,
			stats = $8, cargo = $9, crew = $10, passengers = $11;
	`
	_, err = tx.Exec(ctx, qShip,
		p.Ship.ID,
		p.ID, // linking back to player
		p.Ship.Name,
		p.Ship.ModelName,
		p.Ship.CurrentHull,
		p.Ship.CurrentShield,
		p.Ship.CurrentFuel,
		statsJson,
		cargoJson,
		crewJson,
		passJson,
	)
	if err != nil {
		return fmt.Errorf("failed to save ship: %w", err)
	}

	// 6. Commit Transaction
	return tx.Commit(ctx)
}

// LoadPlayer reconstructs the entire Player > Ship > System graph.
func (r *PostgresRepository) LoadPlayer(ctx context.Context, playerID uuid.UUID) (*entity.Player, error) {
	// Query: Join Players + Ships to get everything in one round-trip
	query := `
		SELECT 
			p.id, p.name, p.credits, p.current_system_id,
			s.id, s.name, s.model_name, 
			s.current_hull, s.current_shield, s.current_fuel,
			s.stats, s.cargo, s.crew, s.passengers
		FROM players p
		JOIN ships s ON p.id = s.player_id
		WHERE p.id = $1
	`

	var p entity.Player
	var s entity.Ship
	var sysID uuid.UUID
	
	// Variables to hold the raw JSON bytes before we unmarshal
	var statsRaw, cargoRaw, crewRaw, passRaw []byte

	err := r.db.QueryRow(ctx, query, playerID).Scan(
		&p.ID, &p.Name, &p.Credits, &sysID,
		&s.ID, &s.Name, &s.ModelName,
		&s.CurrentHull, &s.CurrentShield, &s.CurrentFuel,
		&statsRaw, &cargoRaw, &crewRaw, &passRaw,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Or a specific ErrNotFound
		}
		return nil, fmt.Errorf("failed to load player: %w", err)
	}

	// Hydrate the Ship Struct from JSONB
	if err := json.Unmarshal(statsRaw, &s.Stats); err != nil { return nil, err }
	if err := json.Unmarshal(cargoRaw, &s.Cargo); err != nil { return nil, err }
	if err := json.Unmarshal(crewRaw, &s.Crew); err != nil { return nil, err }
	if err := json.Unmarshal(passRaw, &s.Passengers); err != nil { return nil, err }

	// Handle Arrays being nil if DB was empty JSON '[]' or null
	if s.Cargo == nil { s.Cargo = []entity.Item{} }
	if s.Crew == nil { s.Crew = []entity.Crew{} }
	if s.Passengers == nil { s.Passengers = []entity.Passenger{} }

	// Connect the Ship to the Player
	p.Ship = &s

	// Hydrate the System (Location)
	// We need a separate call to get the System data so the Player has context
	sys, err := r.GetSystem(ctx, sysID)
	if err != nil {
		return nil, fmt.Errorf("failed to load system for player: %w", err)
	}
	p.CurrentSystem = sys

	return &p, nil
}

// --- SYSTEM PERSISTENCE ---

// GetSystem loads a single star system. 
// Cached heavily in a real app, but direct DB lookup for MVP.
func (r *PostgresRepository) GetSystem(ctx context.Context, id uuid.UUID) (*entity.System, error) {
	query := `
		SELECT id, name, x, y, political, economic, social, stats
		FROM systems WHERE id = $1
	`
	var sys entity.System
	var statsRaw []byte
	
	// We scan strings for the Enums for now, or you can cast them if you made them types
	var pol, eco, soc string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sys.ID, &sys.Name, &sys.X, &sys.Y,
		&pol, &eco, &soc, &statsRaw,
	)
	if err != nil {
		return nil, err
	}

	// Map strings back to Enums (assuming simple string conversion works or do a switch)
	// For MVP we just assume the string matches. You might need a helper here.
	// sys.Political = domain.PoliticalStatus(pol) ...

	if err := json.Unmarshal(statsRaw, &sys.Stats); err != nil {
		return nil, err
	}

	return &sys, nil
}

// internal/adapter/repository/postgres.go (Add this method)

// SaveUniverse performs a high-performance batch insert of systems.
func (r *PostgresRepository) SaveUniverse(ctx context.Context, systems []*entity.System) error {
    batch := &pgx.Batch{}

    for _, s := range systems {
        statsJson, _ := json.Marshal(s.Stats)

        // We cast the enums to int for storage (or string if you have a String() method)
        // Adjust these casts based on your domain package definitions.
        politicalVal := fmt.Sprintf("%d", s.Political)
        economicVal := fmt.Sprintf("%d", s.Economic)
        socialVal := fmt.Sprintf("%d", s.Social)

        query := `
            INSERT INTO systems (id, name, x, y, political, economic, social, stats)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            ON CONFLICT (id) DO NOTHING; -- Skip if already seeded
        `
        batch.Queue(query, s.ID, s.Name, s.X, s.Y, politicalVal, economicVal, socialVal, statsJson)
    }

    // Execute the batch
    br := r.db.SendBatch(ctx, batch)
    defer br.Close()

    // Check for errors in the batch execution
    for i := 0; i < len(systems); i++ {
        _, err := br.Exec()
        if err != nil {
            return fmt.Errorf("failed to insert system %d: %w", i, err)
        }
    }

    return nil
}
// LoadUniverse fetches ALL systems from the database to populate the game engine memory.
func (r *PostgresRepository) LoadUniverse(ctx context.Context) (map[uuid.UUID]*entity.System, error) {
    query := `
        SELECT id, name, x, y, political, economic, social, stats 
        FROM systems
    `
    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query universe: %w", err)
    }
    defer rows.Close()

    universe := make(map[uuid.UUID]*entity.System)

    for rows.Next() {
        var s entity.System
        var statsRaw []byte
        
        // Scan into STRINGS first because DB column is TEXT
        var polStr, ecoStr, socStr string

        if err := rows.Scan(&s.ID, &s.Name, &s.X, &s.Y, &polStr, &ecoStr, &socStr, &statsRaw); err != nil {
            return nil, fmt.Errorf("failed to scan system: %w", err)
        }

        // Convert String -> Int
        pVal, _ := strconv.Atoi(polStr)
        eVal, _ := strconv.Atoi(ecoStr)
        sVal, _ := strconv.Atoi(socStr)

        // Cast Int -> Enum
        s.Political = domain.PoliticalStatus(pVal)
        s.Economic = domain.EconomicStatus(eVal)
        s.Social = domain.SocialStatus(sVal)

        // Hydrate Stats from JSONB
        if err := json.Unmarshal(statsRaw, &s.Stats); err != nil {
            return nil, fmt.Errorf("failed to unmarshal stats for system %s: %w", s.ID, err)
        }

        universe[s.ID] = &s
    }

    return universe, nil
}