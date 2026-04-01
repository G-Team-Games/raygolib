package raygolib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Main game interface - use it to implement your game logic
type Game interface {
	Update(dt float32) error
	Draw()
}

// Middleware type - use it to wrap your game with additional functionality 
// (like debug overlay)
type Middleware func(Game) Game

// Game initializer - use it to configure and run your game
type gameInitializer struct {
	game        Game
	cfg         *InitGameConfig
	middlewares []Middleware
}

// InitGame initializes the game with default configuration
func InitGame(game Game) *gameInitializer {
	return InitGameWithConfig(game, nil)
}

func InitGameWithConfig(game Game, cfg *InitGameConfig) *gameInitializer {
	baseCfg := defaultInitGameConfig()
	if cfg != nil {
		baseCfg = mergeInitGameConfigs(baseCfg, cfg)
	}

	return &gameInitializer{
		game: game,
		cfg:  baseCfg,
	}
}

func (gi *gameInitializer) WithMiddleware(m Middleware) *gameInitializer {
	gi.middlewares = append(gi.middlewares, m)
	return gi
}

func (gi *gameInitializer) Run() error {
	rl.InitWindow(int32(gi.cfg.ScreenWidth), int32(gi.cfg.ScreenHeight), gi.cfg.WindowTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(int32(gi.cfg.TargetFPS))

	for _, m := range gi.middlewares {
		gi.game = m(gi.game)
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		if err := gi.game.Update(dt); err != nil {
			return err
		}

		rl.BeginDrawing()
		gi.game.Draw()
		rl.EndDrawing()
	}

	defer rl.CloseWindow()
	return nil
}

// Game config struct
type InitGameConfig struct {
	ScreenWidth  int
	ScreenHeight int
	WindowTitle  string
	TargetFPS    int
}

func defaultInitGameConfig() *InitGameConfig {
	return &InitGameConfig{
		ScreenWidth:  800,
		ScreenHeight: 600,
		WindowTitle:  "Game",
		TargetFPS:    60,
	}
}

func mergeInitGameConfigs(base, other *InitGameConfig) *InitGameConfig {
	if other.ScreenWidth != 0 {
		base.ScreenWidth = other.ScreenWidth
	}
	if other.ScreenHeight != 0 {
		base.ScreenHeight = other.ScreenHeight
	}
	if other.WindowTitle != "" {
		base.WindowTitle = other.WindowTitle
	}
	if other.TargetFPS != 0 {
		base.TargetFPS = other.TargetFPS
	}

	return base
}
