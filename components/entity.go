// components/entity.go
package components

// EntityID represents a unique entity identifier
type EntityID uint32

// EntityManager manages entity IDs and their components
type EntityManager struct {
    nextEntityID    EntityID
    entities        map[EntityID]bool
    componentStores map[ComponentID]map[EntityID]Component
    Registry        *ComponentTypeRegistry // Made public for access from systems
}

// NewEntityManager creates a new entity manager with the given component type registry
func NewEntityManager(registry *ComponentTypeRegistry) *EntityManager {
    return &EntityManager{
        nextEntityID:    1, // Start at 1, reserving 0 for invalid entity
        entities:        make(map[EntityID]bool),
        componentStores: make(map[ComponentID]map[EntityID]Component),
        Registry:        registry,
    }
}

// CreateEntity creates a new entity and returns its ID
func (m *EntityManager) CreateEntity() EntityID {
    id := m.nextEntityID
    m.nextEntityID++
    m.entities[id] = true
    return id
}

// DestroyEntity removes an entity and all its components
func (m *EntityManager) DestroyEntity(entityID EntityID) {
    if !m.entities[entityID] {
        return // Entity doesn't exist
    }
    
    // Remove all components for this entity
    for _, store := range m.componentStores {
        delete(store, entityID)
    }
    
    // Remove the entity
    delete(m.entities, entityID)
}

// AddComponent adds a component to an entity
func (m *EntityManager) AddComponent(entityID EntityID, component Component) {
    if !m.entities[entityID] {
        return // Entity doesn't exist
    }
    
    componentID := component.GetComponentID()
    
    // Create the component store if it doesn't exist
    if _, exists := m.componentStores[componentID]; !exists {
        m.componentStores[componentID] = make(map[EntityID]Component)
    }
    
    // Add the component to the store
    m.componentStores[componentID][entityID] = component
}

// RemoveComponent removes a component from an entity
func (m *EntityManager) RemoveComponent(entityID EntityID, componentID ComponentID) {
    if !m.entities[entityID] {
        return // Entity doesn't exist
    }
    
    store, exists := m.componentStores[componentID]
    if !exists {
        return // Component type doesn't exist
    }
    
    delete(store, entityID)
}

// GetComponent returns a component for an entity if it exists
func (m *EntityManager) GetComponent(entityID EntityID, componentID ComponentID) (Component, bool) {
    if !m.entities[entityID] {
        return nil, false // Entity doesn't exist
    }
    
    store, exists := m.componentStores[componentID]
    if !exists {
        return nil, false // Component type doesn't exist
    }
    
    comp, exists := store[entityID]
    return comp, exists
}

// HasComponent checks if an entity has a specific component
func (m *EntityManager) HasComponent(entityID EntityID, componentID ComponentID) bool {
    if !m.entities[entityID] {
        return false // Entity doesn't exist
    }
    
    store, exists := m.componentStores[componentID]
    if !exists {
        return false // Component type doesn't exist
    }
    
    _, exists = store[entityID]
    return exists
}

// GetEntitiesWithComponent returns all entities that have a specific component
func (m *EntityManager) GetEntitiesWithComponent(componentID ComponentID) []EntityID {
    store, exists := m.componentStores[componentID]
    if !exists {
        return []EntityID{} // Component type doesn't exist
    }
    
    entities := make([]EntityID, 0, len(store))
    for entityID := range store {
        entities = append(entities, entityID)
    }
    
    return entities
}

// GetEntitiesWithComponents returns all entities that have all the specified components
func (m *EntityManager) GetEntitiesWithComponents(componentIDs ...ComponentID) []EntityID {
    if len(componentIDs) == 0 {
        return []EntityID{}
    }
    
    // Start with entities that have the first component
    firstComponentID := componentIDs[0]
    store, exists := m.componentStores[firstComponentID]
    if !exists {
        return []EntityID{} // First component type doesn't exist
    }
    
    // Build the initial set of entities
    entities := make(map[EntityID]bool)
    for entityID := range store {
        entities[entityID] = true
    }
    
    // Filter out entities that don't have all the required components
    for _, componentID := range componentIDs[1:] {
        store, exists := m.componentStores[componentID]
        if !exists {
            return []EntityID{} // One of the component types doesn't exist
        }
        
        // Filter the entities
        for entityID := range entities {
            if _, exists := store[entityID]; !exists {
                delete(entities, entityID)
            }
        }
    }
    
    // Convert map keys to slice
    result := make([]EntityID, 0, len(entities))
    for entityID := range entities {
        result = append(result, entityID)
    }
    
    return result
}

// GetEntityManager returns itself (for compatibility with previous code)
func (m *EntityManager) GetEntityManager() *EntityManager {
    return m
}
