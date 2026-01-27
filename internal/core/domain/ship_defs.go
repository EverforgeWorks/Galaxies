package domain

// 1. CHASSIS (The Noun) - Defines the BASE stats and Slots
type ShipChassis int

const (
	ChassisInterceptor ShipChassis = 0
	ChassisHauler      ShipChassis = 1
	ChassisYacht       ShipChassis = 2
	ChassisCorvette    ShipChassis = 3
	ChassisProspector  ShipChassis = 4
	ChassisCourier     ShipChassis = 5
	ChassisBarge       ShipChassis = 6
	ChassisGunship     ShipChassis = 7
)

func (c ShipChassis) String() string {
	return [...]string{"Interceptor", "Hauler", "Yacht", "Corvette", "Prospector", "Courier", "Barge", "Gunship"}[c]
}

// 2. ORIGIN (The Profession/Location) - Applies SPECIALIZATION Modifiers
type ShipOrigin int

const (
	OriginOuterRim   ShipOrigin = 0
	OriginImperial   ShipOrigin = 1
	OriginVoid       ShipOrigin = 2
	OriginIndustrial ShipOrigin = 3
	OriginClerical   ShipOrigin = 4
	OriginCorporate  ShipOrigin = 5
	OriginScientific ShipOrigin = 6
	OriginSmuggler   ShipOrigin = 7
)

func (o ShipOrigin) String() string {
	return [...]string{"Outer Rim", "Imperial", "Void", "Industrial", "Clerical", "Corporate", "Scientific", "Smuggler"}[o]
}

// 3. QUALIFIER (The Adjective) - Applies GLOBAL QUALITY Multipliers
type ShipQualifier int

const (
	QualHardspace   ShipQualifier = 0
	QualLuxury      ShipQualifier = 1
	QualSurplus     ShipQualifier = 2
	QualPrototype   ShipQualifier = 3
	QualRetrofitted ShipQualifier = 4
	QualRusty       ShipQualifier = 5
	QualReinforced  ShipQualifier = 6
	QualTuned       ShipQualifier = 7
)

func (q ShipQualifier) String() string {
	return [...]string{"Hardspace", "Luxury", "Surplus", "Prototype", "Retrofitted", "Rusty", "Reinforced", "Tuned"}[q]
}

// --- STRUCT DEFINITIONS ---

// ChassisSpec represents the base stats template for a ship hull.
// Renamed from ShipChassis to avoid conflict with the Enum above.
type ChassisSpec struct {
	Name           string
	BaseHull       float64
	BaseShield     float64
	BaseFuel       float64
	FuelEfficiency float64
	BaseCargo      int
	BaseCabins     int
	BaseBunks      int
	SlotsHigh      int
	SlotsMid       int
	SlotsLow       int
	BasePrice      int
}