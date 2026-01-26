package gen

import (
	"math"
	"math/rand"
	"time"

	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// UniverseConfig defines the parameters for a generation run.
type UniverseConfig struct {
	MinX, MaxX   int // e.g. -10, 10
	MinY, MaxY   int // e.g. -10, 10
	SystemCount  int // Target number of systems to spawn
	MinDistance  float64 // Minimum LY between systems to prevent overlap
}

// GeneratedUniverse holds the result of the generation.
type GeneratedUniverse struct {
	Systems []*entity.System
}

// GenerateUniverse creates the initial cluster of systems.
func GenerateUniverse(cfg UniverseConfig) *GeneratedUniverse {
	systems := make([]*entity.System, 0, cfg.SystemCount)
	
	attempts := 0
	maxAttempts := cfg.SystemCount * 50 // Fail-safe to prevent infinite loops

	for len(systems) < cfg.SystemCount && attempts < maxAttempts {
		attempts++

		// 1. Pick a random coordinate within bounds
		width := cfg.MaxX - cfg.MinX + 1
		height := cfg.MaxY - cfg.MinY + 1
		
		x := rand.Intn(width) + cfg.MinX
		y := rand.Intn(height) + cfg.MinY

		// 2. Check for collisions (Minimum Distance)
		if isCrowded(x, y, systems, cfg.MinDistance) {
			continue
		}

		// 3. Generate the System Data
		// Randomize the archetypes
		pol := domain.PoliticalStatus(rand.Intn(20))
		eco := domain.EconomicStatus(rand.Intn(20))
		soc := domain.SocialStatus(rand.Intn(20))

		// Note: We leave Name empty "" so the generator picks one
		sys := GenerateSystem("", x, y, pol, eco, soc)
		
		systems = append(systems, sys)
	}

	return &GeneratedUniverse{
		Systems: systems,
	}
}

// --- Helper Functions ---

// CalculateDistance returns the Euclidean distance between two points (Light Years).
func CalculateDistance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

// isCrowded returns true if a point (x,y) is too close to any existing system.
func isCrowded(x, y int, existing []*entity.System, minDist float64) bool {
	for _, s := range existing {
		dist := CalculateDistance(x, y, s.X, s.Y)
		if dist < minDist {
			return true
		}
	}
	return false
}

// GetSystemsInRange returns a subset of systems reachable from a specific point.
// This is useful for the UI "Scanner" feature.
func GetSystemsInRange(originX, originY int, rangeLY float64, allSystems []*entity.System) []*entity.System {
	var visible []*entity.System
	
	for _, sys := range allSystems {
		// Don't include yourself (distance 0) if you are sitting exactly on a system
		dist := CalculateDistance(originX, originY, sys.X, sys.Y)
		
		if dist <= rangeLY && dist > 0 {
			visible = append(visible, sys)
		}
	}
	return visible
}