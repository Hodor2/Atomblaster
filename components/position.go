// components/position.go
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