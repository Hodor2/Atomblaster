// ui/game_over_screen.go
package ui

import (
	"atomblaster/constants"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// GameOverScreen holds the state needed for rendering the game over screen
type GameOverScreen struct {
	GameScreen        Screen
	Score             *int
	Level             *int
	ElapsedTime       *int64
	ScientistsRescued *int
	TotalScientists   *int
}

// NewGameOverScreen creates a new game over screen instance
func NewGameOverScreen(
	gameScreen Screen,
	score, level *int,
	elapsedTime *int64,
	scientistsRescued, totalScientists *int,
) Screen {
	return &GameOverScreen{
		GameScreen:        gameScreen,
		Score:             score,
		Level:             level,
		ElapsedTime:       elapsedTime,
		ScientistsRescued: scientistsRescued,
		TotalScientists:   totalScientists,
	}
}

// Draw renders the game over screen
func (gos *GameOverScreen) Draw() {
	// Draw the game screen underneath
	gos.GameScreen.Draw()

	// Draw semi-transparent overlay
	rl.DrawRectangle(0, 0, constants.ScreenWidth, constants.ScreenHeight, rl.Fade(rl.Black, 0.7))

	// Draw game over text
	gameOverText := "GAME OVER"
	gameOverSize := int32(60)
	gameOverWidth := rl.MeasureText(gameOverText, gameOverSize)
	rl.DrawText(gameOverText, int32(constants.ScreenWidth)/2-gameOverWidth/2, 120, gameOverSize, rl.Red)

	centerX := constants.ScreenWidth / 2

	// Final Score
	scoreText := fmt.Sprintf("Final Score: %d", *gos.Score)
	scoreWidth := rl.MeasureText(scoreText, 30)
	rl.DrawText(scoreText, int32(centerX)-scoreWidth/2, 200, 30, rl.White)

	// Level reached
	levelText := fmt.Sprintf("Level Reached: %d/%d", *gos.Level, constants.MaxLevel)
	levelWidth := rl.MeasureText(levelText, 24)
	rl.DrawText(levelText, int32(centerX)-levelWidth/2, 240, 24, rl.White)

	// Time played
	timeText := fmt.Sprintf("Time: %02d:%02d", *gos.ElapsedTime/60, *gos.ElapsedTime%60)
	timeWidth := rl.MeasureText(timeText, 24)
	rl.DrawText(timeText, int32(centerX)-timeWidth/2, 270, 24, rl.White)

	// Scientists rescued
	if gos.ScientistsRescued != nil && gos.TotalScientists != nil {
		scienceText := fmt.Sprintf("Scientists Rescued: %d/%d", *gos.ScientistsRescued, *gos.TotalScientists)
		scienceWidth := rl.MeasureText(scienceText, 24)
		rl.DrawText(scienceText, int32(centerX)-scienceWidth/2, 300, 24, rl.Green)

		// Rescue Rate
		var percent int
		if *gos.TotalScientists > 0 {
			percent = (*gos.ScientistsRescued * 100) / *gos.TotalScientists
		}
		percentText := fmt.Sprintf("Rescue Rate: %d%%", percent)
		percentWidth := rl.MeasureText(percentText, 24)
		rl.DrawText(percentText, int32(centerX)-percentWidth/2, 330, 24, rl.Green)

		// Performance message
		var performanceText string
		switch {
		case percent == 100:
			performanceText = "Perfect rescue! All scientists saved!"
		case percent >= 75:
			performanceText = "Great job! Most scientists rescued!"
		case percent >= 50:
			performanceText = "Good work! Half the scientists saved."
		case percent >= 25:
			performanceText = "Some scientists rescued, try harder next time."
		default:
			performanceText = "Few scientists saved. Practice your flying!"
		}
		perfWidth := rl.MeasureText(performanceText, 20)
		rl.DrawText(performanceText, int32(centerX)-perfWidth/2, 360, 20, rl.Yellow)
	}

	// Restart instructions
	restartText := "Press R to Restart or Q to Quit"
	restartWidth := rl.MeasureText(restartText, 24)
	rl.DrawText(restartText, int32(centerX)-restartWidth/2, 420, 24, rl.White)
}

// Update handles input for the game over screen
func (gos *GameOverScreen) Update() bool {
	// This screen doesn't trigger transitions directly (handled in game state)
	return false
}
