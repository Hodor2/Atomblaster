// ui/interfaces.go
package ui

// Model is the interface for all UI screen data models
type Model interface {
    // Any methods that all models should implement
}

// View is the interface for all UI screen views
type View interface {
    Draw()
    SetModel(model Model)
}

// Controller is the interface for all UI screen controllers
type Controller interface {
    HandleInput() bool // Returns true if the screen should change
    SetModel(model Model)
}

// Screen represents a complete UI screen with MVC components
type Screen struct {
    Model      Model
    View       View
    Controller Controller
}

// NewScreen creates a new screen with the given MVC components
func NewScreen(model Model, view View, controller Controller) *Screen {
    view.SetModel(model)
    controller.SetModel(model)
    
    return &Screen{
        Model:      model,
        View:       view,
        Controller: controller,
    }
}

// Draw renders the screen
func (s *Screen) Draw() {
    s.View.Draw()
}

// Update processes input and returns true if the screen should change
func (s *Screen) Update() bool {
    return s.Controller.HandleInput()
}
