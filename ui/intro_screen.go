// Package ui contains UI components for Atom Blaster
package ui

import (
	"atomblaster/constants"
	"atomblaster/entities"
	"math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// IntroScreen shows a short animation of helicopter flying toward scientists
type IntroScreen struct {
	background    rl.Texture2D
	heliTexture   rl.Texture2D
	scientists    []entities.Scientist
	door          rl.Rectangle
	heliPos       rl.Vector2
	animTimer     float32
	animDuration  float32
	showPrompt    bool
	promptTimer   float32
}

// NewIntroScreen creates a new intro screen
func NewIntroScreen(
	background rl.Texture2D,
	heliTexture rl.Texture2D,
) Screen {
	// Create scientists on right side of screen (prison camp)
	scientists := make([]entities.Scientist, 5)
	for i := 0; i < 5; i++ {
		x := float32(constants.ScreenWidth - 150)
		y := float32(200 + i*60)
		scientists[i] = entities.NewScientist(x, y)
	}

	// Create prison door
	door := rl.Rectangle{
		X: float32(constants.ScreenWidth - 100),
		Y: float32(constants.ScreenHeight/2 - 100),
		Width: 60,
		Height: 200,
	}

	// Start helicopter off-screen to the left
	heliStartPos := rl.Vector2{
		X: -float32(constants.HelicopterWidth),
		Y: float32(constants.ScreenHeight / 2),
	}

	return &IntroScreen{
		background:   background,
		heliTexture:  heliTexture,
		scientists:   scientists,
		door:         door,
		heliPos:      heliStartPos,
		animTimer:    0,
		animDuration: 4.0, // 4 seconds animation
		showPrompt:   false,
		promptTimer:  0,
	}
}

// Draw renders the intro animation
func (s *IntroScreen) Draw() {
	// Background
	rl.DrawTexture(s.background, 0, 0, rl.White)

	// Update animation
	dt := rl.GetFrameTime()
	s.animTimer += dt
	s.promptTimer += dt

	// Handle animation phases
	animProgress := s.animTimer / s.animDuration
	if animProgress > 1.0 {
		animProgress = 1.0
		s.showPrompt = true
	}

	// Calculate helicopter position based on animation progress
	targetX := float32(constants.ScreenWidth / 2)
	// Using targetY in the calculation to avoid unused variable warning
	targetY := float32(constants.ScreenHeight / 2)
	s.heliPos.X = -float32(constants.HelicopterWidth) + float32(targetX+float32(constants.HelicopterWidth))*animProgress
	s.heliPos.Y = targetY
	
	// Draw prison fence/building
	rl.DrawRectangle(
		constants.ScreenWidth-200,
		100,
		180,
		constants.ScreenHeight-200,
		rl.Gray,
	)
	
	// Draw fence lines
	for i := 0; i < 10; i++ {
		rl.DrawLine(
			int32(constants.ScreenWidth-200),
			int32(100 + i*40),
			int32(constants.ScreenWidth-20),
			int32(100 + i*40),
			rl.DarkGray,
		)
	}

	// Draw door
	rl.DrawRectangleRec(s.door, rl.DarkGray)
	rl.DrawRectangleLinesEx(s.door, 2, rl.Black)

	// Draw scientists
	for i := range s.scientists {
		// Make scientists wave if helicopter is close
		if animProgress > 0.7 {
			// Customize scientist drawing - simplified for example
			posX := int32(s.scientists[i].Pos.X)
			posY := int32(s.scientists[i].Pos.Y)
			
			// Draw stick figure
			rl.DrawCircle(posX, posY-15, 10, rl.White)     // Head
			rl.DrawLine(posX, posY-5, posX, posY+15, rl.White)  // Body
			
			// Legs
			rl.DrawLine(posX, posY+15, posX-8, posY+30, rl.White)
			rl.DrawLine(posX, posY+15, posX+8, posY+30, rl.White)
			
			// Arms - wave the right arm
			armWaveOffset := int32(5 * math.Sin(float64(s.animTimer * 10)))
			rl.DrawLine(posX, posY, posX-10, posY+5, rl.White)   // Left arm
			rl.DrawLine(posX, posY, posX+10, posY-10-armWaveOffset, rl.White)  // Right arm waving
		} else {
			// Scientists not yet waving - draw normally
			posX := int32(s.scientists[i].Pos.X)
			posY := int32(s.scientists[i].Pos.Y)
			
			// Draw stick figure
			rl.DrawCircle(posX, posY-15, 10, rl.White)     // Head
			rl.DrawLine(posX, posY-5, posX, posY+15, rl.White)  // Body
			rl.DrawLine(posX, posY+15, posX-8, posY+30, rl.White)  // Left leg
			rl.DrawLine(posX, posY+15, posX+8, posY+30, rl.White)  // Right leg
			rl.DrawLine(posX, posY, posX-10, posY+5, rl.White)   // Left arm
			rl.DrawLine(posX, posY, posX+10, posY+5, rl.White)   // Right arm
		}
	}

	// Draw helicopter
	if s.heliTexture.ID > 0 {
		rl.DrawTexturePro(
			s.heliTexture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(s.heliTexture.Width), Height: float32(s.heliTexture.Height)},
			rl.Rectangle{
				X:      s.heliPos.X - float32(constants.HelicopterWidth)/2,
				Y:      s.heliPos.Y - float32(constants.HelicopterHeight)/2,
				Width:  float32(constants.HelicopterWidth),
				Height: float32(constants.HelicopterHeight),
			},
			rl.Vector2{X: float32(constants.HelicopterWidth) / 2, Y: float32(constants.HelicopterHeight) / 2},
			0.0,
			rl.White,
		)
	} else {
		// Fallback to drawing helicopter with primitives (from player.go)
		rl.DrawRectangle(
			int32(s.heliPos.X - float32(constants.HelicopterWidth)/2),
			int32(s.heliPos.Y - float32(constants.HelicopterHeight)/2),
			constants.HelicopterWidth,
			constants.HelicopterHeight,
			rl.DarkGray,
		)
		
		// Cockpit
		rl.DrawCircle(
			int32(s.heliPos.X),
			int32(s.heliPos.Y),
			constants.CockpitRadius,
			rl.SkyBlue,
		)
		
		// Rotor animation
		rotorAngle := float64(s.animTimer * 10.0)
		rotorEndX1 := s.heliPos.X + float32(math.Cos(rotorAngle))*float32(constants.RotorLength)
		rotorEndY1 := s.heliPos.Y - float32(constants.HelicopterHeight)/2 - 5 + float32(math.Sin(rotorAngle))*5
		rotorEndX2 := s.heliPos.X - float32(math.Cos(rotorAngle))*float32(constants.RotorLength)
		rotorEndY2 := s.heliPos.Y - float32(constants.HelicopterHeight)/2 - 5 - float32(math.Sin(rotorAngle))*5
		
		rl.DrawLine(
			int32(rotorEndX1),
			int32(rotorEndY1),
			int32(rotorEndX2),
			int32(rotorEndY2),
			rl.LightGray,
		)
		
		// Rotor center
		rl.DrawCircle(
			int32(s.heliPos.X),
			int32(s.heliPos.Y - float32(constants.HelicopterHeight)/2 - 5),
			3,
			rl.Gray,
		)
		
		// Tail
		rl.DrawRectangle(
			int32(s.heliPos.X + float32(constants.HelicopterWidth)/2 - 5),
			int32(s.heliPos.Y - float32(constants.HelicopterHeight)/4),
			25,
			constants.HelicopterHeight/2,
			rl.Gray,
		)
		
		// Tail rotor
		tailRotorY1 := s.heliPos.Y + float32(math.Sin(float64(rotorAngle*2)))*float32(constants.TailRotorLength)
		tailRotorY2 := s.heliPos.Y - float32(math.Sin(float64(rotorAngle*2)))*float32(constants.TailRotorLength)
		
		rl.DrawLine(
			int32(s.heliPos.X + float32(constants.HelicopterWidth)/2 + 20),
			int32(tailRotorY1),
			int32(s.heliPos.X + float32(constants.HelicopterWidth)/2 + 20),
			int32(tailRotorY2),
			rl.LightGray,
		)
	}

	// Draw intro text
	if animProgress < 0.3 {
		introText := "ATOM BLASTER"
		
		// Fade in effect
		alpha := animProgress / 0.3
		
		textWidth := rl.MeasureText(introText, 60)
		rl.DrawText(
			introText,
			(constants.ScreenWidth - int32(textWidth)) / 2,
			50,
			60,
			rl.ColorAlpha(rl.Red, alpha),
		)
	}
	
	// Draw story text after helicopter gets close
	if animProgress > 0.6 && animProgress < 0.9 {
		storyText := "RESCUE THE SCIENTISTS"
		
		// Fade in/out effect
		alpha := (animProgress - 0.6) / 0.3
		if animProgress > 0.8 {
			alpha = 1.0 - ((animProgress - 0.8) / 0.1)
		}
		
		textWidth := rl.MeasureText(storyText, 40)
		rl.DrawText(
			storyText,
			(constants.ScreenWidth - int32(textWidth)) / 2,
			100,
			40,
			rl.ColorAlpha(rl.White, alpha),
		)
	}

	// Draw prompt when animation completes
	if s.showPrompt {
		promptText := "PRESS ENTER TO START RESCUE MISSION"
		
		// Make text pulse
		alpha := 0.5 + 0.5*float32(math.Sin(float64(s.promptTimer * 4)))
		
		textWidth := rl.MeasureText(promptText, 30)
		rl.DrawText(
			promptText,
			(constants.ScreenWidth - int32(textWidth)) / 2,
			constants.ScreenHeight - 100,
			30,
			rl.ColorAlpha(rl.White, alpha),
		)
	}
}

// Update handles input for the intro screen
func (s *IntroScreen) Update() bool {
	// Skip to title screen when animation is done and ENTER is pressed
	if s.showPrompt && rl.IsKeyPressed(rl.KeyEnter) {
		return true
	}
	return false
}
