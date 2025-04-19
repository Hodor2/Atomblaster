// Package entities contains game entity implementations
package entities

import (
	"math"
	
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Atom represents an enemy atom in the game
type Atom struct {
	Pos       rl.Vector2
	Velocity  rl.Vector2
	Radius    float32
	Color     rl.Color
	HP        int
	Type      int  // 0 = normal, 1 = fast, 2 = big
	Rotation  float32
	SpinSpeed float32
}

// NewAtom creates a new atom enemy with the given position, velocity and type
func NewAtom(pos, vel rl.Vector2, atomType int) Atom {
	// Determine attributes based on type
	var radius float32
	var health int
	var spinSpeed float32
	var color rl.Color
	
	switch atomType {
	case 0: // Normal atom
		radius = 15.0
		health = 2
		spinSpeed = 2.0
		color = rl.Red
	case 1: // Fast atom
		radius = 12.0
		health = 1
		spinSpeed = 4.0
		// Multiply velocity by 1.5 for fast atoms
		vel.X *= 1.5
		vel.Y *= 1.5
		color = rl.Orange
	case 2: // Big atom
		radius = 25.0
		health = 4
		spinSpeed = 1.0
		// Slow down big atoms
		vel.X *= 0.7
		vel.Y *= 0.7
		color = rl.Maroon
	default:
		// Default to normal atom
		radius = 15.0
		health = 2
		spinSpeed = 2.0
		color = rl.Red
	}
	
	return Atom{
		Pos:       pos,
		Velocity:  vel,
		Radius:    radius,
		Color:     color,
		HP:        health,
		Type:      atomType,
		Rotation:  0,
		SpinSpeed: spinSpeed,
	}
}

// Update moves the atom and handles screen edge bouncing
func (a *Atom) Update(dt float32, playerPos rl.Vector2) {
	// Update rotation
	a.Rotation += a.SpinSpeed * dt
	if a.Rotation > 2*math.Pi {
		a.Rotation -= 2 * math.Pi
	}
	
	// Move atom based on velocity
	a.Pos.X += a.Velocity.X * dt
	a.Pos.Y += a.Velocity.Y * dt
	
	// Check for collision with screen edges and bounce
	if a.Pos.X < a.Radius {
		a.Pos.X = a.Radius
		a.Velocity.X = -a.Velocity.X
	}
	if a.Pos.X > 800 - a.Radius {
		a.Pos.X = 800 - a.Radius
		a.Velocity.X = -a.Velocity.X
	}
	if a.Pos.Y < a.Radius {
		a.Pos.Y = a.Radius
		a.Velocity.Y = -a.Velocity.Y
	}
	if a.Pos.Y > 600 - a.Radius {
		a.Pos.Y = 600 - a.Radius
		a.Velocity.Y = -a.Velocity.Y
	}
	
	// For type 1 (fast) atoms, they have a slight homing behavior
	if a.Type == 1 {
		// Calculate direction to player
		dx := playerPos.X - a.Pos.X
		dy := playerPos.Y - a.Pos.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		
		// If player is close, slightly adjust velocity to home in on player
		if distance < 200 {
			// Normalize direction
			if distance > 0 {
				dx /= distance
				dy /= distance
			}
			
			// Gently adjust velocity toward player (weak homing effect)
			homingStrength := float32(50.0)
			a.Velocity.X += dx * homingStrength * dt
			a.Velocity.Y += dy * homingStrength * dt
			
			// Limit velocity to maintain speed cap
			maxSpeed := float32(225.0) // 1.5x base speed
			currentSpeed := float32(math.Sqrt(float64(a.Velocity.X*a.Velocity.X + a.Velocity.Y*a.Velocity.Y)))
			if currentSpeed > maxSpeed {
				a.Velocity.X = (a.Velocity.X / currentSpeed) * maxSpeed
				a.Velocity.Y = (a.Velocity.Y / currentSpeed) * maxSpeed
			}
		}
	}
}

// Draw renders the atom
func (a *Atom) Draw(sprite rl.Texture2D) {
	// Calculate pulse effect for visual interest
	pulseAmount := 0.1 * float32(math.Sin(float64(rl.GetTime()*3.0)))
	drawRadius := a.Radius * (1.0 + pulseAmount)
	
	if sprite.ID > 0 {
		// Draw using sprite
		rl.DrawTexturePro(
			sprite,
			rl.Rectangle{X: 0, Y: 0, Width: float32(sprite.Width), Height: float32(sprite.Height)},
			rl.Rectangle{X: a.Pos.X - drawRadius, Y: a.Pos.Y - drawRadius, Width: drawRadius * 2, Height: drawRadius * 2},
			rl.Vector2{X: drawRadius, Y: drawRadius},
			a.Rotation * 180.0 / math.Pi, // Convert radians to degrees
			a.Color,
		)
	} else {
		// Draw simple atom using circles
		rl.DrawCircleV(a.Pos, drawRadius, a.Color)
		
		// Draw electron rings
		numRings := 1
		if a.Type == 2 {  // More rings for big atoms
			numRings = 2
		}
		
		for i := 0; i < numRings; i++ {
			ringRadius := a.Radius * (0.4 + 0.3*float32(i))
			rl.DrawCircleLines(int32(a.Pos.X), int32(a.Pos.Y), ringRadius, rl.White)
			
			// Draw "electrons" on rings
			electronPos1 := rl.Vector2{
				X: a.Pos.X + ringRadius * float32(math.Cos(float64(a.Rotation + float32(i)*math.Pi/2))),
				Y: a.Pos.Y + ringRadius * float32(math.Sin(float64(a.Rotation + float32(i)*math.Pi/2))),
			}
			
			electronPos2 := rl.Vector2{
				X: a.Pos.X + ringRadius * float32(math.Cos(float64(a.Rotation + float32(i)*math.Pi/2 + math.Pi))),
				Y: a.Pos.Y + ringRadius * float32(math.Sin(float64(a.Rotation + float32(i)*math.Pi/2 + math.Pi))),
			}
			
			rl.DrawCircleV(electronPos1, 3, rl.White)
			rl.DrawCircleV(electronPos2, 3, rl.White)
		}
	}
	
	// For big atoms, show health bar
	if a.Type == 2 {
		maxHealth := 4 // Big atoms have 4 HP
		healthPercent := float32(a.HP) / float32(maxHealth)
		
		barWidth := a.Radius * 2
		barHeight := 4.0
		
		// Background bar
		rl.DrawRectangle(
			int32(a.Pos.X - barWidth/2),
			int32(a.Pos.Y - a.Radius - 10),
			int32(barWidth),
			int32(barHeight),
			rl.DarkGray,
		)
		
		// Health bar
		rl.DrawRectangle(
			int32(a.Pos.X - barWidth/2),
			int32(a.Pos.Y - a.Radius - 10),
			int32(barWidth * healthPercent),
			int32(barHeight),
			rl.Red,
		)
	}
}
