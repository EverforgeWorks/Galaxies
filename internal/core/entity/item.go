package entity

import (
	"galaxies/internal/core/domain"

	"github.com/google/uuid"
)

// Item represents a stack of trade goods.
// In the new simplified model, 1 Quantity = 1 Cargo Slot.
type Item struct {
	ID   uuid.UUID       `json:"id"`
	Name domain.ItemName `json:"name"`

	// Classification
	Category domain.ItemCategory `json:"category"`

	// Value & Properties
	BaseValue int  `json:"base_value"` // Global average price / current market price
	Rarity    int  `json:"rarity"`     // 1 (Common) - 10 (Artifact)
	IsIllegal bool `json:"is_illegal"` // Triggers fines/confiscation

	// State
	Quantity int     `json:"quantity"` // Number of units (Slots used)
	AvgCost  float64 `json:"avg_cost"` // Weighted average cost per unit (0 for loot/market stock)
}

// NewItem creates a fresh stack.
func NewItem(name domain.ItemName, cat domain.ItemCategory, baseVal, rarity, qty int, illegal bool) *Item {
	return &Item{
		ID:        uuid.New(),
		Name:      name,
		Category:  cat,
		BaseValue: baseVal,
		Rarity:    rarity,
		IsIllegal: illegal,
		Quantity:  qty,
		AvgCost:   0, // Default to 0, updated upon purchase
	}
}

// TotalValue calculates the worth of the entire stack at current BaseValue.
func (i *Item) TotalValue() int {
	return i.BaseValue * i.Quantity
}
