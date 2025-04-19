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
