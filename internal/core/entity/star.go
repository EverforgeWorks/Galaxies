package entity

import "github.com/google/uuid"

type Star struct {
	ID   uuid.UUID `json:"id" yaml:"id"`
	Name string    `json:"name" yaml:"name"`
	X    int       `json:"x" yaml:"x"`
	Y    int       `json:"y" yaml:"y"`
	// Future expansion:
	// Resources  []Resource `json:"resources" yaml:"resources"`
	// Planets    []Planet   `json:"planets" yaml:"planets"`
}
