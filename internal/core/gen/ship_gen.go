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

// Map the Enum (Int) to the Spec (Struct)
var chassisSpecs = map[domain.ShipChassis]domain.ChassisSpec{
	domain.ChassisInterceptor: {
		Name: "Interceptor", BaseHull: 800, BaseShield: 500, BaseFuel: 400, FuelEfficiency: 1.2,
		BaseCargo: 10, BaseCabins: 0, BaseBunks: 1,
		SlotsHigh: 3, SlotsMid: 1, SlotsLow: 1, BasePrice: 15000,
	},
	domain.ChassisHauler: {
		Name: "Hauler", BaseHull: 2000, BaseShield: 800, BaseFuel: 2000, FuelEfficiency: 3.5,
		BaseCargo: 500, BaseCabins: 2, BaseBunks: 4,
		SlotsHigh: 0, SlotsMid: 2, SlotsLow: 4, BasePrice: 45000,
	},
	domain.ChassisYacht: {
		Name: "Yacht", BaseHull: 600, BaseShield: 1200, BaseFuel: 800, FuelEfficiency: 1.0,
		BaseCargo: 50, BaseCabins: 8, BaseBunks: 4,
		SlotsHigh: 0, SlotsMid: 4, SlotsLow: 2, BasePrice: 120000,
	},
	domain.ChassisCorvette: {
		Name: "Corvette", BaseHull: 1500, BaseShield: 1500, BaseFuel: 1000, FuelEfficiency: 2.0,
		BaseCargo: 80, BaseCabins: 1, BaseBunks: 6,
		SlotsHigh: 4, SlotsMid: 3, SlotsLow: 2, BasePrice: 250000,
	},
	domain.ChassisProspector: {
		Name: "Prospector", BaseHull: 1200, BaseShield: 600, BaseFuel: 1500, FuelEfficiency: 2.5,
		BaseCargo: 200, BaseCabins: 1, BaseBunks: 2,
		SlotsHigh: 1, SlotsMid: 4, SlotsLow: 3, BasePrice: 60000,
	},
	domain.ChassisCourier: {
		Name: "Courier", BaseHull: 400, BaseShield: 400, BaseFuel: 600, FuelEfficiency: 0.5,
		BaseCargo: 20, BaseCabins: 0, BaseBunks: 1,
		SlotsHigh: 1, SlotsMid: 2, SlotsLow: 1, BasePrice: 35000,
	},
	domain.ChassisBarge: {
		Name: "Barge", BaseHull: 5000, BaseShield: 1000, BaseFuel: 5000, FuelEfficiency: 5.0,
		BaseCargo: 2000, BaseCabins: 5, BaseBunks: 10,
		SlotsHigh: 0, SlotsMid: 1, SlotsLow: 6, BasePrice: 80000,
	},
	domain.ChassisGunship: {
		Name: "Gunship", BaseHull: 2500, BaseShield: 2000, BaseFuel: 800, FuelEfficiency: 4.0,
		BaseCargo: 0, BaseCabins: 0, BaseBunks: 3,
		SlotsHigh: 6, SlotsMid: 2, SlotsLow: 4, BasePrice: 350000,
	},
}

func GenerateShip(name string) *entity.Ship {
	// 1. Randomize Components
	chassisID := domain.ShipChassis(rand.Intn(8))
	originID := domain.ShipOrigin(rand.Intn(8))
	qualID := domain.ShipQualifier(rand.Intn(8))

	spec := chassisSpecs[chassisID]

	// 2. Base Stats
	stats := entity.ShipStats{
		MaxHull:        spec.BaseHull,
		MaxShield:      spec.BaseShield,
		MaxFuel:        spec.BaseFuel,
		FuelEfficiency: spec.FuelEfficiency,
		CargoVolume:    spec.BaseCargo,
		PassengerCabins: spec.BaseCabins,
		CrewBunks:      spec.BaseBunks,
		HighSlots:      spec.SlotsHigh,
		MidSlots:       spec.SlotsMid,
		LowSlots:       spec.SlotsLow,
		Cost:           spec.BasePrice,
		
		// Defaults
		ShieldRegen:   5.0,
		StealthRating: 0.0,
		BaseAccuracy:  1.0,
		DamageBonus:   1.0,
		MaxPowerGrid:  100,
		JumpRange:     20.0,
	}

	// 3. Apply Modifiers (Simplified for MVP)
	applyOriginMods(&stats, originID)
	applyQualifierMods(&stats, qualID)

	// 4. Generate Name if empty
	modelName := originID.String() + " " + spec.Name + " (" + qualID.String() + ")"
	if name == "" {
		name = modelName
	}

	return &entity.Ship{
		ID:            uuid.New(),
		Name:          name,
		ModelName:     modelName,
		CurrentHull:   stats.MaxHull,
		CurrentShield: stats.MaxShield,
		CurrentFuel:   stats.MaxFuel,
		Stats:         stats,
		Cargo:         []entity.Item{},
		Passengers:    []entity.Passenger{},
		Crew:          []entity.Crew{},
	}
}

// Helpers to apply stat changes
func applyOriginMods(s *entity.ShipStats, o domain.ShipOrigin) {
	switch o {
	case domain.OriginOuterRim:
		s.FuelEfficiency *= 0.8 // Better eff
		s.MaxShield *= 0.8
	case domain.OriginImperial:
		s.MaxHull *= 1.2
		s.StealthRating -= 0.1
	case domain.OriginVoid:
		s.StealthRating += 0.3
		s.MaxHull *= 0.7
	case domain.OriginIndustrial:
		s.CargoVolume = int(float64(s.CargoVolume) * 1.5)
		s.MaxFuel *= 0.9
	case domain.OriginCorporate:
		s.MaxPowerGrid += 50
		s.MaxHull *= 0.9
	case domain.OriginScientific:
		s.MaxShield *= 1.3
		s.CargoVolume = int(float64(s.CargoVolume) * 0.8)
	case domain.OriginSmuggler:
		s.StealthRating += 0.2
		s.MaxHull *= 0.8
	}
}

func applyQualifierMods(s *entity.ShipStats, q domain.ShipQualifier) {
	switch q {
	case domain.QualHardspace:
		s.MaxHull *= 1.5
		s.Cost = int(float64(s.Cost) * 1.2)
	case domain.QualLuxury:
		s.Cost *= 3
	case domain.QualSurplus:
		s.Cost = int(float64(s.Cost) * 0.5)
		s.MaxHull *= 0.8
	case domain.QualPrototype:
		s.MaxPowerGrid += 100
		s.MaxShield *= 0.8
	case domain.QualRetrofitted:
		s.LowSlots += 1
		s.MaxHull *= 0.9
	case domain.QualRusty:
		s.MaxHull *= 0.7
		s.MaxShield *= 0.7
		s.Cost = int(float64(s.Cost) * 0.3)
	case domain.QualReinforced:
		s.MaxHull *= 1.3
		s.MaxFuel *= 0.9
	case domain.QualTuned:
		s.FuelEfficiency *= 0.9
		s.MaxHull *= 0.9
	}
}