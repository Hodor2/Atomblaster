// components/velocity.go
package components

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// Velocity component represents an entity's velocity in 2D space
type Velocity struct {
    Value rl.Vector2
    id    ComponentID
}

// NewVelocity creates a new Velocity component
func NewVelocity(dx, dy float32, registry *ComponentTypeRegistry) *Velocity {
    id, _ := registry.GetID("Velocity")
    return &Velocity{
        Value: rl.Vector2{X: dx, Y: dy},
        id:    id,
    }
}

// GetComponentID returns the component's unique ID
func (v *Velocity) GetComponentID() ComponentID {
    return v.id
}