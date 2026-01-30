package data

import (
	"os"

	"galaxies/internal/core/entity"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type UniverseManifest struct {
	Stars []entity.Star `yaml:"stars"`
}

// LoadUniverse is now a "dumb" loader. 
// It relies on the YAML file being the single source of truth.
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
	for _, star := range manifest.Stars {
		// Validation check: ensure gen-universe tool was run
		if star.ID == uuid.Nil {
			// Log a warning or panic in dev mode
			// For now, we just skip invalid entries
			continue
		}
		universeMap[star.ID] = star
	}

	return universeMap, nil
}
