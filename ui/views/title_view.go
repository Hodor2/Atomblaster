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
