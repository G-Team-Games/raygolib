package main

import (
	rgl "github.com/G-Team-Games/raygolib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Debug *rgl.DebugAPI
}

func (g *Game) Update(dt float32) error {
	return nil
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.White)
	rl.DrawText("Press F1 to toggle debug mode", 150, 60, 30, rl.Black)

	// using internal debug API
	g.Debug.Rect(50, 50, 50, 50, rl.Red)
	// using global debug API
	rgl.Debug().Rect(50, 150, 50, 50, rl.Green)
}

// Needed for access to internal debug api (used in middleware)
func (g *Game) SetDebug(d *rgl.DebugAPI) {
	g.Debug = d
}

func main() {
	// With default debug mode config
	// game := rgl.InitGame(&Game{}).WithMiddleware(rgl.DebugMiddleware())

	// but you can change this config if you want
	game := rgl.InitGame(&Game{}).WithMiddleware(
		rgl.DebugMiddlewareWithConfig(rgl.DebugConfig{StartEnabled: true, ToggleKey: rl.KeyF1, ShowFPS: true}),
	)
	if err := game.Run(); err != nil {
		panic(err)
	}
}
