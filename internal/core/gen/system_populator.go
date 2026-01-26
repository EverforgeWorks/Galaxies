package gen

import (
	"galaxies/internal/core"
	"math"
	"math/rand"
	"time"
)

// Base Ranges for random generation
const (
	BasePaxMin = 50
	BasePaxMax = 200
	
	BaseCrewMin = 3
	BaseCrewMax = 12

	BaseVIPMin = 5
	BaseVIPMax = 20

	BaseSlumMin = 100
	BaseSlumMax = 500

	BaseAndroidMin = 5
	BaseAndroidMax = 30

	BasePrisonerMin = 20
	BasePrisonerMax = 100
)

// PopulateSystem takes the calculated Stats (Densities) and rolls the dice
// to determine the actual Counts for this generation cycle.
func PopulateSystem(s *core.System) {
	rand.Seed(time.Now().UnixNano())

	// Helper to calculate count based on Range * Multiplier
	calc := func(min, max int, density float64) int {
		if density <= 0 {
			return 0
		}
		// Roll a random number in the base range
		base := rand.Intn(max-min+1) + min
		// Apply the density multiplier
		result := float64(base) * density
		return int(math.Round(result))
	}

	// 1. Standard Population
	s.Stats.PassengerCount = calc(BasePaxMin, BasePaxMax, s.Stats.PassengerDensity)
	s.Stats.CrewPoolCount = calc(BaseCrewMin, BaseCrewMax, s.Stats.CrewPoolDensity)

	// 2. Special Populations
	// Note: Logic in modifiers.go should have set these Densities > 0 if the booleans are true
	s.Stats.VIPCount = calc(BaseVIPMin, BaseVIPMax, s.Stats.VIPDensity)
	s.Stats.SlumsCount = calc(BaseSlumMin, BaseSlumMax, s.Stats.SlumsDensity)
	s.Stats.AndroidCount = calc(BaseAndroidMin, BaseAndroidMax, s.Stats.AndroidDensity)
	s.Stats.PrisonerCount = calc(BasePrisonerMin, BasePrisonerMax, s.Stats.PrisonerDensity)

	// 3. Safety Cleanups (If flags are false, force count to 0)
	// This handles cases where Density might be high but the facility was removed
	if !s.Stats.HasLuxuryHousing { s.Stats.VIPCount = 0 }
	if !s.Stats.HasSlums { s.Stats.SlumsCount = 0 }
	if !s.Stats.HasAndroidFoundry { s.Stats.AndroidCount = 0 }
	if !s.Stats.HasPrison { s.Stats.PrisonerCount = 0 }
}