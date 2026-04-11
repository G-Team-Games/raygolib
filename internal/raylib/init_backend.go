package raylib

import rl "github.com/gen2brain/raylib-go/raylib"

type InitBackend interface {
	InitWindow(width, height int32, title string)
	CloseWindow()
	SetTargetFPS(fps int32)
	WindowShouldClose() bool
	GetFrameTime() float32
	BeginDrawing()
	EndDrawing()
}

type RlInitBackend struct{}

func NewRlInitBackend() *RlInitBackend {
	return &RlInitBackend{}
}


func (RlInitBackend) InitWindow(width, height int32, title string) {
	rl.InitWindow(width, height, title)
}

func (RlInitBackend) CloseWindow() {
	rl.CloseWindow()
}

func (RlInitBackend) SetTargetFPS(fps int32) {
	rl.SetTargetFPS(fps)
}

func (RlInitBackend) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}

func (RlInitBackend) GetFrameTime() float32 {
	return rl.GetFrameTime()
}

func (RlInitBackend) BeginDrawing() {
	rl.BeginDrawing()
}

func (RlInitBackend) EndDrawing() {
	rl.EndDrawing()
}
