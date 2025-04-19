// components/lifetime.go
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