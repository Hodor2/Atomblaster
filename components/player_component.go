// components/player_component.go
package components

// Player component contains player-specific properties
type Player struct {
    Speed     float32
    HasGun    bool
    IsDashing bool
    DashTimer float32
    id        ComponentID
}

// NewPlayer creates a new Player component
func NewPlayer(speed float32, registry *ComponentTypeRegistry) *Player {
    id, _ := registry.GetID("Player")
    return &Player{
        Speed:     speed,
        HasGun:    false,
        IsDashing: false,
        DashTimer: 0,
        id:        id,
    }
}

// GetComponentID returns the component's unique ID
func (p *Player) GetComponentID() ComponentID {
    return p.id
}

// components/enemy_component.go
package components

// EnemyType represents different types of enemies
type EnemyType int

const (
    NormalAtom EnemyType = iota
    FastAtom
    BigAtom
    Boss
)

// Enemy component contains enemy-specific properties
type Enemy struct {
    Type       EnemyType
    Speed      float32
    FireRate   float32
    SpinSpeed  float32
    Rotation   float32
    id         ComponentID
}

// NewEnemy creates a new Enemy component
func NewEnemy(enemyType EnemyType, speed float32, registry *ComponentTypeRegistry) *Enemy {
    id, _ := registry.GetID("Enemy")
    
    spinSpeed := float32(2.0)
    switch enemyType {
    case FastAtom:
        spinSpeed = 4.0
    case BigAtom:
        spinSpeed = 1.0
    case Boss:
        spinSpeed = 0.5
    }
    
    return &Enemy{
        Type:      enemyType,
        Speed:     speed,
        FireRate:  0,
        SpinSpeed: spinSpeed,
        Rotation:  0,
        id:        id,
    }
}

// GetComponentID returns the component's unique ID
func (e *Enemy) GetComponentID() ComponentID {
    return e.id
}

// components/scientist_component.go
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
        FollowOffset: rl.Vector2{X: float32(rl.GetRandomValue(-15, 15)), Y: float32(rl.GetRandomValue(-15, 15))},
        AnimTimer:   0,
        id:          id,
    }
}

// GetComponentID returns the component's unique ID
func (s *Scientist) GetComponentID() ComponentID {
    return s.id
}
