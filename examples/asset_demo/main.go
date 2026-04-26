package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/G-Team-Games/raygolib"
	rga "github.com/G-Team-Games/raygolib/assets"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	modelPaths = []ModelName{
		_2barrel, Bomb, Cube, Cylinder, Ghost, Plane, Player, Room, Screw,
	}
	currentModelIdx int
	manager         *rga.Manager
	currentModel    *rl.Model
	mu              sync.Mutex

	loadingStatus string
	loading       bool
	nextToLoad    int = -1

	preloadedModel *rl.Model
	preloadErr     error
	preloadDone    bool
	preloading     bool

	cameraAngle float32 = 0
	cameraPitch float32 = 0
	zoomLevel   float32 = 3
)

type modelViewer struct {
	debug *raygolib.DebugAPI
}

func (m *modelViewer) SetDebug(d *raygolib.DebugAPI) {
	m.debug = d
}

func (m *modelViewer) Init() error {
	manager = rga.NewManager()
	loadModel(modelPaths[currentModelIdx])
	return nil
}

func (m *modelViewer) Update(dt float32) error {
	manager.Tick()

	if preloading && preloadDone {
		preloading = false
		if preloadErr != nil {
			loadingStatus = fmt.Sprintf("Failed: %v", preloadErr)
		} else {
			currentModelIdx = nextToLoad
			mu.Lock()
			currentModel = preloadedModel
			mu.Unlock()
			loadingStatus = fmt.Sprintf("Loaded: %s", modelPaths[nextToLoad])
		}
		nextToLoad = -1
		preloadedModel = nil
		preloadErr = nil
		preloadDone = false
	}

	if rl.IsKeyPressed(rl.KeyTab) && !preloading {
		nextToLoad = (currentModelIdx + 1) % len(modelPaths)
		loadingStatus = fmt.Sprintf("Loading %s...", modelPaths[nextToLoad])
		preloading = true
		preloadDone = false

		go func() {
			time.Sleep(2 * time.Second)
			path := string(modelPaths[nextToLoad])
			res, err := manager.GetModel(path)
			if err != nil {
				preloadErr = err
			} else if res != nil {
				preloadedModel = &res.Data
			}
			preloadDone = true
		}()
	}

	updateCamera(dt)

	if rl.IsKeyPressed(rl.KeyR) {
		loadingStatus = "Reloading..."
		manager.ReloadAll(nil)
		loadModel(modelPaths[currentModelIdx])
		loadingStatus = "Reloaded"
	}

	return nil
}

func updateCamera(dt float32) {
	scroll := rl.GetMouseWheelMove()
	if scroll != 0 {
		zoomLevel -= scroll * 0.2
		if zoomLevel < 0.5 {
			zoomLevel = 0.5
		}
	}

	delta := rl.GetMouseDelta()
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		cameraAngle += delta.X * 0.01
		cameraPitch += delta.Y * 0.01
		cameraPitch = clamp(cameraPitch, -math.Pi/2+0.1, math.Pi/2-0.1)
	}

	camera.Position.X = float32(math.Cos(float64(cameraAngle))) * float32(math.Cos(float64(cameraPitch))) * zoomLevel
	camera.Position.Y = float32(math.Sin(float64(cameraPitch))) * zoomLevel
	camera.Position.Z = float32(math.Sin(float64(cameraAngle))) * float32(math.Cos(float64(cameraPitch))) * zoomLevel
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (m *modelViewer) Draw() {
	rl.ClearBackground(rl.NewColor(30, 30, 30, 255))

	rl.BeginMode3D(camera)
	if currentModel != nil {
		rl.DrawModel(*currentModel, rl.Vector3Zero(), 1.0, rl.White)
		rl.DrawModelWires(*currentModel, rl.Vector3Zero(), 1.0, rl.NewColor(80, 80, 80, 255))
	}
	rl.EndMode3D()

	drawHUD()
}

func (m *modelViewer) Close() error {
	manager.ClearAll()
	return nil
}

func loadModel(path ModelName) {
	res, err := manager.GetModel(string(path))
	if err != nil {
		loadingStatus = fmt.Sprintf("Error: %v", err)
		return
	}
	mu.Lock()
	currentModel = &res.Data
	mu.Unlock()
	loadingStatus = fmt.Sprintf("Loaded: %s", path)
}

var camera = rl.Camera3D{
	Position:   rl.NewVector3(2, 2, 2),
	Target:     rl.Vector3Zero(),
	Up:         rl.NewVector3(0, 1, 0),
	Fovy:       60,
	Projection: rl.CameraPerspective,
}

func drawHUD() {
	rl.DrawText(fmt.Sprintf("Model: %s", modelPaths[currentModelIdx]), 10, 10, 20, rl.White)
	rl.DrawText("Controls:", 10, 40, 16, rl.LightGray)
	rl.DrawText("  Right mouse drag - Orbit", 10, 60, 16, rl.LightGray)
	rl.DrawText("  Scroll - Zoom", 10, 80, 16, rl.LightGray)
	rl.DrawText("  Tab - Load next (background)", 10, 100, 16, rl.LightGray)
	rl.DrawText("  R - Reload all", 10, 120, 16, rl.LightGray)
	rl.DrawText(fmt.Sprintf("  [%d/%d]", currentModelIdx+1, len(modelPaths)), 10, 140, 16, rl.LightGray)

	statusColor := rl.Green
	if preloading || loadingStatus == "Reloading..." {
		statusColor = rl.Yellow
	}
	rl.DrawText(loadingStatus, 10, int32(rl.GetScreenHeight())-30, 20, statusColor)
}

func main() {
	reloadFn := func() {
		loadingStatus = "Reloading..."
		manager.ReloadAll(nil)
		loadModel(modelPaths[currentModelIdx])
		loadingStatus = "Reloaded"
	}
	game := &modelViewer{}
	raygolib.InitGameWithConfig(game, &raygolib.InitGameConfig{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		WindowTitle:  "Model Viewer - Asset Demo",
	}).WithMiddleware(raygolib.DebugMiddlewareWithConfig(raygolib.DebugConfig{
		ToggleKey:    rl.KeyF1,
		HotReloadKey: rl.KeyR,
		OnHotReload:  reloadFn,
	})).Run()
}
