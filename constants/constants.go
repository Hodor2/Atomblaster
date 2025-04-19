package constants

// Screen constants
const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	
	// World size (much larger than screen)
	WorldWidth  = 8000
	WorldHeight = 8000
	
	// Visual boundary thickness
	WorldBorderThickness = 50
	
	// Object limits
	MaxEntities = 2000
	MaxFood     = 1000
	MaxPowerUps = 20
	MaxAIs      = 100
	
	// Physics constants
	DefaultFriction    = 0.95
	MaxHelicopterSpeed = 300.0
	
	// Game balance
	InitialPlayerSize  = 20.0
	MaxPlayerSize      = 100.0
	FoodBaseValue      = 5
	PowerUpDuration    = 10.0 // seconds
	
	// Magnet power-up
	MagnetRange        = 200.0
	
	// World generation
	FoodDensity        = 0.1  // Food per 100x100 area
	PowerUpSpawnRate   = 0.05 // Power-ups per second
	
	// Difficulty settings
	DefaultDifficulty  = 1    // 0=Easy, 1=Medium, 2=Hard, 3=Dynamic
)

// Game states
const (
	StateTitle = iota
	StateGame
	StatePause
	StateGameOver
)
