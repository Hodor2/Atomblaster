package util

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Distance calculates the Euclidean distance between two points
func Distance(a, b rl.Vector2) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// NormalizeVector returns a normalized vector (unit vector)
func NormalizeVector(v rl.Vector2) rl.Vector2 {
	mag := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
	if mag > 0 {
		return rl.Vector2{X: v.X / mag, Y: v.Y / mag}
	}
	return v
}

// ClampValue restricts a value to a specific range
func ClampValue(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp performs linear interpolation between two values
func Lerp(start, end, amount float32) float32 {
	return start + amount*(end-start)
}

// RandomRange returns a random float32 between min and max
// Fixed: Changed from 1000 to 999 to prevent exceeding max
func RandomRange(min, max float32) float32 {
	return min + (max-min)*float32(rl.GetRandomValue(0, 999))/999.0
}

// RandomVelocity creates a random velocity vector with the given speed
func RandomVelocity(speed float32) rl.Vector2 {
	angle := float32(rl.GetRandomValue(0, 360)) * rl.Pi / 180.0

	return rl.Vector2{
		X: float32(math.Cos(float64(angle))) * speed,
		Y: float32(math.Sin(float64(angle))) * speed,
	}
}

// PulseValue returns a value that oscillates between min and max over time
// Fixed the type mismatch between float64 and float32
func PulseValue(min, max float32, frequency float32) float32 {
	return min + (max-min)*(0.5+0.5*float32(math.Sin(float64(rl.GetTime())*float64(frequency))))
}
