package gen

import (
	"math/rand"
	"time"

	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateSystem creates a fully realized world based on the 3 pillars of society.
func GenerateSystem(name string, x, y int, pol domain.PoliticalStatus, eco domain.EconomicStatus, soc domain.SocialStatus) *entity.System {
	// 1. Start with the "Average" baseline
	stats := entity.NewDefaultSystemStats()

	// 2. Handle Name Generation
	// If a name wasn't provided (e.g. "", empty string), generate one.
	if name == "" {
		name = GenerateSystemName()
	}

	// 3. Apply Political Layer (Law, Taxes, Facilities, Risk)
	applyPoliticalMods(&stats, pol)

	// 4. Apply Economic Layer (Prices, Market Activity, Wealth)
	applyEconomicMods(&stats, eco)

	// 5. Apply Social Layer (Population Density, Flavor, Special Facilities)
	applySocialMods(&stats, soc)

	sys := &entity.System{
		ID:        uuid.New(),
		Name:      name,
		X:         x,
		Y:         y,
		Political: pol,
		Economic:  eco,
		Social:    soc,
		Stats:     stats,
		// Market initialized empty, populated below
		Market:    []entity.Item{},
	}

	// 👇 NEW: Populate Market based on the System's Traits
	sys.Market = GenerateMarketStock(sys)

	return sys
}

// =============================================================================
// 1. POLITICAL LAYER (The Rules)
// =============================================================================
func applyPoliticalMods(s *entity.SystemStats, p domain.PoliticalStatus) {
	switch p {
	case domain.PolImperialCore:
		s.TaxRate = 0.15
		s.PiracyChance = 0.0
		s.InspectionChance = 0.50
		s.HasShipyard = true
		s.HasOutfitter = true
		s.VIPDensity = 1.0

	case domain.PolFederationOutpost:
		s.TaxRate = 0.05
		s.PiracyChance = 0.10
		s.InspectionChance = 0.15
		s.HasRefueling = true
		s.HasMissionBoard = true

	case domain.PolMartialLaw:
		s.TaxRate = 0.20
		s.PiracyChance = 0.02
		s.InspectionChance = 0.80
		s.BribeCostMult = 3.0
		s.MissionPayMult = 1.2

	case domain.PolCorporateSovereign:
		s.TaxRate = 0.0
		s.DockingFee = 500
		s.MarketBuyMult = 1.2
		s.HasShipyard = true
		s.HasAndroidFoundry = true

	case domain.PolAnarchicFreehold:
		s.TaxRate = 0.0
		s.DockingFee = 0
		s.PiracyChance = 0.30
		s.InspectionChance = 0.0
		s.HasBlackMarket = true

	case domain.PolPirateHaven:
		s.TaxRate = 0.10
		s.PiracyChance = 0.50
		s.InspectionChance = 0.0
		s.HasBlackMarket = true
		s.BlackMarketSellMult = 1.5
		s.HasShipyard = true

	case domain.PolTheocraticRule:
		s.TaxRate = 0.10
		s.HasCantina = false
		s.MissionPayMult = 0.8
		s.PassengerDensity = 2.0

	case domain.PolBureaucraticGridlock:
		s.DockingFee = 200
		s.TaxRate = 0.12
		s.MarketBuyMult = 1.3
		s.MissionPayMult = 0.7

	case domain.PolContestedWarZone:
		s.PiracyChance = 0.20
		s.RepairCostMult = 2.0
		s.FuelCostMult = 3.0
		s.MissionPayMult = 2.0

	case domain.PolDemilitarizedZone:
		s.PiracyChance = 0.0
		s.InspectionChance = 0.05
		s.HasBlackMarket = true

	case domain.PolPuppetState:
		s.TaxRate = 0.25
		s.MarketSellMult = 0.7
		s.SlumsDensity = 1.0

	case domain.PolSyndicateTerritory:
		s.TaxRate = 0.05
		s.BribeCostMult = 0.5
		s.HasBlackMarket = true
		s.HasCantina = true

	case domain.PolIsolationist:
		s.DockingFee = 2000
		s.MarketBuyMult = 2.0
		s.MarketSellMult = 0.5

	case domain.PolRevolutionaryFront:
		s.PiracyChance = 0.15
		s.HasMissionBoard = true
		s.CrewDensity = 2.0

	case domain.PolColonialCharter:
		s.TaxRate = 0.02
		s.MarketBuyMult = 0.9
		s.HasOutfitter = true

	case domain.PolFailedState:
		s.PiracyChance = 0.40
		s.HasRefueling = false
		s.SlumsDensity = 3.0
		s.HasPrison = false

	case domain.PolAIGovernance:
		s.InspectionChance = 1.0
		s.BribeCostMult = 100.0
		s.MarketBuyMult = 1.0
		s.MarketSellMult = 1.0

	case domain.PolExileColony:
		s.PassengerDensity = 0.1
		s.HasMissionBoard = false
		s.HasRefueling = true

	case domain.PolTradeFedNeutrality:
		s.TaxRate = 0.03
		s.MarketBuyMult = 1.05
		s.MarketSellMult = 0.95
		s.HasShipyard = true
		s.HasOutfitter = true

	case domain.PolFeudalDominion:
		s.TaxRate = 0.30
		s.CrewSkillBonus = 2
		s.VIPDensity = 0.5
		s.SlumsDensity = 2.0
	}
}

// =============================================================================
// 2. ECONOMIC LAYER (The Market)
// =============================================================================
func applyEconomicMods(s *entity.SystemStats, e domain.EconomicStatus) {
	switch e {
	case domain.EcoPostScarcity:
		s.MarketBuyMult *= 0.8
		s.PassengerWealth *= 2.0

	case domain.EcoIndustrialBoom:
		s.MarketSellMult *= 1.2
		s.FuelCostMult *= 0.9
		s.MissionPayMult *= 1.5
		s.HasOutfitter = true

	case domain.EcoDepression:
		s.MarketBuyMult *= 0.6
		s.MarketSellMult *= 0.5
		s.MissionPayMult *= 0.5

	case domain.EcoHyperInflation:
		s.MarketBuyMult *= 5.0
		s.MarketSellMult *= 4.0
		s.FuelCostMult *= 3.0

	case domain.EcoResourceRich:
		s.MarketBuyMult *= 0.8
		s.HasRefueling = true

	case domain.EcoFamine:
		s.MarketBuyMult *= 1.5
		s.PassengerDensity *= 0.5

	case domain.EcoTechBottleneck:
		s.ModCostMult *= 3.0
		s.ShipCostMult *= 2.0
		s.HasAndroidFoundry = false

	case domain.EcoBlackMarketHub:
		s.HasBlackMarket = true
		s.BlackMarketBuyMult = 2.0
		s.PiracyChance += 0.1

	case domain.EcoRefuelingDepot:
		s.FuelCostMult = 0.5
		s.HasRefueling = true

	case domain.EcoLuxuryResort:
		s.DockingFee *= 5
		s.PassengerWealth *= 5.0
		s.VIPDensity += 1.0

	case domain.EcoManufacturingHub:
		s.ShipCostMult *= 0.8
		s.ModCostMult *= 0.8
		s.HasShipyard = true

	case domain.EcoTradeEmbargo:
		s.MarketBuyMult *= 3.0
		s.BlackMarketSellMult *= 2.0

	case domain.EcoGoldRush:
		s.PassengerDensity *= 3.0
		// Hiring cost implies scarcity, could map to higher mission pay
		s.MissionPayMult *= 1.2

	case domain.EcoLaborStrike:
		s.FuelCostMult *= 2.0
		s.RepairCostMult *= 3.0
		s.MissionPayMult *= 0.0

	case domain.EcoWarEconomy:
		s.ShipCostMult *= 2.0
		s.ModCostMult *= 2.0
		s.MarketSellMult *= 1.3

	case domain.EcoDepletedWorld:
		s.MarketBuyMult *= 1.2
		s.FuelCostMult *= 2.0
		s.HasRefueling = false

	case domain.EcoAgrarianBreadbasket:
		s.CrewDensity *= 2.0
		s.PassengerWealth *= 0.8

	case domain.EcoCommandEconomy:
		s.MarketBuyMult = 1.0
		s.MarketSellMult = 1.0
		s.BlackMarketBuyMult = 3.0

	case domain.EcoFreePort:
		s.TaxRate = 0.0
		s.InspectionChance = 0.0
		s.HasBlackMarket = true

	case domain.EcoScavengerEconomy:
		s.ModCostMult *= 0.5
		s.RepairCostMult *= 0.5
		s.ShipCostMult *= 0.6
	}
}

// =============================================================================
// 3. SOCIAL LAYER (The People)
// =============================================================================
func applySocialMods(s *entity.SystemStats, soc domain.SocialStatus) {
	switch soc {
	case domain.SocCosmopolitan:
		s.PassengerDensity *= 2.0
		s.CrewDensity *= 1.5
		s.HasCantina = true

	case domain.SocXenophobic:
		s.DockingFee *= 2
		s.CrewDensity *= 0.1
		s.InspectionChance += 0.2

	case domain.SocReligiousPilgrimage:
		s.PassengerDensity *= 3.0
		s.PassengerWealth *= 0.5
		s.HasHospital = true

	case domain.SocPlague:
		s.PassengerDensity = 0.0
		s.DockingFee = 0
		s.HasHospital = true
		s.MissionPayMult *= 2.0

	case domain.SocBrainDrain:
		s.CrewSkillBonus = -2
		s.ModCostMult *= 1.5

	case domain.SocPenalColony:
		s.PassengerDensity = 0.0
		s.HasPrison = true
		s.SlumsDensity = 2.0

	case domain.SocRefugeeCrisis:
		s.PassengerDensity *= 5.0
		s.PassengerWealth = 0.1
		s.SlumsDensity = 3.0

	case domain.SocArtistEnclave:
		s.VIPDensity += 0.5
		s.PassengerWealth *= 1.5

	case domain.SocCyberneticAscension:
		s.HasAndroidFoundry = true
		s.CrewSkillBonus += 3
		s.PassengerDensity *= 0.5

	case domain.SocPreIndustrial:
		s.HasRefueling = false
		s.HasShipyard = false
		s.MarketBuyMult *= 0.5
		s.MarketSellMult *= 0.5

	case domain.SocAcademy:
		s.CrewSkillBonus += 4
		s.MissionPayMult *= 1.1

	case domain.SocSlum:
		s.SlumsDensity = 4.0
		s.PiracyChance += 0.1
		s.CrewDensity *= 3.0

	case domain.SocFrontierSpirit:
		s.CrewDensity *= 1.5
		s.PiracyChance += 0.05
		s.HasMissionBoard = true

	case domain.SocDecadentAristocracy:
		s.VIPDensity = 2.0
		s.VIPWealth = 10.0
		s.MarketBuyMult *= 1.5

	case domain.SocWorkerRebellion:
		s.HasShipyard = false
		s.HasRefueling = false
		s.PiracyChance += 0.2

	case domain.SocCultActivity:
		s.PassengerDensity *= 0.5
		s.MissionPayMult *= 1.2
		s.VIPDensity = 0.2

	case domain.SocScientificExpedition:
		s.HasRefueling = true
		s.PassengerDensity = 0.2
		s.MissionPayMult *= 1.5

	case domain.SocGhostTown:
		s.PassengerDensity = 0.0
		s.CrewDensity = 0.0
		s.HasCantina = false

	case domain.SocGladiatorial:
		s.CrewSkillBonus += 2
		s.MissionPayMult *= 1.2
		s.VIPDensity += 0.5

	case domain.SocHiveMind:
		s.CrewDensity = 0.0
		s.MarketBuyMult = 1.0
		s.InspectionChance = 1.0
	}
}