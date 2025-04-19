package ui

import (
	"github.com/gen2brain/raylib-go/raylib"
	"atomblaster/constants"
	"atomblaster/entities"
)

// Minimap is a UI element that shows a small map of the game world
type Minimap struct {
	Position    rl.Vector2
	Size        rl.Vector2
	BorderColor rl.Color
	MapColor    rl.Color
	PlayerColor rl.Color
	FoodColor   rl.Color
	PowerUpColor rl.Color
	Scale       rl.Vector2
	Opacity     uint8
}

// NewMinimap creates a new minimap UI element
func NewMinimap(worldWidth, worldHeight float32) *Minimap {
	return &Minimap{
		Position:     rl.Vector2{X: constants.ScreenWidth - 150, Y: 20},
		Size:         rl.Vector2{X: 130, Y: 130},
		BorderColor:  rl.Gray,
		MapColor:     rl.Color{R: 20, G: 20, B: 20, A: 200},
		PlayerColor:  rl.Red,
		FoodColor:    rl.Green,
		PowerUpColor: rl.Purple,
		Scale: rl.Vector2{
			X: 130.0 / worldWidth,
			Y: 130.0 / worldHeight,
		},
		Opacity: 180, // Semi-transparent
	}
}

// Draw renders the minimap
func (m *Minimap) Draw(player *entities.Player, camera rl.Camera2D, foodSystem interface{}, powerUpSystem interface{}) {
	// Draw minimap background
	rl.DrawRectangleV(m.Position, m.Size, m.MapColor)
	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      m.Position.X,
		Y:      m.Position.Y,
		Width:  m.Size.X,
		Height: m.Size.Y,
	}, 2, m.BorderColor)
	
	// Draw world boundaries
	borderSize := float32(2)
	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      m.Position.X,
		Y:      m.Position.Y,
		Width:  m.Size.X,
		Height: m.Size.Y,
	}, borderSize, rl.DarkGray)
	
	// Draw food entities (if available)
	if foodSystem != nil {
		if fs, ok := foodSystem.(*systems.FoodGenerator); ok {
			foodColor := rl.Color{
				R: m.FoodColor.R,
				G: m.FoodColor.G,
				B: m.FoodColor.B,
				A: m.Opacity,
			}
			
			// Draw each food item as a tiny dot
			for _, food := range fs.FoodEntities {
				// Convert world position to minimap position
				miniPos := rl.Vector2{
					X: m.Position.X + food.Position.X * m.Scale.X,
					Y: m.Position.Y + food.Position.Y * m.Scale.Y,
				}
				
				// Use different colors based on food type
				color := foodColor
				size := float32(1)
				
				if food.Type == entities.FoodTypePremium {
					color = rl.Blue
					color.A = m.Opacity
					size = 1.5
				} else if food.Type == entities.FoodTypeRare {
					color = rl.Gold
					color.A = m.Opacity
					size = 2
				}
				
				rl.DrawCircleV(miniPos, size, color)
			}
		}
	}
	
	// Draw power-up entities (if available)
	if powerUpSystem != nil {
		if ps, ok := powerUpSystem.(*systems.PowerUpManager); ok {
			powerUpColor := rl.Color{
				R: m.PowerUpColor.R,
				G: m.PowerUpColor.G,
				B: m.PowerUpColor.B,
				A: m.Opacity,
			}
			
			// Draw each power-up as a small triangle
			for _, powerUp := range ps.PowerUps {
				if !powerUp.Active {
					continue
				}
				
				// Convert world position to minimap position
				miniPos := rl.Vector2{
					X: m.Position.X + powerUp.Position.X * m.Scale.X,
					Y: m.Position.Y + powerUp.Position.Y * m.Scale.Y,
				}
				
				// Use different colors based on power-up type
				color := powerUpColor
				size := float32(2.5)
				
				switch powerUp.Type {
				case entities.PowerUpMagnet:
					color = rl.Purple
				case entities.PowerUpSpeed:
					color = rl.Red
				case entities.PowerUpShield:
					color = rl.SkyBlue
				case entities.PowerUpSizeBoost:
					color = rl.Orange
				}
				
				color.A = m.Opacity
				
				// Draw a small triangle
				rl.DrawPoly(miniPos, 3, size, 0, color)
			}
		}
	}
	
	// Draw player position
	playerMiniPos := rl.Vector2{
		X: m.Position.X + player.Position.X * m.Scale.X,
		Y: m.Position.Y + player.Position.Y * m.Scale.Y,
	}
	
	// Draw slightly larger dot for player
	playerColor := rl.Color{
		R: m.PlayerColor.R,
		G: m.PlayerColor.G,
		B: m.PlayerColor.B,
		A: 255, // Player is fully opaque
	}
	rl.DrawCircleV(playerMiniPos, 3, playerColor)
	
	// Draw viewport rectangle
	viewportSize := rl.Vector2{
		X: constants.ScreenWidth * m.Scale.X / camera.Zoom,
		Y: constants.ScreenHeight * m.Scale.Y / camera.Zoom,
	}
	viewportPos := rl.Vector2{
		X: m.Position.X + camera.Target.X * m.Scale.X,
		Y: m.Position.Y + camera.Target.Y * m.Scale.Y,
	}
	
	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      viewportPos.X,
		Y:      viewportPos.Y,
		Width:  viewportSize.X,
		Height: viewportSize.Y,
	}, 1, rl.White)
}

// SetPosition changes the position of the minimap
func (m *Minimap) SetPosition(position rl.Vector2) {
	m.Position = position
}

// SetSize changes the size of the minimap
func (m *Minimap) SetSize(size rl.Vector2) {
	m.Size = size
	
	// Recalculate scale
	m.Scale.X = m.Size.X / constants.WorldWidth
	m.Scale.Y = m.Size.Y / constants.WorldHeight
}

// SetOpacity changes the transparency of the minimap
func (m *Minimap) SetOpacity(opacity uint8) {
	m.Opacity = opacity
}
