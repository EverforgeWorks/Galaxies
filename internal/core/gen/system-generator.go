package gen

import (
	"galaxies/internal/core"
	"galaxies/internal/core/enums"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// GenerateSystemConfig allows you to force specific statuses.
// Pass -1 to any field to make it Random.
type GenerateSystemConfig struct {
	Name      string // Leave empty to auto-generate
	X, Y      int
	Political enums.PoliticalStatus // Pass -1 for Random
	Economic  enums.EconomicStatus  // Pass -1 for Random
	Social    enums.SocialStatus    // Pass -1 for Random
}

func GenerateSystem(config GenerateSystemConfig) *core.System {
	rand.Seed(time.Now().UnixNano())

	// 1. Handle Randomization
	// We cast -1 to the Enum type to check if it's "Unset"
	pol := config.Political
	if int(pol) == -1 {
		pol = enums.PoliticalStatus(rand.Intn(20)) // 0-19
	}

	eco := config.Economic
	if int(eco) == -1 {
		eco = enums.EconomicStatus(rand.Intn(20))
	}

	soc := config.Social
	if int(soc) == -1 {
		soc = enums.SocialStatus(rand.Intn(20))
	}

	// 2. Handle Name
	name := config.Name
	if name == "" {
		name = GenerateName()
	}

	// 3. Create Base Stats
	stats := core.NewDefaultSystemStats()

	// 4. Apply The Modifiers (The "Stacking" Logic)
	// Order matters slightly: We do Political -> Economic -> Social
	ApplyPoliticalMods(&stats, pol)
	ApplyEconomicMods(&stats, eco)
	ApplySocialMods(&stats, soc)

	// 5. Resolve Conflicts / Clamp Values
	// This "decides what to do" if math went weird
	FinalizeStats(&stats)

	// 6. Create the System Struct
	sys := &core.System{
		ID:        uuid.New(),
		Name:      name,
		X:         config.X,
		Y:         config.Y,
		Political: pol,
		Economic:  eco,
		Social:    soc,
		Stats:     stats,
	}

	// 7. POPULATE (Critical Step: Calculates actual counts based on density)
	PopulateSystem(sys)

	return sys
}

// FinalizeStats acts as the safety net for the stacking logic
func FinalizeStats(s *core.SystemStats) {
	// Logic: Tax can never be negative, and rarely above 40% unless specified
	if s.TaxRate < 0 {
		s.TaxRate = 0
	}

	// Logic: Probabilities max out at 100%
	if s.PiracyChance > 1.0 {
		s.PiracyChance = 1.0
	}
	if s.PiracyChance < 0.0 {
		s.PiracyChance = 0.0
	}
	if s.InspectionChance > 1.0 {
		s.InspectionChance = 1.0
	}
	if s.InspectionChance < 0.0 {
		s.InspectionChance = 0.0
	}

	// Logic: Prices can't be zero
	if s.MarketBuyMult < 0.1 {
		s.MarketBuyMult = 0.1
	}
	
	// Logic: Prevent Extreme Deflation (The "14% Price" exploit fix)
	if s.MarketSellMult < 0.25 {
		s.MarketSellMult = 0.25
	}

	// Logic: If a system has NO population, it can't have VIPs or Slums
	if s.PassengerDensity <= 0 {
		s.VIPDensity = 0
		s.SlumsDensity = 0
		s.WantedPassChance = 0
	}

	// Logic: If Inspection Chance is 100%, Black Market is effectively closed (or very hard)
	if s.InspectionChance >= 1.0 {
		s.BlackMarketBuyMult *= 0.5 // High risk, low payout
	}
}