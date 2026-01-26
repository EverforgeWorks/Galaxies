package gen

import (
	"fmt"
	"math/rand"
	"time"
	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateShip creates a unique ship using a strict priority pipeline:
// 1. Chassis (Base Stats / The Noun)
// 2. Origin (Specialization / The Profession)
// 3. Qualifier (Condition / The Adjective)
func GenerateShip(chassis ShipChassis, origin ShipOrigin, qualifier ShipQualifier) *entity.Ship {
	// STEP 1: CHASSIS (Establish the Baseline)
	stats := rollChassisStats(chassis)

	// STEP 2: ORIGIN (Apply Professional Modifications)
	// Origins tend to shift focus (e.g. trading Hull for Speed)
	applyOriginMods(&stats, origin)

	// STEP 3: QUALIFIER (Apply Condition/State)
	// Qualifiers tend to be global multipliers (e.g. Rusty = -20% everything)
	applyQualifierMods(&stats, qualifier)

	// STEP 4: COST (Calculate Final Value)
	stats.Cost = calculateShipValue(stats)

	// Construct Name
	fullName := fmt.Sprintf("%s %s %s", qualifier.String(), origin.String(), chassis.String())

	ship := &entity.Ship{
		ID:            uuid.New(),
		Name:          fullName,
		ModelName:     fullName,
		CurrentHull:   stats.MaxHull,
		CurrentShield: stats.MaxShield,
		CurrentFuel:   stats.MaxFuel,
		Stats:         stats,
	}

	return ship
}

// -----------------------------------------------------------------------------
// HELPER: Variance & Math
// -----------------------------------------------------------------------------

// randRange returns a random float between min and max.
func randRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// applyMult applies a percentage multiplier with slight variance.
// Logic: Multiplicative. Used for "Scale" stats (Hull, Fuel).
// e.g. val=100, mult=1.1 (10% boost) -> returns approx 110.
func applyMult(val float64, mult float64) float64 {
	variance := 0.05 // 5% wiggle room on the modifier
	actualMult := mult + randRange(-variance, variance)
	return val * actualMult
}

// applyFlat adds a fixed amount.
// Logic: Additive. Used for "Count" stats (Slots, Cabins) or "Ratings" (Stealth).
func applyFlatInt(val int, amount int) int {
	// Ensure we don't drop below 0
	res := val + amount
	if res < 0 { return 0 }
	return res
}

func applyFlatFloat(val float64, amount float64) float64 {
	res := val + amount
	if res < 0.0 { return 0.0 }
	return res
}

// -----------------------------------------------------------------------------
// 1. CHASSIS ROLLS (The Base Ranges)
// -----------------------------------------------------------------------------
func rollChassisStats(c ShipChassis) entity.ShipStats {
	// Initialize defaults
	s := entity.ShipStats{
		FuelEfficiency: 1.0, 
		BaseAccuracy:   0.8,
		DamageBonus:    1.0,
		ShieldRegen:    1.0, // Default 1 hp/tick
	}

	switch c {
	case ChassisInterceptor:
		s.MaxHull = randRange(550, 650)
		s.MaxShield = randRange(350, 450)
		s.ShieldRegen = randRange(2.0, 3.0) // Fast Regen
		s.StealthRating = randRange(0.25, 0.35)
		s.CargoVolume = int(randRange(15, 25))
		s.MaxFuel = randRange(180, 220)
		s.JumpRange = randRange(14, 16)
		s.HighSlots = 3; s.MidSlots = 1; s.LowSlots = 1
		s.MaxPowerGrid = int(randRange(140, 160))
		s.CrewBunks = 1; s.PassengerCabins = 0

	case ChassisHauler:
		s.MaxHull = randRange(1400, 1600)
		s.MaxShield = randRange(750, 850)
		s.ShieldRegen = randRange(1.0, 1.5)
		s.StealthRating = 0.0
		s.CargoVolume = int(randRange(750, 850))
		s.MaxFuel = randRange(750, 850)
		s.JumpRange = randRange(22, 28)
		s.HighSlots = 0; s.MidSlots = 2; s.LowSlots = 4
		s.MaxPowerGrid = int(randRange(280, 320))
		s.CrewBunks = 3; s.PassengerCabins = 2

	case ChassisYacht:
		s.MaxHull = randRange(380, 420)
		s.MaxShield = randRange(550, 650)
		s.ShieldRegen = randRange(1.5, 2.5)
		s.StealthRating = randRange(0.05, 0.15)
		s.CargoVolume = int(randRange(40, 60))
		s.MaxFuel = randRange(280, 320)
		s.JumpRange = randRange(28, 32)
		s.HighSlots = 1; s.MidSlots = 3; s.LowSlots = 2
		s.MaxPowerGrid = int(randRange(190, 210))
		s.CrewBunks = 2; s.PassengerCabins = 6

	case ChassisCorvette:
		s.MaxHull = randRange(1100, 1300)
		s.MaxShield = randRange(900, 1100)
		s.ShieldRegen = randRange(3.0, 4.0) // Military Grade
		s.StealthRating = randRange(0.05, 0.15)
		s.CargoVolume = int(randRange(90, 110))
		s.MaxFuel = randRange(380, 420)
		s.JumpRange = randRange(18, 22)
		s.HighSlots = 3; s.MidSlots = 3; s.LowSlots = 2
		s.MaxPowerGrid = int(randRange(430, 470))
		s.CrewBunks = 4; s.PassengerCabins = 0

	case ChassisProspector:
		s.MaxHull = randRange(950, 1050)
		s.MaxShield = randRange(450, 550)
		s.ShieldRegen = randRange(1.0, 1.5)
		s.StealthRating = randRange(0.15, 0.25)
		s.CargoVolume = int(randRange(280, 320))
		s.MaxFuel = randRange(480, 520)
		s.JumpRange = randRange(18, 22)
		s.HighSlots = 2; s.MidSlots = 2; s.LowSlots = 3
		s.MaxPowerGrid = int(randRange(330, 370))
		s.CrewBunks = 2; s.PassengerCabins = 0

	case ChassisCourier:
		s.MaxHull = randRange(180, 220)
		s.MaxShield = randRange(180, 220)
		s.ShieldRegen = randRange(4.0, 5.0) // Very fast regen (hit and run)
		s.StealthRating = randRange(0.45, 0.55)
		s.CargoVolume = int(randRange(35, 45))
		s.MaxFuel = randRange(140, 160)
		s.JumpRange = randRange(32, 38)
		s.HighSlots = 1; s.MidSlots = 1; s.LowSlots = 3
		s.MaxPowerGrid = int(randRange(110, 130))
		s.CrewBunks = 1; s.PassengerCabins = 1

	case ChassisBarge:
		s.MaxHull = randRange(2800, 3200)
		s.MaxShield = randRange(450, 550)
		s.ShieldRegen = randRange(0.5, 0.8) // Very slow
		s.StealthRating = 0.0
		s.CargoVolume = int(randRange(2400, 2600))
		s.MaxFuel = randRange(950, 1050)
		s.JumpRange = randRange(8, 12)
		s.HighSlots = 1; s.MidSlots = 1; s.LowSlots = 5
		s.MaxPowerGrid = int(randRange(380, 420))
		s.CrewBunks = 5; s.PassengerCabins = 0

	case ChassisGunship:
		s.MaxHull = randRange(1900, 2100)
		s.MaxShield = randRange(750, 850)
		s.ShieldRegen = randRange(2.0, 2.5)
		s.StealthRating = 0.0
		s.CargoVolume = int(randRange(70, 90))
		s.MaxFuel = randRange(280, 320)
		s.JumpRange = randRange(10, 14)
		s.HighSlots = 5; s.MidSlots = 2; s.LowSlots = 2
		s.MaxPowerGrid = int(randRange(580, 620))
		s.CrewBunks = 3; s.PassengerCabins = 0
	}
	return s
}

// -----------------------------------------------------------------------------
// 2. ORIGIN MODIFIERS (The Profession)
// -----------------------------------------------------------------------------
func applyOriginMods(s *entity.ShipStats, o ShipOrigin) {
	switch o {
	case OriginOuterRim: // Rugged: Better fuel, worse shield, better hull
		s.FuelEfficiency = applyMult(s.FuelEfficiency, 0.8) 
		s.MaxShield = applyMult(s.MaxShield, 0.8)
		s.MaxHull = applyMult(s.MaxHull, 1.1)

	case OriginImperial: // Military: Heavy, loud, power hungry
		s.MaxHull = applyMult(s.MaxHull, 1.2)
		s.MaxPowerGrid = int(applyMult(float64(s.MaxPowerGrid), 1.15))
		s.StealthRating = 0.0 // OVERRIDE: Imperial ships are never stealthy

	case OriginVoid: // Stealth: Fragile but hidden
		s.StealthRating = applyFlatFloat(s.StealthRating, 0.3) // Additive boost
		s.MaxHull = applyMult(s.MaxHull, 0.7)
		s.MaxShield = applyMult(s.MaxShield, 0.7)

	case OriginIndustrial: // Worker: Cargo focus
		s.CargoVolume = int(applyMult(float64(s.CargoVolume), 1.3))
		s.LowSlots = applyFlatInt(s.LowSlots, 1) // Additive Slot

	case OriginClerical: // Peaceful: Passengers
		s.PassengerCabins = applyFlatInt(s.PassengerCabins, 2)
		s.CrewBunks = applyFlatInt(s.CrewBunks, 2)
		s.HighSlots = applyFlatInt(s.HighSlots, -1) // Penalty

	case OriginCorporate: // Efficient: Cheap mass produce
		s.MaxPowerGrid = int(applyMult(float64(s.MaxPowerGrid), 1.1))
		s.FuelEfficiency = applyMult(s.FuelEfficiency, 0.9)
		s.MaxHull = applyMult(s.MaxHull, 0.9)

	case OriginScientific: // Tech: Sensors/Jump
		s.MaxShield = applyMult(s.MaxShield, 1.25)
		s.JumpRange = applyMult(s.JumpRange, 1.1)
		s.CargoVolume = int(applyMult(float64(s.CargoVolume), 0.8))

	case OriginSmuggler: // Illegal: Fast/Hidden
		s.StealthRating = applyFlatFloat(s.StealthRating, 0.2)
		s.MaxHull = applyMult(s.MaxHull, 0.8)
	}
}

// -----------------------------------------------------------------------------
// 3. QUALIFIER MODIFIERS (The Condition)
// -----------------------------------------------------------------------------
func applyQualifierMods(s *entity.ShipStats, q ShipQualifier) {
	switch q {
	case QualHardspace: // Durable
		s.MaxHull = applyMult(s.MaxHull, 1.3)
		// Tradeoff: Heavier ships jump shorter
		s.JumpRange = applyMult(s.JumpRange, 0.9)

	case QualLuxury: // High End
		s.PassengerCabins = applyFlatInt(s.PassengerCabins, 1)
		s.MaxPowerGrid = int(applyMult(float64(s.MaxPowerGrid), 1.1))

	case QualSurplus: // Cheap / Used
		s.MaxHull = applyMult(s.MaxHull, 0.9)
		s.MaxShield = applyMult(s.MaxShield, 0.8)
		s.DamageBonus = applyMult(s.DamageBonus, 0.9)

	case QualPrototype: // Experimental
		s.MaxPowerGrid = applyFlatInt(s.MaxPowerGrid, 150) // Flat boost
		s.MaxHull = applyMult(s.MaxHull, 0.8)
		s.ShieldRegen = applyMult(s.ShieldRegen, 1.2)

	case QualRetrofitted: // Modded
		s.LowSlots = applyFlatInt(s.LowSlots, 1)
		s.MaxHull = applyMult(s.MaxHull, 0.85)

	case QualRusty: // Garbage
		s.MaxHull = applyMult(s.MaxHull, 0.7)
		s.MaxShield = applyMult(s.MaxShield, 0.6)
		s.BaseAccuracy = applyMult(s.BaseAccuracy, 0.8)
		s.ShieldRegen = applyMult(s.ShieldRegen, 0.5)

	case QualReinforced: // Tank
		s.MaxHull = applyMult(s.MaxHull, 1.4)
		s.JumpRange = applyMult(s.JumpRange, 0.8) // Heavy armor reduces range

	case QualTuned: // Performance
		s.ShieldRegen = applyMult(s.ShieldRegen, 1.2)
		s.BaseAccuracy = applyMult(s.BaseAccuracy, 1.1)
		s.MaxHull = applyMult(s.MaxHull, 0.9)
	}
}

// -----------------------------------------------------------------------------
// 4. COST CALCULATOR
// -----------------------------------------------------------------------------
func calculateShipValue(s entity.ShipStats) int {
	val := 0.0

	val += s.MaxHull * 10
	val += s.MaxShield * 15
	val += s.ShieldRegen * 500
	val += float64(s.CargoVolume) * 50
	val += s.MaxFuel * 5
	val += s.JumpRange * 200
	val += float64(s.MaxPowerGrid) * 20
	
	// Fuel Eff is inverted: Lower is better (0.8 is better than 1.0)
	// We add value if efficiency is < 1.0
	if s.FuelEfficiency < 1.0 {
		val += (1.0 - s.FuelEfficiency) * 5000
	}

	val += float64(s.HighSlots) * 5000
	val += float64(s.MidSlots) * 3000
	val += float64(s.LowSlots) * 2000
	val += float64(s.PassengerCabins) * 1000
	val += float64(s.CrewBunks) * 500

	if s.StealthRating > 0 {
		val += s.StealthRating * 15000
	}
	
	val += s.BaseAccuracy * 1000
	val += s.DamageBonus * 1000

	return int(val)
}