// systems/render_system.go
package systems

import (
    "atomblaster/components"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// RenderSystem handles drawing entities with Position and Sprite components
type RenderSystem struct {
    entityManager *components.EntityManager
    positionID    components.ComponentID
    spriteID      components.ComponentID
    background    rl.Texture2D
    debugMode     bool
}

// NewRenderSystem creates a new render system
func NewRenderSystem(entityManager *components.EntityManager, registry *components.ComponentTypeRegistry, background rl.Texture2D) *RenderSystem {
    positionID, _ := registry.GetID("Position")
    spriteID, _ := registry.GetID("Sprite")
    
    return &RenderSystem{
        entityManager: entityManager,
        positionID:    positionID,
        spriteID:      spriteID,
        background:    background,
        debugMode:     false,
    }
}

// Update is empty for RenderSystem as drawing happens in Draw
func (s *RenderSystem) Update(dt float32) {
    // Rendering system doesn't need to update anything
}

// Draw renders all entities with Position and Sprite components
func (s *RenderSystem) Draw() {
    // Draw background
    rl.DrawTexture(s.background, 0, 0, rl.White)
    
    // Get all entities with both Position and Sprite components
    entities := s.entityManager.GetEntitiesWithComponents(s.positionID, s.spriteID)
    
    // Sort entities by layer if needed (not implemented here)
    
    for _, entityID := range entities {
        // Get components
        posComp, _ := s.entityManager.GetComponent(entityID, s.positionID)
        spriteComp, _ := s.entityManager.GetComponent(entityID, s.spriteID)
        
        position := posComp.(*components.Position)
        sprite := spriteComp.(*components.Sprite)
        
        // Set up destination rectangle
        width := sprite.SourceRect.Width * sprite.Scale
        height := sprite.SourceRect.Height * sprite.Scale
        
        destRect := rl.Rectangle{
            X:      position.Value.X - width/2,
            Y:      position.Value.Y - height/2,
            Width:  width,
            Height: height,
        }
        
        // Calculate origin point for rotation
        origin := rl.Vector2{
            X: width / 2,
            Y: height / 2,
        }
        
        // Draw the sprite
        rl.DrawTexturePro(
            sprite.Texture,
            sprite.SourceRect,
            destRect,
            origin,
            sprite.Rotation,
            sprite.Tint,
        )
        
        // Optional: Draw debug info for entities with colliders
        if s.debugMode {
            s.drawDebugColliders(entityID, position.Value)
        }
    }
    
    // Draw entities without sprites but with special rendering needs
    s.drawSpecialEntities()
}

// drawDebugColliders draws debug visualization for colliders
func (s *RenderSystem) drawDebugColliders(entityID components.EntityID, position rl.Vector2) {
    // Get collider component ID
    colliderID, _ := s.entityManager.Registry.GetID("Collider")
    
    // Check if entity has a collider
    colliderComp, has := s.entityManager.GetComponent(entityID, colliderID)
    if !has {
        return
    }
    
    collider := colliderComp.(*components.Collider)
    
    // Draw different shapes based on collider type
    if collider.Type == components.CircleCollider {
        // Draw circle collider
        rl.DrawCircleLines(
            int32(position.X + collider.Offset.X),
            int32(position.Y + collider.Offset.Y),
            collider.Radius,
            rl.Green,
        )
    } else {
        // Draw rectangle collider
        rect := rl.Rectangle{
            X:      position.X + collider.Offset.X - collider.Width/2,
            Y:      position.Y + collider.Offset.Y - collider.Height/2,
            Width:  collider.Width,
            Height: collider.Height,
        }
        
        rl.DrawRectangleLinesEx(rect, 1, rl.Green)
    }
}

// drawSpecialEntities draws entities that need special rendering logic
func (s *RenderSystem) drawSpecialEntities() {
    // Get tag component ID
    tagID, _ := s.entityManager.Registry.GetID("Tag")
    positionID, _ := s.entityManager.Registry.GetID("Position")
    
    // Draw special entities like rescue zone and door
    specialEntities := s.entityManager.GetEntitiesWithComponents(tagID, positionID)
    for _, entityID := range specialEntities {
        tagComp, _ := s.entityManager.GetComponent(entityID, tagID)
        tag := tagComp.(*components.Tag)
        
        // Draw rescue zone
        if tag.Type == components.RescueZoneTag {
            s.drawRescueZone(entityID)
        } // Draw door
        else if tag.Type == components.DoorTag {
            s.drawDoor(entityID)
        }
    }
}

// drawRescueZone draws a rescue zone entity
func (s *RenderSystem) drawRescueZone(entityID components.EntityID) {
    posComp, _ := s.entityManager.GetComponent(entityID, s.positionID)
    position := posComp.(*components.Position)
    
    // Get collider component for dimensions
    colliderID, _ := s.entityManager.Registry.GetID("Collider")
    colliderComp, has := s.entityManager.GetComponent(entityID, colliderID)
    if !has {
        return
    }
    
    collider := colliderComp.(*components.Collider)
    
    // Draw zone with pulsing effect
    rl.DrawRectangle(
        int32(position.Value.X - collider.Width/2),
        int32(position.Value.Y - collider.Height/2),
        int32(collider.Width),
        int32(collider.Height),
        rl.Fade(rl.Green, 0.3),
    )
    
    // Draw border
    rl.DrawRectangleLinesEx(
        rl.Rectangle{
            X:      position.Value.X - collider.Width/2,
            Y:      position.Value.Y - collider.Height/2,
            Width:  collider.Width,
            Height: collider.Height,
        },
        2, rl.Fade(rl.Green, 0.8),
    )
    
    // Draw "RESCUE ZONE" text
    textWidth := rl.MeasureText("RESCUE ZONE", 20)
    rl.DrawText(
        "RESCUE ZONE",
        int32(position.Value.X - float32(textWidth)/2),
        int32(position.Value.Y - 10),
        20,
        rl.White,
    )
}

// drawDoor draws a door entity
func (s *RenderSystem) drawDoor(entityID components.EntityID) {
    posComp, _ := s.entityManager.GetComponent(entityID, s.positionID)
    position := posComp.(*components.Position)
    
    // Get collider for dimensions
    colliderID, _ := s.entityManager.Registry.GetID("Collider")
    colliderComp, has := s.entityManager.GetComponent(entityID, colliderID)
    if !has {
        return
    }
    
    collider := colliderComp.(*components.Collider)
    
    // Determine door color based on game state
    // For a real implementation, we would check if all requirements are met
    doorColor := rl.Red // Default to closed
    
    // Draw the door
    rl.DrawRectangle(
        int32(position.Value.X - collider.Width/2),
        int32(position.Value.Y - collider.Height/2),
        int32(collider.Width),
        int32(collider.Height),
        doorColor,
    )
}

// ToggleDebugMode turns debug visualization on/off
func (s *RenderSystem) ToggleDebugMode() {
    s.debugMode = !s.debugMode
}

// RequiredComponents returns the component types this system operates on
func (s *RenderSystem) RequiredComponents() []components.ComponentID {
    return []components.ComponentID{s.positionID, s.spriteID}
}
