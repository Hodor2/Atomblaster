// ui/models/title_model.go
package models

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

// TitleModel contains data for the title screen
type TitleModel struct {
    Background   rl.Texture2D
    MenuOptions  []string
    SelectedItem int
}

// NewTitleModel creates a new title screen model
func NewTitleModel(background rl.Texture2D) *TitleModel {
    return &TitleModel{
        Background:   background,
        MenuOptions:  []string{"Start Game", "Instructions", "Exit"},
        SelectedItem: 0,
    }
}

// SelectNextItem moves the selection to the next menu item
func (m *TitleModel) SelectNextItem() {
    m.SelectedItem = (m.SelectedItem + 1) % len(m.MenuOptions)
}

// SelectPreviousItem moves the selection to the previous menu item
func (m *TitleModel) SelectPreviousItem() {
    m.SelectedItem = (m.SelectedItem - 1 + len(m.MenuOptions)) % len(m.MenuOptions)
}

// GetSelectedOption returns the currently selected menu option
func (m *TitleModel) GetSelectedOption() string {
    return m.MenuOptions[m.SelectedItem]
}
