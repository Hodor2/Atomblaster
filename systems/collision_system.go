// systems/collision_system.go
package systems

import (
    "atomblaster/components"
    "atomblaster/constants"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// CollisionSystem handles detection and resolution of collisions between entities
type CollisionSystem struct {
    entityManager *components.EntityManager
    positionID    components.ComponentID
    colliderID    components.ComponentID
    tagID         components.ComponentID
    healthID      components.ComponentID
    playerID      components.ComponentID
    scoreValue    *int // Pointer to the score value in the game state
}

// NewCollisionSystem creates a new collision system
func NewCollisionSystem(entityManager *components.EntityManager, registry *components.ComponentTypeRegistry, score *int) *CollisionSystem {
    positionID, _ := registry.GetID("Position")
    colliderID, _ := registry.GetID("Collider")
    tagID, _ := registry.GetID("Tag")
    healthID, _ := registry.GetID("Health")
    playerID, _ := registry.GetID("Player")
    
    return &CollisionSystem{
        entityManager: entityManager,
        positionID:    positionID,
        colliderID:    colliderID,
        tagID:         tagID,
        healthID:      healthID,
        playerID:      playerID,
        scoreValue:    score,
    }
}

// Update checks for and handles collisions between entities
func (s *CollisionSystem) Update(dt float32) {
    // Get all entities with Position and Collider components
    entities := s.entityManager.GetEntitiesWithComponents(s.positionID, s.colliderID)
    
    // Find the player entity
    var playerEntity components.EntityID
    var playerPosition *components.Position
    var playerCollider *components.Collider
    
    playerEntities := s.entityManager.GetEntitiesWithComponents(s.playerID, s.positionID, s.colliderID)
    if len(playerEntities) > 0 {
        playerEntity = playerEntities[0]
        posComp, _ := s.entityManager.GetComponent(playerEntity, s.positionID)
        collComp, _ := s.entityManager.GetComponent(playerEntity, s.colliderID)
        playerPosition = posComp.(*components.Position)
        playerCollider = collComp.(*components.Collider)
    }
    
    // Process player-specific collisions first (if player exists)
    if playerEntity != 0 && playerPosition != nil && playerCollider != nil {
        s.handlePlayerCollisions(playerEntity, playerPosition, playerCollider, entities)
    }
    
    // Process bullet collisions with enemies
    s.handleBulletCollisions(entities)
    
    // Process other special collisions (scientists, rescue zone, etc.)
    s.handleSpecialCollisions(entities)
}

// handlePlayerCollisions checks for collisions between the player and other entities
func (s *CollisionSystem) handlePlayerCollisions(
    playerEntity components.EntityID,
    playerPos *components.Position,
    playerCollider *components.Collider,
    entities []components.EntityID,
) {
    // Get player health if available
    var playerHealth *components.Health
    if healthComp, has := s.entityManager.GetComponent(playerEntity, s.healthID); has {
        playerHealth = healthComp.(*components.Health)
    }
    
    // Get player component
    playerComp, _ := s.entityManager.GetComponent(playerEntity, s.playerID)
    player := playerComp.(*components.Player)
    
    // Check collisions with all other entities
    for _, entityID := range entities {
        // Skip self
        if entityID == playerEntity {
            continue
        }
        
        // Get entity components
        posComp, _ := s.entityManager.GetComponent(entityID, s.positionID)
        collComp, _ := s.entityManager.GetComponent(entityID, s.colliderID)
        position := posComp.(*components.Position)
        collider := collComp.(*components.Collider)
        
        // Check if entity has a tag
        tagComp, hasTag := s.entityManager.GetComponent(entityID, s.tagID)
        if !hasTag {
            continue
        }
        
        tag := tagComp.(*components.Tag)
        
        // Handle collision based on entity tag
        switch tag.Type {
        case components.EnemyTag:
            // Check for collision between player and enemy
            if s.checkCollision(playerPos.Value, playerCollider, position.Value, collider) {
                // Player hit by enemy
                if playerHealth != nil && !player.IsDashing {
                    playerHealth.TakeDamage(1)
                    
                    // Spawn particle effect
                    s.spawnCollisionParticles(playerPos.Value, 30, rl.Red, 3.0)
                    
                    // Play sound (handled separately)
                    
                    // Remove atom or make player briefly invincible
                    s.entityManager.DestroyEntity(entityID)
                    
                    // Activate dash for brief invincibility
                    player.IsDashing = true
                    player.DashTimer = 0.2
                }
            }
            
        case components.PowerUpTag:
            // Check for collision between player and power-up
            if s.checkCollision(playerPos.Value, playerCollider, position.Value, collider) {
                // Apply power-up effect
                powerUpID, _ := s.entityManager.GetEntityManager().Registry.GetID("PowerUp")
                if powerUpComp, has := s.entityManager.GetComponent(entityID, powerUpID); has {
                    powerUp := powerUpComp.(*components.PowerUp)
                    
                    switch powerUp.Type {
                    case components.PowerUpGun:
                        player.HasGun = true
                        *s.scoreValue += 25
                        s.spawnCollisionParticles(position.Value, 15, rl.Orange, 3.0)
                        
                    case components.PowerUpHealth:
                        if playerHealth != nil {
                            playerHealth.Heal(1)
                            *s.scoreValue += 15
                            s.spawnCollisionParticles(position.Value, 15, rl.Green, 3.0)
                        }
                        
                    case components.PowerUpSpeed:
                        player.Speed += 50
                        *s.scoreValue += 20
                        s.spawnCollisionParticles(position.Value, 15, rl.Purple, 3.0)
                    }
                    
                    // Remove power-up
                    s.entityManager.DestroyEntity(entityID)
                }
            }
            
        case components.DoorTag:
            // Check if player is at the door and conditions are met to advance level
            if s.checkCollision(playerPos.Value, playerCollider, position.Value, collider) {
                // Level completion logic would go here
                // This will be handled by a separate GameStateSystem
            }
        }
    }
}

// handleBulletCollisions checks for collisions between bullets and other entities
func (s *CollisionSystem) handleBulletCollisions(entities []components.EntityID) {
    // Get bullet entities
    bulletEntities := []components.EntityID{}
    for _, entityID := range entities {
        if tagComp, has := s.entityManager.GetComponent(entityID, s.tagID); has {
            tag := tagComp.(*components.Tag)
            if tag.Type == components.BulletTag {
                bulletEntities = append(bulletEntities, entityID)
            }
        }
    }
    
    // Check each bullet against potential targets
    for _, bulletID := range bulletEntities {
        // Get bullet position and collider
        posComp, _ := s.entityManager.GetComponent(bulletID, s.positionID)
        collComp, _ := s.entityManager.GetComponent(bulletID, s.colliderID)
        bulletPos := posComp.(*components.Position)
        bulletCollider := collComp.(*components.Collider)
        
        // Check against all potential targets
        for _, targetID := range entities {
            // Skip self
            if targetID == bulletID {
                continue
            }
            
            // Get target components
            tagComp, hasTag := s.entityManager.GetComponent(targetID, s.tagID)
            if !hasTag {
                continue
            }
            
            tag := tagComp.(*components.Tag)
            
            // Only check collision with enemies
            if tag.Type != components.EnemyTag && tag.Type != components.BossTag {
                continue
            }
            
            posComp, _ := s.entityManager.GetComponent(targetID, s.positionID)
            collComp, _ := s.entityManager.GetComponent(targetID, s.colliderID)
            targetPos := posComp.(*components.Position)
            targetCollider := collComp.(*components.Collider)
            
            // Check for collision
            if s.checkCollision(bulletPos.Value, bulletCollider, targetPos.Value, targetCollider) {
                // Hit detected!
                
                // Check if enemy has health
                if healthComp, has := s.entityManager.GetComponent(targetID, s.healthID); has {
                    health := healthComp.(*components.Health)
                    
                    // Apply damage
                    if !health.TakeDamage(10) {
                        // Enemy defeated
                        *s.scoreValue += 10
                        
                        // Spawn particles
                        s.spawnCollisionParticles(targetPos.Value, 15, rl.Yellow, 2.0)
                        
                        // Check for boss
                        if tag.Type == components.BossTag {
                            *s.scoreValue += 1990 // Total 2000 for boss
                            s.spawnCollisionParticles(targetPos.Value, 50, rl.Orange, 5.0)
                            s.entityManager.DestroyEntity(targetID)
                        } else {
                            // Regular enemy - possibility to spawn power-up
                            if rl.GetRandomValue(0, 100) < 15 {
                                s.spawnPowerUp(targetPos.Value)
                            }
                            
                            // Destroy the enemy
                            s.entityManager.DestroyEntity(targetID)
                        }
                    } else {
                        // Enemy damaged but not defeated
                        *s.scoreValue += 5
                        s.spawnCollisionParticles(targetPos.Value, 5, rl.Yellow, 1.0)
                    }
                } else {
                    // Enemy has no health component - destroy immediately
                    *s.scoreValue += 10
                    s.spawnCollisionParticles(targetPos.Value, 15, rl.Yellow, 2.0)
                    s.entityManager.DestroyEntity(targetID)
                }
                
                // Destroy the bullet
                s.entityManager.DestroyEntity(bulletID)
                
                // Break since this bullet can't hit anything else
                break
            }
        }
    }
}

// handleSpecialCollisions handles special collision types like scientists and rescue zones
func (s *CollisionSystem) handleSpecialCollisions(entities []components.EntityID) {
    // Find the player entity first
    var playerPos *components.Position
    var playerCollider *components.Collider
    
    playerEntities := s.entityManager.GetEntitiesWithComponents(s.playerID, s.positionID, s.colliderID)
    if len(playerEntities) > 0 {
        posComp, _ := s.entityManager.GetComponent(playerEntities[0], s.positionID)
        collComp, _ := s.entityManager.GetComponent(playerEntities[0], s.colliderID)
        playerPos = posComp.(*components.Position)
        playerCollider = collComp.(*components.Collider)
    }
    
    if playerPos == nil || playerCollider == nil {
        return // No player found
    }
    
    // Find all scientists
    scientistID, _ := s.entityManager.GetEntityManager().Registry.GetID("Scientist")
    scientists := s.entityManager.GetEntitiesWithComponents(s.tagID, s.positionID, scientistID)
    
    // Find rescue zone
    var rescueZonePos *components.Position
    var rescueZoneCollider *components.Collider
    
    for _, entityID := range entities {
        if tagComp, has := s.entityManager.GetComponent(entityID, s.tagID); has {
            tag := tagComp.(*components.Tag)
            if tag.Type == components.RescueZoneTag {
                posComp, _ := s.entityManager.GetComponent(entityID, s.positionID)
                collComp, _ := s.entityManager.GetComponent(entityID, s.colliderID)
                rescueZonePos = posComp.(*components.Position)
                rescueZoneCollider = collComp.(*components.Collider)
                break
            }
        }
    }
    
    // Process scientist pickups and rescues
    for _, scientistID := range scientists {
        // Get scientist components
        posComp, _ := s.entityManager.GetComponent(scientistID, s.positionID)
        sciComp, _ := s.entityManager.GetComponent(scientistID, scientistID)
        collComp, _ := s.entityManager.GetComponent(scientistID, s.colliderID)
        
        scientistPos := posComp.(*components.Position)
        scientist := sciComp.(*components.Scientist)
        scientistCollider := collComp.(*components.Collider)
        
        if scientist.State == components.Wandering {
            // Check if player is near scientist to pick up
            distance := rl.Vector2Distance(playerPos.Value, scientistPos.Value)
            if distance < 50.0 {
                scientist.State = components.FollowingPlayer
                // Play pickup sound (handled separately)
            }
        } else if scientist.State == components.FollowingPlayer {
            // Update position to follow player
            scientistPos.Value.X = playerPos.Value.X + scientist.FollowOffset.X
            scientistPos.Value.Y = playerPos.Value.Y + scientist.FollowOffset.Y
            
            // Check if in rescue zone
            if rescueZonePos != nil && rescueZoneCollider != nil {
                if s.checkCollision(scientistPos.Value, scientistCollider, rescueZonePos.Value, rescueZoneCollider) {
                    // Scientist rescued!
                    *s.scoreValue += 100
                    s.spawnCollisionParticles(scientistPos.Value, 20, rl.Green, 3.0)
                    
                    // Play sound (handled separately)
                    
                    // Mark scientist as rescued and remove
                    scientist.State = components.Rescued
                    s.entityManager.DestroyEntity(scientistID)
                }
            }
        }
    }
}

// checkCollision detects if two entities with colliders are intersecting
func (s *CollisionSystem) checkCollision(
    pos1 rl.Vector2, collider1 *components.Collider,
    pos2 rl.Vector2, collider2 *components.Collider,
) bool {
    // Handle circle-circle collision
    if collider1.Type == components.CircleCollider && collider2.Type == components.CircleCollider {
        // Calculate adjusted positions with offsets
        adjustedPos1 := rl.Vector2{
            X: pos1.X + collider1.Offset.X,
            Y: pos1.Y + collider1.Offset.Y,
        }
        
        adjustedPos2 := rl.Vector2{
            X: pos2.X + collider2.Offset.X,
            Y: pos2.Y + collider2.Offset.Y,
        }
        
        // Calculate distance between centers
        distance := rl.Vector2Distance(adjustedPos1, adjustedPos2)
        
        // Check if distance is less than sum of radii
        return distance < (collider1.Radius + collider2.Radius)
    }
    
    // Handle rectangle-rectangle collision
    if collider1.Type == components.RectangleCollider && collider2.Type == components.RectangleCollider {
        rect1 := rl.Rectangle{
            X:      pos1.X + collider1.Offset.X - collider1.Width/2,
            Y:      pos1.Y + collider1.Offset.Y - collider1.Height/2,
            Width:  collider1.Width,
            Height: collider1.Height,
        }
        
        rect2 := rl.Rectangle{
            X:      pos2.X + collider2.Offset.X - collider2.Width/2,
            Y:      pos2.Y + collider2.Offset.Y - collider2.Height/2,
            Width:  collider2.Width,
            Height: collider2.Height,
        }
        
        return rl.CheckCollisionRecs(rect1, rect2)
    }
    
    // Handle circle-rectangle collision
    if collider1.Type == components.CircleCollider && collider2.Type == components.RectangleCollider {
        // Circle
        adjustedPos1 := rl.Vector2{
            X: pos1.X + collider1.Offset.X,
            Y: pos1.Y + collider1.Offset.Y,
        }
        
        // Rectangle
        rect := rl.Rectangle{
            X:      pos2.X + collider2.Offset.X - collider2.Width/2,
            Y:      pos2.Y + collider2.Offset.Y - collider2.Height/2,
            Width:  collider2.Width,
            Height: collider2.Height,
        }
        
        return rl.CheckCollisionCircleRec(adjustedPos1, collider1.Radius, rect)
    }
    
    // Handle rectangle-circle collision
    if collider1.Type == components.RectangleCollider && collider2.Type == components.CircleCollider {
        // Rectangle
        rect := rl.Rectangle{
            X:      pos1.X + collider1.Offset.X - collider1.Width/2,
            Y:      pos1.Y + collider1.Offset.Y - collider1.Height/2,
            Width:  collider1.Width,
            Height: collider1.Height,
        }
        
        // Circle
        adjustedPos2 := rl.Vector2{
            X: pos2.X + collider2.Offset.X,
            Y: pos2.Y + collider2.Offset.Y,
        }
        
        return rl.CheckCollisionCircleRec(adjustedPos2, collider2.Radius, rect)
    }
    
    return false
}

// spawnCollisionParticles creates particle effects for collisions
func (s *CollisionSystem) spawnCollisionParticles(pos rl.Vector2, count int, color rl.Color, size float32) {
    // This would actually create particle entities, but for simplicity,
    // we'll assume a separate ParticleSystem handles this
    // In a real implementation, we'd create entities with Particle components
}

// spawnPowerUp creates a power-up entity at the given position
func (s *CollisionSystem) spawnPowerUp(pos rl.Vector2) {
    // This would actually create a power-up entity with the appropriate components
    // In a real implementation, we'd add that logic here
}

// Draw is empty for CollisionSystem as it doesn't render anything
func (s *CollisionSystem) Draw() {
    // Collision system doesn't need to draw anything
}

// RequiredComponents returns the component types this system operates on
func (s *CollisionSystem) RequiredComponents() []components.ComponentID {
    return []components.ComponentID{s.positionID, s.colliderID, s.tagID}
}