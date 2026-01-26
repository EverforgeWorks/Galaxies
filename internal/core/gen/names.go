package gen

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers = []rune("0123456789")
)

func init() {
	// Seed once on startup
	rand.Seed(time.Now().UnixNano())
}

// GenerateSystemName creates a procedural identifier for a solar system.
// Format: LL-NNNN-L-NNNN (e.g., AB-1234-C-5678)
func GenerateSystemName() string {
	l1 := pick(letters, 2)
	n1 := pick(numbers, 4)
	l2 := pick(letters, 1)
    n2 := pick(numbers, 4)

	return fmt.Sprintf("%s-%s-%s-%s", l1, n1, l2, n2)
}

// GenerateCrewID creates a service identifier for a crew member.
// Format: LL-NNNNNN-L (e.g., KP-850212-Z)
func GenerateCrewID() string {
	l1 := pick(letters, 2)
	n1 := pick(numbers, 6)
	l2 := pick(letters, 1)

	return fmt.Sprintf("%s-%s-%s", l1, n1, l2)
}

// Helper to pick n random characters from a slice
func pick(source []rune, n int) string {
	result := make([]rune, n)
	for i := range result {
		result[i] = source[rand.Intn(len(source))]
	}
	return string(result)
}