package util

import (
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	MaxEntitiesPerNode = 10
	MaxDepth           = 5
)

// QuadtreeItem is an interface for objects that can be stored in quadtree
type QuadtreeItem interface {
	GetBounds() rl.Rectangle
}

// QuadtreeNode represents a node in the quadtree
type QuadtreeNode struct {
	Bounds    rl.Rectangle
	Depth     int
	Entities  []interface{}
	Children  [4]*QuadtreeNode
	Divided   bool
}

// Quadtree is a spatial partitioning data structure for efficient collision detection
type Quadtree struct {
	Root      *QuadtreeNode
	WorldSize rl.Rectangle
}

// NewQuadtree creates a new quadtree with specified world bounds
func NewQuadtree(worldBounds rl.Rectangle) *Quadtree {
	return &Quadtree{
		Root: &QuadtreeNode{
			Bounds:   worldBounds,
			Depth:    0,
			Entities: make([]interface{}, 0, MaxEntitiesPerNode),
			Divided:  false,
		},
		WorldSize: worldBounds,
	}
}

// Clear empties the quadtree
func (q *Quadtree) Clear() {
	q.Root = &QuadtreeNode{
		Bounds:   q.WorldSize,
		Depth:    0,
		Entities: make([]interface{}, 0, MaxEntitiesPerNode),
		Divided:  false,
	}
}

// Insert adds an entity into the quadtree
func (node *QuadtreeNode) Insert(entity interface{}, item QuadtreeItem) bool {
	// Check if entity is within this node's bounds
	entityBounds := item.GetBounds()
	if !CheckRectangleOverlap(node.Bounds, entityBounds) {
		return false
	}
	
	// If not at capacity or max depth, add to this node
	if len(node.Entities) < MaxEntitiesPerNode || node.Depth >= MaxDepth {
		node.Entities = append(node.Entities, entity)
		return true
	}
	
	// Otherwise, subdivide and add to children
	if !node.Divided {
		node.Subdivide()
	}
	
	// Try to insert into children
	inserted := false
	for i := 0; i < 4; i++ {
		if node.Children[i] != nil && node.Children[i].Insert(entity, item) {
			inserted = true
		}
	}
	
	// If doesn't fit in any child (due to bounds checking), add to this node
	if !inserted {
		node.Entities = append(node.Entities, entity)
	}
	
	return true
}

// Subdivide creates four child nodes for this node
func (node *QuadtreeNode) Subdivide() {
	halfWidth := node.Bounds.Width / 2
	halfHeight := node.Bounds.Height / 2
	x := node.Bounds.X
	y := node.Bounds.Y
	nextDepth := node.Depth + 1
	
	// Create four children nodes
	node.Children[0] = &QuadtreeNode{ // Top-left
		Bounds:   rl.Rectangle{X: x, Y: y, Width: halfWidth, Height: halfHeight},
		Depth:    nextDepth,
		Entities: make([]interface{}, 0, MaxEntitiesPerNode),
	}
	
	node.Children[1] = &QuadtreeNode{ // Top-right
		Bounds:   rl.Rectangle{X: x + halfWidth, Y: y, Width: halfWidth, Height: halfHeight},
		Depth:    nextDepth,
		Entities: make([]interface{}, 0, MaxEntitiesPerNode),
	}
	
	node.Children[2] = &QuadtreeNode{ // Bottom-left
		Bounds:   rl.Rectangle{X: x, Y: y + halfHeight, Width: halfWidth, Height: halfHeight},
		Depth:    nextDepth,
		Entities: make([]interface{}, 0, MaxEntitiesPerNode),
	}
	
	node.Children[3] = &QuadtreeNode{ // Bottom-right
		Bounds:   rl.Rectangle{X: x + halfWidth, Y: y + halfHeight, Width: halfWidth, Height: halfHeight},
		Depth:    nextDepth,
		Entities: make([]interface{}, 0, MaxEntitiesPerNode),
	}
	
	node.Divided = true
	
	// Redistribute existing entities among children
	entitiesToRedistribute := node.Entities
	node.Entities = make([]interface{}, 0, MaxEntitiesPerNode)
	
	for _, entity := range entitiesToRedistribute {
		if item, ok := entity.(QuadtreeItem); ok {
			inserted := false
			for i := 0; i < 4; i++ {
				if node.Children[i].Insert(entity, item) {
					inserted = true
					break
				}
			}
			
			// If it doesn't fit in any child, keep it in this node
			if !inserted {
				node.Entities = append(node.Entities, entity)
			}
		} else {
			// If entity doesn't implement QuadtreeItem, keep it here
			node.Entities = append(node.Entities, entity)
		}
	}
}

// Query finds all entities that might collide with the given rectangle
func (node *QuadtreeNode) Query(bounds rl.Rectangle, results *[]interface{}) {
	// If this node doesn't overlap with the query bounds, return
	if !CheckRectangleOverlap(node.Bounds, bounds) {
		return
	}
	
	// Add all entities in this node to the results
	for _, entity := range node.Entities {
		*results = append(*results, entity)
	}
	
	// If this node is divided, query children
	if node.Divided {
		for i := 0; i < 4; i++ {
			if node.Children[i] != nil {
				node.Children[i].Query(bounds, results)
			}
		}
	}
}

// DrawDebug draws the quadtree structure for debugging
func (node *QuadtreeNode) DrawDebug(color rl.Color) {
	// Draw this node's boundary
	rl.DrawRectangleLinesEx(node.Bounds, 1, color)
	
	// Draw children if divided
	if node.Divided {
		for i := 0; i < 4; i++ {
			if node.Children[i] != nil {
				node.Children[i].DrawDebug(color)
			}
		}
	}
}

// CheckRectangleOverlap determines if two rectangles overlap
func CheckRectangleOverlap(a, b rl.Rectangle) bool {
	return !(a.X >= b.X+b.Width || a.X+a.Width <= b.X || 
			 a.Y >= b.Y+b.Height || a.Y+a.Height <= b.Y)
}
