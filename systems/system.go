// systems/system.go
package systems

import (
    "atomblaster/components"
)

// System is the interface for all systems in the ECS architecture
type System interface {
    Update(dt float32)
    Draw()
    RequiredComponents() []components.ComponentID
}

// SystemManager manages and updates all registered systems
type SystemManager struct {
    systems       []System
    entityManager *components.EntityManager
}

// NewSystemManager creates a new system manager
func NewSystemManager(entityManager *components.EntityManager) *SystemManager {
    return &SystemManager{
        systems:       make([]System, 0),
        entityManager: entityManager,
    }
}

// AddSystem adds a system to the manager
func (m *SystemManager) AddSystem(system System) {
    m.systems = append(m.systems, system)
}

// UpdateAll updates all systems
func (m *SystemManager) UpdateAll(dt float32) {
    for _, system := range m.systems {
        system.Update(dt)
    }
}

// DrawAll renders all systems
func (m *SystemManager) DrawAll() {
    for _, system := range m.systems {
        system.Draw()
    }
}

// GetEntityManager returns the entity manager
func (m *SystemManager) GetEntityManager() *components.EntityManager {
    return m.entityManager
}