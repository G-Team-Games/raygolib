package main

import (
	"fmt"
	"math"

	"github.com/G-Team-Games/raygolib"
	rga "github.com/G-Team-Games/raygolib/assets"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var modelPaths = []ModelName{
	_2barrel, Bomb, Cube, Cylinder, Ghost, Plane, Player, Room, Screw,
}

const (
	zoomStep   = 0.2
	minZoom    = 0.5
	rotateStep = 0.01
	pitchPad   = 0.1
)

type modelViewer struct {
	manager         *rga.Manager
	currentModelIdx int
	loadingStatus   string
	cameraAngle     float32
	cameraPitch     float32
	zoomLevel       float32
	camera          rl.Camera3D
}

func newModelViewer() *modelViewer {
	return &modelViewer{
		zoomLevel: 3,
		camera: rl.Camera3D{
			Position:   rl.NewVector3(2, 2, 2),
			Target:     rl.Vector3Zero(),
			Up:         rl.NewVector3(0, 1, 0),
			Fovy:       60,
			Projection: rl.CameraPerspective,
		},
	}
}

func (m *modelViewer) Init() error {
	m.manager = rga.NewManager()
	m.loadModel(modelPaths[m.currentModelIdx])
	return nil
}

func (m *modelViewer) Update(dt float32) error {
	m.manager.Tick()

	if rl.IsKeyPressed(rl.KeyTab) {
		// Loading can be done in goroutines safely, due to queue-based design of the manager.
		// This is not "loading in background" because of raylib operations locked on main
		// // thread, but loading asset in goroutine dont break the game
		go m.loadNextModel()
	}

	m.updateCamera()

	if rl.IsKeyPressed(rl.KeyR) {
		m.reloadAllModels()
	}

	return nil
}

func (m *modelViewer) loadNextModel() {
	m.loadingStatus = "Loading..."
	targetIdx := (m.currentModelIdx + 1) % len(modelPaths)
	path := modelPaths[targetIdx]
	if _, err := m.manager.GetModel(string(path)); err != nil {
		m.loadingStatus = fmt.Sprintf("Failed: %v", err)
		return
	}
	m.currentModelIdx = targetIdx
	m.loadingStatus = fmt.Sprintf("Loaded: %s", path)
}

func (m *modelViewer) updateCamera() {
	scroll := rl.GetMouseWheelMove()
	if scroll != 0 {
		m.zoomLevel -= scroll * zoomStep
		if m.zoomLevel < minZoom {
			m.zoomLevel = minZoom
		}
	}

	delta := rl.GetMouseDelta()
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		m.cameraAngle += delta.X * rotateStep
		m.cameraPitch += delta.Y * rotateStep
		m.cameraPitch = rl.Clamp(m.cameraPitch, -math.Pi/2+pitchPad, math.Pi/2-pitchPad)
	}

	m.camera.Position.X = float32(math.Cos(float64(m.cameraAngle))) * float32(math.Cos(float64(m.cameraPitch))) * m.zoomLevel
	m.camera.Position.Y = float32(math.Sin(float64(m.cameraPitch))) * m.zoomLevel
	m.camera.Position.Z = float32(math.Sin(float64(m.cameraAngle))) * float32(math.Cos(float64(m.cameraPitch))) * m.zoomLevel
}

func (m *modelViewer) Draw() {
	rl.ClearBackground(rl.NewColor(30, 30, 30, 255))

	rl.BeginMode3D(m.camera)
	currentPath := string(modelPaths[m.currentModelIdx])
	// manager.Model return a pointer to the model if it's loaded,
	// or nil if it's still loading or failed
	if current := m.manager.Model(currentPath); current != nil {
		rl.DrawModel(current.Data, rl.Vector3Zero(), 1.0, rl.White)
		rl.DrawModelWiresEx(current.Data, rl.Vector3Zero(), rl.Vector3Zero(), 1.0, rl.NewVector3(1, 1, 1), rl.Black)
	}
	rl.EndMode3D()

	m.drawHUD()
}

func (m *modelViewer) Close() error {
	m.manager.ClearAll()
	return nil
}

func (m *modelViewer) loadModel(path ModelName) {
	_, err := m.manager.GetModel(string(path))
	if err != nil {
		m.loadingStatus = fmt.Sprintf("Error: %v", err)
		return
	}
	m.loadingStatus = fmt.Sprintf("Loaded: %s", path)
}

func (m *modelViewer) reloadAllModels() {
	m.loadingStatus = "Reloading..."
	m.manager.ReloadAll(nil)
	m.loadModel(modelPaths[m.currentModelIdx])
	m.loadingStatus = "Reloaded"
}

func (m *modelViewer) drawHUD() {
	rl.DrawText(fmt.Sprintf("Model: %s", modelPaths[m.currentModelIdx]), 10, 10, 20, rl.White)
	rl.DrawText("Controls:", 10, 40, 16, rl.LightGray)
	rl.DrawText("  Right mouse drag - Orbit", 10, 60, 16, rl.LightGray)
	rl.DrawText("  Scroll - Zoom", 10, 80, 16, rl.LightGray)
	rl.DrawText("  Tab - Load next (background)", 10, 100, 16, rl.LightGray)
	rl.DrawText("  R - Reload all", 10, 120, 16, rl.LightGray)
	rl.DrawText(fmt.Sprintf("  [%d/%d]", m.currentModelIdx+1, len(modelPaths)), 10, 140, 16, rl.LightGray)

	statusColor := rl.Green
	if m.loadingStatus == "Reloading..." {
		statusColor = rl.Yellow
	}
	rl.DrawText(m.loadingStatus, 10, int32(rl.GetScreenHeight())-30, 20, statusColor)
}

func main() {
	game := newModelViewer()
	reloadFn := func() {
		game.reloadAllModels()
	}

	raygolib.InitGameWithConfig(game, &raygolib.InitGameConfig{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		WindowTitle:  "Model Viewer - Asset Demo",
	}).WithMiddleware(raygolib.DebugMiddlewareWithConfig(raygolib.DebugConfig{
		ToggleKey:    rl.KeyF1,
		HotReloadKey: rl.KeyR,
		OnHotReload:  reloadFn,
		StartEnabled: true,
	})).Run()
}
