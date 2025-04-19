package entities

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"math"
)

// PowerUp represents a collectible power-up item
type PowerUp struct {
	Position    rl.Vector2
	Velocity    rl.Vector2
	Type        PowerUpType
	Color       rl.Color
	Size        float32
	Rotation    float32
	Duration    float32 // How long the power-up lasts when collected (in seconds)
	PulseTime   float32
	Lifetime    float32 // How long the power-up exists in the world before disappearing
	MaxLifetime float32 // Maximum lifetime
	Active      bool
}

// NewPowerUp creates a new power-up of the specified type
func NewPowerUp(position rl.Vector2, powerType PowerUpType) *PowerUp {
	var color rl.Color
	var duration float32 = constants.PowerUpDuration
	var size float32 = 12.0 // Base size for all power-ups
	
	// Configure based on type
	switch powerType {
	case PowerUpMagnet:
		color = rl.Purple
		
	case PowerUpSpeed:
		color = rl.Red
		duration = constants.PowerUpDuration * 0.8 // Slightly shorter for balance
		
	case PowerUpShield:
		color = rl.SkyBlue
		duration = constants.PowerUpDuration * 0.5 // Shorter shield duration
		
	case PowerUpSizeBoost:
		color = rl.Orange
		duration = constants.PowerUpDuration * 1.5 // Longer size boost duration
	}
	
	// Random initial velocity
	angle := rl.GetRandomFloat32() * 2 * math.Pi
	speed := 10.0 + rl.GetRandomFloat32() * 20.0
	
	return &PowerUp{
		Position:    position,
		Velocity:    rl.Vector2{
						 X: float32(math.Cos(float64(angle))) * speed,
						 Y: float32(math.Sin(float64(angle))) * speed,
					 },
		Type:        powerType,
		Color:       color,
		Size:        size,
		Rotation:    rl.GetRandomFloat32() * 360,
		Duration:    duration,
		PulseTime:   rl.GetRandomFloat32() * 2 * math.Pi,
		Lifetime:    30.0 + rl.GetRandomFloat32() * 30.0, // 30-60 seconds before disappearing
		MaxLifetime: 30.0 + rl.GetRandomFloat32() * 30.0,
		Active:      true,
	}
}

// NewRandomPowerUp creates a power-up of random type
func NewRandomPowerUp(position rl.Vector2) *PowerUp {
	// Randomly select power-up type
	randomType := PowerUpType(int(rl.GetRandomFloat32() * 4))
	return NewPowerUp(position, randomType)
}

// Update updates the power-up's state
func (p *PowerUp) Update(dt float32, worldBounds rl.Rectangle) {
	// Skip if not active
	if !p.Active {
		return
	}
	
	// Update position based on velocity
	p.Position.X += p.Velocity.X * dt
	p.Position.Y += p.Velocity.Y * dt
	
	// Bounce off world boundaries
	if p.Position.X < worldBounds.X + p.Size {
		p.Position.X = worldBounds.X + p.Size
		p.Velocity.X = -p.Velocity.X * 0.8 // Slight damping on bounce
	} else if p.Position.X > worldBounds.X + worldBounds.Width - p.Size {
		p.Position.X = worldBounds.X + worldBounds.Width - p.Size
		p.Velocity.X = -p.Velocity.X * 0.8
	}
	
	if p.Position.Y < worldBounds.Y + p.Size {
		p.Position.Y = worldBounds.Y + p.Size
		p.Velocity.Y = -p.Velocity.Y * 0.8
	} else if p.Position.Y > worldBounds.Y + worldBounds.Height - p.Size {
		p.Position.Y = worldBounds.Y + worldBounds.Height - p.Size
		p.Velocity.Y = -p.Velocity.Y * 0.8
	}
	
	// Apply drag
	p.Velocity.X *= 0.98
	p.Velocity.Y *= 0.98
	
	// Rotate
	p.Rotation += dt * 45 // 45 degrees per second
	if p.Rotation > 360 {
		p.Rotation -= 360
	}
	
	// Update pulse animation
	p.PulseTime += dt * 3
	if p.PulseTime > 2*math.Pi {
		p.PulseTime -= 2 * math.Pi
	}
	
	// Update lifetime
	p.Lifetime -= dt
	if p.Lifetime <= 0 {
		p.Active = false
	}
}

// Draw renders the power-up
func (p *PowerUp) Draw() {
	if !p.Active {
		return
	}
	
	// Calculate pulse effect (pulsate more rapidly as expiration approaches)
	lifetimeRatio := p.Lifetime / p.MaxLifetime
	pulseSpeed := 1.0 + (1.0 - lifetimeRatio) * 2.0 // Pulse faster when near expiration
	pulseAmount := 0.2 + (1.0 - lifetimeRatio) * 0.3 // Pulse more dramatically when near expiration
	
	pulseFactor := 1.0 - pulseAmount + pulseAmount * float32(math.Sin(float64(p.PulseTime) * pulseSpeed))
	
	// Draw with pulsing size
	drawSize := p.Size * pulseFactor
	
	// Make color pulse alpha when nearing expiration
	drawColor := p.Color
	if p.Lifetime < 5.0 { // Last 5 seconds
		alphaFactor := math.Sin(float64(p.PulseTime) * 2.0)
		minAlpha := uint8(100 + 155 * lifetimeRatio) // Fade from 255 to 100 as time runs out
		
		drawColor.A = uint8(float64(minAlpha) + alphaFactor*float64(255-minAlpha))
	}
	
	// Draw based on power-up type
	switch p.Type {
	case PowerUpMagnet:
		drawMagnetPowerUp(p.Position, drawSize, drawColor, p.Rotation)
	case PowerUpSpeed:
		drawSpeedPowerUp(p.Position, drawSize, drawColor, p.Rotation)
	case PowerUpShield:
		drawShieldPowerUp(p.Position, drawSize, drawColor, p.Rotation)
	case PowerUpSizeBoost:
		drawSizeBoostPowerUp(p.Position, drawSize, drawColor, p.Rotation)
	}
}

// GetBounds returns the collision rectangle for the power-up
func (p *PowerUp) GetBounds() rl.Rectangle {
	return rl.Rectangle{
		X:      p.Position.X - p.Size,
		Y:      p.Position.Y - p.Size,
		Width:  p.Size * 2,
		Height: p.Size * 2,
	}
}

// Drawing functions for different power-up types

func drawMagnetPowerUp(pos rl.Vector2, size float32, color rl.Color, rotation float32) {
	// Draw outer circle
	rl.DrawCircleLinesEx(pos, size, 2.0, color)
	
	// Draw magnet symbol
	radAngle := rotation * math.Pi / 180.0
	
	// North pole
	northX := pos.X + float32(math.Cos(float64(radAngle))) * size * 0.6
	northY := pos.Y + float32(math.Sin(float64(radAngle))) * size * 0.6
	
	// South pole
	southX := pos.X - float32(math.Cos(float64(radAngle))) * size * 0.6
	southY := pos.Y - float32(math.Sin(float64(radAngle))) * size * 0.6
	
	// Draw poles
	rl.DrawRectanglePro(
		rl.Rectangle{X: northX - size*0.3, Y: northY - size*0.3, Width: size*0.6, Height: size*0.6},
		rl.Vector2{X: size*0.3, Y: size*0.3},
		rotation,
		color,
	)
	
	rl.DrawRectanglePro(
		rl.Rectangle{X: southX - size*0.3, Y: southY - size*0.3, Width: size*0.6, Height: size*0.6},
		rl.Vector2{X: size*0.3, Y: size*0.3},
		rotation,
		color,
	)
}

func drawSpeedPowerUp(pos rl.Vector2, size float32, color rl.Color, rotation float32) {
	// Draw outer circle
	rl.DrawCircleLinesEx(pos, size, 2.0, color)
	
	// Draw lightning bolt symbol
	vertices := []rl.Vector2{
		{X: pos.X - size*0.4, Y: pos.Y - size*0.6}, // Top left
		{X: pos.X + size*0.1, Y: pos.Y - size*0.1}, // Middle right
		{X: pos.X - size*0.1, Y: pos.Y - size*0.1}, // Middle left
		{X: pos.X + size*0.4, Y: pos.Y + size*0.6}, // Bottom right
		{X: pos.X,            Y: pos.Y + size*0.1}, // Bottom middle
		{X: pos.X - size*0.2, Y: pos.Y + size*0.1}, // Bottom left
	}
	
	// Rotate vertices around center
	radAngle := rotation * math.Pi / 180.0
	sinRot := float32(math.Sin(float64(radAngle)))
	cosRot := float32(math.Cos(float64(radAngle)))
	
	for i := range vertices {
		// Translate to origin
		x := vertices[i].X - pos.X
		y := vertices[i].Y - pos.Y
		
		// Rotate
		rotX := x*cosRot - y*sinRot
		rotY := x*sinRot + y*cosRot
		
		// Translate back
		vertices[i].X = rotX + pos.X
		vertices[i].Y = rotY + pos.Y
	}
	
	// Draw lightning bolt
	for i := 0; i < len(vertices); i++ {
		j := (i + 1) % len(vertices)
		rl.DrawLineEx(vertices[i], vertices[j], 2.0, color)
	}
}

func drawShieldPowerUp(pos rl.Vector2, size float32, color rl.Color, rotation float32) {
	// Draw outer circle
	rl.DrawCircleLinesEx(pos, size, 2.0, color)
	
	// Draw shield symbol (a rounded rectangle)
	shieldWidth := size * 0.8
	shieldHeight := size * 1.1
	
	// Create shield shape
	shieldRect := rl.Rectangle{
		X: pos.X - shieldWidth/2,
		Y: pos.Y - shieldHeight/2,
		Width: shieldWidth,
		Height: shieldHeight,
	}
	
	// Draw rotated shield
	rl.DrawRectanglePro(
		shieldRect,
		rl.Vector2{X: shieldWidth/2, Y: shieldHeight/2},
		rotation,
		rl.Color{R: color.R, G: color.G, B: color.B, A: 100},
	)
	
	// Draw shield outline
	rl.DrawRectangleLinesEx(
		rl.Rectangle{
			X: pos.X - shieldWidth/2,
			Y: pos.Y - shieldHeight/2,
			Width: shieldWidth,
			Height: shieldHeight,
		},
		2.0,
		color,
	)
}

func drawSizeBoostPowerUp(pos rl.Vector2, size float32, color rl.Color, rotation float32) {
	// Draw outer circle
	rl.DrawCircleLinesEx(pos, size, 2.0, color)
	
	// Draw expand symbol (arrows pointing outward)
	arrowSize := size * 0.6
	
	// Draw four arrows pointing outward
	for i := 0; i < 4; i++ {
		angle := float32(i) * 90.0 + rotation
		radAngle := angle * math.Pi / 180.0
		
		// Arrow start (inner point)
		startX := pos.X + float32(math.Cos(float64(radAngle))) * size * 0.2
		startY := pos.Y + float32(math.Sin(float64(radAngle))) * size * 0.2
		
		// Arrow end (outer point)
		endX := pos.X + float32(math.Cos(float64(radAngle))) * size * 0.8
		endY := pos.Y + float32(math.Sin(float64(radAngle))) * size * 0.8
		
		// Draw arrow line
		rl.DrawLineEx(
			rl.Vector2{X: startX, Y: startY},
			rl.Vector2{X: endX, Y: endY},
			2.0,
			color,
		)
		
		// Draw arrow head
		headSize := size * 0.2
		headAngle1 := radAngle + math.Pi*0.85 // Head angle offset
		headAngle2 := radAngle - math.Pi*0.85
		
		head1X := endX + float32(math.Cos(float64(headAngle1))) * headSize
		head1Y := endY + float32(math.Sin(float64(headAngle1))) * headSize
		
		head2X := endX + float32(math.Cos(float64(headAngle2))) * headSize
		head2Y := endY + float32(math.Sin(float64(headAngle2))) * headSize
		
		rl.DrawLineEx(
			rl.Vector2{X: endX, Y: endY},
			rl.Vector2{X: head1X, Y: head1Y},
			2.0,
			color,
		)
		
		rl.DrawLineEx(
			rl.Vector2{X: endX, Y: endY},
			rl.Vector2{X: head2X, Y: head2Y},
			2.0,
			color,
		)
	}
}
