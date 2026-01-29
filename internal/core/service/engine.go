package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
	"galaxies/internal/core/gen"
)

var (
	ErrPlayerNotOnline   = errors.New("player not active")
	ErrSystemNotFound    = errors.New("system not found")
	ErrInsufficientFunds = errors.New("insufficient credits")
	ErrInsufficientStock = errors.New("market out of stock")
	ErrItemNotFound      = errors.New("item not found")
)

type GameEngine struct {
	mu             sync.RWMutex
	ActivePlayers  map[uuid.UUID]*entity.Player
	Universe       map[uuid.UUID]*entity.System
	repo           *repository.PostgresRepository
	OnPlayerUpdate func(p *entity.Player)
}

func NewGameEngine(repo *repository.PostgresRepository, universe map[uuid.UUID]*entity.System) *GameEngine {
	e := &GameEngine{
		ActivePlayers:  make(map[uuid.UUID]*entity.Player),
		Universe:       universe,
		repo:           repo,
		OnPlayerUpdate: func(p *entity.Player) {},
	}

	// START AUTO-SAVER (Every 1 Minute) as a fallback
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			e.SaveAll()
		}
	}()

	return e
}

func (e *GameEngine) SaveAll() {
	e.mu.RLock()
	players := make([]*entity.Player, 0, len(e.ActivePlayers))
	for _, p := range e.ActivePlayers {
		players = append(players, p)
	}
	e.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, p := range players {
		p.Lock()
		_ = e.repo.SavePlayer(ctx, p)
		p.Unlock()
	}
}

// --- AUTHENTICATION & ONBOARDING ---

func (e *GameEngine) AuthenticateGitHub(ctx context.Context, externalID string, githubName string) (*entity.Player, bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	existing, err := e.repo.GetPlayerByExternalID(ctx, externalID)
	if err == nil && existing != nil {
		if existing.CurrentSystem != nil {
			if sys, ok := e.Universe[existing.CurrentSystem.ID]; ok {
				existing.CurrentSystem = sys
			}
		}
		e.ActivePlayers[existing.ID] = existing
		isNew := (existing.Ship == nil)
		return existing, isNew, nil
	}

	newPlayer := &entity.Player{
		ID:         uuid.New(),
		ExternalID: externalID,
		Name:       githubName,
		Credits:    0,
	}

	// Synchronous Save for Auth
	if err := e.repo.SavePlayer(ctx, newPlayer); err != nil {
		return nil, false, fmt.Errorf("failed to create draft player: %w", err)
	}

	e.ActivePlayers[newPlayer.ID] = newPlayer
	return newPlayer, true, nil
}

func (e *GameEngine) GenerateStarterOptions() []*entity.Ship {
	return []*entity.Ship{
		gen.GenerateShipByChassis(domain.ChassisInterceptor, "Interceptor"),
		gen.GenerateShipByChassis(domain.ChassisCourier, "Courier"),
		gen.GenerateShipByChassis(domain.ChassisProspector, "Prospector"),
	}
}

func (e *GameEngine) CompleteOnboarding(ctx context.Context, playerID uuid.UUID, chosenName string, shipTemplate *entity.Ship) (*entity.Player, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	player, ok := e.ActivePlayers[playerID]
	if !ok {
		return nil, errors.New("session expired")
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9_]{3,18}$`)
	if !validName.MatchString(chosenName) {
		return nil, errors.New("invalid name: must be 3-18 chars, alphanumeric only (no spaces)")
	}

	var startSystem *entity.System
	for _, s := range e.Universe {
		startSystem = s
		break
	}
	if startSystem == nil {
		return nil, errors.New("universe is empty")
	}

	player.Name = chosenName
	player.CurrentSystem = startSystem
	player.Credits = 1000

	shipTemplate.ID = uuid.New()
	player.Ship = shipTemplate

	// Synchronous Save for Onboarding
	if err := e.repo.SavePlayer(ctx, player); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("callsign already taken")
		}
		return nil, fmt.Errorf("failed to register pilot: %w", err)
	}

	return player, nil
}

// --- ECONOMY & TRADE (SYNCHRONOUS) ---

func (e *GameEngine) BuyItem(ctx context.Context, playerID uuid.UUID, itemName domain.ItemName, qty int) error {
	if qty <= 0 {
		return errors.New("invalid quantity")
	}

	e.mu.RLock()
	player, pOk := e.ActivePlayers[playerID]
	e.mu.RUnlock()

	if !pOk {
		return ErrPlayerNotOnline
	}

	// 1. LOCK RESOURCES
	// Critical: Lock System first to prevent Market Race Conditions
	sys := player.CurrentSystem
	sys.Lock()
	defer sys.Unlock()

	// Then Lock Player
	player.Lock()
	defer player.Unlock()
	player.Ship.Lock()
	defer player.Ship.Unlock()

	// 2. VALIDATE STATE
	var marketItem *entity.Item
	marketIdx := -1
	for i, item := range sys.Market {
		if item.Name == itemName {
			marketItem = &sys.Market[i]
			marketIdx = i
			break
		}
	}
	if marketItem == nil {
		return ErrItemNotFound
	}
	if marketItem.Quantity < qty {
		return ErrInsufficientStock
	}

	unitPrice := CalculatePrice(marketItem.BaseValue, sys.Stats.MarketBuyMult, sys.Stats.TaxRate)
	totalCost := unitPrice * qty

	if player.Credits < totalCost {
		return ErrInsufficientFunds
	}
	if !player.Ship.CanFit(qty) {
		return entity.ErrCargoFull
	}

	// 3. MUTATE MEMORY
	player.Credits -= totalCost
	sys.Market[marketIdx].Quantity -= qty

	existingIdx := -1
	for i, invItem := range player.Ship.Cargo {
		if invItem.Name == itemName {
			existingIdx = i
			break
		}
	}

	if existingIdx != -1 {
		stack := &player.Ship.Cargo[existingIdx]
		currentTotalVal := float64(stack.Quantity) * stack.AvgCost
		newTotalVal := float64(qty) * float64(unitPrice)
		stack.Quantity += qty
		stack.AvgCost = (currentTotalVal + newTotalVal) / float64(stack.Quantity)
	} else {
		newItem := *marketItem
		newItem.ID = uuid.New()
		newItem.Quantity = qty
		newItem.AvgCost = float64(unitPrice)
		player.Ship.Cargo = append(player.Ship.Cargo, newItem)
	}

	// 4. SYNCHRONOUS PERSISTENCE
	// We wait for the DB. If this fails, we return an error.
	// (Improvement: We could rollback memory here, but failing securely is the priority)
	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ideally, these would be in a SQL transaction, but sequential is safer than async
	if err := e.repo.SavePlayer(saveCtx, player); err != nil {
		return fmt.Errorf("transaction failed (player save): %w", err)
	}
	if err := e.repo.UpdateMarket(saveCtx, sys.ID, sys.Market); err != nil {
		// Note: If this fails, player lost money but market didn't deduct stock. 
		// This is a known consistency risk without SQL transactions, but preferable to async data loss.
		return fmt.Errorf("transaction failed (market update): %w", err)
	}

	// 5. PUSH UPDATE
	if e.OnPlayerUpdate != nil {
		e.OnPlayerUpdate(player)
	}

	return nil
}

func (e *GameEngine) SellItem(ctx context.Context, playerID uuid.UUID, itemName domain.ItemName, qty int) error {
	if qty <= 0 {
		return errors.New("invalid quantity")
	}

	e.mu.RLock()
	player, pOk := e.ActivePlayers[playerID]
	e.mu.RUnlock()

	if !pOk {
		return ErrPlayerNotOnline
	}

	// 1. LOCK RESOURCES
	sys := player.CurrentSystem
	sys.Lock()
	defer sys.Unlock()

	player.Lock()
	defer player.Unlock()
	player.Ship.Lock()
	defer player.Ship.Unlock()

	// 2. VALIDATE
	var invItem *entity.Item
	invIdx := -1
	for i, item := range player.Ship.Cargo {
		if item.Name == itemName {
			invItem = &player.Ship.Cargo[i]
			invIdx = i
			break
		}
	}
	if invItem == nil {
		return ErrItemNotFound
	}
	if invItem.Quantity < qty {
		return errors.New("not enough items")
	}

	// 3. MUTATE MEMORY
	unitPrice := int(float64(invItem.BaseValue) * sys.Stats.MarketSellMult)
	totalVal := unitPrice * qty

	player.Credits += totalVal
	player.Ship.Cargo[invIdx].Quantity -= qty

	if player.Ship.Cargo[invIdx].Quantity <= 0 {
		lastIdx := len(player.Ship.Cargo) - 1
		player.Ship.Cargo[invIdx] = player.Ship.Cargo[lastIdx]
		player.Ship.Cargo = player.Ship.Cargo[:lastIdx]
	}

	mktIdx := -1
	for i, mItem := range sys.Market {
		if mItem.Name == itemName {
			mktIdx = i
			break
		}
	}
	if mktIdx != -1 {
		sys.Market[mktIdx].Quantity += qty
	} else {
		newItem := *invItem
		newItem.ID = uuid.New()
		newItem.Quantity = qty
		newItem.AvgCost = 0
		sys.Market = append(sys.Market, newItem)
	}

	// 4. SYNCHRONOUS PERSISTENCE
	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.repo.SavePlayer(saveCtx, player); err != nil {
		return fmt.Errorf("transaction failed (player save): %w", err)
	}
	if err := e.repo.UpdateMarket(saveCtx, sys.ID, sys.Market); err != nil {
		return fmt.Errorf("transaction failed (market update): %w", err)
	}

	// 5. PUSH UPDATE
	if e.OnPlayerUpdate != nil {
		e.OnPlayerUpdate(player)
	}

	return nil
}

func CalculatePrice(baseValue int, multiplier float64, taxRate float64) int {
	price := float64(baseValue) * multiplier
	if taxRate > 0 {
		price *= (1 + taxRate)
	}
	return int(math.Ceil(price))
}

// --- NAVIGATION & ACTIONS (SYNCHRONOUS) ---

func (e *GameEngine) Warp(ctx context.Context, playerID uuid.UUID, targetSystemID uuid.UUID) error {
	e.mu.RLock()
	player, ok := e.ActivePlayers[playerID]
	target, sysOk := e.Universe[targetSystemID]
	e.mu.RUnlock()

	if !ok {
		return ErrPlayerNotOnline
	}
	if !sysOk {
		return ErrSystemNotFound
	}
	if player.Ship == nil {
		return errors.New("no ship")
	}

	// Dist Calculation doesn't need locks yet
	dist := entity.CalculateDistance(player.CurrentSystem, target)

	if dist > player.Ship.Stats.JumpRange {
		return fmt.Errorf("target too far (%.1f LY > %.1f LY)", dist, player.Ship.Stats.JumpRange)
	}

	fuelCost := dist * player.Ship.Stats.FuelEfficiency
	
	// Lock for State Mutation
	player.Lock()
	defer player.Unlock()
	player.Ship.Lock()
	defer player.Ship.Unlock()

	if player.Ship.CurrentFuel < fuelCost {
		return errors.New("insufficient fuel")
	}

	player.Ship.CurrentFuel -= fuelCost
	player.CurrentSystem = target

	// Synchronous Save
	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := e.repo.SavePlayer(saveCtx, player); err != nil {
		return fmt.Errorf("jump recorded failed: %w", err)
	}

	if e.OnPlayerUpdate != nil {
		e.OnPlayerUpdate(player)
	}

	return nil
}

func (e *GameEngine) Refuel(ctx context.Context, playerID uuid.UUID) (int, error) {
	e.mu.RLock()
	player, ok := e.ActivePlayers[playerID]
	e.mu.RUnlock()

	if !ok {
		return 0, ErrPlayerNotOnline
	}

	// Lock System (for refueling capability check)
	sys := player.CurrentSystem
	sys.Lock()
	defer sys.Unlock()

	player.Lock()
	defer player.Unlock()
	player.Ship.Lock()
	defer player.Ship.Unlock()

	if !sys.Stats.HasRefueling {
		return 0, errors.New("system has no refueling station")
	}

	missingFuel := player.Ship.Stats.MaxFuel - player.Ship.CurrentFuel
	if missingFuel <= 0 {
		return 0, errors.New("fuel tanks already full")
	}

	basePrice := 1.0
	cost := int(math.Ceil(missingFuel * basePrice * sys.Stats.FuelCostMult))

	if player.Credits < cost {
		return 0, fmt.Errorf("insufficient credits (Need: %d, Have: %d)", cost, player.Credits)
	}

	player.Credits -= cost
	player.Ship.CurrentFuel = player.Ship.Stats.MaxFuel

	// Synchronous Save
	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.repo.SavePlayer(saveCtx, player); err != nil {
		return 0, fmt.Errorf("refuel transaction failed: %w", err)
	}

	if e.OnPlayerUpdate != nil {
		e.OnPlayerUpdate(player)
	}

	return cost, nil
}

func (e *GameEngine) Save(ctx context.Context, playerID uuid.UUID) error {
	e.mu.RLock()
	player, ok := e.ActivePlayers[playerID]
	e.mu.RUnlock()
	if !ok {
		return nil
	}
	
	player.Lock()
	defer player.Unlock()
	return e.repo.SavePlayer(ctx, player)
}

// --- DEBUG / CHEAT ACTIONS (SYNCHRONOUS) ---

func (e *GameEngine) CheatRefuel(ctx context.Context, playerID uuid.UUID) error {
	e.mu.RLock()
	player, ok := e.ActivePlayers[playerID]
	e.mu.RUnlock()
	if !ok {
		return ErrPlayerNotOnline
	}

	player.Ship.Lock()
	player.Ship.CurrentFuel = player.Ship.Stats.MaxFuel
	player.Ship.Unlock()

	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.repo.SavePlayer(saveCtx, player); err != nil {
		return err
	}

	if e.OnPlayerUpdate != nil {
		e.OnPlayerUpdate(player)
	}

	return nil
}

func (e *GameEngine) CheatCredits(ctx context.Context, playerID uuid.UUID, amount int) error {
	e.mu.RLock()
	player, ok := e.ActivePlayers[playerID]
	e.mu.RUnlock()
	if !ok {
		return ErrPlayerNotOnline
	}

	player.Lock()
	player.Credits += amount
	player.Unlock()

	saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.repo.SavePlayer(saveCtx, player); err != nil {
		return err
	}

	if e.OnPlayerUpdate != nil {
		e.OnPlayerUpdate(player)
	}

	return nil
}

// --- READ OPERATIONS ---

func (e *GameEngine) GetPlayer(ctx context.Context, pid uuid.UUID) (*entity.Player, error) {
	e.mu.RLock()
	cached, ok := e.ActivePlayers[pid]
	e.mu.RUnlock()
	
	if ok {
		return cached, nil
	}

	// Not in memory? Check DB.
	p, err := e.repo.GetPlayer(ctx, pid)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("player not found")
	}

	// Re-attach System Pointer if loaded from cold storage
	if p.CurrentSystem != nil {
		if sys, ok := e.Universe[p.CurrentSystem.ID]; ok {
			p.CurrentSystem = sys
		}
	}

	// Cache it
	e.mu.Lock()
	e.ActivePlayers[pid] = p
	e.mu.Unlock()
	
	return p, nil
}

func (e *GameEngine) ScanSystems(centerID uuid.UUID, rangeLY float64) []*entity.System {
	e.mu.RLock()
	center, ok := e.Universe[centerID]
	e.mu.RUnlock()
	if !ok {
		return nil
	}
	var results []*entity.System
	for _, sys := range e.Universe {
		if sys.ID == center.ID {
			continue
		}
		if entity.CalculateDistance(center, sys) <= rangeLY {
			results = append(results, sys)
		}
	}
	return results
}

func (e *GameEngine) GetSystemMarket(ctx context.Context, systemID uuid.UUID) ([]entity.Item, error) {
	// Market data is live on the system struct, but we might want to check the DB
	// For performance, we read from the Repo which likely queries the DB directly
	// or we could return e.Universe[systemID].Market if we trust memory is up to date.
	// Currently, repo.GetSystemMarket queries the DB.
	// If we are doing "Memory Authoritative" for markets, we should read from e.Universe.
	
	e.mu.RLock()
	sys, ok := e.Universe[systemID]
	e.mu.RUnlock()
	
	if ok {
		// Return Memory State (Fastest and consistent with Buy/Sell)
		return sys.Market, nil
	}

	return e.repo.GetSystemMarket(ctx, systemID)
}
