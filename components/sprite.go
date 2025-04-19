// components/sprite.go
package components

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// Sprite component represents a visual representation of an entity
type Sprite struct {
    Texture     rl.Texture2D
    SourceRect  rl.Rectangle
    Rotation    float32
    Scale       float32
    Tint        rl.Color
    id          ComponentID
}

// NewSprite creates a new Sprite component
func NewSprite(texture rl.Texture2D, registry *ComponentTypeRegistry) *Sprite {
    id, _ := registry.GetID("Sprite")
    return &Sprite{
        Texture:    texture,
        SourceRect: rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)},
        Rotation:   0,
        Scale:      1.0,
        Tint:       rl.White,
        id:         id,
    }
}

// GetComponentID returns the component's unique ID
func (s *Sprite) GetComponentID() ComponentID {
    return s.id
}