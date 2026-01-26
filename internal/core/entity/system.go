package entity

import (
	"github.com/google/uuid"
	"galaxies/internal/core/domain"
)

type System struct {
	ID        uuid.UUID
	Name      string
	
	// Coordinates (Light Years or Grid Units)
	X int
	Y int

	// The "DNA" of the system - determining how stats are generated
	Political enums.PoliticalStatus // e.g., "Anarchy", "Democracy"
	Economic  enums.EconomicStatus  // e.g., "Industrial", "Agricultural"
	Social    enums.SocialStatus    // e.g., "Feudal", "Utopian"

	// The computed gameplay modifiers
	Stats SystemStats `json:"stats"`
}

type SystemStats struct {
	// --- ECONOMY & FEES ---
	// Multipliers applied to the Global Base Price of items.
	// Buy = Station sells to player. Sell = Station buys from player.
	MarketBuyMult  float64 `json:"market_buy_mult"`  
	MarketSellMult float64 `json:"market_sell_mult"` 
	
	FuelCostMult   float64 `json:"fuel_cost_mult"`   // Multiplier on global fuel price
	RepairCostMult float64 `json:"repair_cost_mult"` // Multiplier on hull repair/dock services
	DockingFee     int     `json:"docking_fee"`      // Flat credit fee to land
	TaxRate        float64 `json:"tax_rate"`         // % taken from market transactions

	// --- ILLEGAL ACTIVITY ---
	// If BlackMarket exists, these modifiers apply to "Contraband" items.
	BlackMarketBuyMult  float64 `json:"black_market_buy"`
	BlackMarketSellMult float64 `json:"black_market_sell"`
	
	// Risk: 0.0 to 1.0
	PiracyChance     float64 `json:"piracy_chance"`     // Chance of interdiction upon entering
	InspectionChance float64 `json:"inspection_chance"` // Chance of police scan upon docking
	BribeCostMult    float64 `json:"bribe_cost_mult"`   // Multiplier for avoiding fines

	// --- OPPORTUNITIES (Missions & Shipyard) ---
	MissionPayMult float64 `json:"mission_pay"` // "Rich" systems pay more
	ShipCostMult   float64 `json:"ship_cost_mult"`
	ModCostMult    float64 `json:"mod_cost_mult"`

	// --- POPULATION GENERATORS ---
	// "Density" is a multiplier for the number of people generated.
	// "Wealth" is a multiplier for the fares/rewards they offer.
	PassengerDensity float64 `json:"passenger_density"`
	PassengerWealth  float64 `json:"passenger_wealth"`
	
	VIPDensity       float64 `json:"vip_density"`
	VIPWealth        float64 `json:"vip_wealth"`
	
	SlumsDensity     float64 `json:"slums_density"` // Refugees, criminals, cheap labor
	
	// Crew Availability
	CrewDensity      float64 `json:"crew_density"`     // How many recruits appear
	CrewSkillBonus   int     `json:"crew_skill_bonus"` // e.g. +2 levels for "Military" systems

	// --- FACILITIES (Capabilities) ---
	HasShipyard       bool `json:"has_shipyard"`
	HasOutfitter      bool `json:"has_outfitter"`
	HasRefueling      bool `json:"has_refueling"`
	HasBlackMarket    bool `json:"has_black_market"`
	HasMissionBoard   bool `json:"has_mission_board"`
	HasCantina        bool `json:"has_cantina"`       // Recruits Crew
	HasHospital       bool `json:"has_hospital"`      // Heals Character/Crew
	HasPrison         bool `json:"has_prison"`        // Drops off Prisoner passengers
	HasAndroidFoundry bool `json:"has_android_foundry"` // Sells Android Crew
}

func NewDefaultSystemStats() SystemStats {
	return SystemStats{
		// --- ECONOMY & FEES ---
		MarketBuyMult:  1.1,  // Station sells items at 110% of global average (Markup)
		MarketSellMult: 0.9,  // Station buys items at 90% of global average (Spread)
		FuelCostMult:   1.0,  // Standard fuel prices
		RepairCostMult: 1.0,  // Standard repair labor costs
		DockingFee:     100,  // Standard parking ticket
		TaxRate:        0.05, // 5% sales tax on legal transactions

		// --- ILLEGAL ACTIVITY ---
		BlackMarketBuyMult:  1.5, // Buying contraband is expensive (Risk premium)
		BlackMarketSellMult: 0.6, // Fences pay low (they take the risk)
		PiracyChance:        0.05, // 5% chance of ambush when entering system
		InspectionChance:    0.10, // 10% chance of police scan on docking
		BribeCostMult:       1.0,  // Standard bribery rates

		// --- OPPORTUNITIES ---
		MissionPayMult: 1.0,
		ShipCostMult:   1.0,
		ModCostMult:    1.0,

		// --- POPULATION GENERATORS ---
		PassengerDensity: 1.0, // Standard crowd size
		PassengerWealth:  1.0, // Standard fares
		
		// Specials default to 0.0 (Must be enabled by Generators)
		VIPDensity:       0.0, 
		VIPWealth:        5.0, // If they theoretically existed, they'd pay 5x
		SlumsDensity:     0.0, 

		// --- CREW ---
		CrewDensity:    1.0,
		CrewSkillBonus: 1,   // Standard rookies

		// --- FACILITIES (The Baseline Standard) ---
		// Most settled systems have these basics:
		HasRefueling:    true,
		HasMissionBoard: true,
		HasCantina:      true,

		// These are "Premium" facilities, false by default:
		HasShipyard:       false,
		HasOutfitter:      false,
		HasBlackMarket:    false,
		HasHospital:       false,
		HasPrison:         false,
		HasAndroidFoundry: false,
	}
}