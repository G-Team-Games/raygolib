package raygolib

import rl "github.com/gen2brain/raylib-go/raylib"


type DebugAware interface {
	SetDebug(*DebugAPI)
}

type debugWrapper struct {
	next  Game
	debug *DebugAPI
}

func (d *debugWrapper) Update(dt float32) error {
	if rl.IsKeyPressed(rl.KeyF1) {
		d.debug.Toggle()
	}

	return d.next.Update(dt)
}

func (d *debugWrapper) Draw() {
	d.next.Draw()

	if !d.debug.enabled {
		d.debug.draws = nil
		return
	}

	rl.DrawFPS(10, 10)

	d.debug.Flush()
}

type DebugAPI struct {
	enabled bool
	draws   []func()
}

func (d *DebugAPI) Rect(x, y, w, h float32) {
	if !d.enabled {
		return
	}

	d.draws = append(d.draws, func() {
		rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.Red)
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
	return func(next Game) Game {
		debug := &DebugAPI{
			enabled: true,
		}

		if g, ok := next.(DebugAware); ok {
			g.SetDebug(debug)
		}

		return &debugWrapper{
			next:  next,
			debug: debug,
		}
	}
}
