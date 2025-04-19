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
