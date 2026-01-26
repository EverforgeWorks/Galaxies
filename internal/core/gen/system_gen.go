package gen

import (
	"math/rand"
	"time"

	"galaxies/internal/core/entity"
	"galaxies/internal/core/domain"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateSystem creates a fully realized world based on the 3 pillars of society.
func GenerateSystem(name string, x, y int, pol enums.PoliticalStatus, eco enums.EconomicStatus, soc enums.SocialStatus) *entity.System {
	// 1. Start with the "Average" baseline
	stats := entity.NewDefaultSystemStats()

	// 2. Apply Political Layer (Law, Taxes, Facilities, Risk)
	applyPoliticalMods(&stats, pol)

	// 3. Apply Economic Layer (Prices, Market Activity, Wealth)
	applyEconomicMods(&stats, eco)

	// 4. Apply Social Layer (Population Density, Flavor, Special Facilities)
	applySocialMods(&stats, soc)

	return &entity.System{
		ID:        uuid.New(),
		Name:      name,
		X:         x,
		Y:         y,
		Political: pol,
		Economic:  eco,
		Social:    soc,
		Stats:     stats,
	}
}

// =============================================================================
// 1. POLITICAL LAYER (The Rules)
// Controls: Taxes, Security, Facilities, Fees
// =============================================================================
func applyPoliticalMods(s *entity.SystemStats, p enums.PoliticalStatus) {
	switch p {
	case enums.PolImperialCore:
		s.TaxRate = 0.15          // High taxes
		s.PiracyChance = 0.0      // Very safe
		s.InspectionChance = 0.50 // Police state
		s.HasShipyard = true
		s.HasOutfitter = true
		s.VIPDensity = 1.0 // Nobles live here

	case enums.PolFederationOutpost:
		s.TaxRate = 0.05
		s.PiracyChance = 0.10
		s.InspectionChance = 0.15
		s.HasRefueling = true
		s.HasMissionBoard = true

	case enums.PolMartialLaw:
		s.TaxRate = 0.20
		s.PiracyChance = 0.02
		s.InspectionChance = 0.80 // Almost guaranteed scan
		s.BribeCostMult = 3.0     // Soldiers are expensive to bribe
		s.MissionPayMult = 1.2    // Military contracts pay well

	case enums.PolCorporateSovereign:
		s.TaxRate = 0.0 // No tax, but...
		s.DockingFee = 500 // High usage fees
		s.MarketBuyMult = 1.2
		s.HasShipyard = true
		s.HasAndroidFoundry = true

	case enums.PolAnarchicFreehold:
		s.TaxRate = 0.0
		s.DockingFee = 0
		s.PiracyChance = 0.30
		s.InspectionChance = 0.0
		s.HasBlackMarket = true

	case enums.PolPirateHaven:
		s.TaxRate = 0.10 (Protection money)
		s.PiracyChance = 0.50     // Dangerous to enter
		s.InspectionChance = 0.0  // No police
		s.HasBlackMarket = true
		s.BlackMarketSellMult = 1.5 // Best place to sell contraband
		s.HasShipyard = true      // Chop shops

	case enums.PolTheocraticRule:
		s.TaxRate = 0.10 (Tithes)
		s.HasCantina = false      // Alcohol forbidden (Logic handled in event generation)
		s.MissionPayMult = 0.8
		s.PassengerDensity = 2.0  // Pilgrims

	case enums.PolBureaucraticGridlock:
		s.DockingFee = 200
		s.TaxRate = 0.12
		s.MarketBuyMult = 1.3     // Inefficiency increases prices
		s.MissionPayMult = 0.7    // Slow payments

	case enums.PolContestedWarZone:
		s.PiracyChance = 0.20     // Opportunists
		s.RepairCostMult = 2.0    // Mechanics are scarce
		s.FuelCostMult = 3.0      // Fuel is rationed
		s.MissionPayMult = 2.0    // High risk, high reward

	case enums.PolDemilitarizedZone:
		s.PiracyChance = 0.0
		s.InspectionChance = 0.05
		s.HasBlackMarket = true   // Smugglers thrive in the gray area

	case enums.PolPuppetState:
		s.TaxRate = 0.25 (Siphoned off)
		s.MarketSellMult = 0.7
		s.SlumsDensity = 1.0

	case enums.PolSyndicateTerritory:
		s.TaxRate = 0.05
		s.BribeCostMult = 0.5     // Corruption is standard
		s.HasBlackMarket = true
		s.HasCantina = true

	case enums.PolIsolationist:
		s.DockingFee = 2000       // Go away
		s.MarketBuyMult = 2.0     // Imports rare
		s.MarketSellMult = 0.5    // Exports unwanted

	case enums.PolRevolutionaryFront:
		s.PiracyChance = 0.15
		s.HasMissionBoard = true
		s.CrewDensity = 2.0       // Idealists joining up

	case enums.PolColonialCharter:
		s.TaxRate = 0.02          // Incentivized
		s.MarketBuyMult = 0.9     // Subsidized goods
		s.HasOutfitter = true

	case enums.PolFailedState:
		s.PiracyChance = 0.40
		s.HasRefueling = false    // Infrastructure collapsed!
		s.SlumsDensity = 3.0
		s.HasPrison = false       // Prisoners escaped

	case enums.PolAIGovernance:
		s.InspectionChance = 1.0  // 100% scan, instant processing
		s.BribeCostMult = 100.0   // Cannot bribe logic
		s.MarketBuyMult = 1.0     // Perfectly efficient
		s.MarketSellMult = 1.0    // Zero spread

	case enums.PolExileColony:
		s.PassengerDensity = 0.1
		s.HasMissionBoard = false
		s.HasRefueling = true

	case enums.PolTradeFedNeutrality:
		s.TaxRate = 0.03
		s.MarketBuyMult = 1.05
		s.MarketSellMult = 0.95   // Best trade spreads
		s.HasShipyard = true
		s.HasOutfitter = true

	case enums.PolFeudalDominion:
		s.TaxRate = 0.30
		s.CrewSkillBonus = 2      // Warrior caste
		s.VIPDensity = 0.5
		s.SlumsDensity = 2.0
	}
}

// =============================================================================
// 2. ECONOMIC LAYER (The Market)
// Controls: Prices, Multipliers, Wealth
// =============================================================================
func applyEconomicMods(s *entity.SystemStats, e enums.EconomicStatus) {
	switch e {
	case enums.EcoPostScarcity:
		s.MarketBuyMult *= 0.8    // Cheap goods
		s.StandardOfLiving = "High" // Conceptual
		s.PassengerWealth *= 2.0

	case enums.EcoIndustrialBoom:
		s.MarketSellMult *= 1.2   // They need resources BADLY
		s.FuelCostMult *= 0.9
		s.MissionPayMult *= 1.5
		s.HasOutfitter = true

	case enums.EcoDepression:
		s.MarketBuyMult *= 0.6    // Liquidation prices
		s.MarketSellMult *= 0.5   // Nobody has money to buy your stuff
		s.MissionPayMult *= 0.5

	case enums.EcoHyperInflation:
		s.MarketBuyMult *= 5.0    // Prices skyrocketing
		s.MarketSellMult *= 4.0
		s.FuelCostMult *= 3.0

	case enums.EcoResourceRich:
		s.MarketBuyMult *= 0.8    // Raw mats cheap
		s.HasRefueling = true

	case enums.EcoFamine:
		// Food logic handled by Item Categories elsewhere
		s.MarketBuyMult *= 1.5    // General desperation
		s.PassengerDensity *= 0.5 // People dying

	case enums.EcoTechBottleneck:
		s.ModCostMult *= 3.0      // Tech is expensive
		s.ShipCostMult *= 2.0
		s.HasAndroidFoundry = false

	case enums.EcoBlackMarketHub:
		s.HasBlackMarket = true
		s.BlackMarketBuyMult = 2.0 // They pay huge for contraband
		s.PiracyChance += 0.1

	case enums.EcoRefuelingDepot:
		s.FuelCostMult = 0.5      // Dirt cheap fuel
		s.HasRefueling = true

	case enums.EcoLuxuryResort:
		s.DockingFee *= 5
		s.PassengerWealth *= 5.0
		s.VIPDensity += 1.0

	case enums.EcoManufacturingHub:
		s.ShipCostMult *= 0.8     // Ships cheap here
		s.ModCostMult *= 0.8
		s.HasShipyard = true

	case enums.EcoTradeEmbargo:
		s.MarketBuyMult *= 3.0    // Smuggler's paradise
		s.BlackMarketSellMult *= 2.0

	case enums.EcoGoldRush:
		s.PassengerDensity *= 3.0 // Everyone going there
		s.CrewHiringCostMult *= 2.0 // Labor shortage

	case enums.EcoLaborStrike:
		s.RefuelCostMult *= 2.0   // Scabs are expensive
		s.RepairCostMult *= 3.0
		s.MissionPayMult *= 0.0

	case enums.EcoWarEconomy:
		s.ShipCostMult *= 2.0     // Military buying everything
		s.ModCostMult *= 2.0
		s.MarketSellMult *= 1.3

	case enums.EcoDepletedWorld:
		s.MarketBuyMult *= 1.2
		s.FuelCostMult *= 2.0
		s.HasRefueling = false    // Maybe?

	case enums.EcoAgrarianBreadbasket:
		s.CrewDensity *= 2.0      // Farm boys looking for adventure
		s.PassengerWealth *= 0.8

	case enums.EcoCommandEconomy:
		s.MarketBuyMult = 1.0     // Fixed prices
		s.MarketSellMult = 1.0
		s.BlackMarketBuyMult = 3.0 // Illegal goods very valuable

	case enums.EcoFreePort:
		s.TaxRate = 0.0
		s.InspectionChance = 0.0
		s.HasBlackMarket = true

	case enums.EcoScavengerEconomy:
		s.ModCostMult *= 0.5      // Used parts cheap
		s.RepairCostMult *= 0.5   // Janky repairs
		s.ShipCostMult *= 0.6
	}
}

// =============================================================================
// 3. SOCIAL LAYER (The People)
// Controls: Population, Crew, Flavor, Special Facilities
// =============================================================================
func applySocialMods(s *entity.SystemStats, soc enums.SocialStatus) {
	switch soc {
	case enums.SocCosmopolitan:
		s.PassengerDensity *= 2.0
		s.CrewDensity *= 1.5
		s.HasCantina = true

	case enums.SocXenophobic:
		s.DockingFee *= 2
		s.CrewDensity *= 0.1      // Won't work for outsiders
		s.InspectionChance += 0.2

	case enums.SocReligiousPilgrimage:
		s.PassengerDensity *= 3.0
		s.PassengerWealth *= 0.5  // Poor pilgrims
		s.HasHospital = true      // Healing shrines

	case enums.SocPlague:
		s.PassengerDensity = 0.0
		s.DockingFee = 0          // Please help us
		s.HasHospital = true
		s.MissionPayMult *= 2.0   // Hazard pay

	case enums.SocBrainDrain:
		s.CrewSkillBonus = -2
		s.ModCostMult *= 1.5      // No engineers to install stuff

	case enums.SocPenalColony:
		s.PassengerDensity = 0.0
		s.HasPrison = true
		s.SlumsDensity = 2.0

	case enums.SocRefugeeCrisis:
		s.PassengerDensity *= 5.0
		s.PassengerWealth = 0.1   // Broke
		s.SlumsDensity = 3.0

	case enums.SocArtistEnclave:
		s.VIPDensity += 0.5
		s.PassengerWealth *= 1.5

	case enums.SocCyberneticAscension:
		s.HasAndroidFoundry = true
		s.CrewSkillBonus += 3
		s.PassengerDensity *= 0.5

	case enums.SocPreIndustrial:
		s.HasRefueling = false
		s.HasShipyard = false
		s.MarketBuyMult *= 0.5    // Raw goods cheap
		s.MarketSellMult *= 0.5   // Can't afford tech

	case enums.SocAcademy:
		s.CrewSkillBonus += 4     // Experts
		s.CrewHiringCostMult *= 2.0

	case enums.SocSlum:
		s.SlumsDensity = 4.0
		s.PiracyChance += 0.1
		s.CrewDensity *= 3.0      // Desperate for work

	case enums.SocFrontierSpirit:
		s.CrewDensity *= 1.5
		s.PiracyChance += 0.05
		s.HasMissionBoard = true

	case enums.SocDecadentAristocracy:
		s.VIPDensity = 2.0
		s.VIPWealth = 10.0
		s.MarketBuyMult *= 1.5    // Everything overpriced

	case enums.SocWorkerRebellion:
		s.HasShipyard = false     // On strike
		s.HasRefueling = false
		s.PiracyChance += 0.2

	case enums.SocCultActivity:
		s.PassengerDensity *= 0.5
		s.MissionPayMult *= 1.2   // Weird jobs
		s.VIPDensity = 0.2        // Cult leaders

	case enums.SocScientificExpedition:
		s.HasRefueling = true
		s.PassengerDensity = 0.2
		s.MissionPayMult *= 1.5   // Research grants

	case enums.SocGhostTown:
		s.PassengerDensity = 0.0
		s.CrewDensity = 0.0
		s.HasCantina = false

	case enums.SocGladiatorial:
		s.CrewSkillBonus += 2     // Fighters
		s.MissionPayMult *= 1.2
		s.VIPDensity += 0.5       // Spectators

	case enums.SocHiveMind:
		s.CrewDensity = 0.0       // No individuals
		s.MarketBuyMult = 1.0
		s.InspectionChance = 1.0  // We see all
	}
}