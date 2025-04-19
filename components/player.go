// components/player.go
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