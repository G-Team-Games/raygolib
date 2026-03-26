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

	g.Debug.Rect(50, 50, 50, 50)
}

// It has to be
func (g *Game) SetDebug(d *rgl.DebugAPI) {
	g.Debug = d
}

func main() {
	game := &Game{}

	if err := rgl.Run(game, rgl.DebugMiddleware()); err != nil {
		panic(err)
	}
}
