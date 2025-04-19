package ui

import (
	"atomblaster/constants"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// PauseScreen holds the state needed for rendering the pause screen
type PauseScreen struct {
	GameScreen Screen
}

// NewPauseScreen creates a new pause screen instance
func NewPauseScreen(gameScreen Screen) Screen {
	return &PauseScreen{
		GameScreen: gameScreen,
	}
}

// Draw renders the pause screen
func (ps *PauseScreen) Draw() {
	// Draw the game screen underneath
	ps.GameScreen.Draw()

	// Draw semi-transparent overlay
	rl.DrawRectangle(0, 0, constants.ScreenWidth, constants.ScreenHeight, rl.Fade(rl.Black, 0.5))

	// Draw pause text
	pauseText := "GAME PAUSED"
	pauseSize := 50
	pauseWidth := rl.MeasureText(pauseText, int32(pauseSize))
	rl.DrawText(pauseText, int32(constants.ScreenWidth/2-int(pauseWidth)/2), 150, int32(pauseSize), rl.White)

	// Draw instructions
	resumeText := "Press ESC to Resume"
	resumeWidth := rl.MeasureText(resumeText, 24)
	rl.DrawText(resumeText, int32(constants.ScreenWidth/2-int(resumeWidth)/2), 250, 24, rl.Yellow)

	// Show controls reminder
	rl.DrawText("Helicopter Controls:", constants.ScreenWidth/2-150, 300, 20, rl.White)
	rl.DrawText("- WASD or Arrow Keys to fly", constants.ScreenWidth/2-150, 330, 18, rl.LightGray)
	rl.DrawText("- Left Click to fire helicopter weapons", constants.ScreenWidth/2-150, 355, 18, rl.LightGray)
	rl.DrawText("- SPACE for helicopter boost", constants.ScreenWidth/2-150, 380, 18, rl.LightGray)
	rl.DrawText("- Collect orange powerups for weapons", constants.ScreenWidth/2-150, 405, 18, rl.LightGray)
}

// Update handles input for the pause screen
func (ps *PauseScreen) Update() bool {
	// This screen doesn't trigger transitions directly (handled in game state)
	return false
}
