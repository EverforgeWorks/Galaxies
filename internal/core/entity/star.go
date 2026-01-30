package entity

import (
	"github.com/google/uuid"
)

type Star struct {
	ID   uuid.UUID `json:"id" yaml:"-"`    // Ignored in YAML, generated in Go
	Name string    `json:"name" yaml:"-"`  // Ignored in YAML, generated in Go
	X    int       `json:"x" yaml:"x"`
	Y    int       `json:"y" yaml:"y"`
}
