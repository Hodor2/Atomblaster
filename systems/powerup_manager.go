package systems

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/entities"
	"atomblaster/util"
)

// PowerUpManager handles the creation, updating, and removal of power-ups
type PowerUpManager struct {
	PowerUps       []*entities.PowerUp
	PowerUpPool    *util.ObjectPool
	
	// PowerUp generation settings
	SpawnTimer          float32
	BaseSpawnInterval   float32 // Average seconds between spawns
	SpawnIntervalJitter float32 // Random jitter added to spawn interval
	MaxPowerUps         int     // Maximum number of power-ups in world
	
	// World bounds
	WorldBounds rl.Rectangle
}

// NewPowerUpManager creates a new power-up management system
func NewPowerUpManager(worldBounds rl.Rectangle) *PowerUpManager {
	maxPowerUps := constants.MaxPowerUps
	
	powerUpPool := util.NewObjectPool(
		// Create function
		func() interface{} {
			return &entities.PowerUp{}
		},
		// Reset function
		func(obj interface{}) {
			powerUp := obj.(*entities.PowerUp)
			powerUp.Position = rl.Vector2{X: 0, Y: 0}
			powerUp.Velocity = rl.Vector2{X: 0, Y: 0}
			powerUp.Type = 0
			powerUp.Color = rl.White
			powerUp.Size = 12.0
			powerUp.Rotation = 0
			powerUp.Duration = constants.PowerUpDuration
			powerUp.PulseTime = 0
			powerUp.Lifetime = 30.0
			powerUp.MaxLifetime = 30.0
			powerUp.Active = false
		},
		maxPowerUps,
	)
	
	return &PowerUpManager{
		PowerUps:           make([]*entities.PowerUp, 0, maxPowerUps),
		PowerUpPool:        powerUpPool,
		SpawnTimer:         0,
		BaseSpawnInterval:  20.0, // Spawn roughly every 20 seconds
		SpawnIntervalJitter: 10.0, // +/- 10 seconds
		MaxPowerUps:        maxPowerUps,
		WorldBounds:        worldBounds,
	}
}

// Update updates all power-ups and possibly spawns new ones
func (pm *PowerUpManager) Update(dt float32, quadtree *util.Quadtree) {
	// Update spawn timer
	pm.SpawnTimer += dt
	
	// Calculate when next power-up should spawn
	nextSpawnTime := pm.BaseSpawnInterval + (rl.GetRandomFloat32()*2-1) * pm.SpawnIntervalJitter
	
	// Check if it's time to spawn and we have room
	if pm.SpawnTimer >= nextSpawnTime && len(pm.PowerUps) < pm.MaxPowerUps {
		pm.SpawnTimer = 0
		
		// Choose a random position that's at least 10% away from world edges
		padding := pm.WorldBounds.Width * 0.1
		randomX := pm.WorldBounds.X + padding + rl.GetRandomFloat32() * (pm.WorldBounds.Width - padding*2)
		randomY := pm.WorldBounds.Y + padding + rl.GetRandomFloat32() * (pm.WorldBounds.Height - padding*2)
		
		// Create and add the power-up
		powerUp := entities.NewRandomPowerUp(rl.Vector2{X: randomX, Y: randomY})
		pm.PowerUps = append(pm.PowerUps, powerUp)
		
		// Add to quadtree
		quadtree.Root.Insert(powerUp, powerUp)
	}
	
	// Update existing power-ups and remove inactive ones
	for i := len(pm.PowerUps) - 1; i >= 0; i-- {
		pm.PowerUps[i].Update(dt, pm.WorldBounds)
		
		// Remove expired power-ups
		if !pm.PowerUps[i].Active {
			// Return to pool
			pm.PowerUpPool.Return(pm.PowerUps[i])
			
			// Remove from list (using swap-and-pop for efficiency)
			pm.PowerUps[i] = pm.PowerUps[len(pm.PowerUps)-1]
			pm.PowerUps = pm.PowerUps[:len(pm.PowerUps)-1]
		}
	}
}

// Draw renders all active power-ups
func (pm *PowerUpManager) Draw(camera rl.Camera2D) {
	// Get visible bounds with some padding
	viewportBounds := rl.Rectangle{
		X:      camera.Target.X - constants.ScreenWidth/2/camera.Zoom,
		Y:      camera.Target.Y - constants.ScreenHeight/2/camera.Zoom,
		Width:  constants.ScreenWidth/camera.Zoom,
		Height: constants.ScreenHeight/camera.Zoom,
	}
	
	// Add padding for items just offscreen
	padding := float32(100)
	viewportWithPadding := rl.Rectangle{
		X:      viewportBounds.X - padding,
		Y:      viewportBounds.Y - padding,
		Width:  viewportBounds.Width + padding*2,
		Height: viewportBounds.Height + padding*2,
	}
	
	// Draw only power-ups in viewport
	for _, powerUp := range pm.PowerUps {
		// Check if power-up is in view
		powerUpBounds := powerUp.GetBounds()
		
		if util.CheckRectangleOverlap(powerUpBounds, viewportWithPadding) {
			powerUp.Draw()
		}
	}
}

// RemovePowerUp removes a specific power-up from the manager
func (pm *PowerUpManager) RemovePowerUp(powerUp *entities.PowerUp) {
	for i, p := range pm.PowerUps {
		if p == powerUp {
			// Return to pool
			pm.PowerUpPool.Return(powerUp)
			
			// Remove from list
			pm.PowerUps[i] = pm.PowerUps[len(pm.PowerUps)-1]
			pm.PowerUps = pm.PowerUps[:len(pm.PowerUps)-1]
			break
		}
	}
}

// SpawnPowerUp spawns a power-up at a specific position
func (pm *PowerUpManager) SpawnPowerUp(position rl.Vector2, type_ entities.PowerUpType) *entities.PowerUp {
	if len(pm.PowerUps) >= pm.MaxPowerUps {
		return nil
	}
	
	// Create the power-up
	powerUp := entities.NewPowerUp(position, type_)
	pm.PowerUps = append(pm.PowerUps, powerUp)
	
	return powerUp
}

// SpawnRandomPowerUp spawns a random power-up at a specific position
func (pm *PowerUpManager) SpawnRandomPowerUp(position rl.Vector2) *entities.PowerUp {
	if len(pm.PowerUps) >= pm.MaxPowerUps {
		return nil
	}
	
	// Create the power-up
	powerUp := entities.NewRandomPowerUp(position)
	pm.PowerUps = append(pm.PowerUps, powerUp)
	
	return powerUp
}

// GetPowerUpCount returns the number of active power-ups
func (pm *PowerUpManager) GetPowerUpCount() int {
	return len(pm.PowerUps)
}

// GetNearestPowerUp returns the nearest power-up to a position
func (pm *PowerUpManager) GetNearestPowerUp(position rl.Vector2) *entities.PowerUp {
	var nearest *entities.PowerUp = nil
	var nearestDistSq float32 = float32(1e10) // Very large initial value
	
	for _, powerUp := range pm.PowerUps {
		if !powerUp.Active {
			continue
		}
		
		dx := powerUp.Position.X - position.X
		dy := powerUp.Position.Y - position.Y
		distSq := dx*dx + dy*dy
		
		if distSq < nearestDistSq {
			nearest = powerUp
			nearestDistSq = distSq
		}
	}
	
	return nearest
}

// GetNearbyPowerUps returns all power-ups within a radius of a position
func (pm *PowerUpManager) GetNearbyPowerUps(position rl.Vector2, radius float32) []*entities.PowerUp {
	result := make([]*entities.PowerUp, 0)
	
	for _, powerUp := range pm.PowerUps {
		if !powerUp.Active {
			continue
		}
		
		dx := powerUp.Position.X - position.X
		dy := powerUp.Position.Y - position.Y
		distSq := dx*dx + dy*dy
		
		if distSq <= radius*radius {
			result = append(result, powerUp)
		}
	}
	
	return result
}
