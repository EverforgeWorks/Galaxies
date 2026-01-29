package entity

import (
	"encoding/json"
	"math"
	"sync"
	"github.com/google/uuid"
	"galaxies/internal/core/domain"
)

type System struct {
	mu sync.Mutex
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	
	X int `json:"x"`
	Y int `json:"y"`

	Political domain.PoliticalStatus `json:"-"` // Ignored in default marshal (handled by custom)
	Economic  domain.EconomicStatus  `json:"-"`
	Social    domain.SocialStatus    `json:"-"`

	Stats SystemStats `json:"stats"`
	Market []Item	  `json:"market"`
}

// Custom MarshalJSON to convert Enums (Ints) to Strings for the Client
func (s *System) MarshalJSON() ([]byte, error) {
	type Alias System
	return json.Marshal(&struct {
		Political string `json:"political"`
		Economic  string `json:"economic"`
		Social    string `json:"social"`
		*Alias
	}{
		Political: s.Political.String(),
		Economic:  s.Economic.String(),
		Social:    s.Social.String(),
		Alias:     (*Alias)(s),
	})
}

// SystemStats defines all the modifiers for the system
type SystemStats struct {
	// --- ECONOMY & FEES ---
	MarketBuyMult  float64 `json:"market_buy_mult"`  
	MarketSellMult float64 `json:"market_sell_mult"` 
	
	FuelCostMult   float64 `json:"fuel_cost_mult"`   
	RepairCostMult float64 `json:"repair_cost_mult"` 
	DockingFee     int     `json:"docking_fee"`      
	TaxRate        float64 `json:"tax_rate"`         

	// --- ILLEGAL ACTIVITY ---
	BlackMarketBuyMult  float64 `json:"black_market_buy"`
	BlackMarketSellMult float64 `json:"black_market_sell"`
	
	PiracyChance     float64 `json:"piracy_chance"`     
	InspectionChance float64 `json:"inspection_chance"` 
	BribeCostMult    float64 `json:"bribe_cost_mult"`   

	// --- OPPORTUNITIES ---
	MissionPayMult float64 `json:"mission_pay"` 
	ShipCostMult   float64 `json:"ship_cost_mult"`
	ModCostMult    float64 `json:"mod_cost_mult"`

	// --- POPULATION ---
	PassengerDensity float64 `json:"passenger_density"`
	PassengerWealth  float64 `json:"passenger_wealth"`
	
	VIPDensity       float64 `json:"vip_density"`
	VIPWealth        float64 `json:"vip_wealth"`
	SlumsDensity     float64 `json:"slums_density"` 
	
	CrewDensity      float64 `json:"crew_density"`     
	CrewSkillBonus   int     `json:"crew_skill_bonus"` 

	// --- FACILITIES ---
	HasShipyard       bool `json:"has_shipyard"`
	HasOutfitter      bool `json:"has_outfitter"`
	HasRefueling      bool `json:"has_refueling"`
	HasBlackMarket    bool `json:"has_black_market"`
	HasMissionBoard   bool `json:"has_mission_board"`
	HasCantina        bool `json:"has_cantina"`       
	HasHospital       bool `json:"has_hospital"`      
	HasPrison         bool `json:"has_prison"`        
	HasAndroidFoundry bool `json:"has_android_foundry"` 
}

func NewDefaultSystemStats() SystemStats {
	return SystemStats{
		MarketBuyMult:  1.1,  
		MarketSellMult: 0.9,  
		FuelCostMult:   1.0,  
		RepairCostMult: 1.0,  
		DockingFee:     100,  
		TaxRate:        0.05, 

		BlackMarketBuyMult:  1.5, 
		BlackMarketSellMult: 0.6, 
		PiracyChance:        0.05, 
		InspectionChance:    0.10, 
		BribeCostMult:       1.0,  

		MissionPayMult: 1.0,
		ShipCostMult:   1.0,
		ModCostMult:    1.0,

		PassengerDensity: 1.0, 
		PassengerWealth:  1.0, 
		
		VIPDensity:       0.0, 
		VIPWealth:        5.0, 
		SlumsDensity:     0.0, 

		CrewDensity:    1.0,
		CrewSkillBonus: 1,   

		HasRefueling:    true,
		HasMissionBoard: true,
		HasCantina:      true,

		HasShipyard:       false,
		HasOutfitter:      false,
		HasBlackMarket:    false,
		HasHospital:       false,
		HasPrison:         false,
		HasAndroidFoundry: false,
	}
}

// CalculateDistance returns the Euclidean distance between two systems.
func CalculateDistance(a, b *System) float64 {
	return math.Sqrt(math.Pow(float64(b.X-a.X), 2) + math.Pow(float64(b.Y-a.Y), 2))
}

func (s *System) Lock()   { s.mu.Lock() }
func (s *System) Unlock() { s.mu.Unlock() }
