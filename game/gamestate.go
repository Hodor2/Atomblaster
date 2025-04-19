// game/gamestate.go
package game

import (
    "atomblaster/audio"
    "atomblaster/components"
    "atomblaster/constants"
    "atomblaster/systems"
    "atomblaster/ui"
    "atomblaster/ui/controllers"
    "atomblaster/ui/models"
    "atomblaster/ui/views"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// GameState holds the current state of the game
type GameState struct {
    // Game state
    CurrentState   int
    Score          int
    Health         int
    Level          int
    ScientistsRescued int
    TotalScientists   int
    StartTime      int64
    ElapsedTime    int64
    GameOver       bool
    BossDefeated   bool
    IsBossLevel    bool
    
    // ECS Framework
    ComponentRegistry *components.ComponentTypeRegistry
    EntityManager    *components.EntityManager
    SystemManager    *systems.SystemManager
    
    // Systems
    MovementSystem   *systems.MovementSystem
    RenderSystem     *systems.RenderSystem
    CollisionSystem  *systems.CollisionSystem
    InputSystem      *systems.InputSystem
    ParticleSystem   *systems.ParticleSystem
    
    // UI Screens
    IntroScreen     *ui.Screen
    TitleScreen     *ui.Screen
    GameScreen      *ui.Screen
    PauseScreen     *ui.Screen
    GameOverScreen  *ui.Screen
    BossIntroScreen *ui.Screen
    
    // Audio
    Audio *audio.AudioSystem
    
    // Assets
    Background     rl.Texture2D
    PlayerSprite   rl.Texture2D
    EnemySprite    rl.Texture2D
    BulletSprite   rl.Texture2D
    PowerUpSprites [3]rl.Texture2D
}

// NewGameState creates a new game state
func NewGameState(audioSystem *audio.AudioSystem) *GameState {
    g := &GameState{
        CurrentState:      constants.StateIntro,
        Score:             0,
        Health:            3,
        Level:             1,
        ScientistsRescued: 0,
        TotalScientists:   0,
        GameOver:          false,
        IsBossLevel:       false,
        BossDefeated:      false,
        Audio:             audioSystem,
    }
    
    // Initialize assets
    g.initializeAssets()
    
    // Set start time
    g.StartTime = int64(rl.GetTime())
    
    // Initialize ECS framework
    g.initializeECS()
    
    // Set up the first level
    g.initLevel()
    
    // Initialize screens
    g.createScreens()
    
    return g
}

// initializeAssets loads all game textures
func (g *GameState) initializeAssets() {
    g.Background = rl.LoadTexture("assets/background.png")
    g.PlayerSprite = rl.LoadTexture("assets/helicopter.png")
    g.EnemySprite = rl.LoadTexture("assets/atom.png")
    g.BulletSprite = rl.LoadTexture("assets/bullet.png")
    
    // Load power-up sprites
    g.PowerUpSprites[0] = rl.LoadTexture("assets/powerup_weapon.png")
    g.PowerUpSprites[1] = rl.LoadTexture("assets/powerup_health.png")
    g.PowerUpSprites[2] = rl.LoadTexture("assets/powerup_speed.png")
}

// initializeECS sets up the Entity Component System
func (g *GameState) initializeECS() {
    // Create component registry
    g.ComponentRegistry = components.NewComponentTypeRegistry()
    
    // Register component types
    g.ComponentRegistry.Register("Position")
    g.ComponentRegistry.Register("Velocity")
    g.ComponentRegistry.Register("Sprite")
    g.ComponentRegistry.Register("Collider")
    g.ComponentRegistry.Register("Health")
    g.ComponentRegistry.Register("Tag")
    g.ComponentRegistry.Register("PowerUp")
    g.ComponentRegistry.Register("Lifetime")
    g.ComponentRegistry.Register("Player")
    g.ComponentRegistry.Register("Enemy")
    g.ComponentRegistry.Register("Scientist")
    
    // Create entity manager
    g.EntityManager = components.NewEntityManager(g.ComponentRegistry)
    
    // Create system manager
    g.SystemManager = systems.NewSystemManager(g.EntityManager)
    
    // Create systems
    g.MovementSystem = systems.NewMovementSystem(g.EntityManager, g.ComponentRegistry)
    g.RenderSystem = systems.NewRenderSystem(g.EntityManager, g.ComponentRegistry, g.Background)
    g.CollisionSystem = systems.NewCollisionSystem(g.EntityManager, g.ComponentRegistry, &g.Score)
    g.InputSystem = systems.NewInputSystem(g.EntityManager, g.ComponentRegistry, &g.CurrentState, g.Audio)
    g.ParticleSystem = systems.NewParticleSystem(g.EntityManager, g.ComponentRegistry)
    
    // Add systems to system manager in order of execution
    g.SystemManager.AddSystem(g.InputSystem)
    g.SystemManager.AddSystem(g.MovementSystem)
    g.SystemManager.AddSystem(g.CollisionSystem)
    g.SystemManager.AddSystem(g.ParticleSystem)
    g.SystemManager.AddSystem(g.RenderSystem)
}

// initLevel resets the level state and spawns entities
func (g *GameState) initLevel() {
    // Reset player position and properties
    g.createPlayer()
    
    // Reset state
    g.ScientistsRescued = 0
    g.GameOver = false
    g.BossDefeated = false
    
    // Check if this is the boss level (level 5)
    g.IsBossLevel = (g.Level == 5)
    
    if g.IsBossLevel {
        // Create boss
        g.createBoss()
        g.CurrentState = constants.StateBossIntro
    } else {
        // Create normal level entities
        g.createAtoms()
        g.createScientists()
        g.createRescueZone()
    }
    
    // Create door
    g.createDoor()
    
    // Spawn powerups
    g.createPowerUps()
}

// createPlayer creates the player entity with all required components
func (g *GameState) createPlayer() {
    // Create player entity
    playerID := g.EntityManager.CreateEntity()
    
    // Add components
    g.EntityManager.AddComponent(playerID, components.NewPosition(float32(constants.ScreenWidth/4), float32(constants.ScreenHeight/2), g.ComponentRegistry))
    g.EntityManager.AddComponent(playerID, components.NewVelocity(0, 0, g.ComponentRegistry))
    g.EntityManager.AddComponent(playerID, components.NewCircleCollider(constants.HelicopterWidth/2, g.ComponentRegistry))
    g.EntityManager.AddComponent(playerID, components.NewSprite(g.PlayerSprite, g.ComponentRegistry))
    g.EntityManager.AddComponent(playerID, components.NewHealth(g.Health, 10, g.ComponentRegistry))
    g.EntityManager.AddComponent(playerID, components.NewTag(components.PlayerTag, g.ComponentRegistry))
    g.EntityManager.AddComponent(playerID, components.NewPlayer(300, g.ComponentRegistry))
}

// createAtoms creates enemy atom entities
func (g *GameState) createAtoms() {
    // Number of atoms based on level
    numAtoms := g.Level*2 + 3
    if g.IsBossLevel {
        numAtoms = 5 // Fewer atoms in boss level
    }
    
    for i := 0; i < numAtoms; i++ {
        // Random position
        pos := rl.Vector2{
            X: float32(rl.GetRandomValue(constants.ScreenWidth/4, 3*constants.ScreenWidth/4)),
            Y: float32(rl.GetRandomValue(20, constants.ScreenHeight-40)),
        }
        
        // Create random velocity vector based on speed and level
        speed := float32(100+g.Level*10) + float32(rl.GetRandomValue(-20, 20))
        vel := rl.Vector2{
            X: float32(rl.GetRandomValue(-100, 100)) / 100.0 * speed,
            Y: float32(rl.GetRandomValue(-100, 100)) / 100.0 * speed,
        }
        
        // Determine atom type
        atomType := components.NormalAtom
        if g.IsBossLevel && rl.GetRandomValue(0, 1) == 1 {
            atomType = components.FastAtom
        }
        
        // Create atom entity
        atomID := g.EntityManager.CreateEntity()
        
        // Add components
        g.EntityManager.AddComponent(atomID, components.NewPosition(pos.X, pos.Y, g.ComponentRegistry))
        g.EntityManager.AddComponent(atomID, components.NewVelocity(vel.X, vel.Y, g.ComponentRegistry))
        g.EntityManager.AddComponent(atomID, components.NewSprite(g.EnemySprite, g.ComponentRegistry))
        g.EntityManager.AddComponent(atomID, components.NewTag(components.EnemyTag, g.ComponentRegistry))
        g.EntityManager.AddComponent(atomID, components.NewEnemy(atomType, speed, g.ComponentRegistry))
        
        // Add health and collider based on type
        if atomType == components.NormalAtom {
            g.EntityManager.AddComponent(atomID, components.NewCircleCollider(15, g.ComponentRegistry))
            g.EntityManager.AddComponent(atomID, components.NewHealth(2, 2, g.ComponentRegistry))
        } else if atomType == components.FastAtom {
            g.EntityManager.AddComponent(atomID, components.NewCircleCollider(12, g.ComponentRegistry))
            g.EntityManager.AddComponent(atomID, components.NewHealth(1, 1, g.ComponentRegistry))
        } else if atomType == components.BigAtom {
            g.EntityManager.AddComponent(atomID, components.NewCircleCollider(25, g.ComponentRegistry))
            g.EntityManager.AddComponent(atomID, components.NewHealth(4, 4, g.ComponentRegistry))
        }
    }
}

// createBoss creates the boss entity
func (g *GameState) createBoss() {
    // Create boss entity
    bossID := g.EntityManager.CreateEntity()
    
    // Add components
    g.EntityManager.AddComponent(bossID, components.NewPosition(float32(constants.ScreenWidth-200), 150, g.ComponentRegistry))
    g.EntityManager.AddComponent(bossID, components.NewVelocity(0, 0, g.ComponentRegistry))
    g.EntityManager.AddComponent(bossID, components.NewRectangleCollider(80, 40, g.ComponentRegistry))
    g.EntityManager.AddComponent(bossID, components.NewSprite(g.PlayerSprite, g.ComponentRegistry)) // Using player sprite for simplicity
    g.EntityManager.AddComponent(bossID, components.NewHealth(100, 100, g.ComponentRegistry))
    g.EntityManager.AddComponent(bossID, components.NewTag(components.BossTag, g.ComponentRegistry))
    g.EntityManager.AddComponent(bossID, components.NewEnemy(components.Boss, 200, g.ComponentRegistry))
}

// createScientists creates scientist entities to rescue
func (g *GameState) createScientists() {
    // Skip in boss level
    if g.IsBossLevel {
        return
    }
    
    // Number of scientists based on level
    scientistCount := 2 + g.Level
    if scientistCount > 8 {
        scientistCount = 8 // Maximum 8 scientists
    }
    
    g.TotalScientists = scientistCount
    
    for i := 0; i < scientistCount; i++ {
        // Place scientists around the level
        x := constants.ScreenWidth/4 + float32(rl.GetRandomValue(0, int32(constants.ScreenWidth/2)))
        y := constants.ScreenHeight/6 + float32(rl.GetRandomValue(0, int32(constants.ScreenHeight*2/3)))
        
        // Create scientist entity
        scientistID := g.EntityManager.CreateEntity()
        
        // Add components
        g.EntityManager.AddComponent(scientistID, components.NewPosition(x, y, g.ComponentRegistry))
        g.EntityManager.AddComponent(scientistID, components.NewVelocity(0, 0, g.ComponentRegistry))
        g.EntityManager.AddComponent(scientistID, components.NewCircleCollider(15, g.ComponentRegistry))
        g.EntityManager.AddComponent(scientistID, components.NewTag(components.ScientistTag, g.ComponentRegistry))
        g.EntityManager.AddComponent(scientistID, components.NewScientist(g.ComponentRegistry))
    }
}

// createRescueZone creates the rescue zone entity
func (g *GameState) createRescueZone() {
    // Skip in boss level
    if g.IsBossLevel {
        return
    }
    
    // Create rescue zone on the left side
    rescueX := float32(rl.GetRandomValue(50, 200))
    rescueY := constants.ScreenHeight - float32(rl.GetRandomValue(100, 200))
    
    // Create rescue zone entity
    rescueZoneID := g.EntityManager.CreateEntity()
    
    // Add components
    g.EntityManager.AddComponent(rescueZoneID, components.NewPosition(rescueX, rescueY, g.ComponentRegistry))
    g.EntityManager.AddComponent(rescueZoneID, components.NewRectangleCollider(100, 50, g.ComponentRegistry))
    g.EntityManager.AddComponent(rescueZoneID, components.NewTag(components.RescueZoneTag, g.ComponentRegistry))
}

// createDoor creates the door entity (level exit)
func (g *GameState) createDoor() {
    // Create door on the right side
    doorID := g.EntityManager.CreateEntity()
    
    // Add components
    g.EntityManager.AddComponent(doorID, components.NewPosition(
        constants.ScreenWidth - 35,
        constants.ScreenHeight/2,
        g.ComponentRegistry,
    ))
    g.EntityManager.AddComponent(doorID, components.NewRectangleCollider(30, 100, g.ComponentRegistry))
    g.EntityManager.AddComponent(doorID, components.NewTag(components.DoorTag, g.ComponentRegistry))
}

// createPowerUps creates power-up entities
func (g *GameState) createPowerUps() {
    // Create gun power-up if player doesn't have a gun
    playerEntities := g.EntityManager.GetEntitiesWithComponents(
        g.ComponentRegistry.GetIDByName("Player"),
        g.ComponentRegistry.GetIDByName("Tag"),
    )
    
    if len(playerEntities) > 0 {
        playerEntityID := playerEntities[0]
        playerComp, _ := g.EntityManager.GetComponent(
            playerEntityID,
            g.ComponentRegistry.GetIDByName("Player"),
        )
        player := playerComp.(*components.Player)
        
        if !player.HasGun {
            // Create gun power-up
            gunX := float32(rl.GetRandomValue(100, int32(constants.ScreenWidth-100)))
            gunY := float32(rl.GetRandomValue(100, int32(constants.ScreenHeight-100)))
            
            gunID := g.EntityManager.CreateEntity()
            
            g.EntityManager.AddComponent(gunID, components.NewPosition(gunX, gunY, g.ComponentRegistry))
            g.EntityManager.AddComponent(gunID, components.NewCircleCollider(15, g.ComponentRegistry))
            g.EntityManager.AddComponent(gunID, components.NewSprite(g.PowerUpSprites[0], g.ComponentRegistry))
            g.EntityManager.AddComponent(gunID, components.NewTag(components.PowerUpTag, g.ComponentRegistry))
            g.EntityManager.AddComponent(gunID, components.NewPowerUp(components.PowerUpGun, 0, 0, g.ComponentRegistry))
        }
    }
    
    // Health power-up (more common in higher levels and boss level)
    healthChance := int32(50)
    if g.Level > 1 || g.IsBossLevel {
        healthChance = 100 // Always give health in boss level
    }
    
    if rl.GetRandomValue(0, 100) < healthChance {
        healthX := float32(rl.GetRandomValue(100, int32(constants.ScreenWidth-100)))
        healthY := float32(rl.GetRandomValue(100, int32(constants.ScreenHeight-100)))
        
        healthID := g.EntityManager.CreateEntity()
        
        g.EntityManager.AddComponent(healthID, components.NewPosition(healthX, healthY, g.ComponentRegistry))
        g.EntityManager.AddComponent(healthID, components.NewCircleCollider(15, g.ComponentRegistry))
        g.EntityManager.AddComponent(healthID, components.NewSprite(g.PowerUpSprites[1], g.ComponentRegistry))
        g.EntityManager.AddComponent(healthID, components.NewTag(components.PowerUpTag, g.ComponentRegistry))
        g.EntityManager.AddComponent(healthID, components.NewPowerUp(components.PowerUpHealth, 0, 0, g.ComponentRegistry))
    }
    
    // Add speed boost pickup (more common in boss level)
    speedChance := int32(20)
    if g.IsBossLevel {
        speedChance = 50
    }
    
    if rl.GetRandomValue(0, 100) < speedChance {
        speedX := float32(rl.GetRandomValue(50, int32(constants.ScreenWidth-100)))
        speedY := float32(rl.GetRandomValue(50, int32(constants.ScreenHeight-100)))
        
        speedID := g.EntityManager.CreateEntity()
        
        g.EntityManager.AddComponent(speedID, components.NewPosition(speedX, speedY, g.ComponentRegistry))
        g.EntityManager.AddComponent(speedID, components.NewCircleCollider(15, g.ComponentRegistry))
        g.EntityManager.AddComponent(speedID, components.NewSprite(g.PowerUpSprites[2], g.ComponentRegistry))
        g.EntityManager.AddComponent(speedID, components.NewTag(components.PowerUpTag, g.ComponentRegistry))
        g.EntityManager.AddComponent(speedID, components.NewPowerUp(components.PowerUpSpeed, 0, 50, g.ComponentRegistry))
    }
}

// createScreens initializes all game screens using the MVC pattern
func (g *GameState) createScreens() {
    // Create intro screen
    introModel := models.NewIntroModel(g.Background, g.PlayerSprite)
    introView := views.NewIntroView(introModel)
    introController := controllers.NewIntroController(introModel, &g.CurrentState)
    g.IntroScreen = ui.NewScreen(introModel, introView, introController)
    
    // Create title screen
    titleModel := models.NewTitleModel(g.Background)
    titleView := views.NewTitleView(titleModel)
    titleController := controllers.NewTitleController(titleModel, &g.CurrentState)
    g.TitleScreen = ui.NewScreen(titleModel, titleView, titleController)
    
    // Create game screen
    gameModel := models.NewGameModel(
        g.Background,
        g.PlayerSprite,
        g.EnemySprite,
        g.BulletSprite,
        g.PowerUpSprites,
        &g.Score,
        &g.Health,
        &g.Level,
        &g.ScientistsRescued,
        &g.TotalScientists,
        &g.StartTime,
        &g.ElapsedTime,
    )
    gameView := views.NewGameView(gameModel)
    gameController := controllers.NewGameController(gameModel)
    g.GameScreen = ui.NewScreen(gameModel, gameView, gameController)
    
    // Create pause screen
    pauseModel := models.NewPauseModel(gameModel)
    pauseView := views.NewPauseView(pauseModel, gameView)
    pauseController := controllers.NewPauseController(pauseModel, &g.CurrentState, g.ResetGame)
    g.PauseScreen = ui.NewScreen(pauseModel, pauseView, pauseController)
    
    // Create game over screen
    gameOverModel := models.NewGameOverModel(gameModel, false)
    gameOverView := views.NewGameOverView(gameOverModel, gameView)
    gameOverController := controllers.NewGameOverController(gameOverModel, &g.CurrentState, g.ResetGame)
    g.GameOverScreen = ui.NewScreen(gameOverModel, gameOverView, gameOverController)
    
    // Create boss intro screen
    bossIntroModel := models.NewBossIntroModel(g.Background, g.PlayerSprite, g.PlayerSprite)
    bossIntroView := views.NewBossIntroView(bossIntroModel)
    bossIntroController := controllers.NewBossIntroController(bossIntroModel, &g.CurrentState)
    g.BossIntroScreen = ui.NewScreen(bossIntroModel, bossIntroView, bossIntroController)
}

// ResetGame resets the game state to start a new game
func (g *GameState) ResetGame() {
    // Reset game state
    g.Score = 0
    g.Health = 3
    g.Level = 1
    g.ScientistsRescued = 0
    g.TotalScientists = 0
    g.GameOver = false
    g.IsBossLevel = false
    g.BossDefeated = false
    
    // Reset start time
    g.StartTime = int64(rl.GetTime())
    g.ElapsedTime = 0
    
    // Initialize first level
    g.initLevel()
}

// Draw renders the current game state
func (g *GameState) Draw() {
    rl.BeginDrawing()
    rl.ClearBackground(rl.Black)
    
    switch g.CurrentState {
    case constants.StateIntro:
        g.IntroScreen.Draw()
        
    case constants.StateTitle:
        g.TitleScreen.Draw()
        
    case constants.StateBossIntro:
        g.BossIntroScreen.Draw()
        
    case constants.StateGame:
        // The render system handles drawing the game world
        g.RenderSystem.Draw()
        
        // Draw UI overlay
        g.GameScreen.Draw()
        
    case constants.StatePause:
        g.PauseScreen.Draw()
        
    case constants.StateGameOver:
        g.GameOverScreen.Draw()
    }
    
    rl.EndDrawing()
}

// Update updates the game state based on input and time passing
func (g *GameState) Update(dt float32) {
    // Update elapsed time
    g.ElapsedTime = int64(rl.GetTime()) - g.StartTime
    
    // Update based on current state
    switch g.CurrentState {
    case constants.StateIntro:
        if g.IntroScreen.Update() {
            g.CurrentState = constants.StateTitle
        }
        
    case constants.StateTitle:
        if g.TitleScreen.Update() {
            g.CurrentState = constants.StateGame
        }
        
    case constants.StateBossIntro:
        if g.BossIntroScreen.Update() {
            g.CurrentState = constants.StateGame
        }
        
    case constants.StateGame:
        // Update game systems
        g.updateGame(dt)
        
        // Check for game over conditions
        if g.Health <= 0 {
            g.CurrentState = constants.StateGameOver
        }
        
    case constants.StatePause:
        if g.PauseScreen.Update() {
            // Controller handles state changes
        }
        
    case constants.StateGameOver:
        if g.GameOverScreen.Update() {
            // Controller handles restart/quit
        }
    }
}

// updateGame handles all game updates during gameplay
func (g *GameState) updateGame(dt float32) {
    // Update all ECS systems
    g.SystemManager.UpdateAll(dt)
    
    // Update game state based on entity state
    g.updateGameState()
}

// updateGameState updates the game state based on entity state
func (g *GameState) updateGameState() {
    // Update player health
    playerEntities := g.EntityManager.GetEntitiesWithComponents(
        g.ComponentRegistry.GetIDByName("Player"),
        g.ComponentRegistry.GetIDByName("Health"),
    )
    
    if len(playerEntities) > 0 {
        playerEntityID := playerEntities[0]
        healthComp, _ := g.EntityManager.GetComponent(
            playerEntityID,
            g.ComponentRegistry.GetIDByName("Health"),
        )
        health := healthComp.(*components.Health)
        
        g.Health = health.Current
    }
    
    // Update scientists rescued count
    // In a real implementation, we would have a component or system for tracking this
    
    // Check if boss is defeated
    if g.IsBossLevel && !g.BossDefeated {
        bossEntities := g.EntityManager.GetEntitiesWithComponents(
            g.ComponentRegistry.GetIDByName("Tag"),
            g.ComponentRegistry.GetIDByName("Health"),
        )
        
        bossAlive := false
        for _, entityID := range bossEntities {
            tagComp, _ := g.EntityManager.GetComponent(
                entityID,
                g.ComponentRegistry.GetIDByName("Tag"),
            )
            tag := tagComp.(*components.Tag)
            
            if tag.Type == components.BossTag {
                healthComp, _ := g.EntityManager.GetComponent(
                    entityID,
                    g.ComponentRegistry.GetIDByName("Health"),
                )
                health := healthComp.(*components.Health)
                
                if health.Current > 0 {
                    bossAlive = true
                    break
                }
            }
        }
        
        g.BossDefeated = !bossAlive
    }
}
