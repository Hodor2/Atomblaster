// components/component.go
package components

// ComponentID is a unique identifier for component types
type ComponentID uint32

// Component is the basic interface that all components must implement
type Component interface {
    GetComponentID() ComponentID
}

// ComponentTypeRegistry keeps track of registered component types
type ComponentTypeRegistry struct {
    nextID      ComponentID
    componentIDs map[string]ComponentID
}

// NewComponentTypeRegistry creates a new component type registry
func NewComponentTypeRegistry() *ComponentTypeRegistry {
    return &ComponentTypeRegistry{
        nextID:      1, // Start at 1, reserving 0 for invalid ID
        componentIDs: make(map[string]ComponentID),
    }
}

// Register registers a new component type and returns its unique ID
func (r *ComponentTypeRegistry) Register(name string) ComponentID {
    if id, exists := r.componentIDs[name]; exists {
        return id
    }
    
    id := r.nextID
    r.componentIDs[name] = id
    r.nextID++
    
    return id
}

// GetID returns the ID for a registered component type
func (r *ComponentTypeRegistry) GetID(name string) (ComponentID, bool) {
    id, exists := r.componentIDs[name]
    return id, exists
}

// GetIDByName is a helper method to get component ID by name
func (r *ComponentTypeRegistry) GetIDByName(name string) ComponentID {
    id, _ := r.GetID(name)
    return id
}
