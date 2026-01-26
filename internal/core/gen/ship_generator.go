package gen

import (
	"fmt"
	"galaxies/internal/core"
	"galaxies/internal/core/enums"
	"github.com/google/uuid"
	"math/rand"
)

func GenerateShip(qual enums.ShipQualifier, origin enums.ShipOrigin, chassis enums.ShipChassis) *core.Ship {
	// 1. Start with Chassis
	stats := GetBaseChassisStats(chassis)

	// 2. Apply Origin
	ApplyOriginMods(&stats, origin)

	// 3. Apply Qualifier
	ApplyQualifierMods(&stats, qual)

	// 4. Create Name
	// e.g. "Hardspace Imperial Interceptor"
	fullName := fmt.Sprintf("%s %s %s", qual, origin, chassis)

	return &core.Ship{
		ID:          uuid.New(),
		Name:        "Unnamed Vessel",
		ModelName:   fullName,
		Stats:       stats,
		CurrentHull: stats.MaxHull,
		CurrentFuel: stats.MaxFuel,
		// Init other fields...
	}
}

// GenerateRandomShip picks random enums for everything
func GenerateRandomShip() *core.Ship {
	q := enums.ShipQualifier(rand.Intn(8))
	o := enums.ShipOrigin(rand.Intn(8))
	c := enums.ShipChassis(rand.Intn(8))
	
	return GenerateShip(q, o, c)
}