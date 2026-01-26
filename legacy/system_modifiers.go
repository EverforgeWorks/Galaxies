package gen

import (
	"galaxies/internal/core/entity"
	"galaxies/internal/core/enums"
)

// --- POLITICAL MODIFIERS ---
func ApplyPoliticalMods(s *entity.SystemStats, p enums.PoliticalStatus) {
	switch p {
	case enums.PolImperialCore:
		s.TaxRate += 0.20
		s.InspectionChance += 0.80
		s.PiracyChance = 0.0
		s.ShipCostMult *= 1.2
		s.BountyPayMult *= 2.0
		s.HasShipyard = true
		s.HasOutfitter = true

	case enums.PolFederationOutpost:
		s.TaxRate += 0.10
		s.InspectionChance += 0.30
		s.PiracyChance -= 0.02
		s.HasRefueling = true

	case enums.PolMartialLaw:
		s.InspectionChance = 1.0
		s.TaxRate += 0.15
		s.PassengerDensity *= 0.5 // Reduced traffic
		s.MarketBuyMult *= 0.8
		s.ContrabandProfit *= 0.5

	case enums.PolCorporateSovereign:
		s.TaxRate = 0.0
		s.DockingFee = 500
		s.MarketBuyMult *= 1.1
		s.HasOutfitter = true

	case enums.PolAnarchicFreehold:
		s.TaxRate = 0.0
		s.InspectionChance = 0.0
		s.PiracyChance += 0.30
		s.ContrabandProfit *= 1.2
		s.HasMissionBoard = false

	case enums.PolPirateHaven:
		s.TaxRate = 0.0
		s.InspectionChance = 0.0
		s.PiracyChance += 0.50
		s.HasBlackMarket = true
		s.BlackMarketBuyMult = 1.2
		s.DockingFee = 200
		s.BountyPayMult = 0.0

	case enums.PolTheocraticRule:
		s.ContrabandProfit *= 0.8
		s.PassengerWealth *= 0.8
		s.MissionPayMult *= 1.2
		s.HasCantina = false

	case enums.PolBureaucraticGridlock:
		s.DockingFee += 100
		s.MarketSellMult *= 1.1
		s.MissionQuantityMult *= 0.5
		s.BribeCostMult *= 0.5

	case enums.PolContestedWarZone:
		s.PiracyChance += 0.20
		s.FuelCostMult *= 2.0
		s.RepairCostMult *= 2.0
		s.MissionPayMult *= 2.5

	case enums.PolDemilitarizedZone:
		s.InspectionChance += 0.50
		s.ShipCostMult *= 2.0
		s.ContrabandProfit *= 1.5

	case enums.PolPuppetState:
		s.TaxRate += 0.15
		s.MarketBuyMult *= 0.9

	case enums.PolSyndicateTerritory:
		s.HasBlackMarket = true
		s.PiracyChance = 0.05
		s.DockingFee = 300
		s.ContrabandProfit *= 1.1

	case enums.PolIsolationist:
		s.DockingFee = 1000
		s.InspectionChance += 0.50

	case enums.PolRevolutionaryFront:
		s.PiracyChance += 0.15
		s.PassengerDensity += 0.5 // Refugees
		s.HasMissionBoard = true

	case enums.PolColonialCharter:
		s.TaxRate = 0.0
		s.MarketBuyMult *= 1.2
		s.ShipCostMult *= 0.8
		s.CrewHiringCostMult *= 0.5

	case enums.PolFailedState:
		s.HasRefueling = false
		s.HasShipyard = false
		s.PiracyChance += 0.40
		s.MarketSellMult *= 0.5

	case enums.PolAIGovernance:
		s.InspectionChance = 1.0
		s.BribeCostMult = 999.0
		s.TaxRate = 0.10
		s.MarketSellMult = 1.0

	case enums.PolExileColony:
		s.HasPrison = true
		s.PrisonerDensity += 1.0 // Add prisoners
		s.PassengerDensity = 0.0
		s.InspectionChance += 0.60

	case enums.PolTradeFedNeutrality:
		s.TaxRate = 0.02
		s.MarketBuyMult *= 1.05
		s.MarketSellMult *= 0.95

	case enums.PolFeudalDominion:
		s.TaxRate += 0.30
		s.CrewHiringCostMult *= 0.5
		s.VIPDensity += 0.5 // Add Lords/Ladies
		s.HasLuxuryHousing = true
	}
}

// --- ECONOMIC MODIFIERS ---
func ApplyEconomicMods(s *entity.SystemStats, e enums.EconomicStatus) {
	switch e {
	case enums.EcoPostScarcity:
		s.MarketSellMult *= 0.4
		s.MarketBuyMult *= 0.4
		s.PassengerWealth *= 3.0
		s.HasLuxuryHousing = true
		s.VIPDensity += 0.5

	case enums.EcoIndustrialBoom:
		s.MarketBuyMult *= 1.1
		s.HasShipyard = true
		s.ShipCostMult *= 0.9

	case enums.EcoDepression:
		s.MarketBuyMult *= 0.6
		s.MarketSellMult *= 0.6
		s.PassengerDensity += 1.0 // People leaving
		s.CrewHiringCostMult *= 0.5

	case enums.EcoHyperInflation:
		s.MarketBuyMult *= 5.0
		s.MarketSellMult *= 5.0
		s.FuelCostMult *= 5.0
		s.RepairCostMult *= 5.0

	case enums.EcoResourceRich:
		s.MarketSellMult *= 0.7
		s.HasRefueling = true

	case enums.EcoFamine:
		s.MarketBuyMult *= 1.5
		s.PassengerDensity += 0.5
		s.HasSlums = true
		s.SlumsDensity += 1.0

	case enums.EcoTechBottleneck:
		s.ModCostMult *= 2.0
		s.ShipCostMult *= 1.5
		s.AndroidCost *= 2

	case enums.EcoBlackMarketHub:
		s.HasBlackMarket = true
		s.BlackMarketBuyMult = 1.3
		s.ContrabandProfit *= 1.5
		s.InspectionChance -= 0.20

	case enums.EcoRefuelingDepot:
		s.HasRefueling = true
		s.FuelCostMult *= 0.5

	case enums.EcoLuxuryResort:
		s.DockingFee += 400
		s.VIPDensity += 2.0 // Lots of VIPs
		s.HasLuxuryHousing = true
		s.MarketSellMult *= 2.0

	case enums.EcoManufacturingHub:
		s.ShipCostMult *= 0.8
		s.ModCostMult *= 0.8
		s.HasOutfitter = true
		s.HasShipyard = true

	case enums.EcoTradeEmbargo:
		s.MarketSellMult *= 3.0
		s.FuelCostMult *= 2.0

	case enums.EcoGoldRush:
		s.DockingFee += 100
		s.MarketBuyMult *= 1.3
		s.FuelCostMult *= 2.0

	case enums.EcoLaborStrike:
		s.MissionQuantityMult *= 0.1
		s.CrewHiringCostMult *= 3.0

	case enums.EcoWarEconomy:
		s.ShipCostMult *= 1.5
		s.ModCostMult *= 1.5
		s.MissionPayMult *= 1.5

	case enums.EcoDepletedWorld:
		s.MarketSellMult *= 1.5
		s.FuelCostMult *= 2.0
		s.PassengerDensity += 0.5

	case enums.EcoAgrarianBreadbasket:
		s.MarketSellMult *= 0.8
		s.HasRefueling = true

	case enums.EcoCommandEconomy:
		s.MarketBuyMult = 1.0
		s.MarketSellMult = 1.0
		s.TaxRate = 0.0

	case enums.EcoFreePort:
		s.TaxRate = 0.0
		s.DockingFee = 0
		s.InspectionChance = 0.0

	case enums.EcoScavengerEconomy:
		s.ModCostMult *= 0.5
		s.ShipCostMult *= 0.6
		s.RepairCostMult *= 0.5
	}
}

// --- SOCIAL MODIFIERS ---
func ApplySocialMods(s *entity.SystemStats, soc enums.SocialStatus) {
	switch soc {
	case enums.SocCosmopolitan:
		s.PassengerDensity *= 2.0
		s.CrewPoolDensity += 1.0
		s.CrewSkillAvg = 5
		s.HasCantina = true

	case enums.SocXenophobic:
		s.PassengerDensity = 0.0
		s.CrewPoolDensity = 0.0
		s.DockingFee *= 5
		s.MarketBuyMult *= 0.5

	case enums.SocReligiousPilgrimage:
		s.PassengerDensity *= 5.0
		s.PassengerWealth *= 0.5
		s.MissionQuantityMult *= 2.0

	case enums.SocPlague:
		s.HasHospital = true
		s.PassengerDensity = 0.0
		s.DockingFee = 0
		s.MarketBuyMult *= 2.0

	case enums.SocBrainDrain:
		s.CrewSkillAvg = 8
		s.CrewHiringCostMult *= 0.5
		s.CrewPoolDensity += 0.5

	case enums.SocPenalColony:
		s.HasPrison = true
		s.PrisonerDensity += 3.0 // High density
		s.CrewPoolDensity = 0.0
		s.InspectionChance += 0.20

	case enums.SocRefugeeCrisis:
		s.PassengerDensity *= 4.0
		s.PassengerWealth = 0.1
		s.HasSlums = true
		s.SlumsDensity += 4.0

	case enums.SocArtistEnclave:
		s.PassengerWealth *= 1.5
		s.VIPDensity += 0.5
		s.HasLuxuryHousing = true
		s.CrewSkillAvg = 4

	case enums.SocCyberneticAscension:
		s.HasAndroidFoundry = true
		s.AndroidDensity += 1.5
		s.AndroidSkill = 9
		s.CrewPoolDensity = 0.0

	case enums.SocPreIndustrial:
		s.MarketBuyMult *= 0.5
		s.HasShipyard = false
		s.HasRefueling = false

	case enums.SocAcademy:
		s.CrewSkillAvg = 9
		s.CrewHiringCostMult *= 2.0
		s.MissionPayMult *= 1.5

	case enums.SocSlum:
		s.HasSlums = true
		s.SlumsDensity += 2.0
		s.CrewHiringCostMult *= 0.4
		s.BountyPayMult *= 0.5

	case enums.SocFrontierSpirit:
		s.MissionPayMult *= 1.2
		s.CrewSkillAvg = 6
		s.HasRefueling = true

	case enums.SocDecadentAristocracy:
		s.HasLuxuryHousing = true
		s.VIPDensity += 1.5
		s.CrewHiringCostMult *= 2.0
		s.SlumsDensity = 0.0

	case enums.SocWorkerRebellion:
		s.MarketSellMult *= 0.2
		s.HasCantina = true

	case enums.SocCultActivity:
		s.WantedPassChance += 0.10
		s.MissionPayMult *= 0.8
		s.CrewSkillAvg = 2

	case enums.SocScientificExpedition:
		s.MissionPayMult *= 2.0
		s.PassengerWealth *= 2.0
		s.HasOutfitter = true

	case enums.SocGhostTown:
		s.PassengerDensity = 0.0
		s.CrewPoolDensity = 0.0
		s.HasShipyard = false
		s.HasRefueling = false
		s.MarketBuyMult = 0.0
		s.HasMissionBoard = false

	case enums.SocGladiatorial:
		s.HasHospital = true
		s.CrewSkillAvg = 8
		s.BountyPayMult *= 1.5

	case enums.SocHiveMind:
		s.CrewPoolDensity += 2.0
		s.CrewSkillAvg = 10
		s.CrewHiringCostMult *= 0.1
		s.InspectionChance = 1.0
	}
}