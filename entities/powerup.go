// Package entities contains game entity implementations
package entities

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// PowerUp type constants
const (
	PowerUpGun    = 0
	PowerUpHealth = 1
	PowerUpSpeed  = 2
)

// PowerUp represents a collectible item
type PowerUp struct {
	Pos       rl.Vector2
	Radius    float32
	Collected bool
	Type      int // 0 for gun, 1 for health, 2 for speed boost
}

// NewPowerUp creates a new power-up of the specified type
func NewPowerUp(pos rl.Vector2, powerUpType int) PowerUp {
	return PowerUp{
		Pos:       pos,
		Radius:    15,
		Collected: false,
		Type:      powerUpType,
	}
}

// Update handles the power-up's animation and behavior
func (p *PowerUp) Update() {
	// Power-ups don't need much updating, but we could add floating animations here
}

// Draw renders the power-up on screen
func (p *PowerUp) Draw(sprites [3]rl.Texture2D) {
	if p.Collected {
		return
	}
	
	// Make power-ups float a bit
	floatOffset := float32(math.Sin(float64(rl.GetTime()*2.0))) * 5.0
	
	// Choose color based on power-up type
	var color rl.Color
	var symbol string
	
	switch p.Type {
	case PowerUpGun: // Gun
		color = rl.Orange
		symbol = "G"
	case PowerUpHealth: // Health
		color = rl.Green
		symbol = "H"
	case PowerUpSpeed: // Speed
		color = rl.Purple
		symbol = "S"
	}
	
	// Make power-ups glow/pulse
	glowSize := p.Radius * (1.0 + 0.2*float32(math.Sin(float64(rl.GetTime()*3.0))))
	
	// Draw the power-up
	if p.Type < len(sprites) && sprites[p.Type].ID > 0 {
		rl.DrawTexturePro(
			sprites[p.Type],
			rl.Rectangle{X: 0, Y: 0, Width: float32(sprites[p.Type].Width), Height: float32(sprites[p.Type].Height)},
			rl.Rectangle{
				X:      p.Pos.X - p.Radius,
				Y:      p.Pos.Y - p.Radius + floatOffset,
				Width:  p.Radius * 2,
				Height: p.Radius * 2,
			},
			rl.Vector2{X: p.Radius, Y: p.Radius},
			0.0,
			rl.White,
		)
	} else {
		// Fallback to simple shapes with symbols
		rl.DrawCircleV(p.Pos, p.Radius, color)
		rl.DrawText(symbol, int32(p.Pos.X)-5, int32(p.Pos.Y-5+floatOffset), 20, rl.Black)
	}
	
	// Draw glow effect
	rl.DrawCircleLines(int32(p.Pos.X), int32(p.Pos.Y+floatOffset), glowSize, rl.Fade(color, 0.5))
}

// CheckCollision checks if a player has collided with this power-up
func (p *PowerUp) CheckCollision(playerPos rl.Vector2, playerRadius float32) bool {
	if p.Collected {
		return false
	}
	
	dx := p.Pos.X - playerPos.X
	dy := p.Pos.Y - playerPos.Y
	distance := rl.Vector2Length(rl.Vector2{X: dx, Y: dy})
	
	return distance < (p.Radius + playerRadius)
}
