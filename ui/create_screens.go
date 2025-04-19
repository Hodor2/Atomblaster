// ui/create-screens.go
package ui

import (
    "atomblaster/entities"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// CreateScreens initializes all UI screens
func CreateScreens(
    g interface{},
    background rl.Texture2D,
    playerSprite rl.Texture2D,
    enemySprite rl.Texture2D,
    bulletSprite rl.Texture2D,
    powerUpSprites [3]rl.Texture2D,
    player *entities.Player,
    atoms *[]entities.Atom,
    bullets *[]entities.Bullet,
    powerUps *[]entities.PowerUp,
    scientists *[]entities.Scientist,
    rescueZone *entities.RescueZone,
    particleSystem *entities.ParticleSystem,
    messageSystem *FloatingMessageSystem,
    door *rl.Rectangle,
    score *int,
    health *int,
    level *int,
    scientistsRescued *int,
    totalScientists *int,
    startTime *int64,
    elapsedTime *int64,
) (Screen, Screen, Screen, Screen, Screen) {
    // Create intro screen
    introScreen := NewIntroScreen(background, playerSprite)
    
    // Create title screen
    titleScreen := NewTitleScreen(background)

    // Create game screen with scientists and rescue zone
    gameScreen := NewGameScreen(
        background,
        playerSprite,
        enemySprite,
        bulletSprite,
        powerUpSprites,
        player,
        atoms,
        bullets,
        powerUps,
        scientists,
        rescueZone,
        particleSystem,
        messageSystem,
        door,
        score,
        health,
        level,
        scientistsRescued,
        totalScientists,
        startTime,
        elapsedTime,
    )

    // Create pause screen using the game screen as base
    pauseScreen := NewPauseScreen(gameScreen)

    // Create game over screen with all required parameters
    gameOverScreen := NewGameOverScreen(
        gameScreen,
        score,
        level,
        elapsedTime,
        scientistsRescued,
        totalScientists,
    )

    return introScreen, titleScreen, gameScreen, pauseScreen, gameOverScreen
}
