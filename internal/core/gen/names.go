package gen

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
)

// NewSeededGenerator creates a deterministic random source from a string seed.
// Example: NewSeededGenerator("star:0:0") will always return the same sequence.
func NewSeededGenerator(seedInput string) *rand.Rand {
	// 1. Hash the input string to get a consistent byte sequence
	hash := md5.Sum([]byte(seedInput))
	
	// 2. Convert the first 8 bytes of the hash to an int64 seed
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	
	// 3. Create and return a new random source
	return rand.New(rand.NewSource(seed))
}

// GenerateStarName now accepts a specific random source (r).
// If r is nil, it falls back to the global random source (non-deterministic).
func GenerateStarName(r *rand.Rand) string {
	if r == nil {
		r = rand.New(rand.NewSource(rand.Int63()))
	}

	prefixes := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon", "Zeta", "Theta", "Omicron", "Omega", "Prime", "Nova", "Sol", "Vega", "Deneb", "Altair", "Sirius", "Rigel", "Proxima", "Betel", "Antares"}
	suffixes := []string{"Major", "Minor", "Prime", "System", "Cluster", "Nebula", "X", "Y", "Z", "I", "II", "III", "IV", "V", "VI", "VII", "Centauri", "Cygni", "Lyrae", "Orionis"}

	// 50% chance of a simple name (Prefix-Suffix)
	if r.Float32() < 0.5 {
		return fmt.Sprintf("%s %s", 
			prefixes[r.Intn(len(prefixes))], 
			suffixes[r.Intn(len(suffixes))])
	}

	// 50% chance of a catalog-style name (e.g., "X-291-A")
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l1 := letters[r.Intn(len(letters))]
	l2 := letters[r.Intn(len(letters))]
	num := r.Intn(9000) + 1000 // 1000-9999
	l3 := letters[r.Intn(len(letters))]

	return fmt.Sprintf("%c%c-%d-%c", l1, l2, num, l3)
}
