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
            rl.White
        )
    }
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
        rl.Fade(textColor, v.model.Alpha)
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
                rl.Fade(rl.White, v.model.Alpha)
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
            rl.Fade(rl.White, v.model.Alpha)
        )
        
        // Boss helicopter on the right
        rl.DrawTexture(
            v.model.BossSprite,
            int32(constants.ScreenWidth - 250),
            int32(constants.ScreenHeight - 250),
            rl.Fade(rl.White, v.model.Alpha)
        )
        
        // Draw "VS" text in the middle
        vsText := "VS"
        vsWidth := rl.MeasureText(vsText, 70)
        rl.DrawText(
            vsText,
            int32(constants.ScreenWidth/2 - vsWidth/2),
            int32(constants.ScreenHeight - 220),
            70,
            rl.Fade(rl.Red, v.model.Alpha)
        )
    }
}
