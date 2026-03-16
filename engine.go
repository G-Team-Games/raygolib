package raygolib

import rl "github.com/gen2brain/raylib-go/raylib"

func Run(game Game, config Config) error {
	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, config.Title)
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {

		dt := rl.GetFrameTime()

		if err := Step(game, dt); err != nil {
			return err
		}
	}

	return nil
}
func Step(game Game, dt float32) error {
	if err := game.Update(dt); err != nil {
		return err
	}

	game.Render()
	return nil
}
