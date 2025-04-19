// Package entities contains game entity implementations
package entities

import (
	"atomblaster/constants"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Player represents the player character in the game
type Player struct {
	Pos          rl.Vector2
	Radius       float32
	Speed        float32
	HasGun       bool
	DashSpeed    float32
	Dashing      bool
	DashTimer    float32
	DashCooldown float32
	RotorAngle   float32 // Added to animate rotor
}

// NewPlayer creates a new player instance with provided position
func NewPlayer(x, y float32) Player {
	return Player{
		Pos:          rl.Vector2{X: x, Y: y},
		Radius:       constants.HelicopterWidth / 2, // Adjusted radius to match helicopter width
		Speed:        300,
		HasGun:       false,
		DashSpeed:    1000,
		Dashing:      false,
		DashTimer:    0,
		DashCooldown: 0,
		RotorAngle:   0,
	}
}

// Update processes player animation and input
func (p *Player) Update(dt float32) {
	// Animate rotor
	p.RotorAngle += dt * 10.0
	if p.RotorAngle > 2*math.Pi {
		p.RotorAngle = 0
	}

	// Update dash timer and cooldown
	if p.Dashing {
		p.DashTimer -= dt
		if p.DashTimer <= 0 {
			p.Dashing = false
			p.DashCooldown = 1.0 // 1 second cooldown
		}
	}

	if p.DashCooldown > 0 {
		p.DashCooldown -= dt
	}

	// Handle keyboard input for movement
	var dx, dy float32

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		dy -= 1
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		dy += 1
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		dx -= 1
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		dx += 1
	}

	// Normalize diagonal movement
	if dx != 0 && dy != 0 {
		dx *= 0.7071
		dy *= 0.7071
	}

	// Move the player
	p.Move(dx, dy)

	// Keep player within screen bounds
	if p.Pos.X < p.Radius {
		p.Pos.X = p.Radius
	}
	if p.Pos.X > constants.ScreenWidth-p.Radius {
		p.Pos.X = constants.ScreenWidth - p.Radius
	}
	if p.Pos.Y < p.Radius {
		p.Pos.Y = p.Radius
	}
	if p.Pos.Y > constants.ScreenHeight-p.Radius {
		p.Pos.Y = constants.ScreenHeight - p.Radius
	}
}

// Move moves the player by the given delta coordinates
func (p *Player) Move(dx, dy float32) {
	// Calculate speed
	moveSpeed := p.Speed
	if p.Dashing {
		moveSpeed = p.DashSpeed
	}

	// Apply movement with normalized diagonal speed
	p.Pos.X += dx * moveSpeed * rl.GetFrameTime()
	p.Pos.Y += dy * moveSpeed * rl.GetFrameTime()
}

// TryDash attempts to perform a dash move if cooldown allows
func (p *Player) TryDash() bool {
	if !p.Dashing && p.DashCooldown <= 0 {
		p.Dashing = true
		p.DashTimer = 0.2 // Dash lasts 0.2 seconds
		return true
	}
	return false
}

// Draw renders the player on screen
func (p *Player) Draw(sprite rl.Texture2D) {
	if sprite.ID > 0 {
		rl.DrawTexturePro(
			sprite,
			rl.Rectangle{X: 0, Y: 0, Width: float32(sprite.Width), Height: float32(sprite.Height)},
			rl.Rectangle{
				X:      p.Pos.X - p.Radius,
				Y:      p.Pos.Y - p.Radius,
				Width:  p.Radius * 2,
				Height: p.Radius * 2,
			},
			rl.Vector2{X: p.Radius, Y: p.Radius},
			0.0,
			rl.White,
		)
	} else {
		// Draw helicopter using primitives

		// Body
		rl.DrawRectangle(
			int32(p.Pos.X - constants.HelicopterWidth/2),
			int32(p.Pos.Y - constants.HelicopterHeight/2),
			constants.HelicopterWidth,
			constants.HelicopterHeight,
			rl.DarkGray,
		)

		// Cockpit
		rl.DrawCircleV(p.Pos, constants.CockpitRadius, rl.SkyBlue)

		// Main rotor - animate rotation
		rotorEndX1 := p.Pos.X + float32(math.Cos(float64(p.RotorAngle)))*constants.RotorLength
		rotorEndY1 := p.Pos.Y - constants.HelicopterHeight/2 - 5 + float32(math.Sin(float64(p.RotorAngle)))*5
		rotorEndX2 := p.Pos.X - float32(math.Cos(float64(p.RotorAngle)))*constants.RotorLength
		rotorEndY2 := p.Pos.Y - constants.HelicopterHeight/2 - 5 - float32(math.Sin(float64(p.RotorAngle)))*5

		rl.DrawLine(
			int32(rotorEndX1),
			int32(rotorEndY1),
			int32(rotorEndX2),
			int32(rotorEndY2),
			rl.LightGray,
		)

		// Rotor center point
		rl.DrawCircleV(
			rl.Vector2{X: p.Pos.X, Y: p.Pos.Y - constants.HelicopterHeight/2 - 5},
			3,
			rl.Gray,
		)

		// Tail
		rl.DrawRectangle(
			int32(p.Pos.X + constants.HelicopterWidth/2 - 5),
			int32(p.Pos.Y - constants.HelicopterHeight/4),
			25,
			constants.HelicopterHeight/2,
			rl.Gray,
		)

		// Tail rotor - also animated
		tailRotorY1 := p.Pos.Y + float32(math.Sin(float64(p.RotorAngle*2)))*constants.TailRotorLength
		tailRotorY2 := p.Pos.Y - float32(math.Sin(float64(p.RotorAngle*2)))*constants.TailRotorLength

		rl.DrawLine(
			int32(p.Pos.X + constants.HelicopterWidth/2 + 20),
			int32(tailRotorY1),
			int32(p.Pos.X + constants.HelicopterWidth/2 + 20),
			int32(tailRotorY2),
			rl.LightGray,
		)

		// Gun if player has it
		if p.HasGun {
			rl.DrawRectangle(
				int32(p.Pos.X + constants.CockpitRadius + 2),
				int32(p.Pos.Y - 3),
				15,
				6,
				rl.DarkBlue,
			)
		}
	}

	// Draw dash cooldown indicator
	if p.DashCooldown > 0 {
		cooldownRatio := p.DashCooldown / 1.0
		rl.DrawRectangle(
			int32(p.Pos.X-p.Radius),
			int32(p.Pos.Y+constants.HelicopterHeight/2+5),
			int32(2*p.Radius*cooldownRatio),
			3,
			rl.Gray,
		)
	}
}
