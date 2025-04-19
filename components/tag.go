// components/tag.go
package components

// TagType represents different types of entity tags
type TagType int

const (
    PlayerTag TagType = iota
    EnemyTag
    BulletTag
    PowerUpTag
    ScientistTag
    RescueZoneTag
    DoorTag
    BossTag
)

// Tag component identifies the entity type
type Tag struct {
    Type TagType
    id   ComponentID
}

// NewTag creates a new Tag component
func NewTag(tagType TagType, registry *ComponentTypeRegistry) *Tag {
    id, _ := registry.GetID("Tag")
    return &Tag{
        Type: tagType,
        id:   id,
    }
}

// GetComponentID returns the component's unique ID
func (t *Tag) GetComponentID() ComponentID {
    return t.id
}