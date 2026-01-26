package entity

import (
	"github.com/google/uuid"
	"galaxies/internal/core/enums"
)

type System struct {
	ID        uuid.UUID
	Name      string
	X         int
	Y         int
	Political enums.PoliticalStatus
	Economic  enums.EconomicStatus
	Social    enums.SocialStatus
	Stats     SystemStats
}

type SystemStats struct {
	// --- GLOBAL MARKET & ECONOMY ---
	MarketBuyMult    float64 `json:"market_buy_mult"`
	MarketSellMult   float64 `json:"market_sell_mult"`
	FuelCostMult     float64 `json:"fuel_cost_mult"`
	RepairCostMult   float64 `json:"repair_cost_mult"`
	TaxRate          float64 `json:"tax_rate"`
	DockingFee       int     `json:"docking_fee"`

	// --- ILLEGAL MARKET ---
	BlackMarketBuyMult  float64 `json:"black_market_buy"`
	BlackMarketSellMult float64 `json:"black_market_sell"`
	ContrabandProfit    float64 `json:"contraband_profit"`

	// --- SHIPYARD & OUTFITTER ---
	ShipCostMult float64 `json:"ship_cost_mult"`
	ModCostMult  float64 `json:"mod_cost_mult"`

	// --- MISSIONS & BOUNTIES ---
	MissionQuantityMult float64 `json:"mission_quantity"`
	MissionPayMult      float64 `json:"mission_pay"`
	BountyPayMult       float64 `json:"bounty_pay"`

	// --- POPULATION MODIFIERS (The Potential) ---
	// These replace the old integer "Quantity" fields
	PassengerDensity float64 `json:"passenger_density"` // Base 1.0
	PassengerWealth  float64 `json:"passenger_wealth"`
	VIPDensity       float64 `json:"vip_density"`       // Base 0.0 (Special)
	VIPWealth        float64 `json:"vip_wealth"`
	SlumsDensity     float64 `json:"slums_density"`     // Base 0.0
	SlumsWealth      float64 `json:"slums_wealth"`
	AndroidDensity   float64 `json:"android_density"`   // Base 0.0
	AndroidCost      int     `json:"android_cost"`
	AndroidSkill     int     `json:"android_skill"`
	PrisonerDensity  float64 `json:"prisoner_density"`  // Base 0.0
	PrisonerCost     int     `json:"prisoner_cost"`
	PrisonerSkill    int     `json:"prisoner_skill"`

	// --- LIVE POPULATION COUNTS (The Reality) ---
	// These are calculated by the Populator
	PassengerCount int `json:"passenger_count"`
	VIPCount       int `json:"vip_count"`
	SlumsCount     int `json:"slums_count"`
	AndroidCount   int `json:"android_count"`
	PrisonerCount  int `json:"prisoner_count"`

	// --- CREW ---
	CrewPoolDensity    float64 `json:"crew_pool_density"` // Replaces Size
	CrewPoolCount      int     `json:"crew_pool_count"`   // Actual number
	CrewSkillAvg       int     `json:"crew_skill_avg"`
	CrewHiringCostMult float64 `json:"crew_hiring_cost_mult"`

	// --- RISK & LAW ---
	PiracyChance     float64 `json:"piracy_chance"`
	InspectionChance float64 `json:"inspection_chance"`
	BribeCostMult    float64 `json:"bribe_cost_mult"`
	WantedPassChance float64 `json:"wanted_pass_chance"`

	// --- BOOLEAN FLAGS ---
	HasShipyard       bool `json:"has_shipyard"`
	HasOutfitter      bool `json:"has_outfitter"`
	HasCantina        bool `json:"has_cantina"`
	HasHospital       bool `json:"has_hospital"`
	HasBlackMarket    bool `json:"has_black_market"`
	HasRefueling      bool `json:"has_refueling"`
	HasMissionBoard   bool `json:"has_mission_board"`
	HasAndroidFoundry bool `json:"has_android_foundry"`
	HasPrison         bool `json:"has_prison"`
	HasLuxuryHousing  bool `json:"has_luxury_housing"`
	HasSlums          bool `json:"has_slums"`
}

func NewDefaultSystemStats() SystemStats {
	return SystemStats{
		// Global Economy
		MarketBuyMult:  1.0,
		MarketSellMult: 1.0,
		FuelCostMult:   1.0,
		RepairCostMult: 1.0,
		TaxRate:        0.05,
		DockingFee:     50,

		// Population Defaults (Potential)
		PassengerDensity: 1.0, // Normal Traffic
		PassengerWealth:  1.0,
		CrewPoolDensity:  1.0, // Normal Crew availability
		CrewSkillAvg:     3,
		CrewHiringCostMult: 1.0,

		// Special Pops (Default to 0 density)
		VIPDensity:      0.0,
		VIPWealth:       5.0,
		SlumsDensity:    0.0,
		SlumsWealth:     0.2,
		AndroidDensity:  0.0,
		PrisonerDensity: 0.0,
		PrisonerCost:    1000,
		PrisonerSkill:   2,
		
		// Base Stats for Specials
		AndroidCost:  5000,
		AndroidSkill: 5,

		// Risk
		PiracyChance:     0.05,
		InspectionChance: 0.10,
		BribeCostMult:    1.0,
		WantedPassChance: 0.02,
		BlackMarketBuyMult: 0.8,
		BlackMarketSellMult: 1.2,
		ContrabandProfit: 1.5,

		// Facilities
		HasRefueling:    true,
		HasCantina:      true,
		HasMissionBoard: true,
		MissionQuantityMult: 1.0,
		MissionPayMult:      1.0,
		BountyPayMult:       1.0,
		ShipCostMult:        1.0,
		ModCostMult:         1.0,
	}
}