// components/factory.go
package components

import (
    "atomblaster/constants"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// EntityFactory provides convenience functions for creating common game entities
type EntityFactory struct {
    registry *ComponentTypeRegistry
    manager  *EntityManager
    assets   struct {
        playerSprite   rl.Texture2D
        enemySprite    rl.Texture2D
        bulletSprite   rl.Texture2D
        powerUpSprites [3]rl.Texture2D
    }
}

// NewEntityFactory creates a new entity factory
func NewEntityFactory(
    registry *ComponentTypeRegistry,
    manager *EntityManager,
    playerSprite, enemySprite, bulletSprite rl.Texture2D,
    powerUpSprites [3]rl.Texture2D,
) *EntityFactory {
    factory := &EntityFactory{
        registry: registry,
        manager:  manager,
    }
    
    factory.assets.playerSprite = playerSprite
    factory.assets.enemySprite = enemySprite
    factory.assets.bulletSprite = bulletSprite
    factory.assets.powerUpSprites = powerUpSprites
    
    return factory
}

// CreatePlayer creates a player entity
func (f *EntityFactory) CreatePlayer(x, y float32, hasGun bool) EntityID {
    // Create player entity
    playerID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(playerID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(playerID, NewVelocity(0, 0, f.registry))
    f.manager.AddComponent(playerID, NewCircleCollider(constants.HelicopterWidth/2, f.registry))
    f.manager.AddComponent(playerID, NewSprite(f.assets.playerSprite, f.registry))
    f.manager.AddComponent(playerID, NewHealth(3, 10, f.registry))
    f.manager.AddComponent(playerID, NewTag(PlayerTag, f.registry))
    
    player := NewPlayer(300, f.registry)
    player.HasGun = hasGun
    f.manager.AddComponent(playerID, player)
    
    return playerID
}

// CreateAtom creates an enemy atom entity
func (f *EntityFactory) CreateAtom(x, y float32, velX, velY float32, atomType EnemyType, level int) EntityID {
    // Create atom entity
    atomID := f.manager.CreateEntity()
    
    // Get speed based on level and type
    baseSpeed := float32(100 + level*10)
    speed := baseSpeed + float32(rl.GetRandomValue(-20, 20))
    
    if atomType == FastAtom {
        speed *= 1.5
    } else if atomType == BigAtom {
        speed *= 0.7
    }
    
    // Add components
    f.manager.AddComponent(atomID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(atomID, NewVelocity(velX, velY, f.registry))
    f.manager.AddComponent(atomID, NewSprite(f.assets.enemySprite, f.registry))
    f.manager.AddComponent(atomID, NewTag(EnemyTag, f.registry))
    f.manager.AddComponent(atomID, NewEnemy(atomType, speed, f.registry))
    
    // Add health and collider based on type
    if atomType == NormalAtom {
        f.manager.AddComponent(atomID, NewCircleCollider(15, f.registry))
        f.manager.AddComponent(atomID, NewHealth(2, 2, f.registry))
    } else if atomType == FastAtom {
        f.manager.AddComponent(atomID, NewCircleCollider(12, f.registry))
        f.manager.AddComponent(atomID, NewHealth(1, 1, f.registry))
    } else if atomType == BigAtom {
        f.manager.AddComponent(atomID, NewCircleCollider(25, f.registry))
        f.manager.AddComponent(atomID, NewHealth(4, 4, f.registry))
    }
    
    return atomID
}

// CreateBoss creates a boss entity
func (f *EntityFactory) CreateBoss(x, y float32) EntityID {
    // Create boss entity
    bossID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(bossID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(bossID, NewVelocity(0, 0, f.registry))
    f.manager.AddComponent(bossID, NewRectangleCollider(80, 40, f.registry))
    f.manager.AddComponent(bossID, NewSprite(f.assets.playerSprite, f.registry)) // Using player sprite for simplicity
    f.manager.AddComponent(bossID, NewHealth(100, 100, f.registry))
    f.manager.AddComponent(bossID, NewTag(BossTag, f.registry))
    f.manager.AddComponent(bossID, NewEnemy(Boss, 200, f.registry))
    
    return bossID
}

// CreateBullet creates a bullet entity
func (f *EntityFactory) CreateBullet(x, y float32, velX, velY float32, isEnemyBullet bool) EntityID {
    // Create bullet entity
    bulletID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(bulletID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(bulletID, NewVelocity(velX, velY, f.registry))
    f.manager.AddComponent(bulletID, NewCircleCollider(5, f.registry))
    f.manager.AddComponent(bulletID, NewSprite(f.assets.bulletSprite, f.registry))
    f.manager.AddComponent(bulletID, NewTag(BulletTag, f.registry))
    f.manager.AddComponent(bulletID, NewLifetime(constants.BulletLifetime, f.registry))
    
    return bulletID
}

// CreateScientist creates a scientist entity
func (f *EntityFactory) CreateScientist(x, y float32) EntityID {
    // Create scientist entity
    scientistID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(scientistID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(scientistID, NewVelocity(0, 0, f.registry))
    f.manager.AddComponent(scientistID, NewCircleCollider(15, f.registry))
    f.manager.AddComponent(scientistID, NewTag(ScientistTag, f.registry))
    f.manager.AddComponent(scientistID, NewScientist(f.registry))
    
    return scientistID
}

// CreateRescueZone creates a rescue zone entity
func (f *EntityFactory) CreateRescueZone(x, y, width, height float32) EntityID {
    // Create rescue zone entity
    rescueZoneID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(rescueZoneID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(rescueZoneID, NewRectangleCollider(width, height, f.registry))
    f.manager.AddComponent(rescueZoneID, NewTag(RescueZoneTag, f.registry))
    
    return rescueZoneID
}

// CreateDoor creates a door entity (level exit)
func (f *EntityFactory) CreateDoor(x, y, width, height float32) EntityID {
    // Create door entity
    doorID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(doorID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(doorID, NewRectangleCollider(width, height, f.registry))
    f.manager.AddComponent(doorID, NewTag(DoorTag, f.registry))
    
    return doorID
}

// CreatePowerUp creates a power-up entity
func (f *EntityFactory) CreatePowerUp(x, y float32, powerUpType PowerUpType) EntityID {
    // Create power-up entity
    powerUpID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(powerUpID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(powerUpID, NewCircleCollider(15, f.registry))
    
    // Select the appropriate sprite based on power-up type
    spriteIndex := 0
    switch powerUpType {
    case PowerUpGun:
        spriteIndex = 0
    case PowerUpHealth:
        spriteIndex = 1
    case PowerUpSpeed:
        spriteIndex = 2
    }
    
    f.manager.AddComponent(powerUpID, NewSprite(f.assets.powerUpSprites[spriteIndex], f.registry))
    f.manager.AddComponent(powerUpID, NewTag(PowerUpTag, f.registry))
    
    // Create power-up component with appropriate values
    var duration float32 = 0
    var value float32 = 0
    
    switch powerUpType {
    case PowerUpGun:
        // Gun is permanent
        duration = 0
        value = 0
    case PowerUpHealth:
        // Health is instant
        duration = 0
        value = 1
    case PowerUpSpeed:
        // Speed is permanent
        duration = 0
        value = 50
    }
    
    f.manager.AddComponent(powerUpID, NewPowerUp(powerUpType, duration, value, f.registry))
    
    return powerUpID
}

// CreateParticle creates a particle entity
func (f *EntityFactory) CreateParticle(x, y float32, velX, velY float32, lifetime float32, size float32, color rl.Color) EntityID {
    // Create particle entity
    particleID := f.manager.CreateEntity()
    
    // Add components
    f.manager.AddComponent(particleID, NewPosition(x, y, f.registry))
    f.manager.AddComponent(particleID, NewVelocity(velX, velY, f.registry))
    f.manager.AddComponent(particleID, NewLifetime(lifetime, f.registry))
    
    // For a real implementation, we would add a dedicated Particle component
    // But for simplicity in this example, we're skipping that
    
    return particleID
}