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
