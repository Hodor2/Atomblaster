// entities/bullet.go
// Package entities contains game entity implementations
package entities

import (
    "atomblaster/constants"

    rl "github.com/gen2brain/raylib-go/raylib"
)

// Bullet represents a projectile fired by the player
type Bullet struct {
    Pos      rl.Vector2
    Vel      rl.Vector2
    Radius   float32
    Lifetime float32
}

// NewBullet creates a new bullet at playerPos moving in direction.
// It normalizes the direction and applies the game’s BulletSpeed.
func NewBullet(playerPos rl.Vector2, direction rl.Vector2) Bullet {
    // Normalize direction vector
    length := rl.Vector2Length(direction)
    var normalizedDir rl.Vector2
    if length > 0 {
        normalizedDir = rl.Vector2{
            X: direction.X / length,
            Y: direction.Y / length,
        }
    } else {
        normalizedDir = rl.Vector2{X: 1, Y: 0} // Default right if no input
    }

    return Bullet{
        Pos:      playerPos,
        Vel:      rl.Vector2{X: normalizedDir.X * constants.BulletSpeed, Y: normalizedDir.Y * constants.BulletSpeed},
        Radius:   5,
        Lifetime: constants.BulletLifetime,
    }
}

// Update moves the bullet by dt and returns false when it’s gone
func (b *Bullet) Update(dt float32) bool {
    b.Pos.X += b.Vel.X * dt
    b.Pos.Y += b.Vel.Y * dt
    b.Lifetime -= dt

    if b.Lifetime <= 0 ||
        b.Pos.X < 0 || b.Pos.X > constants.ScreenWidth ||
        b.Pos.Y < 0 || b.Pos.Y > constants.ScreenHeight {
        return false
    }
    return true
}

// Draw renders the bullet. If a texture is provided it uses DrawTexturePro,
// otherwise it falls back to a circle + trail.
func (b *Bullet) Draw(sprite rl.Texture2D) {
    if sprite.ID > 0 {
        rl.DrawTexturePro(
            sprite,
            rl.Rectangle{X: 0, Y: 0, Width: float32(sprite.Width), Height: float32(sprite.Height)},
            rl.Rectangle{
                X:      b.Pos.X - b.Radius,
                Y:      b.Pos.Y - b.Radius,
                Width:  b.Radius * 2,
                Height: b.Radius * 2,
            },
            rl.Vector2{X: b.Radius, Y: b.Radius},
            0.0,
            rl.White,
        )
    } else {
        // Fallback to simple circle
        rl.DrawCircleV(b.Pos, b.Radius, rl.Yellow)

        // Draw a small trail behind the bullet
        normalizedVel := rl.Vector2{
            X: b.Vel.X / constants.BulletSpeed,
            Y: b.Vel.Y / constants.BulletSpeed,
        }
        const trailLength = float32(15.0)

        endX := b.Pos.X - normalizedVel.X*trailLength
        endY := b.Pos.Y - normalizedVel.Y*trailLength

        rl.DrawLine(
            int32(b.Pos.X),
            int32(b.Pos.Y),
            int32(endX),
            int32(endY),
            rl.Fade(rl.Orange, 0.7),
        )
    }
}

// CheckCollision returns true if this bullet overlaps the given atom.
func (b *Bullet) CheckCollision(atom Atom) bool {
    dx := b.Pos.X - atom.Pos.X
    dy := b.Pos.Y - atom.Pos.Y
    distance := rl.Vector2Length(rl.Vector2{X: dx, Y: dy})
    return distance < (b.Radius + atom.Radius)
}
