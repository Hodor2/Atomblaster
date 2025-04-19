package util

import (
	"github.com/gen2brain/raylib-go/raylib"
	"math"
	"hash/fnv"
)

// Clamp constrains a value between min and max
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp performs linear interpolation between a and b by t
func Lerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

// Vector2Distance calculates the distance between two Vector2 points
func Vector2Distance(v1, v2 rl.Vector2) float32 {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Vector2Add adds two vectors
func Vector2Add(v1, v2 rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

// Vector2Subtract subtracts v2 from v1
func Vector2Subtract(v1, v2 rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

// Vector2Scale multiplies a vector by a scalar
func Vector2Scale(v rl.Vector2, scale float32) rl.Vector2 {
	return rl.Vector2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

// Vector2Length calculates the length of a vector
func Vector2Length(v rl.Vector2) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Vector2Normalize returns the normalized vector (unit vector)
func Vector2Normalize(v rl.Vector2) rl.Vector2 {
	length := Vector2Length(v)
	if length == 0 {
		return rl.Vector2{X: 0, Y: 0}
	}
	
	return rl.Vector2{
		X: v.X / length,
		Y: v.Y / length,
	}
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float32) float32 {
	return degrees * (math.Pi / 180.0)
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(radians float32) float32 {
	return radians * (180.0 / math.Pi)
}

// AngleBetweenVectors calculates the angle between two vectors in degrees
func AngleBetweenVectors(v1, v2 rl.Vector2) float32 {
	dot := v1.X*v2.X + v1.Y*v2.Y
	len1 := Vector2Length(v1)
	len2 := Vector2Length(v2)
	
	if len1 == 0 || len2 == 0 {
		return 0
	}
	
	angle := float32(math.Acos(float64(dot/(len1*len2))))
	return RadiansToDegrees(angle)
}

// Vector2Rotate rotates a vector by an angle in degrees
func Vector2Rotate(v rl.Vector2, degrees float32) rl.Vector2 {
	radians := DegreesToRadians(degrees)
	cosVal := float32(math.Cos(float64(radians)))
	sinVal := float32(math.Sin(float64(radians)))
	
	return rl.Vector2{
		X: v.X*cosVal - v.Y*sinVal,
		Y: v.X*sinVal + v.Y*cosVal,
	}
}

// LookAtRotation calculates the rotation (in degrees) to look at a target
func LookAtRotation(origin, target rl.Vector2) float32 {
	direction := Vector2Subtract(target, origin)
	return RadiansToDegrees(float32(math.Atan2(float64(direction.Y), float64(direction.X))))
}

// HashString creates a hash from a string
func HashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// RandomInRange returns a random float between min and max
func RandomInRange(min, max float32) float32 {
	return min + (max-min)*rl.GetRandomFloat32()
}

// RandomInCircle returns a random point inside a circle of radius
func RandomInCircle(center rl.Vector2, radius float32) rl.Vector2 {
	angle := RandomInRange(0, 2*math.Pi)
	distance := rl.GetRandomFloat32() * radius
	
	return rl.Vector2{
		X: center.X + float32(math.Cos(float64(angle))) * distance,
		Y: center.Y + float32(math.Sin(float64(angle))) * distance,
	}
}

// ObjectPool is a generic pool for reusing objects
type ObjectPool struct {
	objects     []interface{}
	createFunc  func() interface{}
	resetFunc   func(interface{})
}

// NewObjectPool creates a new object pool
func NewObjectPool(createFunc func() interface{}, resetFunc func(interface{}), initialSize int) *ObjectPool {
	pool := &ObjectPool{
		objects:    make([]interface{}, 0, initialSize),
		createFunc: createFunc,
		resetFunc:  resetFunc,
	}
	
	// Pre-allocate objects
	for i := 0; i < initialSize; i++ {
		pool.objects = append(pool.objects, createFunc())
	}
	
	return pool
}

// Get retrieves an object from the pool or creates a new one
func (p *ObjectPool) Get() interface{} {
	if len(p.objects) == 0 {
		return p.createFunc()
	}
	
	// Remove and return the last object
	lastIndex := len(p.objects) - 1
	obj := p.objects[lastIndex]
	p.objects = p.objects[:lastIndex]
	
	// Reset the object before returning
	if p.resetFunc != nil {
		p.resetFunc(obj)
	}
	
	return obj
}

// Return puts an object back into the pool
func (p *ObjectPool) Return(obj interface{}) {
	p.objects = append(p.objects, obj)
}
