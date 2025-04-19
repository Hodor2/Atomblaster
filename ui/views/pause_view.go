// ui/views/pause_view.go
package views

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// PauseView handles rendering the pause screen
type PauseView struct {
    model    *models.PauseModel
    gameView *GameView
}

// NewPauseView creates a new pause screen view
func NewPauseView(model *models.PauseModel, gameView *GameView) *PauseView {
    return &PauseView{
        model:    model,
        gameView: gameView,
    }
}

// SetModel sets the view's data model
func (v *PauseView) SetModel(model ui.Model) {
    v.model = model.(*models.PauseModel)
}

// Draw renders the pause screen
func (v *PauseView) Draw() {
    // First draw the game screen in the background
    v.gameView.Draw()
    
    // Then draw a semi-transparent overlay
    rl.DrawRectangle(
        0,
        0,
        constants.ScreenWidth,
        constants.ScreenHeight,
        rl.Fade(rl.Black, 0.7),
    )
    
    // Draw pause title
    pauseText := "GAME PAUSED"
    pauseWidth := rl.MeasureText(pauseText, 50)
    
    rl.DrawText(
        pauseText,
        int32(constants.ScreenWidth/2 - pauseWidth/2),
        100,
        50,
        rl.White,
    )
    
    // Draw menu options
    menuY := 250
    menuSpacing := 60
    
    for i, option := range v.model.MenuOptions {
        fontSize := 30
        if i == v.model.SelectedItem {
            fontSize = 35
            rl.DrawText(
                "► ",
                int32(constants.ScreenWidth/2 - 120),
                int32(menuY + i*menuSpacing),
                int32(fontSize),
                rl.Yellow,
            )
            
            rl.DrawText(
                option,
                int32(constants.ScreenWidth/2 - 80),
                int32(menuY + i*menuSpacing),
                int32(fontSize),
                rl.Yellow,
            )
        } else {
            rl.DrawText(
                option,
                int32(constants.ScreenWidth/2 - 80),
                int32(menuY + i*menuSpacing),
                int32(fontSize),
                rl.White,
            )
        }
    }
    
    // Draw controls
    controlsText := "Use Arrow Keys to Navigate, Enter to Select"
    controlsWidth := rl.MeasureText(controlsText, 20)
    
    rl.DrawText(
        controlsText,
        int32(constants.ScreenWidth/2 - controlsWidth/2),
        int32(constants.ScreenHeight - 80),
        20,
        rl.White,
    )
}

// ui/views/game_over_view.go
package views

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    "fmt"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// GameOverView handles rendering the game over screen
type GameOverView struct {
    model    *models.GameOverModel
    gameView *GameView
}

// NewGameOverView creates a new game over screen view
func NewGameOverView(model *models.GameOverModel, gameView *GameView) *GameOverView {
    return &GameOverView{
        model:    model,
        gameView: gameView,
    }
}

// SetModel sets the view's data model
func (v *GameOverView) SetModel(model ui.Model) {
    v.model = model.(*models.GameOverModel)
}

// Draw renders the game over screen
func (v *GameOverView) Draw() {
    // First draw the game screen in the background
    v.gameView.Draw()
    
    // Then draw a semi-transparent overlay
    rl.DrawRectangle(
        0,
        0,
        constants.ScreenWidth,
        constants.ScreenHeight,
        rl.Fade(rl.Black, 0.8),
    )
    
    // Draw game over title
    titleText := "GAME OVER"
    if v.model.PlayerWon {
        titleText = "MISSION COMPLETE!"
    }
    
    titleWidth := rl.MeasureText(titleText, 50)
    
    titleColor := rl.Red
    if v.model.PlayerWon {
        titleColor = rl.Green
    }
    
    rl.DrawText(
        titleText,
        int32(constants.ScreenWidth/2 - titleWidth/2),
        80,
        50,
        titleColor,
    )
    
    // Draw stats
    baseY := 180
    lineSpacing := 40
    
    // Final score
    scoreText := fmt.Sprintf("Final Score: %d", v.model.FinalScore)
    scoreWidth := rl.MeasureText(scoreText, 30)
    rl.DrawText(
        scoreText,
        int32(constants.ScreenWidth/2 - scoreWidth/2),
        int32(baseY),
        30,
        rl.White,
    )
    
    // Levels completed
    levelsText := fmt.Sprintf("Levels Completed: %d/%d", v.model.LevelsComplete, constants.MaxLevel)
    levelsWidth := rl.MeasureText(levelsText, 30)
    rl.DrawText(
        levelsText,
        int32(constants.ScreenWidth/2 - levelsWidth/2),
        int32(baseY + lineSpacing),
        30,
        rl.White,
    )
    
    // Scientists rescued
    scientistsText := fmt.Sprintf("Scientists Rescued: %d/%d", v.model.Scientists, v.model.TotalScientists)
    scientistsWidth := rl.MeasureText(scientistsText, 30)
    rl.DrawText(
        scientistsText,
        int32(constants.ScreenWidth/2 - scientistsWidth/2),
        int32(baseY + 2*lineSpacing),
        30,
        rl.White,
    )
    
    // Time elapsed
    minutes := v.model.TimeElapsed / 60
    seconds := v.model.TimeElapsed % 60
    timeText := fmt.Sprintf("Time: %02d:%02d", minutes, seconds)
    timeWidth := rl.MeasureText(timeText, 30)
    rl.DrawText(
        timeText,
        int32(constants.ScreenWidth/2 - timeWidth/2),
        int32(baseY + 3*lineSpacing),
        30,
        rl.White,
    )
    
    // Draw restart instruction
    restartText := "Press R to Restart, Q to Quit"
    restartWidth := rl.MeasureText(restartText, 25)
    rl.DrawText(
        restartText,
        int32(constants.ScreenWidth/2 - restartWidth/2),
        int32(constants.ScreenHeight - 80),
        25,
        rl.White,
    )
}
