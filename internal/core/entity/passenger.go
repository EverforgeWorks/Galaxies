package entity

import (
	"github.com/google/uuid"
	"galaxies/internal/core/domain"
)

type Passenger struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"` // The ID string (e.g. "TRVL-8592-XK")
	
	// Who they are
	Type domain.PassengerType `json:"type"`
	
	// The Contract
	SourceSystemID uuid.UUID `json:"source_system_id"`
	TargetSystemID uuid.UUID `json:"target_system_id"`
	Fare           int       `json:"fare"`
	
	// Gameplay Modifiers (Active while on board)
	InspectionRisk float64 `json:"inspection_risk"` // Adds to System Inspection chance
	InterdictionRisk float64 `json:"interdiction_risk"` // Adds to Piracy chance (VIPs are targets)
	TimeLimit      int     `json:"time_limit"`      // Jumps remaining before failure (0 = none)
}

// NewPassenger creates a passenger entity.
func NewPassenger(name string, pType domain.PassengerType, source, target uuid.UUID, fare int) *Passenger {
	return &Passenger{
		ID:             uuid.New(),
		Name:           name,
		Type:           pType,
		SourceSystemID: source,
		TargetSystemID: target,
		Fare:           fare,
		InspectionRisk: 0.0,
		InterdictionRisk: 0.0,
		TimeLimit:      0,
	}
}