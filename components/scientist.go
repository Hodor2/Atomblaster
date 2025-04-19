// components/scientist.go
package components

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// ScientistState represents the current state of a scientist
type ScientistState int

const (
    Wandering ScientistState = iota
    FollowingPlayer
    Rescued
)

// Scientist component contains scientist-specific properties
type Scientist struct {
    State       ScientistState
    WanderTimer float32
    WanderDir   rl.Vector2
    FollowOffset rl.Vector2
    AnimTimer   float32
    id          ComponentID
}

// NewScientist creates a new Scientist component
func NewScientist(registry *ComponentTypeRegistry) *Scientist {
    id, _ := registry.GetID("Scientist")
    return &Scientist{
        State:       Wandering,
        WanderTimer: 0,
        WanderDir:   rl.Vector2{X: 0, Y: 0},
        FollowOffset: rl.Vector2{
            X: float32(rl.GetRandomValue(-15, 15)),
            Y: float32(rl.GetRandomValue(-15, 15)),
        },
        AnimTimer:   float32(rl.GetRandomValue(0, 100)) / 100.0, // Random start time
        id:          id,
    }
}

// GetComponentID returns the component's unique ID
func (s *Scientist) GetComponentID() ComponentID {
    return s.id
}