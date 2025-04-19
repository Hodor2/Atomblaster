// Package constants provides game constants for Atom Blaster
package constants

// Screen dimensions
const (
    ScreenWidth  = 800
    ScreenHeight = 600
)

// Game states
const (
    StateIntro = iota  // New intro state
    StateTitle
    StateGame
    StatePause
    StateGameOver
    StateBossIntro     // Add this new state for boss intro
)

// Game parameters
const (
    MaxLevel             = 10
    PlayerInitialHealth  = 3
    FireCooldownDuration = 0.2 // seconds between shots
)

// Helicopter parameters
const (
    HelicopterWidth  = 60
    HelicopterHeight = 30
    CockpitRadius    = 12
    RotorLength      = 40
    TailRotorLength  = 10
)

// Bullet parameters
const (
    // Speed at which bullets travel (pixels per second)
    BulletSpeed = 600.0

    // How many seconds a bullet remains before expiring
    BulletLifetime = 2.0
)
