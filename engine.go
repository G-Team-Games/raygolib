package raygolib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Engine struct
type engine struct {
	game    Game
	modules []Module
}

type Module interface {
	Update(dt float32) error
	Draw()
}

type Option func(*config)

func WithDebugMode(debug bool) Option {
	return func(c *config) {
		c.DebugMode = debug
	}

}

func Run(game Game, config config) error {
	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, config.Title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(config.TargetFPS)

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
