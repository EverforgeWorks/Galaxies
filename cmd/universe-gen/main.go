package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"galaxies/internal/core/entity"
	"galaxies/internal/core/gen"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type UniverseManifest struct {
	Stars []entity.Star `yaml:"stars"`
}

func main() {
	filePath := flag.String("file", "internal/data/universe.yaml", "Path to universe.yaml")
	addCoords := flag.String("add", "", "Comma-separated coordinates to add (e.g. '1,1;2,4;-5,10')")
	flag.Parse()

	// 1. Load Existing
	manifest, err := loadManifest(*filePath)
	if err != nil {
		log.Printf("Creating new manifest (could not load: %v)", err)
		manifest = &UniverseManifest{}
	}

	// 2. Add New Coordinates if requested
	if *addCoords != "" {
		pairs := strings.Split(*addCoords, ";")
		for _, pair := range pairs {
			parts := strings.Split(pair, ",")
			if len(parts) != 2 {
				log.Printf("Skipping invalid coord: %s", pair)
				continue
			}
			x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

			// Check for duplicates
			if exists(manifest, x, y) {
				log.Printf("Star at %d,%d already exists. Skipping.", x, y)
				continue
			}

			manifest.Stars = append(manifest.Stars, entity.Star{X: x, Y: y})
			fmt.Printf("Added placeholder for Star at [%d, %d]\n", x, y)
		}
	}

	// 3. Hydrate (Fill missing IDs/Names)
	// This makes the generation persistent. Once generated, it stays in the YAML.
	count := 0
	for i := range manifest.Stars {
		star := &manifest.Stars[i]
		
		dirty := false
		if star.ID == uuid.Nil {
			// Generate ID based on coordinates (Deterministic Seed)
			seed := fmt.Sprintf("star:%d:%d", star.X, star.Y)
			star.ID = uuid.NewMD5(uuid.NameSpaceOID, []byte(seed))
			dirty = true
		}
		
		if star.Name == "" {
			// Generate Name based on ID
			// We recreate the generator so it's deterministic per star
			seedString := fmt.Sprintf("star:%d:%d", star.X, star.Y)
			rng := gen.NewSeededGenerator(seedString)
			star.Name = gen.GenerateStarName(rng)
			dirty = true
		}

		if dirty {
			count++
		}
	}

	// 4. Save
	if err := saveManifest(*filePath, manifest); err != nil {
		log.Fatalf("Failed to save universe: %v", err)
	}

	fmt.Printf("Universe synchronized. Updated %d stars. Total count: %d\n", count, len(manifest.Stars))
}

func exists(m *UniverseManifest, x, y int) bool {
	for _, s := range m.Stars {
		if s.X == x && s.Y == y {
			return true
		}
	}
	return false
}

func loadManifest(path string) (*UniverseManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m UniverseManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func saveManifest(path string, m *UniverseManifest) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
