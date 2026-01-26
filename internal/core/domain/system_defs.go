package domain

// --- POLITICAL STATUSES (0-19) ---
type PoliticalStatus int

const (
	PolImperialCore      PoliticalStatus = 0
	PolFederationOutpost PoliticalStatus = 1
	PolMartialLaw        PoliticalStatus = 2
	PolCorporateSovereign PoliticalStatus = 3
	PolAnarchicFreehold  PoliticalStatus = 4
	PolPirateHaven       PoliticalStatus = 5
	PolTheocraticRule    PoliticalStatus = 6
	PolBureaucraticGridlock PoliticalStatus = 7
	PolContestedWarZone  PoliticalStatus = 8
	PolDemilitarizedZone PoliticalStatus = 9
	PolPuppetState       PoliticalStatus = 10
	PolSyndicateTerritory PoliticalStatus = 11
	PolIsolationist      PoliticalStatus = 12
	PolRevolutionaryFront PoliticalStatus = 13
	PolColonialCharter   PoliticalStatus = 14
	PolFailedState       PoliticalStatus = 15
	PolAIGovernance      PoliticalStatus = 16
	PolExileColony       PoliticalStatus = 17
	PolTradeFedNeutrality PoliticalStatus = 18
	PolFeudalDominion    PoliticalStatus = 19
)

func (p PoliticalStatus) String() string {
	names := [...]string{
		"Imperial Core", "Federation Outpost", "Martial Law", "Corporate Sovereign",
		"Anarchic Freehold", "Pirate Haven", "Theocratic Rule", "Bureaucratic Gridlock",
		"Contested War Zone", "Demilitarized Zone", "Puppet State", "Syndicate Territory",
		"Isolationist", "Revolutionary Front", "Colonial Charter", "Failed State",
		"AI Governance", "Exile Colony", "Trade Federation Neutrality", "Feudal Dominion",
	}
	if int(p) < len(names) {
		return names[p]
	}
	return "Unknown"
}

// --- ECONOMIC STATUSES (0-19) ---
type EconomicStatus int

const (
	EcoPostScarcity      EconomicStatus = 0
	EcoIndustrialBoom    EconomicStatus = 1
	EcoDepression        EconomicStatus = 2
	EcoHyperInflation    EconomicStatus = 3
	EcoResourceRich      EconomicStatus = 4
	EcoFamine            EconomicStatus = 5
	EcoTechBottleneck    EconomicStatus = 6
	EcoBlackMarketHub    EconomicStatus = 7
	EcoRefuelingDepot    EconomicStatus = 8
	EcoLuxuryResort      EconomicStatus = 9
	EcoManufacturingHub  EconomicStatus = 10
	EcoTradeEmbargo      EconomicStatus = 11
	EcoGoldRush          EconomicStatus = 12
	EcoLaborStrike       EconomicStatus = 13
	EcoWarEconomy        EconomicStatus = 14
	EcoDepletedWorld     EconomicStatus = 15
	EcoAgrarianBreadbasket EconomicStatus = 16
	EcoCommandEconomy    EconomicStatus = 17
	EcoFreePort          EconomicStatus = 18
	EcoScavengerEconomy  EconomicStatus = 19
)

func (e EconomicStatus) String() string {
	names := [...]string{
		"Post-Scarcity Utopia", "Industrial Boom", "Economic Depression", "Hyper-Inflation",
		"Resource Rich", "Famine / Starvation", "Tech Bottleneck", "Black Market Hub",
		"Refueling Depot", "Luxury Resort", "Manufacturing Hub", "Trade Embargo",
		"Gold Rush", "Labor Strike", "War Economy", "Depleted World",
		"Agrarian Breadbasket", "Command Economy", "Free Port", "Scavenger Economy",
	}
	if int(e) < len(names) {
		return names[e]
	}
	return "Unknown"
}

// --- SOCIAL STATUSES (0-19) ---
type SocialStatus int

const (
	SocCosmopolitan      SocialStatus = 0
	SocXenophobic        SocialStatus = 1
	SocReligiousPilgrimage SocialStatus = 2
	SocPlague            SocialStatus = 3
	SocBrainDrain        SocialStatus = 4
	SocPenalColony       SocialStatus = 5
	SocRefugeeCrisis     SocialStatus = 6
	SocArtistEnclave     SocialStatus = 7
	SocCyberneticAscension SocialStatus = 8
	SocPreIndustrial     SocialStatus = 9
	SocAcademy           SocialStatus = 10
	SocSlum              SocialStatus = 11
	SocFrontierSpirit    SocialStatus = 12
	SocDecadentAristocracy SocialStatus = 13
	SocWorkerRebellion   SocialStatus = 14
	SocCultActivity      SocialStatus = 15
	SocScientificExpedition SocialStatus = 16
	SocGhostTown         SocialStatus = 17
	SocGladiatorial      SocialStatus = 18
	SocHiveMind          SocialStatus = 19
)

func (s SocialStatus) String() string {
	names := [...]string{
		"Cosmopolitan Metropolis", "Xenophobic / Closed", "Religious Pilgrimage", "Plague / Quarantine",
		"Brain Drain", "Penal Colony", "Refugee Crisis", "Artist Enclave",
		"Cybernetic Ascension", "Pre-Industrial Society", "Academy / University", "Slum / Overpopulated",
		"Frontier Spirit", "Decadent Aristocracy", "Worker Rebellion", "Cult Activity",
		"Scientific Expedition", "Ghost Town", "Gladiatorial Culture", "Hive Mind",
	}
	if int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}