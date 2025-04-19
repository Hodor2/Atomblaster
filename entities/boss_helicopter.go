// Package entities contains game entity implementations
package entities

import (
	"math"
	
	rl "github.com/gen2brain/raylib-go/raylib"
)

// BossPhase represents different attack patterns of the boss
type BossPhase int

const (
	PhaseCircling BossPhase = iota  // Circle around the arena
	PhaseDashing                     // Quick dashes toward player
	PhaseFiring                      // Stop and fire bullets
)

// BossHelicopter represents the enemy boss helicopter
type BossHelicopter struct {
	Pos          rl.Vector2
	Velocity     rl.Vector2
	Speed        float32
	Health       int
	MaxHealth    int
	Width        float32
	Height       float32
	RotorAngle   float32
	Phase        BossPhase
	PhaseTimer   float32
	AttackTimer  float32
	CircleCenter rl.Vector2
	CircleRadius float32
	CircleAngle  float32
	DashTarget   rl.Vector2
	Invulnerable bool
	InvTimer     float32
	FlashTimer   float32
}

// NewBossHelicopter creates a new boss helicopter
func NewBossHelicopter(screenWidth, screenHeight float32) *BossHelicopter {
	return &BossHelicopter{
		Pos:          rl.Vector2{X: screenWidth / 2, Y: 100},
		Velocity:     rl.Vector2{X: 0, Y: 0},
		Speed:        200,
		Health:       100,
		MaxHealth:    100,
		Width:        80,  // Larger than player helicopter
		Height:       40,
		RotorAngle:   0,
		Phase:        PhaseCircling,
		PhaseTimer:   10.0,        // Time in seconds for each phase
		AttackTimer:  0,
		CircleCenter: rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2},
		CircleRadius: 200,
		CircleAngle:  0,
		Invulnerable: true,        // Start invulnerable during intro
		InvTimer:     2.0,         // Invulnerability time
		FlashTimer:   0,
	}
}

// Update moves and controls the boss helicopter
func (b *BossHelicopter) Update(dt float32, playerPos rl.Vector2, screenWidth, screenHeight float32) bool {
	// Update rotor animation
	b.RotorAngle += dt * 15.0
	if b.RotorAngle > 2*math.Pi {
		b.RotorAngle = 0
	}
	
	// Update invulnerability
	if b.Invulnerable {
		b.InvTimer -= dt
		if b.InvTimer <= 0 {
			b.Invulnerable = false
		}
	}
	
	// Update flash effect when hit
	if b.FlashTimer > 0 {
		b.FlashTimer -= dt
	}
	
	// Update phase timer
	b.PhaseTimer -= dt
	if b.PhaseTimer <= 0 {
		// Switch to next phase
		b.changePhase()
		b.PhaseTimer = 10.0 // Reset timer
	}
	
	// Update boss behavior based on current phase
	switch b.Phase {
	case PhaseCircling:
		b.updateCirclingPhase(dt, playerPos)
	case PhaseDashing:
		b.updateDashingPhase(dt, playerPos)
	case PhaseFiring:
		b.updateFiringPhase(dt, playerPos)
	}
	
	// Constrain to screen bounds
	if b.Pos.X < b.Width/2 {
		b.Pos.X = b.Width / 2
	}
	if b.Pos.X > screenWidth - b.Width/2 {
		b.Pos.X = screenWidth - b.Width/2
	}
	if b.Pos.Y < b.Height/2 {
		b.Pos.Y = b.Height/2
	}
	if b.Pos.Y > screenHeight - b.Height/2 {
		b.Pos.Y = screenHeight - b.Height/2
	}
	
	// Return true if boss is defeated
	return b.Health <= 0
}

// changePhase selects the next attack phase
func (b *BossHelicopter) changePhase() {
	// Cycle through phases
	switch b.Phase {
	case PhaseCircling:
		b.Phase = PhaseDashing
		b.AttackTimer = 1.0 // Initial delay before first dash
	case PhaseDashing:
		b.Phase = PhaseFiring
		b.AttackTimer = 0.5 // Initial delay before firing
	case PhaseFiring:
		b.Phase = PhaseCircling
		b.CircleAngle = float32(math.Atan2(float64(b.Pos.Y-b.CircleCenter.Y), float64(b.Pos.X-b.CircleCenter.X)))
	}
	
	// Brief invulnerability when changing phases
	b.Invulnerable = true
	b.InvTimer = 1.0
}

// updateCirclingPhase handles circling movement pattern
func (b *BossHelicopter) updateCirclingPhase(dt float32, playerPos rl.Vector2) {
	// Update the circle center to follow player somewhat
	b.CircleCenter.X = b.CircleCenter.X*0.99 + playerPos.X*0.01
	b.CircleCenter.Y = b.CircleCenter.Y*0.99 + playerPos.Y*0.01
	
	// Increase circle angle based on time
	b.CircleAngle += dt * 1.0
	
	// Calculate target position on circle
	targetX := b.CircleCenter.X + float32(math.Cos(float64(b.CircleAngle)))*b.CircleRadius
	targetY := b.CircleCenter.Y + float32(math.Sin(float64(b.CircleAngle)))*b.CircleRadius
	
	// Calculate direction to target
	dirX := targetX - b.Pos.X
	dirY := targetY - b.Pos.Y
	
	// Normalize direction
	dist := float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
	if dist > 0 {
		dirX /= dist
		dirY /= dist
	}
	
	// Set velocity
	b.Velocity.X = dirX * b.Speed
	b.Velocity.Y = dirY * b.Speed
	
	// Apply velocity
	b.Pos.X += b.Velocity.X * dt
	b.Pos.Y += b.Velocity.Y * dt
}

// updateDashingPhase handles dashing toward player
func (b *BossHelicopter) updateDashingPhase(dt float32, playerPos rl.Vector2) {
	b.AttackTimer -= dt
	
	if b.AttackTimer <= 0 {
		// Time to start a new dash
		b.DashTarget = playerPos
		
		// Calculate direction to player
		dirX := playerPos.X - b.Pos.X
		dirY := playerPos.Y - b.Pos.Y
		
		// Normalize direction
		dist := float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
		if dist > 0 {
			dirX /= dist
			dirY /= dist
		}
		
		// Set high velocity for dash
		dashSpeed := b.Speed * 3.0
		b.Velocity.X = dirX * dashSpeed
		b.Velocity.Y = dirY * dashSpeed
		
		// Reset attack timer for next dash
		b.AttackTimer = 2.0
	} else if b.AttackTimer > 1.5 {
		// Currently dashing
		b.Pos.X += b.Velocity.X * dt
		b.Pos.Y += b.Velocity.Y * dt
	} else {
		// Slow down after dash
		b.Velocity.X *= 0.95
		b.Velocity.Y *= 0.95
		b.Pos.X += b.Velocity.X * dt
		b.Pos.Y += b.Velocity.Y * dt
	}
}

// updateFiringPhase handles stopping and firing at player
func (b *BossHelicopter) updateFiringPhase(dt float32, playerPos rl.Vector2) {
	b.AttackTimer -= dt
	
	if b.AttackTimer <= 0 {
		// Ready to fire a bullet - this is just a placeholder
		// The actual bullet firing will be handled in the gamestate
		
		// Calculate direction to player for bullet trajectory
		dirX := playerPos.X - b.Pos.X
		dirY := playerPos.Y - b.Pos.Y
		
		// Normalize direction
		dist := float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
		if dist > 0 {
			dirX /= dist
			dirY /= dist
		}
		
		// Set velocity to slightly move while firing
		moveSpeed := b.Speed * 0.5
		b.Velocity.X = dirX * moveSpeed
		b.Velocity.Y = dirY * moveSpeed
		
		// Reset attack timer for next shot
		b.AttackTimer = 0.8
	}
	
	// Move slightly to avoid being a sitting duck
	b.Pos.X += b.Velocity.X * dt
	b.Pos.Y += b.Velocity.Y * dt
	
	// Slow down movement
	b.Velocity.X *= 0.98
	b.Velocity.Y *= 0.98
}

// Draw renders the boss helicopter
func (b *BossHelicopter) Draw() {
	// Skip drawing if dead
	if b.Health <= 0 {
		return
	}
	
	// Determine color based on invulnerability/hit status
	color := rl.Red
	if b.Invulnerable {
		if int(b.InvTimer*10) % 2 == 0 {
			 color = rl.Color{R: 139, G: 0, B: 0, A: 255} // Dark red color
		}
	} else if b.FlashTimer > 0 {
		color = rl.White
	}
	
	// Draw helicopter body
	rl.DrawRectangle(
		int32(b.Pos.X - b.Width/2),
		int32(b.Pos.Y - b.Height/2),
		int32(b.Width),
		int32(b.Height),
		color,
	)
	
	// Draw cockpit
	rl.DrawCircleV(b.Pos, 15, rl.Black)
	
	// Draw main rotor
	rotorEndX1 := b.Pos.X + float32(math.Cos(float64(b.RotorAngle)))*50
	rotorEndY1 := b.Pos.Y - b.Height/2 - 5 + float32(math.Sin(float64(b.RotorAngle)))*5
	rotorEndX2 := b.Pos.X - float32(math.Cos(float64(b.RotorAngle)))*50
	rotorEndY2 := b.Pos.Y - b.Height/2 - 5 - float32(math.Sin(float64(b.RotorAngle)))*5
	
	rl.DrawLine(
		int32(rotorEndX1),
		int32(rotorEndY1),
		int32(rotorEndX2),
		int32(rotorEndY2),
		rl.DarkGray,
	)
	
	// Draw tail
	rl.DrawRectangle(
		int32(b.Pos.X + b.Width/2 - 5),
		int32(b.Pos.Y - b.Height/4),
		int32(30),
		int32(b.Height/2),
		color,
	)
	
	// Draw tail rotor
	tailRotorY1 := b.Pos.Y + float32(math.Sin(float64(b.RotorAngle*2)))*15
	tailRotorY2 := b.Pos.Y - float32(math.Sin(float64(b.RotorAngle*2)))*15
	
	rl.DrawLine(
		int32(b.Pos.X + b.Width/2 + 25),
		int32(tailRotorY1),
		int32(b.Pos.X + b.Width/2 + 25),
		int32(tailRotorY2),
		rl.DarkGray,
	)
	
	// Draw weapon pods under helicopter
	rl.DrawRectangle(
		int32(b.Pos.X - b.Width/3),
		int32(b.Pos.Y + b.Height/2),
		int32(10),
		int32(10),
		rl.DarkGray,
	)
	
	rl.DrawRectangle(
		int32(b.Pos.X + b.Width/3 - 10),
		int32(b.Pos.Y + b.Height/2),
		int32(10),
		int32(10),
		rl.DarkGray,
	)
	
	// Draw health bar
	barWidth := b.Width
	barHeight := 8.0
	healthPercent := float32(b.Health) / float32(b.MaxHealth)
	
	// Background bar
	rl.DrawRectangle(
		int32(b.Pos.X - barWidth/2),
		int32(b.Pos.Y - b.Height/2 - 20),
		int32(barWidth),
		int32(barHeight),
		rl.DarkGray,
	)
	
	// Health bar
	rl.DrawRectangle(
		int32(b.Pos.X - barWidth/2),
		int32(b.Pos.Y - b.Height/2 - 20),
		int32(barWidth * healthPercent),
		int32(barHeight),
		rl.Red,
	)
}

// TakeDamage reduces boss health and returns true if boss was defeated
func (b *BossHelicopter) TakeDamage(amount int) bool {
	// If invulnerable, no damage
	if b.Invulnerable {
		return false
	}
	
	b.Health -= amount
	if b.Health < 0 {
		b.Health = 0
	}
	
	// Set flash effect
	b.FlashTimer = 0.1
	
	return b.Health <= 0
}

// GetCollisionRect returns the collision rectangle for the boss
func (b *BossHelicopter) GetCollisionRect() rl.Rectangle {
	return rl.Rectangle{
		X:      b.Pos.X - b.Width/2,
		Y:      b.Pos.Y - b.Height/2,
		Width:  b.Width,
		Height: b.Height,
	}
}

// FireBullet returns a vector for the bullet's direction (used by gamestate)
func (b *BossHelicopter) FireBullet() rl.Vector2 {
	// This just gives a direction - actual bullet creation happens in gamestate
	// Returns bullet direction: downward and slightly random
	return rl.Vector2{
		X: float32(rl.GetRandomValue(-50, 50)) / 100.0,
		Y: 1.0,
	}
}

// ShouldFireBullet checks if the boss should fire now based on phase
func (b *BossHelicopter) ShouldFireBullet() bool {
	return b.Phase == PhaseFiring && b.AttackTimer < 0.1 && b.AttackTimer > 0
}
