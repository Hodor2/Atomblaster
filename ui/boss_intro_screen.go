// Package ui contains UI components for Atom Blaster
package ui

import (
	"atomblaster/constants"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// BossIntroScreen shows an animation of enemy helicopter killing scientists
type BossIntroScreen struct {
	background    rl.Texture2D
	bossTexture   rl.Texture2D
	playerTexture rl.Texture2D
	scientists    []ScientistAnimInfo
	bossPos       rl.Vector2
	playerPos     rl.Vector2
	animTimer     float32
	animDuration  float32
	showPrompt    bool
	promptTimer   float32
	explosions    []ExplosionInfo
	textAlpha     float32
}

// ScientistAnimInfo holds animation state for a scientist in the intro
type ScientistAnimInfo struct {
	Pos        rl.Vector2
	Dead       bool
	DeathTimer float32
	AnimTimer  float32
}

// ExplosionInfo holds explosion animation data
type ExplosionInfo struct {
	Pos      rl.Vector2
	Size     float32
	Lifetime float32
	MaxSize  float32
}

// NewBossIntroScreen creates a new boss intro screen
func NewBossIntroScreen(
	background rl.Texture2D,
	bossTexture rl.Texture2D,
	playerTexture rl.Texture2D,
) Screen {
	// Create scientists in a group near center
	scientists := make([]ScientistAnimInfo, 8)
	centerX := float32(constants.ScreenWidth / 2)
	centerY := float32(constants.ScreenHeight / 2)

	for i := 0; i < 8; i++ {
		// Position in a rough circle around center
		angle := float32(i) * (2 * math.Pi / 8)
		radius := float32(80)
		x := centerX + float32(math.Cos(float64(angle)))*radius
		y := centerY + float32(math.Sin(float64(angle)))*radius

		scientists[i] = ScientistAnimInfo{
			Pos:        rl.Vector2{X: x, Y: y},
			Dead:       false,
			DeathTimer: 0,
			AnimTimer:  float32(i) * 0.2, // Offset animation timing
		}
	}

	// Start boss helicopter off-screen at the top
	bossStartPos := rl.Vector2{
		X: centerX,
		Y: -50,
	}

	// Start player helicopter off-screen to the left
	playerStartPos := rl.Vector2{
		X: -50,
		Y: centerY + 150,
	}

	return &BossIntroScreen{
		background:    background,
		bossTexture:   bossTexture,
		playerTexture: playerTexture,
		scientists:    scientists,
		bossPos:       bossStartPos,
		playerPos:     playerStartPos,
		animTimer:     0,
		animDuration:  8.0, // 8 seconds animation
		showPrompt:    false,
		promptTimer:   0,
		explosions:    []ExplosionInfo{},
		textAlpha:     0,
	}
}

// Draw renders the boss intro animation
func (s *BossIntroScreen) Draw() {
	// Background
	rl.DrawTexture(s.background, 0, 0, rl.White)

	// Update animation
	dt := rl.GetFrameTime()
	s.animTimer += dt
	s.promptTimer += dt

	// Handle animation phases
	animProgress := s.animTimer / s.animDuration

	// Show darkened background with red tint
	rl.DrawRectangle(0, 0, constants.ScreenWidth, constants.ScreenHeight, rl.Fade(rl.Maroon, 0.2))

	// Boss helicopter movement
	if animProgress < 0.3 {
		// Move in from top
		s.bossPos.Y = -50 + (animProgress/0.3)*250
	} else if animProgress < 0.7 {
		// Circle around scientists
		circleTime := (animProgress - 0.3) / 0.4
		circleAngle := circleTime * 2 * math.Pi
		s.bossPos.X = float32(constants.ScreenWidth/2 + int32(120*math.Cos(float64(circleAngle))))
		s.bossPos.Y = float32(constants.ScreenHeight/2 + int32(80*math.Sin(float64(circleAngle))))

		// Kill scientists as boss circles
		scientistsKilled := int(circleTime * 8)
		for i := 0; i < scientistsKilled && i < 8; i++ {
			if !s.scientists[i].Dead && s.scientists[i].DeathTimer == 0 {
				s.scientists[i].DeathTimer = 0.01 // Start death animation

				// Add explosion
				s.explosions = append(s.explosions, ExplosionInfo{
					Pos:      s.scientists[i].Pos,
					Size:     0,
					Lifetime: 0.5,
					MaxSize:  30,
				})
			}
		}
	} else {
		// Move to hover at right side
		s.bossPos.X = float32(constants.ScreenWidth/2) + (animProgress-0.7)/0.3*300
		s.bossPos.Y = 200
	}

	// Player helicopter movement (entering after scientists die)
	if animProgress > 0.8 {
		playerEntryProgress := (animProgress - 0.8) / 0.2
		s.playerPos.X = -50 + playerEntryProgress*float32(constants.ScreenWidth/2-100)
		s.playerPos.Y = float32(constants.ScreenHeight/2 + 150)
	}

	// Update and draw explosions
	updatedExplosions := []ExplosionInfo{}
	for _, explosion := range s.explosions {
		explosion.Lifetime -= dt
		if explosion.Lifetime > 0 {
			explosion.Size = explosion.MaxSize * (1 - explosion.Lifetime/0.5)

			// Draw explosion
			rl.DrawCircleV(explosion.Pos, explosion.Size, rl.Orange)
			rl.DrawCircleV(explosion.Pos, explosion.Size*0.7, rl.Yellow)

			updatedExplosions = append(updatedExplosions, explosion)
		}
	}
	s.explosions = updatedExplosions

	// Draw scientists
	for i := range s.scientists {
		s.scientists[i].AnimTimer += dt

		// Update death animation
		if s.scientists[i].DeathTimer > 0 {
			s.scientists[i].DeathTimer += dt
			if s.scientists[i].DeathTimer >= 0.5 {
				s.scientists[i].Dead = true
			}
		}

		// Skip if fully dead
		if s.scientists[i].Dead {
			continue
		}

		// Get position
		posX := int32(s.scientists[i].Pos.X)
		posY := int32(s.scientists[i].Pos.Y)

		// Determine color based on death state
		color := rl.White
		if s.scientists[i].DeathTimer > 0 {
			// Fade to red and diminish
			deathProgress := s.scientists[i].DeathTimer / 0.5
			alpha := 1.0 - deathProgress
			color = rl.ColorAlpha(rl.Red, alpha)
		}

		// Draw scientist stick figure
		rl.DrawCircle(posX, posY-15, 10, color)         // Head
		rl.DrawLine(posX, posY-5, posX, posY+15, color) // Body

		// Arms and legs
		if s.scientists[i].DeathTimer == 0 {
			// Regular animation if alive
			armWaveOffset := int32(5 * math.Sin(float64(s.scientists[i].AnimTimer*10)))

			// Legs
			rl.DrawLine(posX, posY+15, posX-8, posY+30, color)
			rl.DrawLine(posX, posY+15, posX+8, posY+30, color)

			// Arms - some waving in panic
			rl.DrawLine(posX, posY, posX-10, posY+5, color)                // Left arm
			rl.DrawLine(posX, posY, posX+10, posY-10-armWaveOffset, color) // Right arm waving
		} else {
			// Death animation - limbs splayed
			deathProgress := s.scientists[i].DeathTimer / 0.5

			// Legs splayed
			legAngle1 := math.Pi/4 + float64(deathProgress)*math.Pi/4
			legAngle2 := -math.Pi/4 - float64(deathProgress)*math.Pi/4
			legLength := 15.0 * (1.0 - float64(deathProgress)*0.5)

			rl.DrawLine(
				posX,
				posY+15,
				posX+int32(math.Cos(legAngle1)*legLength),
				posY+15+int32(math.Sin(legAngle1)*legLength),
				color,
			)
			rl.DrawLine(
				posX,
				posY+15,
				posX+int32(math.Cos(legAngle2)*legLength),
				posY+15+int32(math.Sin(legAngle2)*legLength),
				color,
			)

			// Arms splayed
			armAngle1 := -math.Pi/4 - float64(deathProgress)*math.Pi/4
			armAngle2 := -3*math.Pi/4 + float64(deathProgress)*math.Pi/4
			armLength := 10.0 * (1.0 - float64(deathProgress)*0.5)

			rl.DrawLine(
				posX,
				posY,
				posX+int32(math.Cos(armAngle1)*armLength),
				posY+int32(math.Sin(armAngle1)*armLength),
				color,
			)
			rl.DrawLine(
				posX,
				posY,
				posX+int32(math.Cos(armAngle2)*armLength),
				posY+int32(math.Sin(armAngle2)*armLength),
				color,
			)
		}
	}

	// Draw boss helicopter
	if s.bossTexture.ID > 0 {
		rl.DrawTexturePro(
			s.bossTexture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(s.bossTexture.Width), Height: float32(s.bossTexture.Height)},
			rl.Rectangle{
				X:      s.bossPos.X - 40,
				Y:      s.bossPos.Y - 20,
				Width:  80,
				Height: 40,
			},
			rl.Vector2{X: 40, Y: 20},
			0.0,
			rl.Red,
		)
	} else {
		// Fallback to drawing boss helicopter with primitives
		drawBossHelicopter(s.bossPos, float32(s.animTimer*10.0))
	}

	// Draw player helicopter if it's entered the scene
	if animProgress > 0.8 {
		if s.playerTexture.ID > 0 {
			rl.DrawTexturePro(
				s.playerTexture,
				rl.Rectangle{X: 0, Y: 0, Width: float32(s.playerTexture.Width), Height: float32(s.playerTexture.Height)},
				rl.Rectangle{
					X:      s.playerPos.X - float32(constants.HelicopterWidth)/2,
					Y:      s.playerPos.Y - float32(constants.HelicopterHeight)/2,
					Width:  float32(constants.HelicopterWidth),
					Height: float32(constants.HelicopterHeight),
				},
				rl.Vector2{X: float32(constants.HelicopterWidth) / 2, Y: float32(constants.HelicopterHeight) / 2},
				0.0,
				rl.White,
			)
		} else {
			// Fallback to drawing helicopter with primitives
			drawPlayerHelicopter(s.playerPos, float32(s.animTimer*10.0))
		}
	}

	// Draw text messages based on animation progress
	if animProgress < 0.3 {
		// Intro text
		if s.textAlpha < 1.0 {
			s.textAlpha += dt * 2.0
			if s.textAlpha > 1.0 {
				s.textAlpha = 1.0
			}
		}

		introText := "ENEMY HELICOPTER DETECTED"
		textWidth := rl.MeasureText(introText, 40)
		rl.DrawText(
			introText,
			(constants.ScreenWidth-int32(textWidth))/2,
			50,
			40,
			rl.ColorAlpha(rl.Red, s.textAlpha),
		)
	} else if animProgress > 0.4 && animProgress < 0.7 {
		// Middle text
		middleText := "THE SCIENTISTS ARE BEING KILLED!"
		textWidth := rl.MeasureText(middleText, 30)
		var alpha float32
		if animProgress < 0.5 {
			alpha = float32((animProgress - 0.4) / 0.1)
		} else if animProgress < 0.6 {
			alpha = 1.0
		} else {
			alpha = float32(1.0 - (animProgress-0.6)/0.1)
		}

		rl.DrawText(
			middleText,
			(constants.ScreenWidth-int32(textWidth))/2,
			100,
			30,
			rl.ColorAlpha(rl.White, alpha),
		)
	} else if animProgress > 0.8 {
		// Final battle text
		if animProgress > 0.9 {
			s.showPrompt = true
		}

		finalText := "DEFEAT THE ENEMY HELICOPTER!"
		textWidth := rl.MeasureText(finalText, 40)

		// Dramatic flashing effect
		alpha := 0.5 + 0.5*float32(math.Sin(float64(s.animTimer*8)))

		rl.DrawText(
			finalText,
			(constants.ScreenWidth-int32(textWidth))/2,
			100,
			40,
			rl.ColorAlpha(rl.Red, alpha),
		)
	}

	// Draw prompt when animation completes
	if s.showPrompt {
		promptText := "PRESS ENTER TO BEGIN BATTLE"

		// Make text pulse
		alpha := 0.5 + 0.5*float32(math.Sin(float64(s.promptTimer*4)))

		textWidth := rl.MeasureText(promptText, 30)
		rl.DrawText(
			promptText,
			(constants.ScreenWidth-int32(textWidth))/2,
			constants.ScreenHeight-100,
			30,
			rl.ColorAlpha(rl.White, alpha),
		)
	}
}

// Helper function to draw boss helicopter with primitives
func drawBossHelicopter(pos rl.Vector2, rotorAngle float32) {
	// Body (larger and red)
	rl.DrawRectangle(
		int32(pos.X-40),
		int32(pos.Y-20),
		80,
		40,
		rl.Red,
	)

	// Cockpit
	rl.DrawCircle(
		int32(pos.X),
		int32(pos.Y),
		15,
		rl.Black,
	)

	// Rotor animation
	rotorAngleRad := float64(rotorAngle)
	rotorEndX1 := pos.X + float32(math.Cos(rotorAngleRad))*50
	rotorEndY1 := pos.Y - 25 + float32(math.Sin(rotorAngleRad))*5
	rotorEndX2 := pos.X - float32(math.Cos(rotorAngleRad))*50
	rotorEndY2 := pos.Y - 25 - float32(math.Sin(rotorAngleRad))*5

	rl.DrawLine(
		int32(rotorEndX1),
		int32(rotorEndY1),
		int32(rotorEndX2),
		int32(rotorEndY2),
		rl.DarkGray,
	)

	// Rotor center
	rl.DrawCircle(
		int32(pos.X),
		int32(pos.Y-25),
		3,
		rl.Gray,
	)

	// Tail
	rl.DrawRectangle(
		int32(pos.X+40-5),
		int32(pos.Y-10),
		30,
		20,
		rl.Red,
	)

	// Tail rotor
	tailRotorY1 := pos.Y + float32(math.Sin(float64(rotorAngle*2)))*15
	tailRotorY2 := pos.Y - float32(math.Sin(float64(rotorAngle*2)))*15

	rl.DrawLine(
		int32(pos.X+65),
		int32(tailRotorY1),
		int32(pos.X+65),
		int32(tailRotorY2),
		rl.DarkGray,
	)

	// Add weapon pods under helicopter
	rl.DrawRectangle(
		int32(pos.X-30),
		int32(pos.Y+20),
		15,
		10,
		rl.DarkGray,
	)

	rl.DrawRectangle(
		int32(pos.X+15),
		int32(pos.Y+20),
		15,
		10,
		rl.DarkGray,
	)
}

// Helper function to draw player helicopter with primitives
func drawPlayerHelicopter(pos rl.Vector2, rotorAngle float32) {
	// Draw helicopter body
	rl.DrawRectangle(
		int32(pos.X-30),
		int32(pos.Y-15),
		60,
		30,
		rl.DarkGray,
	)

	// Cockpit
	rl.DrawCircleV(pos, 12, rl.SkyBlue)

	// Main rotor - animate rotation
	rotorEndX1 := pos.X + float32(math.Cos(float64(rotorAngle)))*40
	rotorEndY1 := pos.Y - 15 - 5 + float32(math.Sin(float64(rotorAngle)))*5
	rotorEndX2 := pos.X - float32(math.Cos(float64(rotorAngle)))*40
	rotorEndY2 := pos.Y - 15 - 5 - float32(math.Sin(float64(rotorAngle)))*5

	rl.DrawLine(
		int32(rotorEndX1),
		int32(rotorEndY1),
		int32(rotorEndX2),
		int32(rotorEndY2),
		rl.LightGray,
	)

	// Rotor center point
	rl.DrawCircleV(
		rl.Vector2{X: pos.X, Y: pos.Y - 15 - 5},
		3,
		rl.Gray,
	)

	// Tail
	rl.DrawRectangle(
		int32(pos.X+30-5),
		int32(pos.Y-7),
		25,
		15,
		rl.Gray,
	)

	// Tail rotor
	tailRotorY1 := pos.Y + float32(math.Sin(float64(rotorAngle*2)))*10
	tailRotorY2 := pos.Y - float32(math.Sin(float64(rotorAngle*2)))*10

	rl.DrawLine(
		int32(pos.X+30+20),
		int32(tailRotorY1),
		int32(pos.X+30+20),
		int32(tailRotorY2),
		rl.LightGray,
	)
}

// Update handles input for the boss intro screen
func (s *BossIntroScreen) Update() bool {
	// Complete animation if time is up
	if s.animTimer >= s.animDuration {
		s.showPrompt = true
	}

	// Skip to boss battle when animation is done and ENTER is pressed
	if s.showPrompt && rl.IsKeyPressed(rl.KeyEnter) {
		return true
	}

	// Allow skipping the intro with ESC
	if rl.IsKeyPressed(rl.KeyEscape) {
		return true
	}

	return false
}
