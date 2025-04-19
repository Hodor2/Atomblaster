package entities

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/util"
	"math"
	"time"
)

// PowerUpType defines different types of power-ups
type PowerUpType int

const (
	PowerUpMagnet PowerUpType = iota
	PowerUpSpeed
	PowerUpShield
	PowerUpSizeBoost
)

// PowerUpEffect tracks active power-up duration
type PowerUpEffect struct {
	Type         PowerUpType
	RemainingTime float32
	TotalDuration float32
}

// Player represents the player-controlled helicopter
type Player struct {
	Position       rl.Vector2
	Velocity       rl.Vector2
	Rotation       float32  // Current rotation in degrees
	TargetRotation float32  // Desired rotation in degrees
	Scale          float32  // Visual scale
	Size           float32  // Collision size
	Color          rl.Color
	
	// Physics parameters
	Acceleration   float32
	MaxSpeed       float32
	RotationSpeed  float32
	Friction       float32
	
	// Visual components
	MainRotor      struct {
		Rotation     float32
		RotationSpeed float32
		Size         float32
		Offset       rl.Vector2
	}
	TailRotor      struct {
		Rotation     float32
		RotationSpeed float32
		Size         float32
		Offset       rl.Vector2
	}
	
	// Particle effects
	ExhaustParticles *ParticleEmitter
	
	// Game state
	Health         float32
	Score          int
	BaseSize       float32  // Starting/minimum size
	MaxSize        float32  // Maximum possible size
	GrowthMultiplier float32 // How quickly player grows when eating food
	TotalFoodCollected int
	GameTime       float32
	IsInvulnerable bool
	IsDead         bool
	LastDeathTime  time.Time
	
	// PowerUp effects
	ActivePowerUps    map[PowerUpType]*PowerUpEffect
	MagnetRange       float32  // Range to attract food
	SpeedMultiplier   float32  // Multiplier for movement speed
	HasShield         bool     // Whether shield is active
	SizeBoostFactor   float32  // Multiplier for size growth from food
	
	// Debug 
	DebugMode      bool
}

// NewPlayer creates a new player helicopter
func NewPlayer() *Player {
	player := &Player{
		Position:       rl.Vector2{X: constants.WorldWidth / 2, Y: constants.WorldHeight / 2},
		Velocity:       rl.Vector2{X: 0, Y: 0},
		Rotation:       0,
		TargetRotation: 0,
		Scale:          1.0,
		Size:           constants.InitialPlayerSize,
		Color:          rl.Red,
		
		// Physics
		Acceleration:   500.0,
		MaxSpeed:       constants.MaxHelicopterSpeed,
		RotationSpeed:  5.0,
		Friction:       constants.DefaultFriction,
		
		// Game state
		Health:         100.0,
		Score:          0,
		BaseSize:       constants.InitialPlayerSize,
		MaxSize:        constants.MaxPlayerSize,
		GrowthMultiplier: 0.2,
		TotalFoodCollected: 0,
		GameTime:       0,
		IsInvulnerable: false,
		IsDead:         false,
		LastDeathTime:  time.Now(),
		
		// PowerUp default values
		ActivePowerUps:   make(map[PowerUpType]*PowerUpEffect),
		MagnetRange:      0, // No magnetic attraction by default
		SpeedMultiplier:  1.0, // Normal speed
		HasShield:        false,
		SizeBoostFactor:  1.0, // Normal size growth
		
		DebugMode:      false,
	}
	
	// Set up rotors
	player.MainRotor.Rotation = 0
	player.MainRotor.RotationSpeed = 15.0
	player.MainRotor.Size = 30.0
	player.MainRotor.Offset = rl.Vector2{X: 0, Y: 0}
	
	player.TailRotor.Rotation = 0
	player.TailRotor.RotationSpeed = 25.0
	player.TailRotor.Size = 15.0
	player.TailRotor.Offset = rl.Vector2{X: -18, Y: 0}
	
	// Setup particle emitter for exhaust
	player.ExhaustParticles = NewParticleEmitter(
		rl.Vector2{X: -20, Y: 0},  // Offset from helicopter center
		rl.Gray,                   // Base color
		0.5,                       // Particle size
		20,                        // Particles per second
		2.0,                       // Particle lifetime
	)
	
	return player
}

// Update updates the player state
func (p *Player) Update(dt float32) {
	// Update game time
	p.GameTime += dt
	
	// Update power-up timers
	p.UpdatePowerUps(dt)
	
	// Handle input if not dead
	if !p.IsDead {
		p.HandleInput(dt)
	}
	
	// Apply physics
	p.UpdatePhysics(dt)
	
	// Update rotors
	speed := util.Vector2Length(p.Velocity)
	p.MainRotor.Rotation += p.MainRotor.RotationSpeed * speed/p.MaxSpeed * dt * 360
	p.TailRotor.Rotation += p.TailRotor.RotationSpeed * speed/p.MaxSpeed * dt * 360
	
	// Update exhaust particles
	exhaustPos := p.Position
	exhaustOffset := rl.Vector2{
		X: float32(math.Cos(float64(p.Rotation) * math.Pi / 180)) * -20,
		Y: float32(math.Sin(float64(p.Rotation) * math.Pi / 180)) * -20,
	}
	exhaustPos.X += exhaustOffset.X
	exhaustPos.Y += exhaustOffset.Y
	
	p.ExhaustParticles.Position = exhaustPos
	p.ExhaustParticles.EmissionRate = 5 + speed/p.MaxSpeed*15 // More particles at higher speed
	p.ExhaustParticles.Update(dt)
	
	// Check for respawn after death
	if p.IsDead && time.Since(p.LastDeathTime).Seconds() > 2.0 {
		p.Respawn()
	}
}

// HandleInput processes player input
func (p *Player) HandleInput(dt float32) {
	// Movement direction from input
	moveDir := rl.Vector2{X: 0, Y: 0}
	
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		moveDir.Y = -1
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		moveDir.Y = 1
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		moveDir.X = -1
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		moveDir.X = 1
	}
	
	// Normalize direction if moving diagonally
	length := float32(math.Sqrt(float64(moveDir.X*moveDir.X + moveDir.Y*moveDir.Y)))
	if length > 0 {
		moveDir.X /= length
		moveDir.Y /= length
		
		// Update target rotation based on movement direction
		p.TargetRotation = float32(math.Atan2(float64(moveDir.Y), float64(moveDir.X))) * 180 / math.Pi
	}
	
	// Smoothly rotate toward target direction
	angleDiff := p.TargetRotation - p.Rotation
	
	// Normalize angle to [-180, 180]
	for angleDiff > 180 {
		angleDiff -= 360
	}
	for angleDiff < -180 {
		angleDiff += 360
	}
	
	// Apply rotation with smoothing
	p.Rotation += angleDiff * p.RotationSpeed * dt
	
	// Apply acceleration in moving direction
	if length > 0 {
		p.Velocity.X += moveDir.X * p.Acceleration * dt
		p.Velocity.Y += moveDir.Y * p.Acceleration * dt
	}
}

// UpdatePhysics applies physics to the player
func (p *Player) UpdatePhysics(dt float32) {
	// Apply friction
	p.Velocity.X *= p.Friction
	p.Velocity.Y *= p.Friction
	
	// Apply speed multiplier from power-up
	finalVelocity := p.Velocity
	if p.SpeedMultiplier != 1.0 {
		finalVelocity = util.Vector2Scale(p.Velocity, p.SpeedMultiplier)
	}
	
	// Limit top speed
	speed := util.Vector2Length(finalVelocity)
	maxAdjustedSpeed := p.MaxSpeed * p.SpeedMultiplier
	if speed > maxAdjustedSpeed {
		finalVelocity = util.Vector2Scale(
			util.Vector2Normalize(finalVelocity),
			maxAdjustedSpeed,
		)
	}
	
	// Update position
	p.Position.X += finalVelocity.X * dt
	p.Position.Y += finalVelocity.Y * dt
	
	// Constrain to world bounds
	p.Position.X = util.Clamp(p.Position.X, p.Size, constants.WorldWidth-p.Size)
	p.Position.Y = util.Clamp(p.Position.Y, p.Size, constants.WorldHeight-p.Size)
}

// Draw renders the player helicopter
func (p *Player) Draw() {
	// Skip drawing if dead
	if p.IsDead {
		return
	}
	
	// Draw exhaust particles first (behind helicopter)
	p.ExhaustParticles.Draw()
	
	// Draw helicopter body
	bodyRect := rl.Rectangle{
		X:      p.Position.X - p.Size,
		Y:      p.Position.Y - p.Size/2,
		Width:  p.Size * 2,
		Height: p.Size,
	}
	
	// Draw rotated body
	rl.DrawRectanglePro(bodyRect, 
					   rl.Vector2{X: p.Size, Y: p.Size/2},
					   p.Rotation,
					   p.Color)
	
	// Draw main rotor
	rotorPos := p.Position
	rotorLength := p.MainRotor.Size
	
	// Main rotor horizontal line
	rl.DrawLineEx(
		rl.Vector2{X: rotorPos.X - rotorLength * float32(math.Cos(float64(p.MainRotor.Rotation)*math.Pi/180)), 
				  Y: rotorPos.Y - rotorLength * float32(math.Sin(float64(p.MainRotor.Rotation)*math.Pi/180))},
		rl.Vector2{X: rotorPos.X + rotorLength * float32(math.Cos(float64(p.MainRotor.Rotation)*math.Pi/180)), 
				  Y: rotorPos.Y + rotorLength * float32(math.Sin(float64(p.MainRotor.Rotation)*math.Pi/180))},
		2.0,
		rl.White,
	)
	
	// Main rotor vertical line (perpendicular)
	rl.DrawLineEx(
		rl.Vector2{X: rotorPos.X - rotorLength * float32(math.Cos(float64(p.MainRotor.Rotation+90)*math.Pi/180)), 
				  Y: rotorPos.Y - rotorLength * float32(math.Sin(float64(p.MainRotor.Rotation+90)*math.Pi/180))},
		rl.Vector2{X: rotorPos.X + rotorLength * float32(math.Cos(float64(p.MainRotor.Rotation+90)*math.Pi/180)), 
				  Y: rotorPos.Y + rotorLength * float32(math.Sin(float64(p.MainRotor.Rotation+90)*math.Pi/180))},
		2.0,
		rl.White,
	)
	
	// Draw tail rotor
	tailOffset := rl.Vector2{
		X: -p.Size * 1.5 * float32(math.Cos(float64(p.Rotation)*math.Pi/180)),
		Y: -p.Size * 1.5 * float32(math.Sin(float64(p.Rotation)*math.Pi/180)),
	}
	
	tailRotorPos := rl.Vector2{
		X: p.Position.X + tailOffset.X,
		Y: p.Position.Y + tailOffset.Y,
	}
	
	tailRotorLength := p.TailRotor.Size / 2
	
	rl.DrawLineEx(
		rl.Vector2{X: tailRotorPos.X - tailRotorLength * float32(math.Cos(float64(p.TailRotor.Rotation)*math.Pi/180)), 
				  Y: tailRotorPos.Y - tailRotorLength * float32(math.Sin(float64(p.TailRotor.Rotation)*math.Pi/180))},
		rl.Vector2{X: tailRotorPos.X + tailRotorLength * float32(math.Cos(float64(p.TailRotor.Rotation)*math.Pi/180)), 
				  Y: tailRotorPos.Y + tailRotorLength * float32(math.Sin(float64(p.TailRotor.Rotation)*math.Pi/180