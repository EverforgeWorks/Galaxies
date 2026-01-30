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

// LoadUniverse reads coordinates and generates STABLE IDs based on location
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

		// 1. Generate Deterministic ID from Coordinates (e.g., "star:0:0")
		// This ensures the ID is the same in YAML, DB, and RAM.
		seed := fmt.Sprintf("star:%d:%d", star.X, star.Y)
		star.ID = uuid.NewMD5(uuid.NameSpaceOID, []byte(seed))

		// 2. Generate Deterministic Name
		// (Optional: If the YAML already has a name, you could skip this)
		if star.Name == "" {
			rng := gen.NewSeededGenerator(seed)
			star.Name = gen.GenerateStarName(rng)
		}
		
		universeMap[star.ID] = *star
	}

	return universeMap, nil
}
