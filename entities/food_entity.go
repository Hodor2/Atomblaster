package entities

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"math"
)

// FoodType defines different food variants
type FoodType int

const (
	FoodTypeBasic FoodType = iota
	FoodTypePremium
	FoodTypeRare
)

// Food represents a collectible item that increases player size/score
type Food struct {
	Position   rl.Vector2
	Velocity   rl.Vector2 // For moving food
	Color      rl.Color
	Size       float32
	Value      int
	Type       FoodType
	PulseTime  float32
	Rotation   float32
	HasPhysics bool // Whether food moves or stays stationary
}

// NewFood creates a basic food item
func NewFood(position rl.Vector2) *Food {
	// Basic small food
	size := 3.0 + rl.GetRandomFloat32() * 2.0
	value := constants.FoodBaseValue
	
	return &Food{
		Position:   position,
		Velocity:   rl.Vector2{X: 0, Y: 0},
		Color:      rl.Green,
		Size:       size,
		Value:      value,
		Type:       FoodTypeBasic,
		PulseTime:  rl.GetRandomFloat32() * 2 * math.Pi, // Random start phase
		Rotation:   rl.GetRandomFloat32() * 360,
		HasPhysics: false,
	}
}

// NewRandomFood creates a random food item with varying type/value
func NewRandomFood(position rl.Vector2) *Food {
	// Determine food type based on rarity
	foodTypeRoll := rl.GetRandomFloat32()
	var foodType FoodType
	var color rl.Color
	var size float32
	var value int
	var hasPhysics bool
	
	switch {
	case foodTypeRoll < 0.02: // 2% chance for rare food
		foodType = FoodTypeRare
		color = rl.Gold
		size = 8.0 + rl.GetRandomFloat32() * 2.0
		value = constants.FoodBaseValue * 5 + int(rl.GetRandomFloat32() * 15)
		hasPhysics = true // Rare food moves around
		
	case foodTypeRoll < 0.15: // 13% chance for premium food
		foodType = FoodTypePremium
		color = rl.Blue
		size = 5.0 + rl.GetRandomFloat32() * 3.0
		value = constants.FoodBaseValue * 2 + int(rl.GetRandomFloat32() * 8)
		hasPhysics = rl.GetRandomFloat32() < 0.3 // 30% chance to have physics
		
	default: // 85% common food
		foodType = FoodTypeBasic
		color = rl.Green
		size = 3.0 + rl.GetRandomFloat32() * 2.0
		value = constants.FoodBaseValue
		hasPhysics = false
	}
	
	// Set velocity for moving food
	var velocity rl.Vector2
	if hasPhysics {
		angle := rl.GetRandomFloat32() * 2 * math.Pi
		speed := 20.0 + rl.GetRandomFloat32() * 30.0
		velocity = rl.Vector2{
			X: float32(math.Cos(float64(angle))) * speed,
			Y: float32(math.Sin(float64(angle))) * speed,
		}
	}
	
	return &Food{
		Position:   position,
		Velocity:   velocity,
		Color:      color,
		Size:       size,
		Value:      value,
		Type:       foodType,
		PulseTime:  rl.GetRandomFloat32() * 2 * math.Pi, // Random start phase
		Rotation:   rl.GetRandomFloat32() * 360,
		HasPhysics: hasPhysics,
	}
}

// Update updates the food's state
func (f *Food) Update(dt float32, worldBounds rl.Rectangle) {
	// Animation update
	f.PulseTime += dt * 2
	if f.PulseTime > 2*math.Pi {
		f.PulseTime -= 2 * math.Pi
	}
	
	// Rotation update
	rotationSpeed := 30.0
	if f.Type == FoodTypePremium {
		rotationSpeed = 60.0
	} else if f.Type == FoodTypeRare {
		rotationSpeed = 90.0
	}
	
	f.Rotation += dt * rotationSpeed
	if f.Rotation > 360 {
		f.Rotation -= 360
	}
	
	// Physics update (for moving food)
	if f.HasPhysics {
		f.Position.X += f.Velocity.X * dt
		f.Position.Y += f.Velocity.Y * dt
		
		// Bounce off world boundaries
		if f.Position.X < worldBounds.X + f.Size {
			f.Position.X = worldBounds.X + f.Size
			f.Velocity.X = -f.Velocity.X
		} else if f.Position.X > worldBounds.X + worldBounds.Width - f.Size {
			f.Position.X = worldBounds.X + worldBounds.Width - f.Size
			f.Velocity.X = -f.Velocity.X
		}
		
		if f.Position.Y < worldBounds.Y + f.Size {
			f.Position.Y = worldBounds.Y + f.Size
			f.Velocity.Y = -f.Velocity.Y
		} else if f.Position.Y > worldBounds.Y + worldBounds.Height - f.Size {
			f.Position.Y = worldBounds.Y + worldBounds.Height - f.Size
			f.Velocity.Y = -f.Velocity.Y
		}
		
		// Apply slight drag to eventually slow down
		f.Velocity.X *= 0.99
		f.Velocity.Y *= 0.99
	}
}

// Draw renders the food item
func (f *Food) Draw() {
	// Get pulse factor (between 0.85 and 1.15)
	pulseFactor := 0.85 + 0.15*float32(math.Sin(float64(f.PulseTime)))
	
	// Draw food with pulsing size
	drawSize := f.Size * pulseFactor
	
	// Number of sides varies by food type
	sides := 5 // Pentagon for basic
	if f.Type == FoodTypePremium {
		sides = 6 // Hexagon for premium
	} else if f.Type == FoodTypeRare {
		sides = 8 // Octagon for rare
	}
	
	// Draw a polygon with rotation
	centerX := f.Position.X
	centerY := f.Position.Y
	
	// Calculate all vertex positions
	vertices := make([]rl.Vector2, sides)
	for i := 0; i < sides; i++ {
		angle := f.Rotation*math.Pi/180 + float32(i)*2*math.Pi/float32(sides)
		vertices[i] = rl.Vector2{
			X: centerX + drawSize * float32(math.Cos(