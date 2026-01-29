package gen

import (
	"math/rand"

	"galaxies/internal/core/entity"
)

// GenerateMarketStock populates a system with items based on its socio-political tags.
func GenerateMarketStock(sys *entity.System) []entity.Item {
	var market []entity.Item

	for itemName, template := range GlobalItemTemplates {
		// Check if this system produces this item
		rules, exists := GlobalProductionMap[itemName]
		if !exists {
			continue
		}

		produces := false
		isBonus := false

		for _, rule := range rules {
			// Check Political (-1 means Any)
			pMatch := rule.Pol == -1 || rule.Pol == int(sys.Political)
			// Check Economic
			eMatch := rule.Eco == -1 || rule.Eco == int(sys.Economic)
			// Check Social
			sMatch := rule.Soc == -1 || rule.Soc == int(sys.Social)

			if pMatch && eMatch && sMatch {
				produces = true
				if rule.Bonus {
					isBonus = true
				}
				break // Found a valid production rule, no need to check others
			}
		}

		if produces {
			// Calculate Quantity based on Rarity (Lower rarity = Higher quantity)
			// Base: (11 - Rarity) * 50
			// Example: Rarity 1 (Common) = 500 units
			// Example: Rarity 10 (Artifact) = 50 units
			baseQty := (11 - template.Rarity) * 50
			
			// Variance +/- 20%
			variance := float64(baseQty) * 0.2
			qty := baseQty + int(rand.Float64()*variance*2-variance)

			// Bonus Multiplier (Specialized economies produce double)
			if isBonus {
				qty *= 2
			}

			// Ensure at least 1 exists if it matches rules
			if qty < 1 {
				qty = 1
			}

			// Create the Item Entity
			item := entity.NewItem(
				template.Name,
				template.Category,
				template.BaseValue,
				template.Rarity,
				qty,
				template.IsIllegal,
			)
			market = append(market, *item)
		}
	}

	return market
}