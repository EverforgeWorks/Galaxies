package gen

import (
	"galaxies/internal/core/entity"
	"galaxies/internal/core/enums"
)

// 1. BASE STATS (The Chassis)
func GetBaseChassisStats(c enums.ShipChassis) entity.ShipStats {
	// Default Baseline
	s := entity.ShipStats{
		MaxHull: 100, MaxShield: 100, MaxFuel: 100,
		Speed: 1.0, FuelEfficiency: 10.0,
		MaxCargo: 10, MaxPassengers: 0, MaxCrew: 2,
		HighSlots: 1, MidSlots: 1, LowSlots: 1,
		PowerGridOutput: 100, CPUOutput: 100,
		SensorRange: 5.0,
	}

	switch c {
	case enums.ChassisInterceptor:
		s.Speed = 2.5
		s.HighSlots = 3
		s.MaxHull = 80
		s.EvasionRating = 0.3
		s.MaxCargo = 5

	case enums.ChassisHauler:
		s.Speed = 0.6
		s.MaxCargo = 200
		s.LowSlots = 4 // Good for Cargo Expanders
		s.MaxHull = 200
		s.FuelEfficiency = 5.0 // Efficient

	case enums.ChassisYacht:
		s.Speed = 1.5
		s.MaxPassengers = 10
		s.MaxShield = 200
		s.HighSlots = 0 // No guns usually

	case enums.ChassisCorvette:
		s.MaxHull = 300
		s.HighSlots = 4
		s.MidSlots = 3
		s.PowerGridOutput = 200
		s.MaxCrew = 10

	case enums.ChassisCourier:
		s.Speed = 3.5 // Very fast
		s.MaxHull = 50
		s.StealthRating = 0.2
	}
	return s
}

// 2. ORIGIN MODIFIERS (The Profession)
func ApplyOriginMods(s *entity.ShipStats, o enums.ShipOrigin) {
	switch o {
	case enums.OriginOuterRim:
		s.FuelEfficiency *= 0.8 // 20% Better
		s.MaxHull += 50
		s.MaxShield -= 20
		s.SensorRange += 2.0

	case enums.OriginImperial:
		s.ArmorRating += 5.0
		s.ThermalHandling += 0.2
		s.StealthRating -= 0.1 // Loud and proud

	case enums.OriginVoid:
		s.StealthRating += 0.3
		s.MaxShield += 50
		s.MaxHull -= 30
		s.ThermalHandling += 0.1

	case enums.OriginIndustrial:
		s.MaxCargo += 50
		s.LowSlots += 1
		s.Speed *= 0.9

	case enums.OriginClerical:
		s.MaxPassengers += 5
		s.MaxCrew += 5
		s.HighSlots -= 1 // Less weapons
		s.MidSlots += 1  // More scanners/utility

	case enums.OriginSmuggler:
		s.StealthRating += 0.2
		s.Speed *= 1.1
		s.MaxCargo -= 10 // Hidden compartments take space
	}
}

// 3. QUALIFIER MODIFIERS (The Adjective)
func ApplyQualifierMods(s *entity.ShipStats, q enums.ShipQualifier) {
	switch q {
	case enums.QualHardspace:
		s.MaxHull *= 1.5
		s.Speed *= 0.8
		s.ArmorRating += 2.0

	case enums.QualLuxury:
		s.MaxPassengers += 2
		s.MaxShield *= 1.2
		s.Cost *= 3.0 // We need a Cost field in ShipStats or Item
	
	case enums.QualRusty:
		s.MaxHull *= 0.7
		s.Speed *= 0.9
		s.PowerGridOutput = int(float64(s.PowerGridOutput) * 0.8)
		s.Cost *= 0.5

	case enums.QualPrototype:
		s.PowerGridOutput = int(float64(s.PowerGridOutput) * 1.5)
		s.CPUOutput = int(float64(s.CPUOutput) * 1.5)
		s.MaxHull *= 0.8 // Experimental materials

	case enums.QualTuned:
		s.Speed *= 1.2
		s.EvasionRating += 0.1
		s.FuelEfficiency *= 1.2 // Burns more fuel
	}
}