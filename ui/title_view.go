// ui/views/title_view.go
package views

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// TitleView handles rendering the title screen
type TitleView struct {
    model *models.TitleModel
}

// NewTitleView creates a new title screen view
func NewTitleView(model *models.TitleModel) *TitleView {
    return &TitleView{
        model: model,
    }
}

// SetModel sets the view's data model
func (v *TitleView) SetModel(model ui.Model) {
    v.model = model.(*models.TitleModel)
}

// Draw renders the title screen
func (v *TitleView) Draw() {
    // Draw background
    rl.DrawTexture(v.model.Background, 0, 0, rl.White)
    
    // Draw title
    titleText := "ATOM BLASTER"
    titleFontSize := 60
    titleWidth := rl.MeasureText(titleText, int32(titleFontSize))
    
    rl.DrawText(
        titleText,
        int32(constants.ScreenWidth/2 - titleWidth/2),
        100,
        int32(titleFontSize),
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
    
    // Draw footer text
    footerText := "© 2023 ATOM BLASTER TEAM"
    footerWidth := rl.MeasureText(footerText, 20)
    
    rl.DrawText(
        footerText,
        int32(constants.ScreenWidth/2 - footerWidth/2),
        int32(constants.ScreenHeight - 50),
        20,
        rl.LightGray,
    )
    
    // Draw instruction
    instructionText := "Use Arrow Keys to Navigate, Enter to Select"
    instructionWidth := rl.MeasureText(instructionText, 20)
    
    rl.DrawText(
        instructionText,
        int32(constants.ScreenWidth/2 - instructionWidth/2),
        int32(constants.ScreenHeight - 80),
        20,
        rl.White,
    )
}

// ui/views/intro_view.go
package views

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// IntroView handles rendering the intro screen
type IntroView struct {
    model *models.IntroModel
}

// NewIntroView creates a new intro screen view
func NewIntroView(model *models.IntroModel) *IntroView {
    return &IntroView{
        model: model,
    }
}

// SetModel sets the view's data model
func (v *IntroView) SetModel(model ui.Model) {
    v.model = model.(*models.IntroModel)
}

// Draw renders the intro screen
func (v *IntroView) Draw() {
    // Draw background with alpha fade
    rl.DrawTexture(v.model.Background, 0, 0, rl.Fade(rl.White, v.model.Alpha))
    
    // Draw story text with fade effect
    storyText := []string{
        "In a world ravaged by atomic disasters...",
        "Scientists are trapped in dangerous radiation zones.",
        "You're their only hope for rescue.",
        "Pilot your helicopter, avoid radioactive atoms,",
        "and bring the scientists to safety!",
        "",
        "Press SPACE to continue..."
    }
    
    baseY := 150
    
    for i, line := range storyText {
        textWidth := rl.MeasureText(line, 30)
        rl.DrawText(
            line,
            int32(constants.ScreenWidth/2 - textWidth/2),
            int32(baseY + i*40),
            30,
            rl.Fade(rl.White, v.model.Alpha),
        )
    }
    
    // Draw player sprite
    if v.model.Timer > 1.5 {
        rl.DrawTexture(
            v.model.PlayerSprite,
            int32(constants.ScreenWidth/2 - v.model.PlayerSprite.Width/2),
            int32(constants.ScreenHeight - 200),
            rl.White,
        )
    }
}

// ui/views/game_view.go
package views

import (
    "atomblaster/ui"
    "atomblaster/ui/models"
    "fmt"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// GameView handles rendering the game screen
type GameView struct {
    model *models.GameModel
}

// NewGameView creates a new game screen view
func NewGameView(model *models.GameModel) *GameView {
    return &GameView{
        model: model,
    }
}

// SetModel sets the view's data model
func (v *GameView) SetModel(model ui.Model) {
    v.model = model.(*models.GameModel)
}

// Draw renders the game screen
func (v *GameView) Draw() {
    // Note: The actual game entities are drawn by the RenderSystem
    // This view only draws UI overlays like score, health, etc.
    
    // Draw UI elements
    rl.DrawText(
        fmt.Sprintf("SCORE: %d", *v.model.Score),
        10,
        10,
        20,
        rl.White,
    )
    
    rl.DrawText(
        fmt.Sprintf("HEALTH: %d", *v.model.Health),
        10,
        40,
        20,
        rl.White,
    )
    
    rl.DrawText(
        fmt.Sprintf("LEVEL: %d", *v.model.Level),
        10,
        70,
        20,
        rl.White,
    )
    
    // Draw scientists rescued counter
    rl.DrawText(
        fmt.Sprintf("SCIENTISTS: %d/%d", *v.model.ScientistsRescued, *v.model.TotalScientists),
        10,
        100,
        20,
        rl.White,
    )
    
    // Calculate and display elapsed time
    minutes := *v.model.ElapsedTime / 60
    seconds := *v.model.ElapsedTime % 60
    
    rl.DrawText(
        fmt.Sprintf("TIME: %02d:%02d", minutes, seconds),
        10,
        130,
        20,
        rl.White,
    )
}

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

// ui/views/boss_intro_view.go
package views

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// BossIntroView handles rendering the boss introduction screen
type BossIntroView struct {
    model *models.BossIntroModel
}

// NewBossIntroView creates a new boss intro screen view
func NewBossIntroView(model *models.BossIntroModel) *BossIntroView {
    return &BossIntroView{
        model: model,
    }
}

// SetModel sets the view's data model
func (v *BossIntroView) SetModel(model ui.Model) {
    v.model = model.(*models.BossIntroModel)
}

// Draw renders the boss intro screen
func (v *BossIntroView) Draw() {
    // Draw background with alpha fade
    rl.DrawTexture(v.model.Background, 0, 0, rl.Fade(rl.White, v.model.Alpha))
    
    // Draw warning text
    warningText := "WARNING: BOSS APPROACHING!"
    warningWidth := rl.MeasureText(warningText, 50)
    
    // Make text blink if enough time has passed
    textColor := rl.Red
    if v.model.Timer > 0.5 && int(v.model.Timer*5)%2 == 0 {
        textColor = rl.Yellow
    }
    
    rl.DrawText(
        warningText,
        int32(constants.ScreenWidth/2 - warningWidth/2),
        100,
        50,
        rl.Fade(textColor, v.model.Alpha),
    )
    
    // Draw boss description
    if v.model.Timer > 1.0 {
        descriptionText := []string{
            "A hostile combat helicopter has been detected!",
            "It's heavily armed and extremely dangerous.",
            "Defeat it to complete the level and advance.",
            "",
            "Good luck, pilot. You'll need it.",
            "",
            "Press SPACE to continue..."
        }
        
        baseY := 200
        
        for i, line := range descriptionText {
            textWidth := rl.MeasureText(line, 25)
            rl.DrawText(
                line,
                int32(constants.ScreenWidth/2 - textWidth/2),
                int32(baseY + i*40),
                25,
                rl.Fade(rl.White, v.model.Alpha),
            )
        }
    }
    
    // Draw player and boss sprites
    if v.model.Timer > 2.0 {
        // Player helicopter on the left
        rl.DrawTexture(
            v.model.PlayerSprite,
            150,
            int32(constants.ScreenHeight - 200),
            rl.Fade(rl.White, v.model.Alpha),
        )
        
        // Boss helicopter on the right
        rl.DrawTexture(
            v.model.BossSprite,
            int32(constants.ScreenWidth - 250),
            int32(constants.ScreenHeight - 250),
            rl.Fade(rl.White, v.model.Alpha),
        )
        
        // Draw "VS" text in the middle
        vsText := "VS"
        vsWidth := rl.MeasureText(vsText, 70)
        rl.DrawText(
            vsText,
            int32(constants.ScreenWidth/2 - vsWidth/2),
            int32(constants.ScreenHeight - 220),
            70,
            rl.Fade(rl.Red, v.model.Alpha),
        )
    }
}