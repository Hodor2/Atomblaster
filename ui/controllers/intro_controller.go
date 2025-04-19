// ui/controllers/intro_controller.go
package controllers

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// IntroController handles input for the intro screen
type IntroController struct {
    model        *models.IntroModel
    currentState *int
}

// NewIntroController creates a new intro screen controller
func NewIntroController(model *models.IntroModel, currentState *int) *IntroController {
    return &IntroController{
        model:        model,
        currentState: currentState,
    }
}

// SetModel sets the controller's data model
func (c *IntroController) SetModel(model ui.Model) {
    c.model = model.(*models.IntroModel)
}

// HandleInput processes input for the intro screen
func (c *IntroController) HandleInput() bool {
    // Update the intro animation
    c.model.Update(rl.GetFrameTime())
    
    // Allow skipping the intro after a short delay
    if c.model.Timer > 1.0 && (rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter)) {
        *c.currentState = constants.StateTitle
        return true
    }
    
    // Auto-progress after a certain time
    if c.model.Timer > 10.0 {
        *c.currentState = constants.StateTitle
        return true
    }
    
    return false
}

// ui/controllers/game_controller.go
package controllers

import (
    "atomblaster/ui"
    "atomblaster/ui/models"
)

// GameController handles input for the main game screen
type GameController struct {
    model *models.GameModel
}

// NewGameController creates a new game screen controller
func NewGameController(model *models.GameModel) *GameController {
    return &GameController{
        model: model,
    }
}

// SetModel sets the controller's data model
func (c *GameController) SetModel(model ui.Model) {
    c.model = model.(*models.GameModel)
}

// HandleInput processes input for the game screen
// The actual game input is handled by the InputSystem in ECS
func (c *GameController) HandleInput() bool {
    // Nothing to do here, as input is handled by the InputSystem
    return false
}
