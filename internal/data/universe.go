package data

import (
	"fmt"
	"os"

	"galaxies/internal/core/entity"
	"galaxies/internal/core/gen"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type UniverseManifest struct {
	Stars []entity.Star `yaml:"stars"`
}

// LoadUniverse reads coordinates and generates DETERMINISTIC names/IDs
func LoadUniverse(path string) (map[uuid.UUID]entity.Star, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest UniverseManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	universeMap := make(map[uuid.UUID]entity.Star)
	for i := range manifest.Stars {
		star := &manifest.Stars[i]

		// 1. Create a Seed String from Coordinates
		// This string is the "DNA" of the star. (e.g., "star:0:0")
		seedString := fmt.Sprintf("star:%d:%d", star.X, star.Y)

		// 2. Generate Deterministic UUID
		// We use uuid.NewMD5 to generate a UUID based on the seed string.
		// This ensures ID is always the same for these coordinates.
		star.ID = uuid.NewMD5(uuid.NameSpaceOID, []byte(seedString))

		// 3. Generate Deterministic Name
		// We initialize your new generator with the same seed string.
		rng := gen.NewSeededGenerator(seedString)
		star.Name = gen.GenerateStarName(rng)
		
		universeMap[star.ID] = *star
	}

	return universeMap, nil
}
