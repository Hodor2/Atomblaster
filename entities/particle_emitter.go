package entities

import (
	"github.com/gen2brain/raylib-go/raylib"
	"math"
)

// Particle represents a single particle in an effect
type Particle struct {
	Position  rl.Vector2
	Velocity  rl.Vector2
	Color     rl.Color
	Size      float32
	Lifetime  float32 // Remaining lifetime in seconds
	MaxLife   float32 // Total lifetime in seconds
	Active    bool
}

// ParticleEmitter generates and manages particles
type ParticleEmitter struct {
	Position        rl.Vector2
	BaseColor       rl.Color
	ParticleSize    float32
	EmissionRate    float32 // Particles per second
	ParticleLife    float32 // Seconds per particle
	ParticleSpread  float32 // How wide particles spread
	GravityEffect   float32 // How much gravity affects particles
	Particles       []Particle
	TimeSinceLastEmit float32
}

// NewParticleEmitter creates a new particle emitter
func NewParticleEmitter(position rl.Vector2, color rl.Color, size float32, 
						emissionRate float32, lifetime float32) *ParticleEmitter {
	
	return &ParticleEmitter{
		Position:       position,
		BaseColor:      color,
		ParticleSize:   size,
		EmissionRate:   emissionRate,
		ParticleLife:   lifetime,
		ParticleSpread: math.Pi, // Default spread is 180 degrees
		GravityEffect:  0,       // Default no gravity
		Particles:      make([]Particle, 100), // Pre-allocate 100 particles
		TimeSinceLastEmit: 0,
	}
}

// Update updates all particles and emits new ones
func (e *ParticleEmitter) Update(dt float32) {
	// Update existing particles
	for i := range e.Particles {
		if e.Particles[i].Active {
			// Update position
			e.Particles[i].Position.X += e.Particles[i].Velocity.X * dt
			e.Particles[i].Position.Y += e.Particles[i].Velocity.Y * dt
			
			// Apply gravity if enabled
			if e.GravityEffect != 0 {
				e.Particles[i].Velocity.Y += e.GravityEffect * dt
			}
			
			// Update lifetime
			e.Particles[i].Lifetime -= dt
			
			// Fade out particle as it ages
			alpha := e.Particles[i].Lifetime / e.Particles[i].MaxLife
			e.Particles[i].Color.A = uint8(255 * alpha)
			
			// Deactivate expired particles
			if e.Particles[i].Lifetime <= 0 {
				e.Particles[i].Active = false
			}
		}
	}
	
	// Emit new particles
	e.TimeSinceLastEmit += dt
	particlesToEmit := e.EmissionRate * dt
	
	// Handle fractional particles with a time accumulator
	for particlesToEmit > 0 {
		if int(particlesToEmit) > 0 || rl.GetRandomFloat32() < particlesToEmit {
			// Find inactive particle
			for i := range e.Particles {
				if !e.Particles[i].Active {
					// Random velocity in specified spread angle
					angle := math.Pi + (rl.GetRandomFloat32() - 0.5) * e.ParticleSpread
					speed := 20.0 + rl.GetRandomFloat32() * 10.0
					
					// Add color variation
					r := int(e.BaseColor.R) + (int(rl.GetRandomFloat32() * 40) - 20)
					g := int(e.BaseColor.G) + (int(rl.GetRandomFloat32() * 40) - 20)
					b := int(e.BaseColor.B) + (int(rl.GetRandomFloat32() * 40) - 20)
					
					// Clamp color values
					if r < 0 { r = 0 } else if r > 255 { r = 255 }
					if g < 0 { g = 0 } else if g > 255 { g = 255 }
					if b < 0 { b = 0 } else if b > 255 { b = 255 }
					
					// Create particle
					lifetime := e.ParticleLife * (0.8 + rl.GetRandomFloat32() * 0.4)
					e.Particles[i] = Particle{
						Position: e.Position,
						Velocity: rl.Vector2{
							X: float32(math.Cos(float64(angle))) * speed,
							Y: float32(math.Sin(float64(angle))) * speed,
						},
						Color:    rl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 255},
						Size:     e.ParticleSize * (0.7 + rl.GetRandomFloat32() * 0.6),
						Lifetime: lifetime,
						MaxLife:  lifetime,
						Active:   true,
					}
					break
				}
			}
		}
		particlesToEmit--
	}
}

// Draw renders all active particles
func (e *ParticleEmitter) Draw() {
	for i := range e.Particles {
		if e.Particles[i].Active {
			rl.DrawCircleV(e.Particles[i].Position, e.Particles[i].Size, e.Particles[i].Color)
		}
	}
}

// SetPosition updates the emitter position
func (e *ParticleEmitter) SetPosition(pos rl.Vector2) {
	e.Position = pos
}

// SetColor updates the base particle color
func (e *ParticleEmitter) SetColor(color rl.Color) {
	e.BaseColor = color
}

// SetEmissionRate updates how many particles are emitted per second
func (e *ParticleEmitter) SetEmissionRate(rate float32) {
	e.EmissionRate = rate
}

// SetParticleLife updates how long particles live
func (e *ParticleEmitter) SetParticleLife(lifetime float32) {
	e.ParticleLife = lifetime
}

// SetGravity updates the gravity effect on particles
func (e *ParticleEmitter) SetGravity(gravity float32) {
	e.GravityEffect = gravity
}

// SetSpread updates the angle range particles are emitted in (in radians)
func (e *ParticleEmitter) SetSpread(spread float32) {
	e.ParticleSpread = spread
}

// CreateExplosion returns a pre-configured particle emitter for an explosion
func CreateExplosion(position rl.Vector2, size float32, color rl.Color) *ParticleEmitter {
	emitter := NewParticleEmitter(
		position,
		color,
		size/10, // Particle size relative to explosion size
		200,     // Many particles
		0.5,     // Short lifetime
	)
	
	// Configure for explosion
	emitter.ParticleSpread = 2 * math.Pi // Full 360 degree spread
	emitter.GravityEffect = 50           // Some gravity
	
	return emitter
}

// CreateSpark returns a pre-configured particle emitter for sparks
func CreateSpark(position rl.Vector2, direction float32, color rl.Color) *ParticleEmitter {
	emitter := NewParticleEmitter(
		position, 
		color,
		1.0,  // Small particles
		50,   // Medium emission rate
		0.3,  // Very short lifetime
	)
	
	// Configure for directional sparks
	emitter.ParticleSpread = math.Pi / 6 // Narrow spread (30 degrees)
	emitter.SetGravity(20)
	
	return emitter
}

// Burst emits a specified number of particles immediately
func (e *ParticleEmitter) Burst(count int) {
	for i := 0; i < count; i++ {
		// Find inactive particle
		for j := range e.Particles {
			if !e.Particles[j].Active {
				// Random velocity in a circle
				angle := rl.GetRandomFloat32() * 2 * math.Pi
				speed := 20.0 + rl.GetRandomFloat32() * 40.0
				
				// Add color variation
				r := int(e.BaseColor.R) + (int(rl.GetRandomFloat32() * 40) - 20)
				g := int(e.BaseColor.G) + (int(rl.GetRandomFloat32() * 40) - 20)
				b := int(e.BaseColor.B) + (int(rl.GetRandomFloat32() * 40) - 20)
				
				// Clamp color values
				if r < 0 { r = 0 } else if r > 255 { r = 255 }
				if g < 0 { g = 0 } else if g > 255 { g = 255 }
				if b < 0 { b = 0 } else if b > 255 { b = 255 }
				
				// Create particle
				lifetime := e.ParticleLife * (0.8 + rl.GetRandomFloat32() * 0.4)
				e.Particles[j] = Particle{
					Position: e.Position,
					Velocity: rl.Vector2{
						X: float32(math.Cos(float64(angle))) * speed,
						Y: float32(math.Sin(float64(angle))) * speed,
					},
					Color:    rl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 255},
					Size:     e.ParticleSize * (0.7 + rl.GetRandomFloat32() * 0.6),
					Lifetime: lifetime,
					MaxLife:  lifetime,
					Active:   true,
				}
				break
			}
		}
	}
}
