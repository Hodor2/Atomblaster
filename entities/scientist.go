package entities

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ScientistState represents different states a scientist can be in
type ScientistState int

const (
	Wandering ScientistState = iota
	FollowingPlayer
	Rescued
)

// Scientist represents a scientist that needs to be rescued
type Scientist struct {
	Pos            rl.Vector2
	Velocity       rl.Vector2
	Size           float32
	State          ScientistState
	WanderTimer    float32
	WanderDir      rl.Vector2
	RescueTimer    float32
	FollowOffset   rl.Vector2
	IsBeingRescued bool
	AnimTimer      float32 // Added for animation timing
}

// NewScientist creates a new scientist at the given position
func NewScientist(x, y float32) Scientist {
	return Scientist{
		Pos:            rl.Vector2{X: x, Y: y},
		Velocity:       rl.Vector2{X: 0, Y: 0},
		Size:           15,
		State:          Wandering,
		WanderTimer:    0,
		WanderDir:      rl.Vector2{X: 0, Y: 0},
		RescueTimer:    0,
		FollowOffset:   rl.Vector2{X: float32(rand.Intn(30) - 15), Y: float32(rand.Intn(30) - 15)},
		IsBeingRescued: false,
		AnimTimer:      float32(rand.Float64() * 10.0), // Random start time for animation variation
	}
}

// Update updates the scientist's position and state
func (s *Scientist) Update(dt float32, playerPos rl.Vector2, rescueZoneRect rl.Rectangle) bool {
	// Update animation timer
	s.AnimTimer += dt
	
	// Return immediately if already rescued
	if s.State == Rescued {
		return true
	}
	
	// Check if scientist is in rescue zone while following player
	if s.State == FollowingPlayer {
		sciRect := s.GetRectangle()
		if rl.CheckCollisionRecs(sciRect, rescueZoneRect) {
			s.State = Rescued
			return true
		}
	}
	
	switch s.State {
	case Wandering:
		// Random wandering behavior
		s.WanderTimer -= dt
		if s.WanderTimer <= 0 {
			// Change direction every few seconds
			s.WanderTimer = rand.Float32() * 2.0 + 1.0
			s.WanderDir = rl.Vector2{
				X: rand.Float32()*2.0 - 1.0,
				Y: rand.Float32()*2.0 - 1.0,
			}
			// Normalize the direction
			length := float32(math.Sqrt(float64(s.WanderDir.X*s.WanderDir.X + s.WanderDir.Y*s.WanderDir.Y)))
			if length > 0 {
				s.WanderDir.X /= length
				s.WanderDir.Y /= length
			}
		}
		
		// Move in the wander direction
		s.Velocity.X = s.WanderDir.X * 30.0
		s.Velocity.Y = s.WanderDir.Y * 30.0
		
		// Check if player is nearby to switch to following
		distToPlayer := rl.Vector2Distance(s.Pos, playerPos)
		if distToPlayer < 100 {
			s.State = FollowingPlayer
		}
		
	case FollowingPlayer:
		// Follow the player with unique offset
		targetX := playerPos.X + s.FollowOffset.X
		targetY := playerPos.Y + s.FollowOffset.Y
		
		// Calculate direction to target
		dirX := targetX - s.Pos.X
		dirY := targetY - s.Pos.Y
		
		// Normalize direction
		length := float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
		if length > 0 {
			dirX /= length
			dirY /= length
		}
		
		// Set velocity to move toward player
		s.Velocity.X = dirX * 80.0
		s.Velocity.Y = dirY * 80.0
	}
	
	// Apply movement
	s.Pos.X += s.Velocity.X * dt
	s.Pos.Y += s.Velocity.Y * dt
	
	// Keep within screen bounds (assuming constants exist elsewhere)
	if s.Pos.X < s.Size {
		s.Pos.X = s.Size
		s.WanderDir.X *= -1
	}
	if s.Pos.X > 800 - s.Size {
		s.Pos.X = 800 - s.Size
		s.WanderDir.X *= -1
	}
	if s.Pos.Y < s.Size {
		s.Pos.Y = s.Size
		s.WanderDir.Y *= -1
	}
	if s.Pos.Y > 600 - s.Size {
		s.Pos.Y = 600 - s.Size
		s.WanderDir.Y *= -1
	}
	
	return false // Not rescued yet
}

// Draw renders the scientist on screen using the stick figure style from intro screen
func (s *Scientist) Draw() {
	if s.State == Rescued {
		return // Don't draw if already rescued
	}
	
	// Get position coordinates
	posX := int32(s.Pos.X)
	posY := int32(s.Pos.Y)
	
	// Draw stick figure with improved appearance
	
	// Head
	headColor := rl.White
	if s.State == FollowingPlayer {
		headColor = rl.Green // Green when following player
	}
	rl.DrawCircle(posX, posY-15, 10, headColor)
	
	// Body
	rl.DrawLine(posX, posY-5, posX, posY+15, headColor)
	
	// Legs
	rl.DrawLine(posX, posY+15, posX-8, posY+30, headColor)
	rl.DrawLine(posX, posY+15, posX+8, posY+30, headColor)
	
	// Arms - with animation
	if s.State == FollowingPlayer {
		// Waving animation for right arm when following player
		armWaveOffset := int32(5 * math.Sin(float64(s.AnimTimer * 10)))
		rl.DrawLine(posX, posY, posX-10, posY+5, headColor)  // Left arm
		rl.DrawLine(posX, posY, posX+10, posY-10-armWaveOffset, headColor)  // Right arm waving
	} else {
		// Normal arms when wandering
		armSwing := float32(math.Sin(float64(s.AnimTimer * 3))) * 3
		rl.DrawLine(posX, posY, posX-10, posY+int32(5+armSwing), headColor)   // Left arm
		rl.DrawLine(posX, posY, posX+10, posY+int32(5-armSwing), headColor)   // Right arm
	}
	
	// Face expression
	eyeOffsetY := int32(-15)
	if s.State == FollowingPlayer {
		// Happy face
		rl.DrawCircle(posX-4, posY+eyeOffsetY-2, 2, rl.Black)  // Left eye
		rl.DrawCircle(posX+4, posY+eyeOffsetY-2, 2, rl.Black)  // Right eye
		
		// Smile
		rl.DrawCircleLines(
			posX,
			posY+eyeOffsetY+4,
			4,
			rl.Black,
		)
	} else {
		// Neutral face
		rl.DrawCircle(posX-4, posY+eyeOffsetY-1, 2, rl.Black)  // Left eye
		rl.DrawCircle(posX+4, posY+eyeOffsetY-1, 2, rl.Black)  // Right eye
		
		// Straight mouth
		rl.DrawLine(
			posX-4, 
			posY+eyeOffsetY+4,
			posX+4,
			posY+eyeOffsetY+4,
			rl.Black,
		)
	}
}

// GetRectangle returns a rectangle representing the scientist
func (s *Scientist) GetRectangle() rl.Rectangle {
	return rl.Rectangle{
		X:      s.Pos.X - s.Size,
		Y:      s.Pos.Y - s.Size,
		Width:  s.Size * 2,
		Height: s.Size * 2,
	}
}

// IsNearPlayer checks if the scientist is near the player
func (s *Scientist) IsNearPlayer(playerPos rl.Vector2, distance float32) bool {
	return rl.Vector2Distance(s.Pos, playerPos) < distance
}
