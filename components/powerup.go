// components/powerup.go
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