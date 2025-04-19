// systems/particle_system.go
package systems

import (
    "atomblaster/components"
    "math"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// ParticleSystem manages visual effect particles
type ParticleSystem struct {
    entityManager *components.EntityManager
    positionID    components.ComponentID
    velocityID    components.ComponentID
    lifetimeID    components.ComponentID
    particleQueue []ParticleSpawnRequest
}

// ParticleSpawnRequest represents a request to spawn particles
type ParticleSpawnRequest struct {
    Position rl.Vector2
    Count    int
    Color    rl.Color
    Size     float32
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(entityManager *components.EntityManager, registry *components.ComponentTypeRegistry) *ParticleSystem {
    positionID, _ := registry.GetID("Position")
    velocityID, _ := registry.GetID("Velocity")
    lifetimeID, _ := registry.GetID("Lifetime")
    
    return &ParticleSystem{
        entityManager: entityManager,
        positionID:    positionID,
        velocityID:    velocityID,
        lifetimeID:    lifetimeID,
        particleQueue: make([]ParticleSpawnRequest, 0),
    }
}

// Update updates all particle entities and spawns new particles
func (s *ParticleSystem) Update(dt float32) {
    // Process particles in the queue
    for _, request := range s.particleQueue {
        s.spawnParticles(request.Position, request.Count, request.Color, request.Size)
    }
    
    // Clear the queue
    s.particleQueue = s.particleQueue[:0]
    
    // Get all entities with both Position, Velocity and Lifetime components
    entities := s.entityManager.GetEntitiesWithComponents(s.positionID, s.velocityID, s.lifetimeID)
    
    // Update each particle
    for _, entityID := range entities {
        // Get components
        lifetimeComp, _ := s.entityManager.GetComponent(entityID, s.lifetimeID)
        lifetime := lifetimeComp.(*components.Lifetime)
        
        // Update lifetime
        lifetime.Remaining -= dt
        
        // Check if particle has expired
        if lifetime.Remaining <= 0 {
            s.entityManager.DestroyEntity(entityID)
        }
    }
}

// RequestParticles adds a particle spawn request to the queue
func (s *ParticleSystem) RequestParticles(pos rl.Vector2, count int, color rl.Color, size float32) {
    s.particleQueue = append(s.particleQueue, ParticleSpawnRequest{
        Position: pos,
        Count:    count,
        Color:    color,
        Size:     size,
    })
}

// spawnParticles creates particle entities at the given position
func (s *ParticleSystem) spawnParticles(pos rl.Vector2, count int, color rl.Color, size float32) {
    for i := 0; i < count; i++ {
        // Random velocity in all directions
        angle := float32(2.0 * math.Pi * float64(i) / float64(count)) + float32(rl.GetRandomValue(-30, 30)) * 0.01
        particleSpeed := size * (0.7 + float32(rl.GetRandomValue(0, 60))/100.0)
        
        vel := rl.Vector2{
            X: float32(math.Cos(float64(angle))) * particleSpeed,
            Y: float32(math.Sin(float64(angle))) * particleSpeed,
        }
        
        // Random size for variety and actually use it
        particleSize := float32(rl.GetRandomValue(1, 5))
        
        // Slightly randomize color and actually use it
        colorVar := int(rl.GetRandomValue(-20, 20))
        randomizedColor := rl.Color{
            R: s.clampUint8(int(color.R) + colorVar),
            G: s.clampUint8(int(color.G) + colorVar),
            B: s.clampUint8(int(color.B) + colorVar),
            A: color.A,
        }
        
        // Add a small random offset to position
        offsetPos := rl.Vector2{
            X: pos.X + float32(rl.GetRandomValue(-5, 5)),
            Y: pos.Y + float32(rl.GetRandomValue(-5, 5)),
        }
        
        // Create particle entity
        entityID := s.entityManager.CreateEntity()
        