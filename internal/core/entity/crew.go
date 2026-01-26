package entity

import (
	"github.com/google/uuid"
	"galaxies/internal/core/domain"
)

type Crew struct {
	ID   uuid.UUID       `json:"id"`
	Name string          `json:"name"`
	
	Role domain.CrewRole `json:"role"`
	Type domain.CrewType `json:"type"`

	Skill  int `json:"skill"`  // 1-10
	Salary int `json:"salary"` // Credits per cycle
}

// NewCrew creates a hired hand.
func NewCrew(name string, role domain.CrewRole, cType domain.CrewType, skill, salary int) *Crew {
	return &Crew{
		ID:     uuid.New(),
		Name:   name,
		Role:   role,
		Type:   cType,
		Skill:  skill,
		Salary: salary,
	}
}