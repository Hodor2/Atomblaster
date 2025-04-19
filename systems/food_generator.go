package systems

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/entities"
	"atomblaster/util"
	"math"
)

// FoodGenerator manages the creation and lifecycle of food entities
type FoodGenerator struct {
	FoodEntities []*entities.Food
	FoodPool     *util.ObjectPool
	
	// Food generation parameters
	MaxFood           int
	SpawnTimer        float32
	BaseSpawnRate     float32 // Base food items per second
	SpawnRateVariance float32 // Random variance in spawn rate
	CurrentSpawnRate  float32 // Current adjusted spawn rate
	
	// World parameters
	WorldBounds rl.Rectangle
	
	// Food cluster generation
	ClusterTimer      float32
	ClusterInterval   float32
	ClusterChance     float32
}

// NewFoodGenerator creates a new food generation system
func NewFoodGenerator(worldBounds rl.Rectangle) *FoodGenerator {
	maxFood := constants.MaxFood
	
	// Create pool for food entities
	foodPool := util.NewObjectPool(
		// Create function
		func() interface{} { 
			return &entities.Food{} 
		},
		// Reset function
		func(obj interface{}) {
			food := obj.(*entities.Food)
			food.Position = rl.Vector2{X: 0, Y: 0}
			food.Velocity = rl.Vector2{X: 0, Y: 0}
			food.Color = rl.Green
			food.Size = 3.0
			food.Value = constants.FoodBaseValue
			food.Type = entities.FoodTypeBasic
			food.PulseTime = 0
			food.Rotation = 0
			food.HasPhysics = false
		},
		maxFood / 2, // Pre-allocate half of max capacity
	)
	
	return &FoodGenerator{
		FoodEntities:     make([]*entities.Food, 0, maxFood),
		FoodPool:         foodPool,
		MaxFood:          maxFood,
		SpawnTimer:       0,
		BaseSpawnRate:    8.0,  // 8 food items per second
		SpawnRateVariance: 3.0, // +/- 3 variance
		CurrentSpawnRate: 8.0,  // Initial rate
		WorldBounds:      worldBounds,
		ClusterTimer:     0,
		ClusterInterval:  5.0,  // Check for cluster generation every 5 seconds
		ClusterChance:    0.3,  // 30% chance to spawn a cluster
	}
}

// Initialize with some starting food
func (fg *FoodGenerator) Initialize(amount int) {
	// Ensure we don't exceed the maximum
	if amount > fg.MaxFood {
		amount = fg.MaxFood
	}
	
	// Create initial food items distributed across the world
	for i := 0; i < amount; i++ {
		randomX := fg.WorldBounds.X + rl.GetRandomFloat32() * fg.WorldBounds.Width
		randomY := fg.WorldBounds.Y + rl.GetRandomFloat32() * fg.WorldBounds.Height
		
		// Create more premium/rare food initially for an exciting start
		var food *entities.Food
		
		if i % 20 == 0 { // 5% rare
			food = entities.NewRandomFood(rl.Vector2{X: randomX, Y: randomY})
			if food.Type != entities.FoodTypeRare {
				food.Type = entities.FoodTypeRare
				food.Color = rl.Gold
				food.Size = 8.0 + rl.GetRandomFloat32() * 2.0
				food.Value = constants.FoodBaseValue * 5
				food.HasPhysics = true
			}
		} else if i % 5 == 0 { // 20% premium
			food = entities.NewRandomFood(rl.Vector2{X: randomX, Y: randomY})
			if food.Type != entities.FoodTypePremium {
				food.Type = entities.FoodTypePremium
				food.Color = rl.Blue
				food.Size = 5.0 + rl.GetRandomFloat32() * 3.0
				food.Value = constants.FoodBaseValue * 2
				food.HasPhysics = rl.GetRandomFloat32() < 0.3
			}
		} else { // 75% basic
			food = entities.NewFood(rl.Vector2{X: randomX, Y: randomY})
		}
		
		fg.FoodEntities = append(fg.FoodEntities, food)
	}
}

// Update updates the food generation system
func (fg *FoodGenerator) Update(dt float32, quadtree *util.Quadtree) {
	// Update spawn timer
	fg.SpawnTimer += dt
	
	// Vary spawn rate slightly to create natural feeling
	fg.CurrentSpawnRate = fg.BaseSpawnRate + (rl.GetRandomFloat32()*2-1) * fg.SpawnRateVariance
	
	// Check if it's time to spawn new food
	spawnInterval := 1.0 / fg.CurrentSpawnRate
	
	// Spawn new food if below maximum
	for fg.SpawnTimer >= spawnInterval && len(fg.FoodEntities) < fg.MaxFood {
		fg.SpawnTimer -= spawnInterval
		
		// Create new food at random position
		randomX := fg.WorldBounds.X + rl.GetRandomFloat32() * fg.WorldBounds.Width
		randomY := fg.WorldBounds.Y + rl.GetRandomFloat32() * fg.WorldBounds.Height
		
		food := entities.NewRandomFood(rl.Vector2{X: randomX, Y: randomY})
		fg.FoodEntities = append(fg.FoodEntities, food)
		
		// Add to quadtree for collision detection
		quadtree.Root.Insert(food, food)
	}
	
	// Handle cluster generation
	fg.ClusterTimer += dt
	if fg.ClusterTimer >= fg.ClusterInterval {
		fg.ClusterTimer = 0
		
		// Randomly decide whether to generate a cluster
		if rl.GetRandomFloat32() < fg.ClusterChance && len(fg.FoodEntities) < fg.MaxFood - 20 {
			fg.GenerateCluster(quadtree)
		}
	}
	
	// Update all food entities
	for _, food := range fg.FoodEntities {
		food.Update(dt, fg.WorldBounds)
	}
}

// Generate a cluster of food in one area
func (fg *FoodGenerator) GenerateCluster(quadtree *util.Quadtree) {
	// Choose a random point for the cluster center
	centerX := fg.WorldBounds.X + rl.GetRandomFloat32() * fg.WorldBounds.Width
	centerY := fg.WorldBounds.Y + rl.GetRandomFloat32() * fg.WorldBounds.Height
	center := rl.Vector2{X: centerX, Y: centerY}
	
	// Determine cluster parameters
	clusterSize := 10 + int(rl.GetRandomFloat32() * 15) // 10-25 food items
	clusterRadius := 50.0 + rl.GetRandomFloat32() * 100.0 // 50-150 radius
	
	// Higher chance of premium/rare food in clusters
	premiumChance := 0.3  // 30% chance of premium
	rareChance := 0.08    // 8% chance of rare
	
	// Generate food in the cluster
	for i := 0; i < clusterSize; i++ {
		// Make sure we don't exceed max food
		if len(fg.FoodEntities) >= fg.MaxFood {
			break
		}
		
		// Random position within the cluster radius
		angle := rl.GetRandomFloat32() * 2 * math.Pi
		distance := rl.GetRandomFloat32() * clusterRadius
		
		position := rl.Vector2{
			X: center.X + float32(math.Cos(float64(angle))) * distance,
			Y: center.Y + float32(math.Sin(float64(angle))) * distance,
		}
		
		// Create food with higher chance of premium/rare
		var food *entities.Food
		
		typeRoll := rl.GetRandomFloat32()
		if typeRoll < rareChance {
			// Create a rare food
			food = entities.NewRandomFood(position)
			food.Type = entities.FoodTypeRare
			food.Color = rl.Gold
			food.Size = 8.0 + rl.GetRandomFloat32() * 2.0
			food.Value = constants.FoodBaseValue * 5 + int(rl.GetRandomFloat32() * 15)
			food.HasPhysics = true
			
			// Set velocity for rare food
			angle := rl.GetRandomFloat32() * 2 * math.Pi
			speed := 20.0 + rl.GetRandomFloat32() * 30.0
			food.Velocity = rl.Vector2{
				X: float32(math.Cos(float64(angle))) * speed,
				Y: float32(math.Sin(float64(angle))) * speed,
			}
		} else if typeRoll < rareChance + premiumChance {
			// Create a premium food
			food = entities.NewRandomFood(position)
			food.Type = entities.FoodTypePremium
			food.Color = rl.Blue
			food.Size = 5.0 + rl.GetRandomFloat32() * 3.0
			food.Value = constants.FoodBaseValue * 2 + int(rl.GetRandomFloat32() * 8)
			food.HasPhysics = rl.GetRandomFloat32() < 0.3
			
			if food.HasPhysics {
				angle := rl.GetRandomFloat32() * 2 * math.Pi
				speed := 15.0 + rl.GetRandomFloat32() * 20.0
				food.Velocity = rl.Vector2{
					X: float32(math.Cos(float64(angle))) * speed,
					Y: float32(math.Sin(float64(angle))) * speed,
				}
			}
		} else {
			// Create a basic food
			food = entities.NewFood(position)
		}
		
		fg.FoodEntities = append(fg.FoodEntities, food)
		quadtree.Root.Insert(food, food)
	}
	
	// Create a small particle effect to highlight the new cluster
	// This would be implemented if we had a central effect system
}

// RemoveFood removes a food entity from the list
func (fg *FoodGenerator) RemoveFood(food *entities.Food) {
	// Find and remove food from list
	for i, f := range fg.FoodEntities {
		if f == food {
			// Return to pool for potential reuse
			fg.FoodPool.Return(food)
			
			// Remove from list (order not important, so use swap-and-pop)
			lastIdx := len(fg.FoodEntities) - 1
			fg.FoodEntities[i] = fg.FoodEntities[lastIdx]
			fg.FoodEntities = fg.FoodEntities[:lastIdx]
			break
		}
	}
}

// GetFoodAt returns a food entity at the given position if one exists
func (fg *FoodGenerator) GetFoodAt(position rl.Vector2, radius float32) *entities.Food {
	for _, food := range fg.FoodEntities {
		dx := food.Position.X - position.X
		dy := food.Position.Y - position.Y
		distSq := dx*dx + dy*dy
		
		if distSq <= radius*radius {
			return food
		}
	}
	
	return nil
}

// GetNearbyFood returns all food entities within a radius of the position
func (fg *FoodGenerator) GetNearbyFood(position rl.Vector2, radius float32) []*entities.Food {
	result := make([]*entities.Food, 0)
	
	for _, food := range