// ui/draw.go
package ui

import (
	"atomblaster/constants"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawGameState renders the current game state
func DrawGameState(currentState int, screens map[int]Screen, messageSystem *FloatingMessageSystem) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	// Draw the appropriate screen based on game state
	if screen, ok := screens[currentState]; ok {
		screen.Draw()
	}

	// Draw floating messages in game state
	if currentState == constants.StateGame {
		messageSystem.Draw()
	}

	rl.EndDrawing()
}