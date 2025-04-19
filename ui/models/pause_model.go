// ui/models/pause_model.go
package models

// PauseModel contains data for the pause screen
type PauseModel struct {
    GameModel    *GameModel
    MenuOptions  []string
    SelectedItem int
}

// NewPauseModel creates a new pause screen model
func NewPauseModel(gameModel *GameModel) *PauseModel {
    return &PauseModel{
        GameModel:    gameModel,
        MenuOptions:  []string{"Resume", "Restart", "Quit"},
        SelectedItem: 0,
    }
}

// SelectNextItem moves the selection to the next menu item
func (m *PauseModel) SelectNextItem() {
    m.SelectedItem = (m.SelectedItem + 1) % len(m.MenuOptions)
}

// SelectPreviousItem moves the selection to the previous menu item
func (m *PauseModel) SelectPreviousItem() {
    m.SelectedItem = (m.SelectedItem - 1 + len(m.MenuOptions)) % len(m.MenuOptions)
}

// GetSelectedOption returns the currently selected menu option
func (m *PauseModel) GetSelectedOption() string {
    return m.MenuOptions[m.SelectedItem]
}

// ui/models/game_over_model.go
package models

// GameOverModel contains data for the game over screen
type GameOverModel struct {
    GameModel      *GameModel
    PlayerWon      bool
    FinalScore     int
    LevelsComplete int
    TimeElapsed    int64
    Scientists     int
    TotalScientists int
}

// NewGameOverModel creates a new game over screen model
func NewGameOverModel(
    gameModel *GameModel,
    playerWon bool,
) *GameOverModel {
    return &GameOverModel{
        GameModel:      gameModel,
        PlayerWon:      playerWon,
        FinalScore:     *gameModel.Score,
        LevelsComplete: *gameModel.Level,
        TimeElapsed:    *gameModel.ElapsedTime,
        Scientists:     *gameModel.ScientistsRescued,
        TotalScientists: *gameModel.TotalScientists,
    }
}

// ui/models/boss_intro_model.go
package models

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// BossIntroModel contains data for the boss introduction screen
type BossIntroModel struct {
    Background   rl.Texture2D
    PlayerSprite rl.Texture2D
    BossSprite   rl.Texture2D
    Timer        float32
    Alpha        float32
}

// NewBossIntroModel creates a new boss intro screen model
func NewBossIntroModel(background, playerSprite, bossSprite rl.Texture2D) *BossIntroModel {
    return &BossIntroModel{
        Background:   background,
        PlayerSprite: playerSprite,
        BossSprite:   bossSprite,
        Timer:        0,
        Alpha:        0,
    }
}

// Update advances the boss intro animation
func (m *BossIntroModel) Update(dt float32) {
    m.Timer += dt
    
    // Fade in over the first 1 second
    if m.Timer < 1.0 {
        m.Alpha = m.Timer
    } else {
        m.Alpha = 1.0
    }
}
