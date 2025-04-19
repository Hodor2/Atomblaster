// systems/input_system.go
package systems

import (
    "atomblaster/components"
    "atomblaster/constants"
    "math"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// InputSystem handles player input and controls
type InputSystem struct {
    entityManager *components.EntityManager
    positionID    components.ComponentID
    velocityID    components.ComponentID
    playerID      components.ComponentID
    colliderID    components.ComponentID
    tagID         components.ComponentID
    lifetimeID    components.ComponentID
    fireCooldown  float32
    currentState  *int
    audio         interface{} // Would be a proper AudioSystem in the real implementation
}

// NewInputSystem creates a new input system
func NewInputSystem(
    entityManager *components.EntityManager,
    registry *components.ComponentTypeRegistry,
    currentState *int,
    audio interface{},
) *InputSystem {
    positionID, _ := registry.GetID("Position")
    velocityID, _ := registry.GetID("Velocity")
    playerID, _ := registry.GetID("Player")
    colliderID, _ := registry.GetID("Collider")
    tagID, _ := registry.GetID("Tag")
    lifetimeID, _ := registry.GetID("Lifetime")
    
    return &InputSystem{
        entityManager: entityManager,
        positionID:    positionID,
        velocityID:    velocityID,
        playerID:      playerID,
        colliderID:    colliderID,
        tagID:         tagID,
        lifetimeID:    lifetimeID,
        fireCooldown:  0,
        currentState:  currentState,
        audio:         audio,
    }
}

// Update processes input and updates entity states accordingly
func (s *InputSystem) Update(dt float32) {
    // Handle game state transitions first
    s.handleStateTransitions()
    
    // Only process gameplay input if we're in the game state
    if *s.currentState != constants.StateGame {
        return
    }
    
    // Find the player entity
    playerEntities := s.entityManager.GetEntitiesWithComponents(s.playerID, s.positionID, s.velocityID)
    if len(playerEntities) == 0 {
        return // No player found
    }
    
    playerEntity := playerEntities[0]
    
    // Get components
    posComp, _ := s.entityManager.GetComponent(playerEntity, s.positionID)
    velComp, _ := s.entityManager.GetComponent(playerEntity, s.velocityID)
    playerComp, _ := s.entityManager.GetComponent(playerEntity, s.playerID)
    
    position := posComp.(*components.Position)
    velocity := velComp.(*components.Velocity)
    player := playerComp.(*components.Player)
    
    // Handle keyboard movement
    s.handleMovementInput(player, velocity, dt)
    
    // Update dash state
    if player.IsDashing {
        player.DashTimer -= dt
        if player.DashTimer <= 0 {
            player.IsDashing = false
        }
    }
    
    // Handle shooting
    s.handleShootingInput(player, position, dt)
}

// handleStateTransitions processes inputs for changing game states
func (s *InputSystem) handleStateTransitions() {
    switch *s.currentState {
    case constants.StateIntro:
        // Handled by the IntroController
        
    case constants.StateTitle:
        // Handled by the TitleController
        
    case constants.StateBossIntro:
        // Handled by the BossIntroController
        
    case constants.StateGame:
        // Escape to pause
        if rl.IsKeyPressed(rl.KeyEscape) {
            *s.currentState = constants.StatePause
        }
        
    case constants.StatePause:
        // Handled by the PauseController
        
    case constants.StateGameOver:
        // Handled by the GameOverController
    }
}

// handleMovementInput processes keyboard input for movement
func (s *InputSystem) handleMovementInput(player *components.Player, velocity *components.Velocity, dt float32) {
    // Reset velocity
    velocity.Value.X = 0
    velocity.Value.Y = 0
    
    // Process movement keys
    var dx, dy float32
    
    if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
        dy -= 1
    }
    if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
        dy += 1
    }
    if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
        dx -= 1
    }
    if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
        dx += 1
    }
    
    // Normalize diagonal movement
    if dx != 0 && dy != 0 {
        length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
        dx /= length
        dy /= length
    }
    
    // Apply movement speed
    moveSpeed := player.Speed
    if player.IsDashing {
        moveSpeed = player.Speed * 3 // Dash is 3x normal speed
    }
    
    velocity.Value.X = dx * moveSpeed
    velocity.Value.Y = dy * moveSpeed
    
    // Process dash input (space key)
    if rl.IsKeyPressed(rl.KeySpace) && !player.IsDashing {
        player.IsDashing = true
        player.DashTimer = 0.2
    }
}

// handleShootingInput processes mouse input for shooting
func (s *InputSystem) handleShootingInput(player *components.Player, position *components.Position, dt float32) {
    // Update cooldown timer
    s.fireCooldown -= dt
    
    // Only allow shooting if player has a gun and cooldown is expired
    if player.HasGun && s.fireCooldown <= 0 && rl.IsMouseButtonDown(rl.MouseLeftButton) {
        // Reset cooldown
        s.fireCooldown = constants.FireCooldownDuration
        
        // Create bullet entity at player position
        s.spawnBullet(position.Value)
        
        // Play sound (audio would be handled separately)
    }
}

// spawnBullet creates a new bullet entity
func (s *InputSystem) spawnBullet(playerPos rl.Vector2) {
    // Create bullet entity
    entityID := s.entityManager.CreateEntity()
    
    // Calculate bullet direction based on mouse position
    mousePos := rl.GetMousePosition()
    dir := rl.Vector2Subtract(mousePos, playerPos)
    
    // Normalize direction
    dir = rl.Vector2Normalize(dir)
    
    // Set bullet velocity
    bulletSpeed := float32(constants.BulletSpeed)
    bulletVel := rl.Vector2Scale(dir, bulletSpeed)
    
    // Add components to the entity
    s.entityManager.AddComponent(entityID, components.NewPosition(playerPos.X, playerPos.Y, s.entityManager.Registry))
    s.entityManager.AddComponent(entityID, components.NewVelocity(bulletVel.X, bulletVel.Y, s.entityManager.Registry))
    s.entityManager.AddComponent(entityID, components.NewCircleCollider(5, s.entityManager.Registry))
    s.entityManager.AddComponent(entityID, components.NewTag(components.BulletTag, s.entityManager.Registry))
    s.entityManager.AddComponent(entityID, components.NewLifetime(constants.BulletLifetime, s.entityManager.Registry))
}

// Draw is empty for InputSystem as it doesn't render anything
func (s *InputSystem) Draw() {
    // Input system doesn't need to draw anything
}

// RequiredComponents returns the component types this system operates on
func (s *InputSystem) RequiredComponents() []components.ComponentID {
    return []components.ComponentID{s.playerID, s.positionID, s.velocityID}
}
