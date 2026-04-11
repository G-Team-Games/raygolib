package raygolib

import (
	irl "github.com/G-Team-Games/raygolib/internal/raylib"
)

// Main game interface - used to implement game logic
type Game interface {
	Update(dt float32) error
	Draw()
}

// Middleware type - used to wrap game with additional functionality
// (like debug overlay)
type Middleware func(Game) Game

// Game initializer - used to configure and run game
type gameInitializer struct {
	game        Game
	cfg         *InitGameConfig
	middlewares []Middleware
	backend     irl.InitBackend
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
		game:    game,
		cfg:     baseCfg,
		backend: irl.NewRlInitBackend(),
	}
}

func (gi *gameInitializer) WithMiddleware(m Middleware) *gameInitializer {
	gi.middlewares = append(gi.middlewares, m)
	return gi
}

func (gi *gameInitializer) Run() error {
	gi.backend.InitWindow(int32(gi.cfg.ScreenWidth), int32(gi.cfg.ScreenHeight), gi.cfg.WindowTitle)
	defer gi.backend.CloseWindow()

	gi.backend.SetTargetFPS(int32(gi.cfg.TargetFPS))

	for _, m := range gi.middlewares {
		gi.game = m(gi.game)
	}

	for !gi.backend.WindowShouldClose() {
		dt := gi.backend.GetFrameTime()
		if err := gi.game.Update(dt); err != nil {
			return err
		}

		gi.backend.BeginDrawing()
		gi.game.Draw()
		gi.backend.EndDrawing()
	}

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
