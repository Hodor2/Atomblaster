// game/state.go
// Package game contains shared helpers for game logic
package game

import (
    "atomblaster/constants"
    "atomblaster/entities"
    rl "github.com/gen2brain/raylib-go/raylib"
)

// getHelicopterRect returns a rectangle representing the helicopter's collision area
func getHelicopterRect(player entities.Player) rl.Rectangle {
    return rl.Rectangle{
        X:      player.Pos.X - constants.HelicopterWidth/2,
        Y:      player.Pos.Y - constants.HelicopterHeight/2,
        Width:  constants.HelicopterWidth,
        Height: constants.HelicopterHeight,
    }
}
