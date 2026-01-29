package gen

import (
	"math"
	"math/rand"
	"time"

	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var chassisSpecs = map[domain.ShipChassis]domain.ChassisSpec{
	domain.ChassisInterceptor: {
		Name: "Interceptor", BaseHull: 800, BaseShield: 500,
		BaseFuel: 40, FuelEfficiency: 1.0, JumpRange: 4.0,
		BaseCargo:  10,
		BaseCabins: 0, BaseBunks: 1,
		SlotsHigh: 3, SlotsMid: 1, SlotsLow: 1, BasePrice: 15000,
	},
	domain.ChassisCorvette: {
		Name: "Corvette", BaseHull: 1500, BaseShield: 1500,
		BaseFuel: 50, FuelEfficiency: 1.2, JumpRange: 6.0,
		BaseCargo:  15,
		BaseCabins: 1, BaseBunks: 6,
		SlotsHigh: 4, SlotsMid: 3, SlotsLow: 2, BasePrice: 250000,
	},
	domain.ChassisGunship: {
		Name: "Gunship", BaseHull: 2500, BaseShield: 2000,
		BaseFuel: 80, FuelEfficiency: 1.6, JumpRange: 7.0,
		BaseCargo:  5,
		BaseCabins: 0, BaseBunks: 3,
		SlotsHigh: 6, SlotsMid: 2, SlotsLow: 4, BasePrice: 350000,
	},
	domain.ChassisHauler: {
		Name: "Hauler", BaseHull: 2000, BaseShield: 800,
		BaseFuel: 100, FuelEfficiency: 1.5, JumpRange: 10.0,
		BaseCargo:  35,
		BaseCabins: 2, BaseBunks: 4,
		SlotsHigh: 0, SlotsMid: 2, SlotsLow: 4, BasePrice: 45000,
	},
	domain.ChassisProspector: {
		Name: "Prospector", BaseHull: 1200, BaseShield: 600,
		BaseFuel: 100, FuelEfficiency: 1.1, JumpRange: 12.0,
		BaseCargo:  25,
		BaseCabins: 1, BaseBunks: 2,
		SlotsHigh: 1, SlotsMid: 4, SlotsLow: 3, BasePrice: 60000,
	},
	domain.ChassisBarge: {
		Name: "Barge", BaseHull: 5000, BaseShield: 1000,
		BaseFuel: 400, FuelEfficiency: 5.0, JumpRange: 8.0,
		BaseCargo:  50,
		BaseCabins: 5, BaseBunks: 10,
		SlotsHigh: 0, SlotsMid: 1, SlotsLow: 6, BasePrice: 80000,
	},
	domain.ChassisCourier: {
		Name: "Courier", BaseHull: 400, BaseShield: 350,
		BaseFuel: 70, FuelEfficiency: 0.5, JumpRange: 18.0,
		BaseCargo:  12,
		BaseCabins: 0, BaseBunks: 1,
		SlotsHigh: 1, SlotsMid: 2, SlotsLow: 1, BasePrice: 35000,
	},
	domain.ChassisYacht: {
		Name: "Yacht", BaseHull: 600, BaseShield: 1200,
		BaseFuel: 250, FuelEfficiency: 3.0, JumpRange: 25.0,
		BaseCargo:  20,
		BaseCabins: 8, BaseBunks: 4,
		SlotsHigh: 0, SlotsMid: 4, SlotsLow: 2, BasePrice: 120000,
	},
}

func GenerateShip(name string) *entity.Ship {
	randomChassis := domain.ShipChassis(rand.Intn(8))
	return GenerateShipByChassis(randomChassis, name)
}

func GenerateShipByChassis(chassisID domain.ShipChassis, name string) *entity.Ship {
	originID := domain.ShipOrigin(rand.Intn(8))
	qualID := domain.ShipQualifier(rand.Intn(8))

	spec := chassisSpecs[chassisID]

	stats := entity.ShipStats{
		MaxHull:         spec.BaseHull,
		MaxShield:       spec.BaseShield,
		MaxFuel:         spec.BaseFuel,
		FuelEfficiency:  spec.FuelEfficiency,
		JumpRange:       spec.JumpRange,
		CargoVolume:     spec.BaseCargo,
		PassengerCabins: spec.BaseCabins,
		CrewBunks:       spec.BaseBunks,
		HighSlots:       spec.SlotsHigh,
		MidSlots:        spec.SlotsMid,
		LowSlots:        spec.SlotsLow,
		Cost:            spec.BasePrice,
		ShieldRegen:     5.0,
		StealthRating:   0.0,
		BaseAccuracy:    1.0,
		DamageBonus:     1.0,
		MaxPowerGrid:    100,
	}

	applyOriginMods(&stats, originID)
	applyQualifierMods(&stats, qualID)

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

func round(val float64) float64 {
	return math.Round(val)
}

func applyOriginMods(s *entity.ShipStats, o domain.ShipOrigin) {
	switch o {
	case domain.OriginOuterRim:
		s.FuelEfficiency *= 0.85
		s.LowSlots += 1
		s.StealthRating += 0.2
		s.MaxShield = round(s.MaxShield * 0.8) // FIXED
		s.BaseAccuracy *= 0.9

	case domain.OriginImperial:
		s.MaxHull = round(s.MaxHull * 1.3) // FIXED
		s.DamageBonus *= 1.15
		s.HighSlots += 1
		s.StealthRating -= 0.3
		s.FuelEfficiency *= 1.2

	case domain.OriginVoid:
		s.ShieldRegen *= 1.4
		s.JumpRange += 4.0
		s.StealthRating += 0.3
		s.MaxHull = round(s.MaxHull * 0.7) // FIXED
		s.Cost = int(float64(s.Cost) * 1.5)

	case domain.OriginIndustrial:
		s.CargoVolume = int(float64(s.CargoVolume) * 1.5)
		s.MaxHull = round(s.MaxHull * 1.2) // FIXED
		s.LowSlots += 1
		s.JumpRange -= 2.0
		s.ShieldRegen *= 0.8

	case domain.OriginCorporate:
		s.MaxPowerGrid += 50
		s.MidSlots += 1
		s.BaseAccuracy *= 1.1
		s.MaxHull = round(s.MaxHull * 0.9) // FIXED
		s.Cost = int(float64(s.Cost) * 1.3)

	case domain.OriginScientific:
		s.JumpRange += 5.0
		s.MaxShield = round(s.MaxShield * 1.2) // FIXED
		s.MidSlots += 2
		s.DamageBonus *= 0.8
		s.CargoVolume = int(float64(s.CargoVolume) * 0.7)

	case domain.OriginSmuggler:
		s.StealthRating += 0.4
		s.FuelEfficiency *= 0.9
		s.JumpRange += 2.0
		s.MaxHull = round(s.MaxHull * 0.8) // FIXED
		s.HighSlots = max(0, s.HighSlots-1)

	case domain.OriginClerical:
		s.PassengerCabins += 3
		s.MaxShield = round(s.MaxShield * 1.15) // FIXED
		s.CrewBunks += 2
		s.DamageBonus *= 0.7
		s.HighSlots = max(0, s.HighSlots-1)
	}
}

func applyQualifierMods(s *entity.ShipStats, q domain.ShipQualifier) {
	switch q {
	case domain.QualHardspace:
		s.MaxHull = round(s.MaxHull * 1.2) // FIXED
		s.MaxFuel = round(s.MaxFuel * 1.2) // FIXED
		s.LowSlots += 1
		s.JumpRange -= 1.5
		s.ShieldRegen *= 0.9

	case domain.QualLuxury:
		s.PassengerCabins += 2
		s.MaxShield = round(s.MaxShield * 1.2) // FIXED
		s.MaxPowerGrid += 30
		s.Cost *= 4
		s.CargoVolume = int(float64(s.CargoVolume) * 0.5)

	case domain.QualSurplus:
		s.Cost = int(float64(s.Cost) * 0.5)
		s.HighSlots += 1
		s.MaxHull = round(s.MaxHull * 1.1) // FIXED
		s.FuelEfficiency *= 1.3
		s.StealthRating -= 0.2

	case domain.QualPrototype:
		s.MaxPowerGrid += 100
		s.DamageBonus *= 1.2
		s.JumpRange += 3.0
		s.Cost *= 3.0
		s.MaxHull = round(s.MaxHull * 0.8) // FIXED

	case domain.QualRetrofitted:
		s.LowSlots += 1
		s.MidSlots += 1
		s.CargoVolume = int(float64(s.CargoVolume) * 1.2)
		s.MaxHull = round(s.MaxHull * 0.9) // FIXED
		s.BaseAccuracy *= 0.9

	case domain.QualRusty:
		s.Cost = int(float64(s.Cost) * 0.25)
		s.StealthRating += 0.2
		s.LowSlots += 1
		s.MaxShield = round(s.MaxShield * 0.6) // FIXED
		s.ShieldRegen *= 0.5

	case domain.QualReinforced:
		s.MaxHull = round(s.MaxHull * 1.4)     // FIXED
		s.MaxShield = round(s.MaxShield * 1.1) // FIXED
		s.LowSlots += 1
		s.JumpRange -= 2.5
		s.FuelEfficiency *= 1.2

	case domain.QualTuned:
		s.FuelEfficiency *= 0.7
		s.JumpRange += 2.0
		s.BaseAccuracy *= 1.1
		s.MaxHull = round(s.MaxHull * 0.8) // FIXED
		s.Cost = int(float64(s.Cost) * 1.5)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
