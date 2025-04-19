// ui/game_screen.go
package ui

import (
	"fmt"

	"atomblaster/entities"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GameScreen holds the state needed for rendering the main game.
type GameScreen struct {
	Background     rl.Texture2D
	PlayerSprite   rl.Texture2D
	EnemySprite    rl.Texture2D
	BulletSprite   rl.Texture2D
	PowerUpSprites [3]rl.Texture2D

	Player         *entities.Player
	Atoms          *[]entities.Atom
	Bullets        *[]entities.Bullet
	PowerUps       *[]entities.PowerUp
	Scientists     *[]entities.Scientist
	RescueZone     *entities.RescueZone
	ParticleSystem *entities.ParticleSystem
	MessageSystem  *FloatingMessageSystem
	Door           *rl.Rectangle

	Score             *int
	Health            *int
	Level             *int
	StartTime         *int64
	ElapsedTime       *int64
	ScientistsRescued *int
	TotalScientists   *int
}

// NewGameScreen creates and returns a Screen for the main game.
func NewGameScreen(
	bg, playerSprite, enemySprite, bulletSprite rl.Texture2D,
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
	score, health, level *int,
	scientistsRescued, totalScientists *int,
	startTime, elapsedTime *int64,
) Screen {
	return &GameScreen{
		Background:        bg,
		PlayerSprite:      playerSprite,
		EnemySprite:       enemySprite,
		BulletSprite:      bulletSprite,
		PowerUpSprites:    powerUpSprites,
		Player:            player,
		Atoms:             atoms,
		Bullets:           bullets,
		PowerUps:          powerUps,
		Scientists:        scientists,
		RescueZone:        rescueZone,
		ParticleSystem:    particleSystem,
		MessageSystem:     messageSystem,
		Door:              door,
		Score:             score,
		Health:            health,
		Level:             level,
		StartTime:         startTime,
		ElapsedTime:       elapsedTime,
		ScientistsRescued: scientistsRescued,
		TotalScientists:   totalScientists,
	}
}

// Draw renders the game screen. It satisfies the ui.Screen interface.
func (s *GameScreen) Draw() {
	// Draw background
	rl.DrawTexture(s.Background, 0, 0, rl.White)

	// Draw rescue zone
	s.RescueZone.Draw()
	
	// Draw door
	doorColor := rl.Red
	if *s.ScientistsRescued >= *s.TotalScientists {
		doorColor = rl.Green // Door opens when all scientists are rescued
	}
	rl.DrawRectangleRec(*s.Door, doorColor)

	// Draw player
	s.Player.Draw(s.PlayerSprite)

	// Draw atoms
	for _, atom := range *s.Atoms {
		atom.Draw(s.EnemySprite)
	}

	// Draw bullets
	for _, bullet := range *s.Bullets {
		bullet.Draw(s.BulletSprite)
	}

	// Draw power‑ups (pass the full slice, not a single texture)
	for _, pu := range *s.PowerUps {
		pu.Draw(s.PowerUpSprites)
	}

	// Draw scientists
	for _, sci := range *s.Scientists {
		sci.Draw()
	}

	// Draw particles and floating messages
	s.ParticleSystem.Draw()
	s.MessageSystem.Draw()

	// Draw HUD: score, health, level, timer
	hudY := int32(10)
	rl.DrawText(fmt.Sprintf("Score: %d", *s.Score), 10, hudY, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Health: %d", *s.Health), 200, hudY, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Level: %d", *s.Level), 400, hudY, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Time: %d", *s.ElapsedTime), 600, hudY, 20, rl.White)
	
	// Add scientists rescued counter
	rl.DrawText(fmt.Sprintf("Scientists: %d/%d", *s.ScientistsRescued, *s.TotalScientists), 10, hudY+30, 20, rl.Green)
}

// Update handles input for the game screen - nothing to do here as game input
// is handled by the game state
func (s *GameScreen) Update() bool {
	// This screen doesn't trigger transitions directly
	return false
}
