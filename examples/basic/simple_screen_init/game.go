package main

import (
	rgl "github.com/G-Team-Games/raygolib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct{}

func (g *Game) Update(dt float32) error {
	return nil
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.White)
	rl.DrawText("Hello World!", 20, 20, 30, rl.Black)
}

func main() {
	game := &Game{}

	if err := rgl.Run(game, rgl.DefaultConfig); err != nil {
		panic(err)
	}
}
