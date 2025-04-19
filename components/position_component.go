// components/position_component.go
package components

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// Position component represents an entity's position in 2D space
type Position struct {
    Value rl.Vector2
    id    ComponentID
}

// NewPosition creates a new Position component
func NewPosition(x, y float32, registry *ComponentTypeRegistry) *Position {
    id, _ := registry.GetID("Position")
    return &Position{
        Value: rl.Vector2{X: x, Y: y},
        id:    id,
    }
}

// GetComponentID returns the component's unique ID
func (p *Position) GetComponentID() ComponentID {
    return p.id
}

// components/velocity_component.go
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
