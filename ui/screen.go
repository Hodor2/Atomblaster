// Package ui contains UI components for Atom Blaster
package ui

// Screen is an interface for all screens (intro, title, game, pause, game-over)
type Screen interface {
    // Draw renders the screen
    Draw()
    
    // Update processes user input and returns true if the screen should transition
    // to the next screen
    Update() bool
}