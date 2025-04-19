// ui/controllers/pause_controller.go
package controllers

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// PauseController handles input for the pause screen
type PauseController struct {
    model        *models.PauseModel
    currentState *int
    resetGame    func()
}

// NewPauseController creates a new pause screen controller
func NewPauseController(model *models.PauseModel, currentState *int, resetGame func()) *PauseController {
    return &PauseController{
        model:        model,
        currentState: currentState,
        resetGame:    resetGame,
    }
}

// SetModel sets the controller's data model
func (c *PauseController) SetModel(model ui.Model) {
    c.model = model.(*models.PauseModel)
}

// HandleInput processes input for the pause screen
func (c *PauseController) HandleInput() bool {
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
        case "Resume":
            *c.currentState = constants.StateGame
            return true
            
        case "Restart":
            c.resetGame()
            *c.currentState = constants.StateGame
            return true
            
        case "Quit":
            rl.CloseWindow()
            return true
        }
    }
    
    // Escape also resumes
    if rl.IsKeyPressed(rl.KeyEscape) {
        *c.currentState = constants.StateGame
        return true
    }
    
    return false
}

// ui/controllers/game_over_controller.go
package controllers

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// GameOverController handles input for the game over screen
type GameOverController struct {
    model        *models.GameOverModel
    currentState *int
    resetGame    func()
}

// NewGameOverController creates a new game over screen controller
func NewGameOverController(model *models.GameOverModel, currentState *int, resetGame func()) *GameOverController {
    return &GameOverController{
        model:        model,
        currentState: currentState,
        resetGame:    resetGame,
    }
}

// SetModel sets the controller's data model
func (c *GameOverController) SetModel(model ui.Model) {
    c.model = model.(*models.GameOverModel)
}

// HandleInput processes input for the game over screen
func (c *GameOverController) HandleInput() bool {
    // Check for restart
    if rl.IsKeyPressed(rl.KeyR) {
        c.resetGame()
        *c.currentState = constants.StateGame
        return true
    }
    
    // Check for quit
    if rl.IsKeyPressed(rl.KeyQ) {
        rl.CloseWindow()
        return true
    }
    
    return false
}

// ui/controllers/boss_intro_controller.go
package controllers

import (
    "atomblaster/constants"
    "atomblaster/ui"
    "atomblaster/ui/models"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// BossIntroController handles input for the boss intro screen
type BossIntroController struct {
    model        *models.BossIntroModel
    currentState *int
}

// NewBossIntroController creates a new boss intro screen controller
func NewBossIntroController(model *models.BossIntroModel, currentState *int) *BossIntroController {
    return &BossIntroController{
        model:        model,
        currentState: currentState,
    }
}

// SetModel sets the controller's data model
func (c *BossIntroController) SetModel(model ui.Model) {
    c.model = model.(*models.BossIntroModel)
}

// HandleInput processes input for the boss intro screen
func (c *BossIntroController) HandleInput() bool {
    // Update the boss intro animation
    c.model.Update(rl.GetFrameTime())
    
    // Allow skipping the boss intro after a short delay
    if c.model.Timer > 1.0 && (rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter)) {
        *c.currentState = constants.StateGame
        return true
    }
    
    // Auto-progress after a certain time
    if c.model.Timer > 8.0 {
        *c.currentState = constants.StateGame
        return true
    }
    
    return false
}
