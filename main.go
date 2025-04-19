package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/game"
	"math/rand"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Initialize window
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagMsaa4xHint)
	rl.InitWindow(constants.ScreenWidth, constants.ScreenHeight, "Atomblaster 2.5D")
	rl.SetTargetFPS(60)
	
	// Initialize audio device
	rl.InitAudioDevice()
	
	// Create game state
	gameState := game.New()
	
	// Main game loop
	for !rl.WindowShouldClose() {
		// Calculate delta time
		dt := rl.GetFrameTime()
		
		// Update game state
		gameState.Update(dt)
		
		// Draw game
		gameState.Draw()
	}
	
	// Clean up
	gameState.Cleanup()
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
