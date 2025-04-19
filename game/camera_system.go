package game

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/entities"
	"atomblaster/util"
)

// Camera handles the view into the game world
type Camera struct {
	// Target entity to follow
	Target   *entities.Player
	
	// Camera properties
	Position rl.Vector2
	Rotation float32
	Zoom     float32
	
	// Smoothing parameters
	SmoothingFactor float32
	EdgePadding     float32
	
	// Camera effects
	ShakeAmount     float32
	ShakeDecay      float32
	ShakeDuration   float32
	
	// World boundaries
	WorldBounds     rl.Rectangle
}

// NewCamera creates a camera that follows a target entity
func NewCamera(target *entities.Player, worldBounds rl.Rectangle) *Camera {
	return &Camera{
		Target:          target,
		Position:        rl.Vector2{X: 0, Y: 0},
		Rotation:        0.0,
		Zoom:            1.0,
		SmoothingFactor: 0.1,  // Lower = smoother camera
		EdgePadding:     100,  // Keep player this far from screen edge
		ShakeAmount:     0.0,
		ShakeDecay:      5.0,  // How quickly shake fades
		ShakeDuration:   0.0,
		WorldBounds:     worldBounds,
	}
}

// Update updates the camera position
func (c *Camera) Update(dt float32) {
	// Skip if no target
	if c.Target == nil {
		return
	}
	
	// Calculate target position (centered on player)
	targetPos := rl.Vector2{
		X: c.Target.Position.X - constants.ScreenWidth/2/c.Zoom,
		Y: c.Target.Position.Y - constants.ScreenHeight/2/c.Zoom,
	}
	
	// Apply smoothing with lerp
	c.Position.X = util.Lerp(c.Position.X, targetPos.X, c.SmoothingFactor)
	c.Position.Y = util.Lerp(c.Position.Y, targetPos.Y, c.SmoothingFactor)
	
	// Constrain camera to world bounds
	c.Position.X = util.Clamp(c.Position.X, c.WorldBounds.X, 
							 c.WorldBounds.X + c.WorldBounds.Width - constants.ScreenWidth/c.Zoom)
	c.Position.Y = util.Clamp(c.Position.Y, c.WorldBounds.Y, 
							 c.WorldBounds.Y + c.WorldBounds.Height - constants.ScreenHeight/c.Zoom)
	
	// Update camera shake effect
	if c.ShakeDuration > 0 {
		c.ShakeDuration -= dt
		
		// Apply random offset for shake
		if c.ShakeDuration > 0 {
			c.Position.X += (rl.GetRandomFloat32() * 2 - 1) * c.ShakeAmount
			c.Position.Y += (rl.GetRandomFloat32() * 2 - 1) * c.ShakeAmount
			
			// Reduce shake amount over time
			c.ShakeAmount -= c.ShakeDecay * dt
			if c.ShakeAmount < 0 {
				c.ShakeAmount = 0
			}
		} else {
			// Reset shake when duration expires
			c.ShakeDuration = 0
			c.ShakeAmount = 0
		}
	}
}

// GetRLCamera2D returns a raylib Camera2D for rendering
func (c *Camera) GetRLCamera2D() rl.Camera2D {
	return rl.Camera2D{
		Target:   c.Position,
		Offset:   rl.Vector2{X: 0, Y: 0},
		Rotation: c.Rotation,
		Zoom:     c.Zoom,
	}
}

// StartShake initiates a camera shake effect
func (c *Camera) StartShake(amount float32, duration float32) {
	c.ShakeAmount = amount
	c.ShakeDuration = duration
}

// SetZoom sets the camera zoom level
func (c *Camera) SetZoom(zoom float32) {
	// Constrain to reasonable values
	c.Zoom = util.Clamp(zoom, 0.25, 3.0)
}

// ZoomIn increases camera zoom
func (c *Camera) ZoomIn(amount float32) {
	c.SetZoom(c.Zoom + amount)
}

// ZoomOut decreases camera zoom
func (c *Camera) ZoomOut(amount float32) {
	c.SetZoom(c.Zoom - amount)
}

// SetTarget changes the entity the camera follows
func (c *Camera) SetTarget(target *entities.Player) {
	c.Target = target
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldPos rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: (worldPos.X - c.Position.X) * c.Zoom,
		Y: (worldPos.Y - c.Position.Y) * c.Zoom,
	}
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenPos rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: screenPos.X/c.Zoom + c.Position.X,
		Y: screenPos.Y/c.Zoom + c.Position.Y,
	}
}

// GetVisibleBounds returns the rectangle of the world currently visible
func (c *Camera) GetVisibleBounds() rl.Rectangle {
	return rl.Rectangle{
		X:      c.Position.X,
		Y:      c.Position.Y,
		Width:  constants.ScreenWidth / c.Zoom,
		Height: constants.ScreenHeight / c.Zoom,
	}
}
