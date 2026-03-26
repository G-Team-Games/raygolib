package raygolib

import rl "github.com/gen2brain/raylib-go/raylib"

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
