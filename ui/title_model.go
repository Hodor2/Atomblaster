// ui/models/title_model.go
package models

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// TitleModel contains data for the title screen
type TitleModel struct {
    Background   rl.Texture2D
    MenuOptions  []string
    SelectedItem int
}

// NewTitleModel creates a new title screen model
func NewTitleModel(background rl.Texture2D) *TitleModel {
    return &TitleModel{
        Background:   background,
        MenuOptions:  []string{"Start Game", "Instructions", "Exit"},
        SelectedItem: 0,
    }
}

// SelectNextItem moves the selection to the next menu item
func (m *TitleModel) SelectNextItem() {
    m.SelectedItem = (m.SelectedItem + 1) % len(m.MenuOptions)
}

// SelectPreviousItem moves the selection to the previous menu item
func (m *TitleModel) SelectPreviousItem() {
    m.SelectedItem = (m.SelectedItem - 1 + len(m.MenuOptions)) % len(m.MenuOptions)
}

// GetSelectedOption returns the currently selected menu option
func (m *TitleModel) GetSelectedOption() string {
    return m.MenuOptions[m.SelectedItem]
}

// ui/models/intro_model.go
package models

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// IntroModel contains data for the intro screen
type IntroModel struct {
    Background   rl.Texture2D
    PlayerSprite rl.Texture2D
    Timer        float32
    Alpha        float32 // For fade effects
}

// NewIntroModel creates a new intro screen model
func NewIntroModel(background, playerSprite rl.Texture2D) *IntroModel {
    return &IntroModel{
        Background:   background,
        PlayerSprite: playerSprite,
        Timer:        0,
        Alpha:        0,
    }
}

// Update advances the intro animation
func (m *IntroModel) Update(dt float32) {
    m.Timer += dt
    
    // Fade in over the first 1 second
    if m.Timer < 1.0 {
        m.Alpha = m.Timer
    } else {
        m.Alpha = 1.0
    }
}

// ui/models/game_model.go
package models

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// GameModel contains data for the game screen
type GameModel struct {
    Background        rl.Texture2D
    PlayerSprite      rl.Texture2D
    EnemySprite       rl.Texture2D
    BulletSprite      rl.Texture2D
    PowerUpSprites    [3]rl.Texture2D
    Score             *int
    Health            *int
    Level             *int
    ScientistsRescued *int
    TotalScientists   *int
    StartTime         *int64
    ElapsedTime       *int64
}

// NewGameModel creates a new game screen model
func NewGameModel(
    background, playerSprite, enemySprite, bulletSprite rl.Texture2D,
    powerUpSprites [3]rl.Texture2D,
    score, health, level, scientistsRescued, totalScientists *int,
    startTime, elapsedTime *int64,
) *GameModel {
    return &GameModel{
        Background:        background,
        PlayerSprite:      playerSprite,
        EnemySprite:       enemySprite,
        BulletSprite:      bulletSprite,
        PowerUpSprites:    powerUpSprites,
        Score:             score,
        Health:            health,
        Level:             level,
        ScientistsRescued: scientistsRescued,
        TotalScientists:   totalScientists,
        StartTime:         startTime,
        ElapsedTime:       elapsedTime,
    }
}

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