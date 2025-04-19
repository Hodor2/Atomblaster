// Package entities contains game entity implementations
package entities

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Particle represents a visual effect particle
type Particle struct {
	Pos      rl.Vector2
	Vel      rl.Vector2
	LifeTime float32
	Color    rl.Color
	Size     float32
}

// NewParticle creates a new particle
func NewParticle(pos rl.Vector2, vel rl.Vector2, lifetime float32, color rl.Color, size float32) Particle {
	return Particle{
		Pos:      pos,
		Vel:      vel,
		LifeTime: lifetime,
		Color:    color,
		Size:     size,
	}
}

// Update moves the particle and returns true if it's still active
func (p *Particle) Update(dt float32) bool {
	p.Pos.X += p.Vel.X * dt
	p.Pos.Y += p.Vel.Y * dt
	p.LifeTime -= dt

	if p.LifeTime <= 0 {
		return false
	}

	// Make particles get smaller as they age
	p.Size = p.Size * (p.LifeTime / 3.0)

	return true
}

// Draw renders the particle on screen
func (p *Particle) Draw() {
	// For small particles, just draw a pixel
	if p.Size <= 1.0 {
		rl.DrawPixelV(p.Pos, p.Color)
	} else {
		// For larger particles, draw a circle
		rl.DrawCircleV(p.Pos, p.Size, p.Color)
	}
}

// ParticleSystem manages multiple particles
type ParticleSystem struct {
	Particles []Particle
}

// NewParticleSystem creates a new particle system
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		Particles: make([]Particle, 0, 100), // Pre-allocate space for particles
	}
}

// Update updates all particles in the system
func (ps *ParticleSystem) Update(dt float32) {
	// Use a more efficient in-place filtering approach
	newSize := 0
	for i := range ps.Particles {
		if ps.Particles[i].Update(dt) {
			ps.Particles[newSize] = ps.Particles[i]
			newSize++
		}
	}
	ps.Particles = ps.Particles[:newSize]
}

// Draw renders all particles
func (ps *ParticleSystem) Draw() {
	for i := range ps.Particles {
		ps.Particles[i].Draw()
	}
}

// SpawnParticles creates particles at the given position
func (ps *ParticleSystem) SpawnParticles(pos rl.Vector2, count int, color rl.Color, speed float32) {
	for i := 0; i < count; i++ {
		// Random velocity in all directions
		angle := float32(2.0 * math.Pi * float64(i) / float64(count)) 
		particleSpeed := speed * (0.7 + float32(rl.GetRandomValue(0, 60))/100.0)

		vel := rl.Vector2{
			X: float32(math.Cos(float64(angle))) * particleSpeed,
			Y: float32(math.Sin(float64(angle))) * particleSpeed,
		}

		// Random size for variety
		size := float32(rl.GetRandomValue(1, 5))

		// Slightly randomize color
		colorVar := int32(rl.GetRandomValue(-20, 20))
		randomizedColor := rl.Color{
			R: clampUint8(int(color.R) + int(colorVar)),
			G: clampUint8(int(color.G) + int(colorVar)),
			B: clampUint8(int(color.B) + int(colorVar)),
			A: color.A,
		}

		// Add a small random offset to position
		offsetPos := rl.Vector2{
			X: pos.X + float32(rl.GetRandomValue(-5, 5)),
			Y: pos.Y + float32(rl.GetRandomValue(-5, 5)),
		}

		ps.Particles = append(ps.Particles, NewParticle(
			offsetPos,
			vel,
			1.0 * (0.7 + float32(rl.GetRandomValue(0, 60))/100.0), // Randomize lifetime a bit
			randomizedColor,
			size,
		))
	}
}

// Helper function to clamp uint8 values
func clampUint8(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
