// ui/views/view.go
package views

import "atomblaster/ui/models"

// View is the interface for all UI views
type View interface {
    Draw()
    SetModel(model models.Model)
}