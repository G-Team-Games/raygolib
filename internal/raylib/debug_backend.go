package raylib

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DebugBackend interface {
	IsKeyPressed(key int32) bool
	DrawFPS(x, y int32)
	DrawRectangle(x, y, w, h int32, color color.RGBA)
}

type RlDebugBackend struct{}

func NewRlDebugBackend() *RlDebugBackend {
	return &RlDebugBackend{}
}

func (RlDebugBackend) IsKeyPressed(key int32) bool {
	return rl.IsKeyPressed(key)
}

func (RlDebugBackend) DrawFPS(x, y int32) {
	rl.DrawFPS(x, y)
}

func (RlDebugBackend) DrawRectangle(x, y, w, h int32, color rl.Color) {
	rl.DrawRectangle(x, y, w, h, color)
}
