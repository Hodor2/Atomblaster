// main.go
package main

import (
	"atomblaster/audio"
	"atomblaster/game"
	"atomblaster/constants"
	
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize window
	rl.InitWindow(constants.ScreenWidth, constants.ScreenHeight, "Atom Blaster")
	rl.SetTargetFPS(60)
	
	// Initialize audio
	rl.InitAudioDevice()
	audioSystem := audio.NewAudioSystem()
	
	// Initialize game
	gameState := game.NewGameState(audioSystem)
	
	// Main game loop
	for !rl.WindowShouldClose() {
		// Get frame time
		dt := rl.GetFrameTime()
		
		// Update game
		gameState.Update(dt)
		
		// Draw game
		gameState.Draw()
		
		// Update audio system
		audioSystem.Update()
	}
	
	// Cleanup
	rl.CloseAudioDevice()
	rl.CloseWindow()
}