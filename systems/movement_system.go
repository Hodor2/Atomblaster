// systems/movement_system.go
package systems

import (
    "atomblaster/components"
    "atomblaster/constants"
)

// MovementSystem handles updating positions based on velocities
type MovementSystem struct {
    entityManager *components.EntityManager
    positionID    components.ComponentID
    velocityID    components.ComponentID
}

// NewMovementSystem creates a new movement system
func NewMovementSystem(entityManager *components.EntityManager, registry *components.ComponentTypeRegistry) *MovementSystem {
    positionID, _ := registry.GetID("Position")
    velocityID, _ := registry.GetID("Velocity")
    
    return &MovementSystem{
        entityManager: entityManager,
        positionID:    positionID,
        velocityID:    velocityID,
    }
}

// Update moves all entities that have both Position and Velocity components
func (s *MovementSystem) Update(dt float32) {
    // Get all entities with both Position and Velocity components
    entities := s.entityManager.GetEntitiesWithComponents(s.positionID, s.velocityID)
    
    for _, entityID := range entities {
        // Get components
        posComp, _ := s.entityManager.GetComponent(entityID, s.positionID)
        velComp, _ := s.entityManager.GetComponent(entityID, s.velocityID)
        
        position := posComp.(*components.Position)
        velocity := velComp.(*components.Velocity)
        
        // Update position based on velocity
        position.Value.X += velocity.Value.X * dt
        position.Value.Y += velocity.Value.Y * dt
        
        // Optional: Handle screen bounds
        // This keeps entities within the screen bounds, with a small margin
        // You might want different behavior for some entities
        
        // Check if entity has a collider component to determine bounds
        colliderID, _ := s.entityManager.GetEntityManager().Registry.GetID("Collider")
        if colliderComp, has := s.entityManager.GetComponent(entityID, colliderID); has {
            collider := colliderComp.(*components.Collider)
            
            var margin float32
            
            if collider.Type == components.CircleCollider {
                margin = collider.Radius
            } else {
                margin = (collider.Width + collider.Height) / 4
            }
            
            // Constrain position to screen bounds (with margin)
            if position.Value.X < margin {
                position.Value.X = margin
                velocity.Value.X *= -1 // Bounce
            }
            
            if position.Value.X > constants.ScreenWidth-margin {
                position.Value.X = constants.ScreenWidth - margin
                velocity.Value.X *= -1 // Bounce
            }
            
            if position.Value.Y < margin {
                position.Value.Y = margin
                velocity.Value.Y *= -1 // Bounce
            }
            
            if position.Value.Y > constants.ScreenHeight-margin {
                position.Value.Y = constants.ScreenHeight - margin
                velocity.Value.Y *= -1 // Bounce
            }
        }
    }
}

// Draw is empty for MovementSystem as it doesn't render anything
func (s *MovementSystem) Draw() {
    // Movement system doesn't need to draw anything
}

// RequiredComponents returns the component types this system operates on
func (s *MovementSystem) RequiredComponents() []components.ComponentID {
    return []components.ComponentID{s.positionID, s.velocityID}
}