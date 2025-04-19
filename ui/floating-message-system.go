// ui/floating-message-system.go
package ui

import (
    "math"

    rl "github.com/gen2brain/raylib-go/raylib"
)

type FloatingMessage struct {
    Text      string
    Pos       rl.Vector2
    StartTime float32
    Duration  float32
    Alpha     float32
    Scale     float32
    Velocity  rl.Vector2
}

type FloatingMessageSystem struct {
    Messages []FloatingMessage
}

func NewFloatingMessageSystem() *FloatingMessageSystem {
    return &FloatingMessageSystem{Messages: []FloatingMessage{}}
}

func (ms *FloatingMessageSystem) AddMessage(text string, pos rl.Vector2, duration float32) {
    msg := FloatingMessage{
        Text:      text,
        Pos:       pos,
        StartTime: float32(rl.GetTime()),
        Duration:  duration,
        Alpha:     1.0,
        Scale:     1.0,
        Velocity:  rl.Vector2{X: 0, Y: -30},
    }
    ms.Messages = append(ms.Messages, msg)
}

func (ms *FloatingMessageSystem) Update() {
    currentTime := float32(rl.GetTime())
    dt := rl.GetFrameTime()

    var active []FloatingMessage
    for _, msg := range ms.Messages {
        elapsed := currentTime - msg.StartTime
        if elapsed >= msg.Duration {
            continue
        }
        remaining := msg.Duration - elapsed
        alpha := float32(1.0)
        if remaining < 0.5 {
            alpha = remaining * 2.0
        }
        progress := elapsed / msg.Duration
        newPos := rl.Vector2{
            X: msg.Pos.X + msg.Velocity.X*dt,
            Y: msg.Pos.Y + msg.Velocity.Y*dt,
        }
        active = append(active, FloatingMessage{
            Text:      msg.Text,
            Pos:       newPos,
            StartTime: msg.StartTime,
            Duration:  msg.Duration,
            Alpha:     alpha,
            Scale:     1 + float32(math.Sin(float64(progress*math.Pi)))*0.2,
            Velocity:  msg.Velocity,
        })
    }
    ms.Messages = active
}

func (ms *FloatingMessageSystem) Draw() {
    currentTime := float32(rl.GetTime())
    for _, msg := range ms.Messages {
        elapsed := currentTime - msg.StartTime
        if elapsed >= msg.Duration {
            continue
        }
        fontSize := int32(24 * msg.Scale)
        textWidth := rl.MeasureText(msg.Text, fontSize)
        col := rl.White
        col.A = uint8(255 * msg.Alpha)
        shadow := rl.Black
        shadow.A = uint8(128 * msg.Alpha)

        // Draw drop shadow
        rl.DrawText(
            msg.Text,
            int32(msg.Pos.X)-textWidth/2+2,
            int32(msg.Pos.Y)+2,
            fontSize,
            shadow,
        )
        // Draw main text
        rl.DrawText(
            msg.Text,
            int32(msg.Pos.X)-textWidth/2,
            int32(msg.Pos.Y),
            fontSize,
            col,
        )
    }
}
