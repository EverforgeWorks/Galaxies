package gen

import "galaxies/internal/core/domain"

// ProductionRule defines a set of conditions required for a system to produce an item.
// If a field is -1, it means "Any" (don't care).
type ProductionRule struct {
	Pol   int // Cast to domain.PoliticalStatus
	Eco   int // Cast to domain.EconomicStatus
	Soc   int // Cast to domain.SocialStatus
	Bonus bool // If true, this represents a high-yield/specialized source
}

// --- Helpers for readability ---
const Any = -1

func P(p domain.PoliticalStatus) ProductionRule {
	return ProductionRule{Pol: int(p), Eco: Any, Soc: Any}
}

func E(e domain.EconomicStatus) ProductionRule {
	return ProductionRule{Pol: Any, Eco: int(e), Soc: Any}
}

func S(s domain.SocialStatus) ProductionRule {
	return ProductionRule{Pol: Any, Eco: Any, Soc: int(s)}
}

func ComboPE(p domain.PoliticalStatus, e domain.EconomicStatus) ProductionRule {
	return ProductionRule{Pol: int(p), Eco: int(e), Soc: Any, Bonus: true}
}

func ComboES(e domain.EconomicStatus, s domain.SocialStatus) ProductionRule {
	return ProductionRule{Pol: Any, Eco: int(e), Soc: int(s), Bonus: true}
}

func ComboPS(p domain.PoliticalStatus, s domain.SocialStatus) ProductionRule {
	return ProductionRule{Pol: int(p), Eco: Any, Soc: int(s), Bonus: true}
}

// GlobalProductionMap links Items to the System States that manufacture/sell them.
var GlobalProductionMap = map[domain.ItemName][]ProductionRule{

	// =============================================================================
	// SECTOR 1: RAW MATERIALS
	// Produced by: Resource Rich, Frontier, Mining Colonies
	// =============================================================================
	domain.ItemHydrogenFuelCells: {
		E(domain.EcoRefuelingDepot), E(domain.EcoIndustrialBoom), P(domain.PolFederationOutpost),
		S(domain.SocScientificExpedition), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolIsolationist, domain.EcoRefuelingDepot),
	},
	domain.ItemIronOre: {
		E(domain.EcoResourceRich), S(domain.SocFrontierSpirit), P(domain.PolColonialCharter),
		E(domain.EcoScavengerEconomy), S(domain.SocPreIndustrial),
		ComboES(domain.EcoResourceRich, domain.SocWorkerRebellion),
	},
	domain.ItemCopperIngots: {
		E(domain.EcoResourceRich), P(domain.PolColonialCharter), S(domain.SocFrontierSpirit),
		E(domain.EcoManufacturingHub), S(domain.SocPreIndustrial),
		ComboPE(domain.PolFailedState, domain.EcoScavengerEconomy),
	},
	domain.ItemTitaniumPlating: {
		E(domain.EcoManufacturingHub), P(domain.PolMartialLaw), E(domain.EcoWarEconomy),
		P(domain.PolCorporateSovereign), S(domain.SocCyberneticAscension),
		ComboPE(domain.PolImperialCore, domain.EcoWarEconomy),
	},
	domain.ItemSiliconWafers: {
		E(domain.EcoTechBottleneck), P(domain.PolCorporateSovereign), S(domain.SocCyberneticAscension),
		E(domain.EcoPostScarcity), P(domain.PolAIGovernance),
		ComboES(domain.EcoManufacturingHub, domain.SocAcademy),
	},
	domain.ItemCarbonFiber: {
		E(domain.EcoManufacturingHub), P(domain.PolCorporateSovereign), E(domain.EcoIndustrialBoom),
		S(domain.SocCosmopolitan), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolFederationOutpost, domain.EcoIndustrialBoom),
	},
	domain.ItemRawBauxite: {
		E(domain.EcoResourceRich), S(domain.SocFrontierSpirit), P(domain.PolColonialCharter),
		S(domain.SocPreIndustrial), P(domain.PolFeudalDominion),
		ComboES(domain.EcoDepletedWorld, domain.SocSlum),
	},
	domain.ItemLithiumSalts: {
		E(domain.EcoResourceRich), P(domain.PolExileColony), S(domain.SocScientificExpedition),
		E(domain.EcoRefuelingDepot), P(domain.PolAnarchicFreehold),
		ComboPE(domain.PolDemilitarizedZone, domain.EcoResourceRich),
	},
	domain.ItemCobalt: {
		E(domain.EcoResourceRich), P(domain.PolMartialLaw), S(domain.SocWorkerRebellion),
		P(domain.PolColonialCharter), E(domain.EcoWarEconomy),
		ComboES(domain.EcoScavengerEconomy, domain.SocGhostTown),
	},
	domain.ItemTungstenCarbide: {
		E(domain.EcoManufacturingHub), P(domain.PolImperialCore), E(domain.EcoIndustrialBoom),
		S(domain.SocCyberneticAscension), P(domain.PolCorporateSovereign),
		ComboPE(domain.PolContestedWarZone, domain.EcoWarEconomy),
	},
	domain.ItemGoldBullion: {
		E(domain.EcoGoldRush), P(domain.PolImperialCore), S(domain.SocDecadentAristocracy),
		E(domain.EcoLuxuryResort), P(domain.PolSyndicateTerritory),
		ComboPE(domain.PolPirateHaven, domain.EcoGoldRush),
	},
	domain.ItemPlatinumCrystals: {
		E(domain.EcoGoldRush), P(domain.PolCorporateSovereign), E(domain.EcoLuxuryResort),
		S(domain.SocArtistEnclave), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoResourceRich, domain.SocScientificExpedition),
	},
	domain.ItemUraniumOre: {
		E(domain.EcoResourceRich), P(domain.PolMartialLaw), E(domain.EcoWarEconomy),
		P(domain.PolFailedState), S(domain.SocPreIndustrial),
		ComboPE(domain.PolDemilitarizedZone, domain.EcoScavengerEconomy),
	},
	domain.ItemPlutoniumRods: {
		P(domain.PolMartialLaw), P(domain.PolImperialCore), E(domain.EcoWarEconomy),
		P(domain.PolAIGovernance), S(domain.SocCyberneticAscension),
		ComboES(domain.EcoTechBottleneck, domain.SocScientificExpedition),
	},
	domain.ItemHelium3: {
		E(domain.EcoRefuelingDepot), S(domain.SocScientificExpedition), P(domain.PolFederationOutpost),
		P(domain.PolIsolationist), E(domain.EcoResourceRich),
		ComboPE(domain.PolExileColony, domain.EcoRefuelingDepot),
	},
	domain.ItemDeuterium: {
		E(domain.EcoRefuelingDepot), P(domain.PolTradeFedNeutrality), E(domain.EcoPostScarcity),
		S(domain.SocScientificExpedition), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoIndustrialBoom, domain.SocCosmopolitan),
	},
	domain.ItemWaterIce: {
		E(domain.EcoResourceRich), P(domain.PolExileColony), S(domain.SocFrontierSpirit),
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFederationOutpost),
		ComboPE(domain.PolIsolationist, domain.EcoDepletedWorld),
	},
	domain.ItemOxygenCanisters: {
		E(domain.EcoRefuelingDepot), S(domain.SocScientificExpedition), P(domain.PolFederationOutpost),
		E(domain.EcoAgrarianBreadbasket), S(domain.SocRefugeeCrisis),
		ComboES(domain.EcoDepletedWorld, domain.SocGhostTown),
	},
	domain.ItemNitrogenTanks: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolColonialCharter), E(domain.EcoResourceRich),
		S(domain.SocScientificExpedition), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolTerraforming, domain.EcoIndustrialBoom), // Implicit terraforming logic
	},
	domain.ItemRawBiomass: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocPreIndustrial), P(domain.PolFeudalDominion),
		E(domain.EcoFamine), S(domain.SocRefugeeCrisis),
		ComboPE(domain.PolFailedState, domain.EcoDepletedWorld),
	},
	domain.ItemScrapMetal: {
		E(domain.EcoScavengerEconomy), P(domain.PolFailedState), S(domain.SocGhostTown),
		E(domain.EcoDepletedWorld), P(domain.PolAnarchicFreehold),
		ComboES(domain.EcoWarEconomy, domain.SocSlum),
	},
	domain.ItemPolymerResins: {
		E(domain.EcoManufacturingHub), P(domain.PolCorporateSovereign), E(domain.EcoIndustrialBoom),
		S(domain.SocCosmopolitan), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolAIGovernance, domain.EcoPostScarcity),
	},
	domain.ItemGrapheneSheets: {
		E(domain.EcoTechBottleneck), P(domain.PolAIGovernance), S(domain.SocAcademy),
		E(domain.EcoPostScarcity), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoManufacturingHub, domain.SocCyberneticAscension),
	},
	domain.ItemObsidianGlass: {
		E(domain.EcoResourceRich), S(domain.SocArtistEnclave), P(domain.PolTheocraticRule),
		E(domain.EcoLuxuryResort), S(domain.SocPreIndustrial),
		ComboPE(domain.PolExileColony, domain.EcoResourceRich),
	},
	domain.ItemLiquidMethane: {
		E(domain.EcoRefuelingDepot), P(domain.PolExileColony), E(domain.EcoResourceRich),
		S(domain.SocFrontierSpirit), P(domain.PolColonialCharter),
		ComboES(domain.EcoDepletedWorld, domain.SocScientificExpedition),
	},

	// =============================================================================
	// SECTOR 2: INDUSTRIAL
	// Produced by: Manufacturing Hubs, Boom Economies, Corporate States
	// =============================================================================
	domain.ItemSteelBeams: {
		E(domain.EcoManufacturingHub), E(domain.EcoIndustrialBoom), P(domain.PolCorporateSovereign),
		S(domain.SocWorkerRebellion), P(domain.PolColonialCharter),
		ComboES(domain.EcoIndustrialBoom, domain.SocSlum),
	},
	domain.ItemCeramicComposites: {
		E(domain.EcoManufacturingHub), P(domain.PolMartialLaw), E(domain.EcoWarEconomy),
		P(domain.PolTradeFedNeutrality), S(domain.SocAcademy),
		ComboPE(domain.PolImperialCore, domain.EcoIndustrialBoom),
	},
	domain.ItemInsulatedWiring: {
		E(domain.EcoManufacturingHub), S(domain.SocWorkerRebellion), P(domain.PolPuppetState),
		E(domain.EcoIndustrialBoom), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoScavengerEconomy, domain.SocSlum),
	},
	domain.ItemHydraulicFluids: {
		E(domain.EcoIndustrialBoom), P(domain.PolFederationOutpost), E(domain.EcoRefuelingDepot),
		P(domain.PolTradeFedNeutrality), S(domain.SocWorkerRebellion),
		ComboPE(domain.PolCorporateSovereign, domain.EcoManufacturingHub),
	},
	domain.ItemMagneticCoils: {
		E(domain.EcoManufacturingHub), P(domain.PolAIGovernance), E(domain.EcoTechBottleneck),
		S(domain.SocCyberneticAscension), P(domain.PolImperialCore),
		ComboES(domain.EcoPostScarcity, domain.SocAcademy),
	},
	domain.ItemSolarPanels: {
		E(domain.EcoPostScarcity), S(domain.SocScientificExpedition), P(domain.PolTradeFedNeutrality),
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFederationOutpost),
		ComboPE(domain.PolAIGovernance, domain.EcoResourceRich),
	},
	domain.ItemRadiationShielding: {
		E(domain.EcoWarEconomy), P(domain.PolMartialLaw), S(domain.SocScientificExpedition),
		E(domain.EcoDepletedWorld), P(domain.PolFailedState),
		ComboPE(domain.PolExileColony, domain.EcoResourceRich),
	},
	domain.ItemTransparisteel: {
		E(domain.EcoManufacturingHub), P(domain.PolImperialCore), E(domain.EcoLuxuryResort),
		S(domain.SocCosmopolitan), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoIndustrialBoom, domain.SocArtistEnclave),
	},
	domain.ItemHullPlating: {
		E(domain.EcoManufacturingHub), P(domain.PolMartialLaw), E(domain.EcoWarEconomy),
		P(domain.PolImperialCore), S(domain.SocGladiatorial),
		ComboPE(domain.PolContestedWarZone, domain.EcoScavengerEconomy),
	},
	domain.ItemBallBearings: {
		E(domain.EcoManufacturingHub), P(domain.PolPuppetState), S(domain.SocWorkerRebellion),
		E(domain.EcoIndustrialBoom), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoScavengerEconomy, domain.SocSlum),
	},
	domain.ItemIndustrialLubricants: {
		E(domain.EcoIndustrialBoom), E(domain.EcoRefuelingDepot), P(domain.PolTradeFedNeutrality),
		P(domain.PolColonialCharter), S(domain.SocWorkerRebellion),
		ComboPE(domain.PolCorporateSovereign, domain.EcoResourceRich),
	},
	domain.ItemFusionRegulators: {
		E(domain.EcoTechBottleneck), P(domain.PolImperialCore), S(domain.SocAcademy),
		E(domain.EcoPostScarcity), P(domain.PolAIGovernance),
		ComboES(domain.EcoManufacturingHub, domain.SocScientificExpedition),
	},
	domain.ItemCircuitBoards: {
		E(domain.EcoManufacturingHub), P(domain.PolCorporateSovereign), E(domain.EcoTechBottleneck),
		S(domain.SocCyberneticAscension), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoIndustrialBoom, domain.SocCosmopolitan),
	},
	domain.ItemMicroprocessors: {
		E(domain.EcoTechBottleneck), P(domain.PolCorporateSovereign), S(domain.SocAcademy),
		P(domain.PolAIGovernance), E(domain.EcoPostScarcity),
		ComboES(domain.EcoManufacturingHub, domain.SocCyberneticAscension),
	},
	domain.ItemOpticalFibers: {
		E(domain.EcoManufacturingHub), S(domain.SocCosmopolitan), P(domain.PolCorporateSovereign),
		E(domain.EcoPostScarcity), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolImperialCore, domain.EcoTechBottleneck),
	},
	domain.ItemNanotubes: {
		E(domain.EcoTechBottleneck), P(domain.PolAIGovernance), S(domain.SocScientificExpedition),
		E(domain.EcoPostScarcity), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoManufacturingHub, domain.SocAcademy),
	},
	domain.ItemConcreteSlabs: {
		E(domain.EcoIndustrialBoom), P(domain.PolColonialCharter), P(domain.PolPuppetState),
		E(domain.EcoWarEconomy), S(domain.SocSlum),
		ComboPE(domain.PolFailedState, domain.EcoScavengerEconomy),
	},
	domain.ItemPrefabHabitationUnits: {
		E(domain.EcoIndustrialBoom), P(domain.PolColonialCharter), S(domain.SocRefugeeCrisis),
		E(domain.EcoGoldRush), P(domain.PolFederationOutpost),
		ComboPE(domain.PolExileColony, domain.EcoDepletedWorld),
	},
	domain.ItemLifeSupportFilters: {
		E(domain.EcoManufacturingHub), P(domain.PolFederationOutpost), E(domain.EcoDepletedWorld),
		S(domain.SocRefugeeCrisis), P(domain.PolIsolationist),
		ComboES(domain.EcoScavengerEconomy, domain.SocGhostTown),
	},
	domain.ItemWasteRecyclers: {
		E(domain.EcoDepletedWorld), P(domain.PolFailedState), S(domain.SocSlum),
		E(domain.EcoScavengerEconomy), P(domain.PolIsolationist),
		ComboPE(domain.PolAIGovernance, domain.EcoPostScarcity), // Efficient recycling
	},
	domain.ItemDrillingHeads: {
		E(domain.EcoManufacturingHub), P(domain.PolColonialCharter), E(domain.EcoGoldRush),
		E(domain.EcoResourceRich), S(domain.SocFrontierSpirit),
		ComboPE(domain.PolCorporateSovereign, domain.EcoResourceRich),
	},
	domain.ItemConveyorBelts: {
		E(domain.EcoIndustrialBoom), P(domain.PolCorporateSovereign), E(domain.EcoManufacturingHub),
		S(domain.SocWorkerRebellion), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoResourceRich, domain.SocPreIndustrial),
	},
	domain.ItemGravPlates: {
		E(domain.EcoTechBottleneck), P(domain.PolImperialCore), S(domain.SocScientificExpedition),
		P(domain.PolAIGovernance), E(domain.EcoPostScarcity),
		ComboPE(domain.PolCorporateSovereign, domain.EcoLuxuryResort),
	},
	domain.ItemDockingClamps: {
		E(domain.EcoIndustrialBoom), P(domain.PolFederationOutpost), P(domain.PolTradeFedNeutrality),
		E(domain.EcoRefuelingDepot), S(domain.SocFrontierSpirit),
		ComboES(domain.EcoManufacturingHub, domain.SocCosmopolitan),
	},
	domain.ItemAirLockSeals: {
		E(domain.EcoManufacturingHub), P(domain.PolFederationOutpost), S(domain.SocFrontierSpirit),
		P(domain.PolColonialCharter), E(domain.EcoRefuelingDepot),
		ComboPE(domain.PolIsolationist, domain.EcoDepletedWorld),
	},

	// =============================================================================
	// SECTOR 3: CONSUMABLES
	// Produced by: Agrarian Worlds, Hydroponics, Relief Camps
	// =============================================================================
	domain.ItemNutrientPaste: {
		E(domain.EcoFamine), S(domain.SocRefugeeCrisis), P(domain.PolMartialLaw),
		E(domain.EcoDepression), P(domain.PolFailedState),
		ComboES(domain.EcoIndustrialBoom, domain.SocSlum),
	},
	domain.ItemSyntheticMeat: {
		E(domain.EcoManufacturingHub), P(domain.PolCorporateSovereign), S(domain.SocCosmopolitan),
		E(domain.EcoPostScarcity), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolAIGovernance, domain.EcoAgrarianBreadbasket),
	},
	domain.ItemGrainSacks: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFeudalDominion), S(domain.SocPreIndustrial),
		P(domain.PolColonialCharter), E(domain.EcoResourceRich),
		ComboPE(domain.PolIsolationist, domain.EcoAgrarianBreadbasket),
	},
	domain.ItemPurifiedWater: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFederationOutpost), E(domain.EcoLuxuryResort),
		S(domain.SocReligiousPilgrimage), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoRefuelingDepot, domain.SocScientificExpedition),
	},
	domain.ItemDehydratedVegetables: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocFrontierSpirit), P(domain.PolColonialCharter),
		E(domain.EcoWarEconomy), P(domain.PolExileColony),
		ComboES(domain.EcoRefuelingDepot, domain.SocScientificExpedition),
	},
	domain.ItemSoyProtein: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolCorporateSovereign), S(domain.SocSlum),
		E(domain.EcoFamine), P(domain.PolPuppetState),
		ComboPE(domain.PolFederationOutpost, domain.EcoDepletedWorld),
	},
	domain.ItemAlgaeMash: {
		E(domain.EcoDepletedWorld), S(domain.SocSlum), P(domain.PolFailedState),
		E(domain.EcoScavengerEconomy), P(domain.PolExileColony),
		ComboPE(domain.PolIsolationist, domain.EcoFamine),
	},
	domain.ItemCoffeeBeans: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocCosmopolitan), P(domain.PolTradeFedNeutrality),
		S(domain.SocArtistEnclave), E(domain.EcoLuxuryResort),
		ComboPE(domain.PolColonialCharter, domain.EcoAgrarianBreadbasket),
	},
	domain.ItemTeaLeaves: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFeudalDominion), S(domain.SocReligiousPilgrimage),
		P(domain.PolIsolationist), E(domain.EcoLuxuryResort),
		ComboES(domain.EcoAgrarianBreadbasket, domain.SocArtistEnclave),
	},
	domain.ItemSpices: {
		E(domain.EcoResourceRich), P(domain.PolFeudalDominion), S(domain.SocFrontierSpirit),
		E(domain.EcoLuxuryResort), E(domain.EcoGoldRush),
		ComboPE(domain.PolColonialCharter, domain.EcoAgrarianBreadbasket),
	},
	domain.ItemExoticFruits: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocArtistEnclave), E(domain.EcoLuxuryResort),
		P(domain.PolImperialCore), S(domain.SocReligiousPilgrimage),
		ComboES(domain.EcoPostScarcity, domain.SocCosmopolitan),
	},
	domain.ItemLivestock: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFeudalDominion), S(domain.SocPreIndustrial),
		P(domain.PolColonialCharter), S(domain.SocFrontierSpirit),
		ComboPE(domain.PolIsolationist, domain.EcoResourceRich),
	},
	domain.ItemFertilizer: {
		E(domain.EcoManufacturingHub), P(domain.PolCorporateSovereign), E(domain.EcoIndustrialBoom),
		S(domain.SocPreIndustrial), P(domain.PolColonialCharter),
		ComboES(domain.EcoScavengerEconomy, domain.SocSlum),
	},
	domain.ItemSeeds: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocScientificExpedition), P(domain.PolColonialCharter),
		P(domain.PolFeudalDominion), E(domain.EcoResourceRich),
		ComboPE(domain.PolTerraforming, domain.EcoPostScarcity),
	},
	domain.ItemAntibiotics: {
		S(domain.SocPlague), P(domain.PolImperialCore), S(domain.SocScientificExpedition),
		E(domain.EcoPostScarcity), P(domain.PolFederationOutpost),
		ComboES(domain.EcoTechBottleneck, domain.SocRefugeeCrisis),
	},
	domain.ItemFirstAidKits: {
		P(domain.PolFederationOutpost), E(domain.EcoWarEconomy), S(domain.SocRefugeeCrisis),
		P(domain.PolMartialLaw), S(domain.SocPlague),
		ComboPE(domain.PolTradeFedNeutrality, domain.EcoManufacturingHub),
	},
	domain.ItemBandages: {
		E(domain.EcoManufacturingHub), S(domain.SocRefugeeCrisis), P(domain.PolFailedState),
		E(domain.EcoDepression), S(domain.SocSlum),
		ComboPE(domain.PolFeudalDominion, domain.EcoScavengerEconomy),
	},
	domain.ItemVitamins: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocCosmopolitan), P(domain.PolCorporateSovereign),
		E(domain.EcoPostScarcity), S(domain.SocScientificExpedition),
		ComboES(domain.EcoFamine, domain.SocRefugeeCrisis), // Aid shipments
	},
	domain.ItemStimPacks: {
		P(domain.PolMartialLaw), E(domain.EcoWarEconomy), S(domain.SocGladiatorial),
		P(domain.PolSyndicateTerritory), S(domain.SocWorkerRebellion),
		ComboPE(domain.PolContestedWarZone, domain.EcoBlackMarketHub),
	},
	domain.ItemAlcohol: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocCosmopolitan), P(domain.PolAnarchicFreehold),
		S(domain.SocFrontierSpirit), P(domain.PolFeudalDominion),
		ComboES(domain.EcoLaborStrike, domain.SocSlum),
	},
	domain.ItemDistilledSpirits: {
		P(domain.PolImperialCore), S(domain.SocDecadentAristocracy), E(domain.EcoLuxuryResort),
		P(domain.PolSyndicateTerritory), S(domain.SocArtistEnclave),
		ComboPE(domain.PolPirateHaven, domain.EcoBlackMarketHub),
	},
	domain.ItemCigarettes: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocWorkerRebellion), P(domain.PolCorporateSovereign),
		S(domain.SocSlum), E(domain.EcoIndustrialBoom),
		ComboES(domain.EcoDepression, domain.SocRefugeeCrisis),
	},
	domain.ItemChocolate: {
		E(domain.EcoAgrarianBreadbasket), E(domain.EcoLuxuryResort), S(domain.SocCosmopolitan),
		P(domain.PolImperialCore), S(domain.SocDecadentAristocracy),
		ComboPE(domain.PolTradeFedNeutrality, domain.EcoPostScarcity),
	},
	domain.ItemHoney: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFeudalDominion), S(domain.SocPreIndustrial),
		E(domain.EcoLuxuryResort), P(domain.PolColonialCharter),
		ComboES(domain.EcoResourceRich, domain.SocArtistEnclave),
	},
	domain.ItemDairyProducts: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolColonialCharter), S(domain.SocFrontierSpirit),
		P(domain.PolFeudalDominion), S(domain.SocPreIndustrial),
		ComboPE(domain.PolFederationOutpost, domain.EcoAgrarianBreadbasket),
	},

	// =============================================================================
	// SECTOR 4: TECHNOLOGY
	// Produced by: Tech Worlds, AI Governance, Research Stations
	// =============================================================================
	domain.ItemQuantumProcessors: {
		P(domain.PolAIGovernance), E(domain.EcoPostScarcity), S(domain.SocAcademy),
		P(domain.PolImperialCore), S(domain.SocCyberneticAscension),
		ComboPE(domain.PolCorporateSovereign, domain.EcoTechBottleneck),
	},
	domain.ItemAICoresBasic: {
		P(domain.PolCorporateSovereign), E(domain.EcoManufacturingHub), S(domain.SocCosmopolitan),
		P(domain.PolTradeFedNeutrality), E(domain.EcoIndustrialBoom),
		ComboES(domain.EcoPostScarcity, domain.SocWorkerRebellion),
	},
	domain.ItemAICoresAdvanced: {
		P(domain.PolAIGovernance), S(domain.SocCyberneticAscension), E(domain.EcoPostScarcity),
		S(domain.SocAcademy), P(domain.PolImperialCore),
		ComboPE(domain.PolCorporateSovereign, domain.EcoTechBottleneck),
	},
	domain.ItemHolographicProjectors: {
		E(domain.EcoManufacturingHub), S(domain.SocCosmopolitan), E(domain.EcoLuxuryResort),
		S(domain.SocArtistEnclave), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoPostScarcity, domain.SocDecadentAristocracy),
	},
	domain.ItemNavComputerChips: {
		E(domain.EcoManufacturingHub), P(domain.PolFederationOutpost), E(domain.EcoTechBottleneck),
		P(domain.PolTradeFedNeutrality), S(domain.SocAcademy),
		ComboPE(domain.PolAIGovernance, domain.EcoIndustrialBoom),
	},
	domain.ItemSensorArrays: {
		S(domain.SocScientificExpedition), P(domain.PolMartialLaw), E(domain.EcoTechBottleneck),
		P(domain.PolFederationOutpost), S(domain.SocFrontierSpirit),
		ComboPE(domain.PolDemilitarizedZone, domain.EcoBlackMarketHub),
	},
	domain.ItemTargetingLogicBoards: {
		P(domain.PolMartialLaw), E(domain.EcoWarEconomy), P(domain.PolImperialCore),
		P(domain.PolCorporateSovereign), S(domain.SocGladiatorial),
		ComboES(domain.EcoTechBottleneck, domain.SocWorkerRebellion),
	},
	domain.ItemCommunicationRelays: {
		E(domain.EcoManufacturingHub), P(domain.PolFederationOutpost), S(domain.SocCosmopolitan),
		E(domain.EcoTechBottleneck), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolIsolationist, domain.EcoBlackMarketHub),
	},
	domain.ItemTranslationMatrix: {
		S(domain.SocCosmopolitan), P(domain.PolTradeFedNeutrality), S(domain.SocAcademy),
		P(domain.PolFederationOutpost), E(domain.EcoLuxuryResort),
		ComboES(domain.EcoPostScarcity, domain.SocArtistEnclave),
	},
	domain.ItemCryptographicKeys: {
		P(domain.PolCorporateSovereign), S(domain.SocCyberneticAscension), E(domain.EcoBlackMarketHub),
		P(domain.PolMartialLaw), P(domain.PolSyndicateTerritory),
		ComboES(domain.EcoTechBottleneck, domain.SocAcademy),
	},
	domain.ItemDataStorageDrives: {
		E(domain.EcoManufacturingHub), P(domain.PolCorporateSovereign), E(domain.EcoTechBottleneck),
		S(domain.SocAcademy), P(domain.PolTradeFedNeutrality),
		ComboPE(domain.PolAIGovernance, domain.EcoPostScarcity),
	},
	domain.ItemServerRacks: {
		P(domain.PolCorporateSovereign), E(domain.EcoTechBottleneck), S(domain.SocAcademy),
		P(domain.PolAIGovernance), E(domain.EcoManufacturingHub),
		ComboES(domain.EcoPostScarcity, domain.SocCyberneticAscension),
	},
	domain.ItemRoboticArms: {
		E(domain.EcoManufacturingHub), E(domain.EcoIndustrialBoom), P(domain.PolCorporateSovereign),
		S(domain.SocWorkerRebellion), P(domain.PolAIGovernance),
		ComboPE(domain.PolColonialCharter, domain.EcoTechBottleneck),
	},
	domain.ItemDroneChassis: {
		E(domain.EcoManufacturingHub), P(domain.PolMartialLaw), E(domain.EcoWarEconomy),
		P(domain.PolCorporateSovereign), S(domain.SocCyberneticAscension),
		ComboES(domain.EcoIndustrialBoom, domain.SocFrontierSpirit),
	},
	domain.ItemCyberneticLimbs: {
		S(domain.SocCyberneticAscension), P(domain.PolCorporateSovereign), E(domain.EcoWarEconomy),
		S(domain.SocGladiatorial), P(domain.PolImperialCore),
		ComboES(domain.EcoManufacturingHub, domain.SocPlague),
	},
	domain.ItemOcularImplants: {
		S(domain.SocCyberneticAscension), P(domain.PolCorporateSovereign), S(domain.SocAcademy),
		E(domain.EcoTechBottleneck), P(domain.PolSyndicateTerritory),
		ComboPE(domain.PolAIGovernance, domain.EcoLuxuryResort),
	},
	domain.ItemNeuralInterfaces: {
		S(domain.SocCyberneticAscension), P(domain.PolAIGovernance), E(domain.EcoPostScarcity),
		P(domain.PolImperialCore), S(domain.SocAcademy),
		ComboES(domain.EcoTechBottleneck, domain.SocDecadentAristocracy),
	},
	domain.ItemMedicalScanners: {
		S(domain.SocPlague), P(domain.PolImperialCore), E(domain.EcoPostScarcity),
		S(domain.SocScientificExpedition), P(domain.PolFederationOutpost),
		ComboPE(domain.PolTheocraticRule, domain.EcoLuxuryResort),
	},
	domain.ItemTerraformingMicrobes: {
		S(domain.SocScientificExpedition), P(domain.PolCorporateSovereign), E(domain.EcoPostScarcity),
		P(domain.PolColonialCharter), E(domain.EcoAgrarianBreadbasket),
		ComboPE(domain.PolTerraforming, domain.EcoTechBottleneck),
	},
	domain.ItemAtmosphericProcessors: {
		E(domain.EcoManufacturingHub), P(domain.PolColonialCharter), S(domain.SocScientificExpedition),
		E(domain.EcoIndustrialBoom), P(domain.PolTerraforming),
		ComboES(domain.EcoDepletedWorld, domain.SocFrontierSpirit),
	},
	domain.ItemEncryptionBreakers: {
		E(domain.EcoBlackMarketHub), P(domain.PolSyndicateTerritory), P(domain.PolAnarchicFreehold),
		S(domain.SocCyberneticAscension), P(domain.PolRevolutionaryFront),
		ComboPE(domain.PolPirateHaven, domain.EcoTechBottleneck),
	},
	domain.ItemSecurityFirewalls: {
		P(domain.PolCorporateSovereign), P(domain.PolImperialCore), P(domain.PolMartialLaw),
		E(domain.EcoTechBottleneck), S(domain.SocAcademy),
		ComboES(domain.EcoPostScarcity, domain.SocCyberneticAscension),
	},
	domain.ItemVRHeadsets: {
		E(domain.EcoManufacturingHub), S(domain.SocCosmopolitan), E(domain.EcoLuxuryResort),
		P(domain.PolCorporateSovereign), S(domain.SocDecadentAristocracy),
		ComboPE(domain.PolTradeFedNeutrality, domain.EcoPostScarcity),
	},
	domain.ItemPowerInverters: {
		E(domain.EcoManufacturingHub), E(domain.EcoIndustrialBoom), P(domain.PolFederationOutpost),
		P(domain.PolColonialCharter), E(domain.EcoRefuelingDepot),
		ComboES(domain.EcoScavengerEconomy, domain.SocSlum),
	},
	domain.ItemBatteryBanks: {
		E(domain.EcoManufacturingHub), P(domain.PolTradeFedNeutrality), E(domain.EcoIndustrialBoom),
		S(domain.SocScientificExpedition), P(domain.PolColonialCharter),
		ComboPE(domain.PolFederationOutpost, domain.EcoRefuelingDepot),
	},

	// =============================================================================
	// SECTOR 5: LUXURY
	// Produced by: Artist Enclaves, Imperial Cores, Rich Worlds
	// =============================================================================
	domain.ItemSilk: {
		E(domain.EcoAgrarianBreadbasket), P(domain.PolFeudalDominion), S(domain.SocArtistEnclave),
		E(domain.EcoLuxuryResort), P(domain.PolIsolationist),
		ComboES(domain.EcoAgrarianBreadbasket, domain.SocDecadentAristocracy),
	},
	domain.ItemDesignerClothing: {
		E(domain.EcoLuxuryResort), S(domain.SocCosmopolitan), P(domain.PolImperialCore),
		S(domain.SocArtistEnclave), S(domain.SocDecadentAristocracy),
		ComboPE(domain.PolCorporateSovereign, domain.EcoPostScarcity),
	},
	domain.ItemJewelry: {
		E(domain.EcoLuxuryResort), P(domain.PolImperialCore), E(domain.EcoGoldRush),
		S(domain.SocDecadentAristocracy), P(domain.PolSyndicateTerritory),
		ComboES(domain.EcoResourceRich, domain.SocArtistEnclave),
	},
	domain.ItemGemstones: {
		E(domain.EcoResourceRich), E(domain.EcoGoldRush), P(domain.PolColonialCharter),
		P(domain.PolFeudalDominion), S(domain.SocFrontierSpirit),
		ComboPE(domain.PolExileColony, domain.EcoLuxuryResort),
	},
	domain.ItemSculptures: {
		S(domain.SocArtistEnclave), P(domain.PolImperialCore), E(domain.EcoLuxuryResort),
		S(domain.SocReligiousPilgrimage), P(domain.PolTheocraticRule),
		ComboPE(domain.PolFeudalDominion, domain.EcoPostScarcity),
	},
	domain.ItemPaintings: {
		S(domain.SocArtistEnclave), S(domain.SocDecadentAristocracy), E(domain.EcoLuxuryResort),
		P(domain.PolImperialCore), S(domain.SocCosmopolitan),
		ComboES(domain.EcoPostScarcity, domain.SocReligiousPilgrimage),
	},
	domain.ItemAncientArtifacts: {
		E(domain.EcoScavengerEconomy), P(domain.PolExileColony), S(domain.SocScientificExpedition),
		S(domain.SocCultActivity), P(domain.PolAnarchicFreehold),
		ComboPE(domain.PolFailedState, domain.EcoResourceRich),
	},
	domain.ItemRareBooks: {
		S(domain.SocAcademy), P(domain.PolImperialCore), S(domain.SocReligiousPilgrimage),
		E(domain.EcoPostScarcity), P(domain.PolIsolationist),
		ComboES(domain.EcoLuxuryResort, domain.SocDecadentAristocracy),
	},
	domain.ItemVinylRecords: {
		S(domain.SocArtistEnclave), S(domain.SocCosmopolitan), P(domain.PolTradeFedNeutrality),
		E(domain.EcoScavengerEconomy), P(domain.PolAnarchicFreehold),
		ComboPE(domain.PolExileColony, domain.EcoLuxuryResort),
	},
	domain.ItemFineWine: {
		E(domain.EcoLuxuryResort), P(domain.PolFeudalDominion), S(domain.SocDecadentAristocracy),
		P(domain.PolImperialCore), S(domain.SocReligiousPilgrimage),
		ComboES(domain.EcoAgrarianBreadbasket, domain.SocCosmopolitan),
	},
	domain.ItemAgedWhiskey: {
		E(domain.EcoLuxuryResort), S(domain.SocFrontierSpirit), P(domain.PolColonialCharter),
		S(domain.SocDecadentAristocracy), P(domain.PolSyndicateTerritory),
		ComboES(domain.EcoAgrarianBreadbasket, domain.SocArtistEnclave),
	},
	domain.ItemPerfume: {
		E(domain.EcoLuxuryResort), S(domain.SocArtistEnclave), P(domain.PolImperialCore),
		S(domain.SocDecadentAristocracy), E(domain.EcoAgrarianBreadbasket),
		ComboPE(domain.PolCorporateSovereign, domain.EcoPostScarcity),
	},
	domain.ItemExoticPets: {
		E(domain.EcoResourceRich), S(domain.SocFrontierSpirit), E(domain.EcoLuxuryResort),
		P(domain.PolSyndicateTerritory), P(domain.PolColonialCharter),
		ComboES(domain.EcoGoldRush, domain.SocDecadentAristocracy),
	},
	domain.ItemMusicalInstruments: {
		S(domain.SocArtistEnclave), S(domain.SocCosmopolitan), E(domain.EcoLuxuryResort),
		P(domain.PolImperialCore), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoManufacturingHub, domain.SocDecadentAristocracy),
	},
	domain.ItemCeremonialRobes: {
		P(domain.PolTheocraticRule), S(domain.SocReligiousPilgrimage), P(domain.PolFeudalDominion),
		S(domain.SocCultActivity), P(domain.PolImperialCore),
		ComboES(domain.EcoLuxuryResort, domain.SocArtistEnclave),
	},
	domain.ItemRelics: {
		S(domain.SocReligiousPilgrimage), P(domain.PolTheocraticRule), E(domain.EcoScavengerEconomy),
		P(domain.PolFailedState), S(domain.SocCultActivity),
		ComboPE(domain.PolIsolationist, domain.EcoResourceRich),
	},
	domain.ItemTrophies: {
		S(domain.SocGladiatorial), P(domain.PolFeudalDominion), E(domain.EcoWarEconomy),
		S(domain.SocFrontierSpirit), P(domain.PolMartialLaw),
		ComboES(domain.EcoLuxuryResort, domain.SocDecadentAristocracy),
	},
	domain.ItemVRExperiences: {
		S(domain.SocCosmopolitan), E(domain.EcoPostScarcity), P(domain.PolCorporateSovereign),
		E(domain.EcoLuxuryResort), S(domain.SocCyberneticAscension),
		ComboPE(domain.PolTradeFedNeutrality, domain.EcoTechBottleneck),
	},
	domain.ItemPersonalButlers: {
		P(domain.PolImperialCore), E(domain.EcoPostScarcity), P(domain.PolCorporateSovereign),
		S(domain.SocDecadentAristocracy), P(domain.PolAIGovernance),
		ComboPE(domain.PolPuppetState, domain.EcoLuxuryResort),
	},
	domain.ItemRestoredClassicCars: {
		E(domain.EcoLuxuryResort), S(domain.SocArtistEnclave), P(domain.PolImperialCore),
		S(domain.SocCosmopolitan), E(domain.EcoScavengerEconomy),
		ComboPE(domain.PolCorporateSovereign, domain.EcoPostScarcity),
	},
	domain.ItemGoldPlatedWeapons: {
		P(domain.PolSyndicateTerritory), P(domain.PolPirateHaven), S(domain.SocDecadentAristocracy),
		E(domain.EcoLuxuryResort), P(domain.PolMartialLaw),
		ComboES(domain.EcoWarEconomy, domain.SocGladiatorial),
	},
	domain.ItemAlienPottery: {
		E(domain.EcoScavengerEconomy), S(domain.SocScientificExpedition), P(domain.PolExileColony),
		E(domain.EcoResourceRich), P(domain.PolIsolationist),
		ComboPE(domain.PolFailedState, domain.EcoLuxuryResort),
	},
	domain.ItemGeneEditedFlowers: {
		E(domain.EcoAgrarianBreadbasket), S(domain.SocScientificExpedition), E(domain.EcoLuxuryResort),
		P(domain.PolImperialCore), S(domain.SocArtistEnclave),
		ComboPE(domain.PolCorporateSovereign, domain.EcoPostScarcity),
	},
	domain.ItemCaviar: {
		E(domain.EcoLuxuryResort), P(domain.PolImperialCore), S(domain.SocDecadentAristocracy),
		E(domain.EcoAgrarianBreadbasket), P(domain.PolTradeFedNeutrality),
		ComboES(domain.EcoPostScarcity, domain.SocCosmopolitan),
	},
	domain.ItemTruffles: {
		E(domain.EcoAgrarianBreadbasket), E(domain.EcoLuxuryResort), P(domain.PolFeudalDominion),
		S(domain.SocPreIndustrial), P(domain.PolIsolationist),
		ComboPE(domain.PolImperialCore, domain.EcoResourceRich),
	},

	// =============================================================================
	// SECTOR 6: CONTRABAND
	// Produced by: Pirate Havens, Black Markets, Anarchic Zones
	// =============================================================================
	domain.ItemNarcotics: {
		P(domain.PolPirateHaven), E(domain.EcoBlackMarketHub), S(domain.SocSlum),
		P(domain.PolSyndicateTerritory), P(domain.PolAnarchicFreehold),
		ComboES(domain.EcoLuxuryResort, domain.SocDecadentAristocracy),
	},
	domain.ItemStolenIDChips: {
		P(domain.PolSyndicateTerritory), E(domain.EcoBlackMarketHub), S(domain.SocCyberneticAscension),
		P(domain.PolAnarchicFreehold), S(domain.SocSlum),
		ComboPE(domain.PolCorporateSovereign, domain.EcoDepression),
	},
	domain.ItemHackedCreditSticks: {
		P(domain.PolSyndicateTerritory), E(domain.EcoBlackMarketHub), S(domain.SocCosmopolitan),
		P(domain.PolAnarchicFreehold), S(domain.SocAcademy),
		ComboES(domain.EcoTechBottleneck, domain.SocSlum),
	},
	domain.ItemUnregisteredWeapons: {
		P(domain.PolAnarchicFreehold), P(domain.PolRevolutionaryFront), E(domain.EcoWarEconomy),
		P(domain.PolPirateHaven), P(domain.PolSyndicateTerritory),
		ComboPE(domain.PolContestedWarZone, domain.EcoScavengerEconomy),
	},
	domain.ItemCombatStims: {
		P(domain.PolMartialLaw), E(domain.EcoWarEconomy), S(domain.SocGladiatorial),
		P(domain.PolPirateHaven), E(domain.EcoBlackMarketHub),
		ComboES(domain.EcoIndustrialBoom, domain.SocWorkerRebellion),
	},
	domain.ItemSlaveCollars: {
		P(domain.PolPirateHaven), P(domain.PolFeudalDominion), P(domain.PolFailedState),
		E(domain.EcoBlackMarketHub), P(domain.PolSyndicateTerritory),
		ComboPE(domain.PolPenalColony, domain.EcoScavengerEconomy),
	},
	domain.ItemCounterfeitCurrency: {
		P(domain.PolSyndicateTerritory), E(domain.EcoBlackMarketHub), P(domain.PolAnarchicFreehold),
		S(domain.SocSlum), E(domain.EcoDepression),
		ComboES(domain.EcoIndustrialBoom, domain.SocCosmopolitan),
	},
	domain.ItemBootlegSoftware: {
		P(domain.PolAnarchicFreehold), S(domain.SocAcademy), E(domain.EcoBlackMarketHub),
		P(domain.PolSyndicateTerritory), S(domain.SocCyberneticAscension),
		ComboPE(domain.PolTradeFedNeutrality, domain.EcoTechBottleneck),
	},
	domain.ItemAILimitersDisabled: {
		P(domain.PolAIGovernance), S(domain.SocCyberneticAscension), E(domain.EcoBlackMarketHub),
		P(domain.PolPirateHaven), S(domain.SocScientificExpedition),
		ComboES(domain.EcoTechBottleneck, domain.SocWorkerRebellion),
	},
	domain.ItemHarvestedOrgans: {
		P(domain.PolFailedState), S(domain.SocRefugeeCrisis), E(domain.EcoBlackMarketHub),
		P(domain.PolPirateHaven), S(domain.SocPlague),
		ComboPE(domain.PolPenalColony, domain.EcoDepression),
	},
	domain.ItemEndangeredSpecies: {
		P(domain.PolPirateHaven), E(domain.EcoBlackMarketHub), P(domain.PolColonialCharter),
		E(domain.EcoResourceRich), P(domain.PolIsolationist),
		ComboES(domain.EcoLuxuryResort, domain.SocDecadentAristocracy),
	},
	domain.ItemRadioactiveWaste: {
		E(domain.EcoIndustrialBoom), E(domain.EcoWarEconomy), P(domain.PolFailedState),
		E(domain.EcoDepletedWorld), E(domain.EcoScavengerEconomy),
		ComboPE(domain.PolCorporateSovereign, domain.EcoManufacturingHub),
	},
	domain.ItemExplosives: {
		E(domain.EcoWarEconomy), P(domain.PolRevolutionaryFront), E(domain.EcoResourceRich),
		P(domain.PolMartialLaw), P(domain.PolPirateHaven),
		ComboES(domain.EcoIndustrialBoom, domain.SocWorkerRebellion),
	},
	domain.ItemPoison: {
		E(domain.EcoResourceRich), P(domain.PolSyndicateTerritory), S(domain.SocScientificExpedition),
		E(domain.EcoBlackMarketHub), P(domain.PolIsolationist),
		ComboPE(domain.PolFeudalDominion, domain.EcoAgrarianBreadbasket),
	},
	domain.ItemSpyware: {
		P(domain.PolCorporateSovereign), P(domain.PolImperialCore), E(domain.EcoBlackMarketHub),
		S(domain.SocAcademy), P(domain.PolSyndicateTerritory),
		ComboES(domain.EcoTechBottleneck, domain.SocCosmopolitan),
	},
	domain.ItemTortureDevices: {
		P(domain.PolMartialLaw), P(domain.PolPirateHaven), P(domain.PolFeudalDominion),
		E(domain.EcoBlackMarketHub), P(domain.PolFailedState),
		ComboPE(domain.PolPenalColony, domain.EcoDepression),
	},
	domain.ItemRelicFragments: {
		P(domain.PolExileColony), E(domain.EcoScavengerEconomy), S(domain.SocCultActivity),
		P(domain.PolAnarchicFreehold), S(domain.SocScientificExpedition),
		ComboES(domain.EcoResourceRich, domain.SocFrontierSpirit),
	},
	domain.ItemMindWipeDrugs: {
		P(domain.PolSyndicateTerritory), S(domain.SocCyberneticAscension), E(domain.EcoBlackMarketHub),
		P(domain.PolMartialLaw), P(domain.PolCorporateSovereign),
		ComboPE(domain.PolPenalColony, domain.EcoTechBottleneck),
	},
	domain.ItemFugitiveBiometrics: {
		P(domain.PolSyndicateTerritory), E(domain.EcoBlackMarketHub), S(domain.SocSlum),
		P(domain.PolPirateHaven), P(domain.PolPenalColony),
		ComboES(domain.EcoDepression, domain.SocCosmopolitan),
	},
	domain.ItemRedMercury: {
		P(domain.PolFailedState), E(domain.EcoScavengerEconomy), E(domain.EcoWarEconomy),
		P(domain.PolPirateHaven), E(domain.EcoBlackMarketHub),
		ComboPE(domain.PolDemilitarizedZone, domain.EcoDepletedWorld),
	},
	domain.ItemSentientAICode: {
		P(domain.PolAIGovernance), S(domain.SocCyberneticAscension), E(domain.EcoBlackMarketHub),
		S(domain.SocScientificExpedition), P(domain.PolCorporateSovereign),
		ComboES(domain.EcoTechBottleneck, domain.SocCultActivity),
	},
	domain.ItemCloneVats: {
		S(domain.SocCyberneticAscension), E(domain.EcoBlackMarketHub), P(domain.PolCorporateSovereign),
		S(domain.SocScientificExpedition), P(domain.PolAIGovernance),
		ComboPE(domain.PolSyndicateTerritory, domain.EcoPostScarcity),
	},
	domain.ItemGeneWarfareCanisters: {
		P(domain.PolFailedState), E(domain.EcoWarEconomy), E(domain.EcoBlackMarketHub),
		P(domain.PolMartialLaw), S(domain.SocPlague),
		ComboPE(domain.PolIsolationist, domain.EcoTechBottleneck),
	},
	domain.ItemPsychoActiveSpores: {
		E(domain.EcoResourceRich), S(domain.SocCultActivity), E(domain.EcoBlackMarketHub),
		P(domain.PolExileColony), S(domain.SocFrontierSpirit),
		ComboES(domain.EcoFamine, domain.SocReligiousPilgrimage),
	},
	domain.ItemTheGoodStuff: {
		P(domain.PolPirateHaven), S(domain.SocDecadentAristocracy), E(domain.EcoBlackMarketHub),
		P(domain.PolSyndicateTerritory), S(domain.SocArtistEnclave),
		ComboPE(domain.PolAnarchicFreehold, domain.EcoLuxuryResort),
	},
}