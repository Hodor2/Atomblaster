// ui/controllers/controller.go
package controllers

import "atomblaster/ui/models"

// Controller is the interface for all UI controllers
type Controller interface {
    HandleInput() bool // Returns true if the screen should change
    SetModel(model models.Model)
}