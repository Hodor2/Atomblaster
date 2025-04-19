// components/powerup_component.go
package components

// PowerUpType represents different types of power-ups
type PowerUpType int

const (
    PowerUpGun PowerUpType = iota
    PowerUpHealth
    PowerUpSpeed
)

// PowerUp component represents a power-up effect
type PowerUp struct {
    Type     PowerUpType
    Duration float32  // Duration in seconds (0 for permanent)
    Value    float32  // Value of the power-up effect
    id       ComponentID
}

// NewPowerUp creates a new PowerUp component
func NewPowerUp(powerUpType PowerUpType, duration, value float32, registry *ComponentTypeRegistry) *PowerUp {
    id, _ := registry.GetID("PowerUp")
    return &PowerUp{
        Type:     powerUpType,
        Duration: duration,
        Value:    value,
        id:       id,
    }
}

// GetComponentID returns the component's unique ID
func (p *PowerUp) GetComponentID() ComponentID {
    return p.id
}

// components/lifetime_component.go
package components

// Lifetime component gives an entity a limited lifespan
type Lifetime struct {
    Remaining float32
    id        ComponentID
}

// NewLifetime creates a new Lifetime component
func NewLifetime(duration float32, registry *ComponentTypeRegistry) *Lifetime {
    id, _ := registry.GetID("Lifetime")
    return &Lifetime{
        Remaining: duration,
        id:        id,
    }
}

// GetComponentID returns the component's unique ID
func (l *Lifetime) GetComponentID() ComponentID {
    return l.id
}
