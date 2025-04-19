// Package entities contains game entity implementations
package entities

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// RescueZone represents a safe area where scientists need to be brought
type RescueZone struct {
	Pos         rl.Vector2
	Size        rl.Vector2
	Color       rl.Color
	PulseAmount float32
	Timer       float32
}

// NewRescueZone creates a new rescue zone
func NewRescueZone(x, y, width, height float32) RescueZone {
	return RescueZone{
		Pos:         rl.Vector2{X: x, Y: y},
		Size:        rl.Vector2{X: width, Y: height},
		Color:       rl.Green,
		PulseAmount: 0,
		Timer:       0,
	}
}

// Update animates the rescue zone
func (rz *RescueZone) Update(dt float32) {
	// Update pulse effect
	rz.Timer += dt
	rz.PulseAmount = float32(math.Sin(float64(rz.Timer * 2)))
}

// Draw renders the rescue zone
func (rz *RescueZone) Draw() {
	// Calculate pulse size
	pulseSize := 5.0 * rz.PulseAmount
	
	// Update pulsing effect
	rz.Timer += rl.GetFrameTime()
	rz.PulseAmount = float32(math.Sin(float64(rz.Timer * 2)))
	
	// Draw zone with pulsing effect
	rl.DrawRectangleV(
		rl.Vector2{
			X: rz.Pos.X - pulseSize,
			Y: rz.Pos.Y - pulseSize,
		},
		rl.Vector2{
			X: rz.Size.X + pulseSize*2,
			Y: rz.Size.Y + pulseSize*2,
		},
		rl.Fade(rz.Color, 0.3),
	)
	
	// Draw border
	rl.DrawRectangleLinesEx(
		rl.Rectangle{
			X:      rz.Pos.X - pulseSize,
			Y:      rz.Pos.Y - pulseSize,
			Width:  rz.Size.X + pulseSize*2,
			Height: rz.Size.Y + pulseSize*2,
		},
		2, rl.Fade(rz.Color, 0.8),
	)
	
	// Draw "RESCUE ZONE" text
	textWidth := rl.MeasureText("RESCUE ZONE", 20)
	rl.DrawText(
		"RESCUE ZONE",
		int32(rz.Pos.X + rz.Size.X/2 - float32(textWidth)/2),
		int32(rz.Pos.Y + rz.Size.Y/2 - 10),
		20,
		rl.White,
	)
}

// GetRectangle returns the rectangle representing the rescue zone
func (rz *RescueZone) GetRectangle() rl.Rectangle {
	return rl.Rectangle{
		X:      rz.Pos.X,
		Y:      rz.Pos.Y,
		Width:  rz.Size.X,
		Height: rz.Size.Y,
	}
}

// CheckHelicopterInZone determines if the helicopter is inside the rescue zone
func (rz *RescueZone) CheckHelicopterInZone(helicopterRect rl.Rectangle) bool {
	return rl.CheckCollisionRecs(helicopterRect, rz.GetRectangle())
}

// IsHelicopterNear checks if the helicopter is within a certain distance of the rescue zone
func (rz *RescueZone) IsHelicopterNear(helicopterRect rl.Rectangle, distance float32) bool {
	// Get center points
	helicopterCenter := rl.Vector2{
		X: helicopterRect.X + helicopterRect.Width/2,
		Y: helicopterRect.Y + helicopterRect.Height/2,
	}
	
	zoneCenter := rl.Vector2{
		X: rz.Pos.X + rz.Size.X/2,
		Y: rz.Pos.Y + rz.Size.Y/2,
	}
	
	// Check distance between centers
	return rl.Vector2Distance(helicopterCenter, zoneCenter) < distance
}
