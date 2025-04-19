// game/game.go
package game

import (
	"atomblaster/audio"
)

// GameWrapper wraps the GameState and provides compatibility with the main function
type GameWrapper struct {
	state *GameState
	audio *audio.AudioSystem
}

// New creates a new game instance for use with the main function
func New() *GameWrapper {
	// Create audio system
	audioSystem := audio.NewAudioSystem()

	// Create game state
	gameState := NewGameState(audioSystem)

	// Return wrapper with both components
	return &GameWrapper{
		state: gameState,
		audio: audioSystem,
	}
}

// Update updates the game state
func (g *GameWrapper) Update(dt float32) {
	// Pass the delta time to the GameState's Update method
	g.state.Update(dt)
}

// UpdateAudio updates the audio system
func (g *GameWrapper) UpdateAudio() {
	g.audio.Update()
}

// Draw renders the game
func (g *GameWrapper) Draw() {
	g.state.Draw()
}

// Cleanup releases resources
func (g *GameWrapper) Cleanup() {
	g.audio.Cleanup()
}
