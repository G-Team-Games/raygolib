package raygolib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Run(game Game, middlewares ...Middleware) error {
	rl.InitWindow(800, 600, "Test")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for _, m := range middlewares {
		game = m(game)
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		if err := update(game, dt); err != nil {
			return err
		}
		draw(game)
	}

	return nil
}

func update(game Game, dt float32) error {
	return game.Update(dt)
}

func draw(game Game) {
	rl.BeginDrawing()
	defer rl.EndDrawing()
	game.Draw()
}
