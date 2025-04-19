// components/health.go
package components

// Health component represents an entity's health points
type Health struct {
    Current int
    Max     int
    id      ComponentID
}

// NewHealth creates a new Health component
func NewHealth(current, max int, registry *ComponentTypeRegistry) *Health {
    id, _ := registry.GetID("Health")
    return &Health{
        Current: current,
        Max:     max,
        id:      id,
    }
}

// GetComponentID returns the component's unique ID
func (h *Health) GetComponentID() ComponentID {
    return h.id
}

// TakeDamage reduces the entity's health and returns true if the entity is still alive
func (h *Health) TakeDamage(amount int) bool {
    h.Current -= amount
    if h.Current < 0 {
        h.Current = 0
    }
    return h.Current > 0
}

// Heal increases the entity's health up to its maximum
func (h *Health) Heal(amount int) {
    h.Current += amount
    if h.Current > h.Max {
        h.Current = h.Max
    }
}