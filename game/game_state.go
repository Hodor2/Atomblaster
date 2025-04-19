package game

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/entities"
	"atomblaster/systems"
	"atomblaster/ui"
	"atomblaster/util"
	"fmt"
	"math"
)

// Game represents the main game state
type Game struct {
	// Main player
	Player      *entities.Player
	
	// Systems
	Camera      *Camera
	Quadtree    *util.Quadtree
	FoodSystem  *systems.FoodGenerator
	PowerUpSystem *systems.PowerUpManager
	
	// UI elements
	Minimap     *ui.Minimap
	
	// Game state tracking
	Score       int
	GameTime    float32
	State       int // Current game state (menu, playing, game over)
	
	// Performance stats
	FrameTime   float32
	EntityCount int
	
	// Debug
	DebugMode   bool
}

// New creates a new game instance
func New() *Game {
	// Create world bounds
	worldBounds := rl.Rectangle{
		X: 0,
		Y: 0,
		Width: constants.WorldWidth,
		Height: constants.WorldHeight,
	}

	// Initialize player
	player := entities.NewPlayer()
	
	// Initialize camera
	camera := NewCamera(player, worldBounds)
	
	// Create quadtree for collision detection
	quadtree := util.NewQuadtree(worldBounds)
	
	// Create food system
	foodSystem := systems.NewFoodGenerator(worldBounds)
	
	// Initialize with starting food
	foodSystem.Initialize(300)
	
	// Create power-up system
	powerUpSystem := systems.NewPowerUpManager(worldBounds)
	
	// Create minimap
	minimap := ui.NewMinimap(constants.WorldWidth, constants.WorldHeight)
	
	return &Game{
		Player:        player,
		Camera:        camera,
		Quadtree:      quadtree,
		FoodSystem:    foodSystem,
		PowerUpSystem: powerUpSystem,
		Minimap:       minimap,
		Score:         0,
		GameTime:      0,
		State:         constants.StateGame, // Start in game state
		FrameTime:     0,
		EntityCount:   0,
		DebugMode:     false,
	}
}

// Update updates the game state
func (g *Game) Update(dt float32) {
	// Update game time
	g.GameTime += dt
	
	// Store frame time for stats
	g.FrameTime = dt
	
	// Process different game states
	switch g.State {
	case constants.StateTitle:
		// Title screen logic
		if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) {
			g.State = constants.StateGame
		}
		
	case constants.StatePause:
		// Pause screen logic
		if rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeyEscape) {
			g.State = constants.StateGame
		}
		
	case constants.StateGameOver:
		// Game over screen logic
		if rl.IsKeyPressed(rl.KeyR) {
			g.Reset()
			g.State = constants.StateGame
		}
		
	case constants.StateGame:
		// Check for pause
		if rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeyEscape) {
			g.State = constants.StatePause
			return
		}
		
		// Check for toggling debug mode
		if rl.Is