package ui

import (
	"math"
	"atomblaster/constants"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TitleScreen holds any state or assets needed for the title screen
type TitleScreen struct {
	Background rl.Texture2D
	timer      float32
}

// NewTitleScreen creates a new title screen instance
func NewTitleScreen(bg rl.Texture2D) Screen {
	return &TitleScreen{
		Background: bg,
		timer: 0,
	}
}

// Draw renders the title screen
func (ts *TitleScreen) Draw() {
	// Update timer for animations
	ts.timer += rl.GetFrameTime()
	
	// Draw background
	if ts.Background.ID > 0 {
		// Draw darkened background
		rl.DrawTexturePro(
			ts.Background,
			rl.Rectangle{X: 0, Y: 0, Width: float32(ts.Background.Width), Height: float32(ts.Background.Height)},
			rl.Rectangle{X: 0, Y: 0, Width: float32(constants.ScreenWidth), Height: float32(constants.ScreenHeight)},
			rl.Vector2{X: 0, Y: 0},
			0.0,
			rl.Fade(rl.White, 0.5),
		)
	}

	// Draw title
	titleText := "ATOM BLASTER"
	titleSize := 50
	titleWidth := rl.MeasureText(titleText, int32(titleSize))

	// Pulse effect
	pulse := 0.8 + 0.2*float32(math.Sin(float64(ts.timer * 2.0)))
	
	// Shadow effect
	rl.DrawText(titleText, int32(constants.ScreenWidth/2-int(titleWidth)/2)+4, 104, int32(titleSize), rl.Black)
	rl.DrawText(titleText, int32(constants.ScreenWidth/2-int(titleWidth)/2), 100, int32(titleSize), rl.ColorAlpha(rl.Red, pulse))

	// Subtitle
	subtitleText := "A 2D Roguelike Shooter"
	subtitleSize := 30
	subtitleWidth := rl.MeasureText(subtitleText, int32(subtitleSize))
	rl.DrawText(subtitleText, int32(constants.ScreenWidth/2-int(subtitleWidth)/2), 160, int32(subtitleSize), rl.White)

	// Pulsing "Press Enter to Start" message
	startText := "Press ENTER to Start"
	startSize := 24
	startWidth := rl.MeasureText(startText, int32(startSize))

	// Pulsing effect - blink on and off
	if int(ts.timer*2.0) % 2 == 0 {
		rl.DrawText(startText, int32(constants.ScreenWidth/2-int(startWidth)/2), 300, int32(startSize), rl.Yellow)
	}

	// Game controls explanation
	rl.DrawText("Controls:", 100, 380, 20, rl.White)
	rl.DrawText("- WASD or Arrow Keys to move", 120, 410, 18, rl.LightGray)
	rl.DrawText("- Left Click to shoot (once you have a gun)", 120, 435, 18, rl.LightGray)
	rl.DrawText("- SPACE to dash (with cooldown)", 120, 460, 18, rl.LightGray)
	rl.DrawText("- ESC to pause", 120, 485, 18, rl.LightGray)

	// Game objective
	rl.DrawText("Objective:", constants.ScreenWidth-320, 380, 20, rl.White)
	rl.DrawText("- Clear all enemies", constants.ScreenWidth-300, 410, 18, rl.LightGray)
	rl.DrawText("- Reach the exit door", constants.ScreenWidth-300, 435, 18, rl.LightGray)
	rl.DrawText("- Collect power-ups", constants.ScreenWidth-300, 460, 18, rl.LightGray)
	rl.DrawText("- Survive all 10 levels", constants.ScreenWidth-300, 485, 18, rl.LightGray)

	// Version info
	rl.DrawText("Version 2.0", 20, constants.ScreenHeight-30, 16, rl.DarkGray)
}

// Update processes user input for the title screen
func (ts *TitleScreen) Update() bool {
	// Return true to transition to the next screen when ENTER is pressed
	return rl.IsKeyPressed(rl.KeyEnter)
}
