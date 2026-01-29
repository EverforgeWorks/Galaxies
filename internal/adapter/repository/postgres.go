package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// --- PLAYER & SHIP PERSISTENCE ---

func (r *PostgresRepository) SavePlayer(ctx context.Context, p *entity.Player) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var sysID *uuid.UUID
	if p.CurrentSystem != nil {
		sysID = &p.CurrentSystem.ID
	}

	var extID *string
	if p.ExternalID != "" {
		extID = &p.ExternalID
	}

	// 1. Save Player (Included is_admin)
	qPlayer := `
        INSERT INTO players (id, name, external_id, credits, current_system_id, is_admin, last_login)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        ON CONFLICT (id) DO UPDATE 
        SET name = $2, external_id = $3, credits = $4, current_system_id = $5, is_admin = $6, last_login = NOW();
    `
	_, err = tx.Exec(ctx, qPlayer, p.ID, p.Name, extID, p.Credits, sysID, p.IsAdmin)
	if err != nil {
		return fmt.Errorf("db error saving player: %w", err)
	}

	// 2. Save Ship
	if p.Ship != nil {
		statsJson, _ := json.Marshal(p.Ship.Stats)
		crewJson, _ := json.Marshal(p.Ship.Crew)
		passJson, _ := json.Marshal(p.Ship.Passengers)

		qShip := `
            INSERT INTO ships (
                id, player_id, name, model_name, 
                current_hull, current_shield, current_fuel,
                stats, crew, passengers
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            ON CONFLICT (id) DO UPDATE 
            SET 
                name = $3, model_name = $4,
                current_hull = $5, current_shield = $6, current_fuel = $7,
                stats = $8, crew = $9, passengers = $10;
        `
		_, err = tx.Exec(ctx, qShip,
			p.Ship.ID,
			p.ID,
			p.Ship.Name,
			p.Ship.ModelName,
			p.Ship.CurrentHull,
			p.Ship.CurrentShield,
			p.Ship.CurrentFuel,
			statsJson,
			crewJson,
			passJson,
		)
		if err != nil {
			return fmt.Errorf("db error saving ship: %w", err)
		}

		// 3. Save Inventory
		_, err = tx.Exec(ctx, "DELETE FROM ship_inventory WHERE ship_id = $1", p.Ship.ID)
		if err != nil {
			return fmt.Errorf("failed to clear inventory: %w", err)
		}

		for _, item := range p.Ship.Cargo {
			_, err = tx.Exec(ctx, `
				INSERT INTO ship_inventory (id, ship_id, name, category, base_value, rarity, is_illegal, quantity, avg_cost)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, item.ID, p.Ship.ID, item.Name, item.Category, item.BaseValue, item.Rarity, item.IsIllegal, item.Quantity, item.AvgCost)
			if err != nil {
				return fmt.Errorf("failed to insert inventory item: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetPlayer(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	query := `
        SELECT 
            p.id, p.name, p.external_id, p.credits, p.current_system_id, p.is_admin,
            s.id, s.name, s.model_name, 
            s.current_hull, s.current_shield, s.current_fuel,
            s.stats, s.crew, s.passengers
        FROM players p
        LEFT JOIN ships s ON p.id = s.player_id
        WHERE p.id = $1
    `
	return r.scanPlayer(ctx, query, id)
}

func (r *PostgresRepository) GetPlayerByExternalID(ctx context.Context, externalID string) (*entity.Player, error) {
	query := `
        SELECT 
            p.id, p.name, p.external_id, p.credits, p.current_system_id, p.is_admin,
            s.id, s.name, s.model_name, 
            s.current_hull, s.current_shield, s.current_fuel,
            s.stats, s.crew, s.passengers
        FROM players p
        LEFT JOIN ships s ON p.id = s.player_id
        WHERE p.external_id = $1
    `
	return r.scanPlayer(ctx, query, externalID)
}

func (r *PostgresRepository) scanPlayer(ctx context.Context, query string, arg interface{}) (*entity.Player, error) {
	var p entity.Player
	var s entity.Ship
	var sysID *uuid.UUID
	var extID *string

	var sID *uuid.UUID
	var sName, sModel *string
	var sHull, sShield, sFuel *float64
	var statsRaw, crewRaw, passRaw []byte

	err := r.db.QueryRow(ctx, query, arg).Scan(
		&p.ID, &p.Name, &extID, &p.Credits, &sysID, &p.IsAdmin, // Added IsAdmin
		&sID, &sName, &sModel,
		&sHull, &sShield, &sFuel,
		&statsRaw, &crewRaw, &passRaw,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load player: %w", err)
	}

	if extID != nil {
		p.ExternalID = *extID
	}

	if sID != nil {
		s.ID = *sID
		s.Name = *sName
		s.ModelName = *sModel
		s.CurrentHull = *sHull
		s.CurrentShield = *sShield
		s.CurrentFuel = *sFuel
		_ = json.Unmarshal(statsRaw, &s.Stats)
		_ = json.Unmarshal(crewRaw, &s.Crew)
		_ = json.Unmarshal(passRaw, &s.Passengers)
		if s.Crew == nil {
			s.Crew = []entity.Crew{}
		}
		if s.Passengers == nil {
			s.Passengers = []entity.Passenger{}
		}

		invRows, err := r.db.Query(ctx, `
			SELECT id, name, category, base_value, rarity, is_illegal, quantity, avg_cost
			FROM ship_inventory WHERE ship_id = $1
		`, s.ID)
		if err == nil {
			var cargo []entity.Item
			for invRows.Next() {
				var i entity.Item
				invRows.Scan(&i.ID, &i.Name, &i.Category, &i.BaseValue, &i.Rarity, &i.IsIllegal, &i.Quantity, &i.AvgCost)
				cargo = append(cargo, i)
			}
			invRows.Close()
			s.Cargo = cargo
		} else {
			s.Cargo = []entity.Item{}
		}
		p.Ship = &s
	}
	if sysID != nil {
		p.CurrentSystem = &entity.System{ID: *sysID}
	}
	return &p, nil
}

// --- SYSTEM PERSISTENCE ---

func (r *PostgresRepository) LoadUniverse(ctx context.Context) (map[uuid.UUID]*entity.System, error) {
	query := `SELECT id, name, x, y, political, economic, social, stats FROM systems`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query universe: %w", err)
	}

	universe := make(map[uuid.UUID]*entity.System)
	var loadedSystems []*entity.System

	for rows.Next() {
		var s entity.System
		var statsRaw []byte
		var polStr, ecoStr, socStr *string
		if err := rows.Scan(&s.ID, &s.Name, &s.X, &s.Y, &polStr, &ecoStr, &socStr, &statsRaw); err != nil {
			continue
		}
		if polStr != nil {
			val, _ := strconv.Atoi(*polStr)
			s.Political = domain.PoliticalStatus(val)
		}
		if ecoStr != nil {
			val, _ := strconv.Atoi(*ecoStr)
			s.Economic = domain.EconomicStatus(val)
		}
		if socStr != nil {
			val, _ := strconv.Atoi(*socStr)
			s.Social = domain.SocialStatus(val)
		}
		_ = json.Unmarshal(statsRaw, &s.Stats)
		s.Market = []entity.Item{}
		loadedSystems = append(loadedSystems, &s)
	}
	rows.Close()

	for _, sys := range loadedSystems {
		market, err := r.GetSystemMarket(ctx, sys.ID)
		if err == nil {
			sys.Market = market
		}
		universe[sys.ID] = sys
	}
	return universe, nil
}

func (r *PostgresRepository) SaveUniverse(ctx context.Context, systems []*entity.System) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, sys := range systems {
		statsJson, _ := json.Marshal(sys.Stats)
		_, err := tx.Exec(ctx, `
			INSERT INTO systems (id, name, x, y, political, economic, social, stats)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO NOTHING;
		`, sys.ID, sys.Name, sys.X, sys.Y,
			fmt.Sprintf("%d", sys.Political), fmt.Sprintf("%d", sys.Economic), fmt.Sprintf("%d", sys.Social),
			statsJson,
		)
		if err != nil {
			return err
		}
		for _, item := range sys.Market {
			_, err := tx.Exec(ctx, `
				INSERT INTO market_listings (id, system_id, name, category, base_value, rarity, is_illegal, quantity)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			`, item.ID, sys.ID, item.Name, item.Category, item.BaseValue, item.Rarity, item.IsIllegal, item.Quantity)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) UpdateMarket(ctx context.Context, systemID uuid.UUID, market []entity.Item) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "DELETE FROM market_listings WHERE system_id = $1", systemID)
	if err != nil {
		return err
	}
	for _, item := range market {
		_, err := tx.Exec(ctx, `
			INSERT INTO market_listings (id, system_id, name, category, base_value, rarity, is_illegal, quantity)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, item.ID, systemID, item.Name, item.Category, item.BaseValue, item.Rarity, item.IsIllegal, item.Quantity)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetSystemMarket(ctx context.Context, systemID uuid.UUID) ([]entity.Item, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, category, base_value, rarity, is_illegal, quantity
		FROM market_listings WHERE system_id = $1
	`, systemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Item
	for rows.Next() {
		var i entity.Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Category, &i.BaseValue, &i.Rarity, &i.IsIllegal, &i.Quantity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

// RawSQL allows execution of arbitrary commands from the admin console
func (r *PostgresRepository) RawSQL(ctx context.Context, query string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	fieldDescriptions := rows.FieldDescriptions()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, fd := range fieldDescriptions {
			rowMap[fd.Name] = values[i]
		}
		results = append(results, rowMap)
	}

	return results, nil
}
