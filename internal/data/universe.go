package data

import (
	"os"

	"galaxies/internal/core/entity"
	"galaxies/internal/core/gen"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type UniverseManifest struct {
	Stars []entity.Star `yaml:"stars"`
}

// LoadUniverse reads coordinates from YAML and generates names/IDs for each star
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
		// Generate unique identity for this star
		manifest.Stars[i].ID = uuid.New()
		manifest.Stars[i].Name = gen.GenerateStarName()
		
		universeMap[manifest.Stars[i].ID] = manifest.Stars[i]
	}

	return universeMap, nil
}
