// components/collider.go
package components

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// ColliderType defines the type of collider shape
type ColliderType int

const (
    CircleCollider ColliderType = iota
    RectangleCollider
)

// Collider component represents a collision area for an entity
type Collider struct {
    Type        ColliderType
    Radius      float32      // Used for circle colliders
    Width       float32      // Used for rectangle colliders
    Height      float32      // Used for rectangle colliders
    Offset      rl.Vector2   // Offset from the entity's position
    IsTrigger   bool         // If true, doesn't cause physical collision
    id          ComponentID
}

// NewCircleCollider creates a new circle collider component
func NewCircleCollider(radius float32, registry *ComponentTypeRegistry) *Collider {
    id, _ := registry.GetID("Collider")
    return &Collider{
        Type:      CircleCollider,
        Radius:    radius,
        Offset:    rl.Vector2{X: 0, Y: 0},
        IsTrigger: false,
        id:        id,
    }
}

// NewRectangleCollider creates a new rectangle collider component
func NewRectangleCollider(width, height float32, registry *ComponentTypeRegistry) *Collider {
    id, _ := registry.GetID("Collider")
    return &Collider{
        Type:      RectangleCollider,
        Width:     width,
        Height:    height,
        Offset:    rl.Vector2{X: 0, Y: 0},
        IsTrigger: false,
        id:        id,
    }
}

// GetComponentID returns the component's unique ID
func (c *Collider) GetComponentID() ComponentID {
    return c.id
}

// GetBounds returns the collider's bounds as a rectangle, based on the entity's position
func (c *Collider) GetBounds(position rl.Vector2) rl.Rectangle {
    switch c.Type {
    case RectangleCollider:
        return rl.Rectangle{
            X:      position.X + c.Offset.X - c.Width/2,
            Y:      position.Y + c.Offset.Y - c.Height/2,
            Width:  c.Width,
            Height: c.Height,
        }
    case CircleCollider:
        // Approximating the circle with a square for simplicity
        return rl.Rectangle{
            X:      position.X + c.Offset.X - c.Radius,
            Y:      position.Y + c.Offset.Y - c.Radius,
            Width:  c.Radius * 2,
            Height: c.Radius * 2,
        }
    default:
        return rl.Rectangle{}
    }
}