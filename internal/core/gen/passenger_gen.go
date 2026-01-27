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

// GeneratePassenger creates a person looking for a ride based on the local system conditions.
// sourceSys: Where they are now.
// targetSysID: Where they want to go (The caller needs to pick a valid system ID).
// distance: The distance to that target (used to calc fare).
func GeneratePassenger(sourceSys *entity.System, targetSysID uuid.UUID, distance float64) *entity.Passenger {
	
	// 1. Determine Type based on System Flavor
	pType := determinePassengerType(sourceSys)

	// 2. Calculate Fare & Risk
	// Base: 15 Credits per Light Year
	baseRate := 15.0
	
	// Modifiers
	wealthMult := sourceSys.Stats.PassengerWealth
	typeMult := 1.0
	
	// Risks (0.0 - 1.0)
	inspectionRisk := 0.0 // Police attention
	interdictionRisk := 0.0 // Pirate attention (Kidnapping target)
	timeLimit := 0

	switch pType {
	case domain.PaxTourist:
		typeMult = 1.0
	case domain.PaxWorker:
		typeMult = 0.4 // Very cheap, usually bulk transport
	case domain.PaxBusiness:
		typeMult = 2.5
		timeLimit = int(distance/10) + 4 // Strict deadline
	case domain.PaxVIP:
		typeMult = 10.0 // Huge payout
		interdictionRisk = 0.15 // Pirates want them
		inspectionRisk = 0.05 // Paparazzi/Security checks
	case domain.PaxRefugee:
		typeMult = 0.1 // Charity work
	case domain.PaxPrisoner:
		typeMult = 4.0 // Hazard pay
		inspectionRisk = 0.2 // Police are watching
	case domain.PaxFugitive:
		typeMult = 25.0 // "Batshit" payout
		inspectionRisk = 0.6 // Huge risk of getting caught/fined
		timeLimit = int(distance/10) + 2 // They are running for their life
	}

	// 3. Final Fare Calculation
	// (Distance * Base * Wealth * Type) + Random Variance
	calculatedFare := int(distance * baseRate * wealthMult * typeMult)

	// 4. Generate the ID string
	idString := GeneratePassengerID()

	// 5. Build Entity
	pass := entity.NewPassenger(
		idString,
		pType,
		sourceSys.ID,
		targetSysID,
		calculatedFare,
	)
	
	// Apply calculated risks
	pass.InspectionRisk = inspectionRisk
	pass.InterdictionRisk = interdictionRisk
	pass.TimeLimit = timeLimit

	return pass
}

// determinePassengerType uses the System Stats to weight the RNG.
func determinePassengerType(sys *entity.System) domain.PassengerType {
	roll := rand.Float64()

	// 1. Check for SPECIAL Conditions first
	
	// Prison Systems spawn Prisoners (30% chance)
	if sys.Stats.HasPrison && roll < 0.30 {
		return domain.PaxPrisoner
	}

	// War Zones / Slums / Failed States spawn Refugees (40% chance)
	if sys.Stats.SlumsDensity > 2.0 && roll < 0.40 {
		return domain.PaxRefugee
	}

	// Luxury Resorts / Imperial Cores spawn VIPs (25% chance)
	if sys.Stats.VIPDensity > 1.0 && roll < 0.25 {
		return domain.PaxVIP
	}

	// 2. Standard Distribution for normal rolls
	// Re-roll for standard distribution if special didn't hit
	roll = rand.Float64()

	if roll < 0.05 {
		return domain.PaxFugitive // 5% chance anywhere to find a criminal
	} else if roll < 0.25 {
		return domain.PaxBusiness
	} else if roll < 0.60 {
		return domain.PaxWorker
	} else {
		return domain.PaxTourist
	}
}