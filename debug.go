package raygolib

import (
	"image/color"

	irl "github.com/G-Team-Games/raygolib/internal/raylib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type DebugAware interface {
	SetDebug(*DebugAPI)
}

type DebugConfig struct {
	StartEnabled bool
	ToggleKey    int32
	ShowFPS      bool
	HotReloadKey int32
	OnHotReload  func()
}

func defaultDebugConfig() DebugConfig {
	return DebugConfig{
		StartEnabled: false,
		ToggleKey:    rl.KeyF1,
		ShowFPS:      true,
	}
}

type debugWrapper struct {
	next         Game
	debug        *DebugAPI
	toggleKey    int32
	showFPS      bool
	hotReloadKey int32
	onHotReload  func()
	backend      irl.DebugBackend
}

func (d *debugWrapper) Init() error {
	return d.next.Init()
}

func (d *debugWrapper) Close() error {
	return d.next.Close()
}

func (d *debugWrapper) Unwrap() Game {
	return d.next
}

func (d *debugWrapper) Update(dt float32) error {
	if d.toggleKey != 0 && d.backend.IsKeyPressed(d.toggleKey) {
		d.debug.Toggle()
	}

	if d.hotReloadKey != 0 && d.onHotReload != nil && d.backend.IsKeyPressed(d.hotReloadKey) {
		d.onHotReload()
	}

	return d.next.Update(dt)
}

func (d *debugWrapper) Draw() {
	d.next.Draw()

	if !d.debug.enabled {
		d.debug.draws = nil
		return
	}

	if d.showFPS {
		d.backend.DrawFPS(10, 10)
	}

	d.debug.Flush()
}

// DEBUG API --------------------------------------------------------------- //

type DebugAPI struct {
	enabled bool
	draws   []func()
	backend irl.DebugBackend
}

var debugInstance *DebugAPI // nil until middleware is applied

func Debug() *DebugAPI {
	if debugInstance == nil {
		panic("Apply debug middleware first to use debug API")
	}
	return debugInstance
}

func (d *DebugAPI) Rect(x, y, w, h float32, color color.RGBA) {
	if !d.enabled {
		return
	}

	d.draws = append(d.draws, func() {
		d.backend.DrawRectangle(int32(x), int32(y), int32(w), int32(h), color)
	})
}

func (d *DebugAPI) Flush() {
	for _, draw := range d.draws {
		draw()
	}
	d.draws = nil
}

func (d *DebugAPI) Toggle() {
	d.enabled = !d.enabled
}

func DebugMiddleware() Middleware {
	return DebugMiddlewareWithConfig(defaultDebugConfig())
}

func DebugMiddlewareWithConfig(cfg DebugConfig) Middleware {
	return func(next Game) Game {
		debugBackend := irl.NewRlDebugBackend()

		debug := &DebugAPI{
			enabled: cfg.StartEnabled,
			backend: debugBackend,
		}
		debugInstance = debug

		// Walk the middleware chain to find the DebugAware implementation
		current := next
		for current != nil {
			if g, ok := current.(DebugAware); ok {
				g.SetDebug(debug)
				break
			}
			if w, ok := current.(Wrapper); ok {
				current = w.Unwrap()
			} else {
				current = nil
			}
		}

		return &debugWrapper{
			next:         next,
			debug:        debug,
			toggleKey:    cfg.ToggleKey,
			showFPS:      cfg.ShowFPS,
			hotReloadKey: cfg.HotReloadKey,
			onHotReload:  cfg.OnHotReload,
			backend:      debugBackend,
		}
	}
}
