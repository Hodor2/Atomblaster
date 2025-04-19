// ui/controllers/title_controller.go
package controllers

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// TitleController handles input for the title screen
type TitleController struct {
    model        *models.TitleModel
    currentState *int
}

// NewTitleController creates a new title screen controller
func NewTitleController(model *models.TitleModel, currentState *int) *TitleController {
    return &TitleController{
        model:        model,
        currentState: currentState,
    }
}

// SetModel sets the controller's data model
func (c *TitleController) SetModel(model ui.Model) {
    c.model = model.(*models.TitleModel)
}

// HandleInput processes input for the title screen
func (c *TitleController) HandleInput() bool {
    // Handle menu navigation
    if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
        c.model.SelectPreviousItem()
    }
    
    if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
        c.model.SelectNextItem()
    }
    
    // Handle menu selection
    if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
        selectedOption := c.model.GetSelectedOption()
        
        switch selectedOption {
        case "Start Game":
            *c.currentState = constants.StateGame
            return true
            
        case "Instructions":
            // Could show instructions screen
            // For simplicity, we'll just start the game
            *c.currentState = constants.StateIntro
            return true
            
        case "Exit":
            rl.CloseWindow()
            return true
        }
    }
    
    return false
}
