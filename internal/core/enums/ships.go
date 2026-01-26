package enums

// 1. CHASSIS (The Noun) - Defines the BASE stats and Slots
type ShipChassis int

const (
	ChassisInterceptor ShipChassis = 0 // Fast, Weapon Heavy, Tiny Cargo
	ChassisHauler      ShipChassis = 1 // Slow, Massive Cargo, Weak Defenses
	ChassisYacht       ShipChassis = 2 // Passenger focused, Fast, Weak Hull
	ChassisCorvette    ShipChassis = 3 // Balanced Military, Good Hull
	ChassisProspector  ShipChassis = 4 // Industrial, Mining Slots, Good Sensors
	ChassisCourier     ShipChassis = 5 // Extreme Speed, Small Cargo, Paper Hull
	ChassisBarge       ShipChassis = 6 // Massive Hull, No Speed, Industrial
	ChassisGunship     ShipChassis = 7 // Flying Tank, Slow, Heavy Weapons
)

func (c ShipChassis) String() string {
	return [...]string{"Interceptor", "Hauler", "Yacht", "Corvette", "Prospector", "Courier", "Barge", "Gunship"}[c]
}

// 2. ORIGIN (The Profession/Location) - Applies SPECIALIZATION Modifiers
type ShipOrigin int

const (
	OriginOuterRim   ShipOrigin = 0 // Rugged: +Fuel Eff, -Shields
	OriginImperial   ShipOrigin = 1 // Military: +Armor, +Heat Handling, -Stealth
	OriginVoid       ShipOrigin = 2 // Stealth: +Stealth, -Hull
	OriginIndustrial ShipOrigin = 3 // Working: +Cargo, -Speed
	OriginClerical   ShipOrigin = 4 // Peaceful: +Pass/Crew, -Weapons
	OriginCorporate  ShipOrigin = 5 // Efficient: +CPU, +Power, -Hull (Cheap build)
	OriginScientific ShipOrigin = 6 // Tech: +Sensors, +Shields, -Cargo
	OriginSmuggler   ShipOrigin = 7 // Illegal: +Speed, +Stealth, -Armor
)

func (o ShipOrigin) String() string {
	return [...]string{"Outer Rim", "Imperial", "Void", "Industrial", "Clerical", "Corporate", "Scientific", "Smuggler"}[o]
}

// 3. QUALIFIER (The Adjective) - Applies GLOBAL QUALITY Multipliers
type ShipQualifier int

const (
	QualHardspace   ShipQualifier = 0 // Durable: ++Hull, --Speed
	QualLuxury      ShipQualifier = 1 // Expensive: ++PassengerWealth, ++Cost
	QualSurplus     ShipQualifier = 2 // Cheap: --Cost, --Reliability
	QualPrototype   ShipQualifier = 3 // Experimental: ++Power/CPU, --Reliability
	QualRetrofitted ShipQualifier = 4 // Modded: +Slots, -Hull
	QualRusty       ShipQualifier = 5 // Bad: --Stats everywhere, Cheap
	QualReinforced  ShipQualifier = 6 // Tanky: +Armor, -Agility
	QualTuned       ShipQualifier = 7 // Fast: +Speed, +Evasion, -Hull
)

func (q ShipQualifier) String() string {
	return [...]string{"Hardspace", "Luxury", "Surplus", "Prototype", "Retrofitted", "Rusty", "Reinforced", "Tuned"}[q]
}