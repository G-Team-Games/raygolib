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

// type DebugModule struct {
// 	enabled bool
// }

// func NewDebugModule() *DebugModule {
// 	return &DebugModule{enabled: true}
// }

// func (d *DebugModule) Update(dt float32) error {
// 	// np. toggle
// 	return nil
// }

// func (d *DebugModule) Draw() {
// 	if !d.enabled {
// 		return
// 	}

// 	// DrawFPS itd.
// }
