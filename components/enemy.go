// components/enemy.go
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