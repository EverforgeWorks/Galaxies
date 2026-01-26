package core

import (
	"math"
)

// ShipStats placeholders (we'll expand this later when we do Ships)
type ShipStats struct {
	FuelEfficiency float64 // Fuel units burned per LY (Lower is better)
	Speed          float64 // LY per hour
	MaxRange       float64 // Max fuel capacity
}

// Default generic ship for testing
var DefaultShip = ShipStats{
	FuelEfficiency: 10.0, // Burns 10 fuel per LY
	Speed:          1.0,  // 1 LY per hour (Slow)
	MaxRange:       200.0, // 20 LY range
}

// NavigationResult holds the pre-flight calculation
type NavigationResult struct {
	Distance     float64
	FuelRequired float64
	TimeHours    float64
	IsReachable  bool
}

// CalculateJump computes the physics of moving from A to B
func CalculateJump(fromX, fromY, toX, toY int, ship ShipStats) NavigationResult {
	// 1. Distance Formula (Pythagorean theorem)
	deltaX := float64(toX - fromX)
	deltaY := float64(toY - fromY)
	dist := math.Sqrt((deltaX * deltaX) + (deltaY * deltaY))

	// 2. Fuel Math
	fuel := dist * ship.FuelEfficiency

	// 3. Time Math
	time := dist / ship.Speed

	return NavigationResult{
		Distance:     dist,
		FuelRequired: fuel,
		TimeHours:    time,
		IsReachable:  dist <= ship.MaxRange, // Simple check against max fuel tank
	}
}