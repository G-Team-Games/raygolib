package raygolib

import (
	"image/color"
	"testing"

	"github.com/G-Team-Games/raygolib/internal/testutils/mocks"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type fakeGame struct {
	debug       *DebugAPI
	updateCalls int
	drawCalls   int
}

func (f *fakeGame) Init() error {
	return nil
}

func (f *fakeGame) Close() error {
	return nil
}

func (f *fakeGame) Update(dt float32) error {
	f.updateCalls++
	return nil
}

func (f *fakeGame) Draw() {
	f.drawCalls++
}

func (f *fakeGame) SetDebug(d *DebugAPI) {
	f.debug = d
}

func withDebugInstance(t *testing.T) func() {
	t.Helper()
	prev := debugInstance
	return func() {
		debugInstance = prev
	}
}

func TestDebugAPISuite(t *testing.T) {
	t.Run("Rect disabled does not queue", func(t *testing.T) {
		d := &DebugAPI{enabled: false}

		d.Rect(10, 10, 20, 20, color.RGBA{R: 255, A: 255})

		if got := len(d.draws); got != 0 {
			t.Fatalf("expected no draws queued when disabled, got %d", got)
		}
	})

	t.Run("Rect enabled queues", func(t *testing.T) {
		d := &DebugAPI{enabled: true}

		d.Rect(10, 10, 20, 20, color.RGBA{G: 255, A: 255})

		if got := len(d.draws); got != 1 {
			t.Fatalf("expected 1 draw queued when enabled, got %d", got)
		}
	})

	t.Run("Flush clears queue and draws", func(t *testing.T) {
		backend := &mocks.DebugBackendMock{}
		d := &DebugAPI{
			enabled: true,
			backend: backend,
		}

		d.Rect(0, 0, 1, 1, color.RGBA{B: 255, A: 255})
		d.Rect(1, 1, 2, 2, color.RGBA{R: 128, G: 128, A: 255})

		d.Flush()

		if backend.DrawRectangleCalls != 2 {
			t.Fatalf("expected 2 rectangle draws, got %d", backend.DrawRectangleCalls)
		}
		if got := len(d.draws); got != 0 {
			t.Fatalf("expected draws cleared after flush, got %d", got)
		}
	})

	t.Run("Toggle flips enabled state", func(t *testing.T) {
		d := &DebugAPI{enabled: false}

		d.Toggle()
		if !d.enabled {
			t.Fatalf("expected enabled after toggle")
		}

		d.Toggle()
		if d.enabled {
			t.Fatalf("expected disabled after second toggle")
		}
	})
}

func TestDebugMiddlewareSuite(t *testing.T) {
	t.Run("injects and configures", func(t *testing.T) {
		restore := withDebugInstance(t)
		defer restore()

		cfg := DebugConfig{
			StartEnabled: true,
			ToggleKey:    rl.KeyF2,
			ShowFPS:      false,
		}

		game := &fakeGame{}
		wrapped := DebugMiddlewareWithConfig(cfg)(game)

		if game.debug == nil {
			t.Fatalf("expected DebugAPI to be injected into DebugAware game")
		}
		if Debug() != game.debug {
			t.Fatalf("expected Debug() to return injected debug instance")
		}

		dw, ok := wrapped.(*debugWrapper)
		if !ok {
			t.Fatalf("expected wrapped game to be *debugWrapper, got %T", wrapped)
		}

		if dw.debug != game.debug {
			t.Fatalf("expected wrapper debug to match injected debug instance")
		}
		if dw.debug.enabled != cfg.StartEnabled {
			t.Fatalf("expected StartEnabled=%v, got %v", cfg.StartEnabled, dw.debug.enabled)
		}
		if dw.toggleKey != cfg.ToggleKey {
			t.Fatalf("expected ToggleKey=%v, got %v", cfg.ToggleKey, dw.toggleKey)
		}
		if dw.showFPS != cfg.ShowFPS {
			t.Fatalf("expected ShowFPS=%v, got %v", cfg.ShowFPS, dw.showFPS)
		}
	})

	t.Run("Update toggles on key press", func(t *testing.T) {
		restore := withDebugInstance(t)
		defer restore()

		backend := &mocks.DebugBackendMock{
			IsKeyPressedFunc: func(key int32) bool {
				return key == rl.KeyF1
			},
		}

		game := &fakeGame{}
		cfg := DebugConfig{
			StartEnabled: false,
			ToggleKey:    rl.KeyF1,
			ShowFPS:      true,
		}
		wrapped := DebugMiddlewareWithConfig(cfg)(game)

		dw := wrapped.(*debugWrapper)
		dw.backend = backend

		if dw.debug.enabled {
			t.Fatalf("expected debug disabled initially")
		}

		if err := dw.Update(0.016); err != nil {
			t.Fatalf("unexpected error in Update: %v", err)
		}

		if !dw.debug.enabled {
			t.Fatalf("expected debug enabled after key press")
		}
	})

	t.Run("Draw calls FPS when enabled", func(t *testing.T) {
		restore := withDebugInstance(t)
		defer restore()

		backend := &mocks.DebugBackendMock{}

		game := &fakeGame{}
		cfg := DebugConfig{
			StartEnabled: true,
			ToggleKey:    rl.KeyF1,
			ShowFPS:      true,
		}
		wrapped := DebugMiddlewareWithConfig(cfg)(game)

		dw := wrapped.(*debugWrapper)
		dw.backend = backend

		dw.Draw()

		if backend.DrawFPSCalls != 1 {
			t.Fatalf("expected DrawFPS to be called once, got %d", backend.DrawFPSCalls)
		}
	})

	t.Run("Draw skips FPS when disabled", func(t *testing.T) {
		restore := withDebugInstance(t)
		defer restore()

		backend := &mocks.DebugBackendMock{}

		game := &fakeGame{}
		cfg := DebugConfig{
			StartEnabled: false,
			ToggleKey:    rl.KeyF1,
			ShowFPS:      true,
		}
		wrapped := DebugMiddlewareWithConfig(cfg)(game)

		dw := wrapped.(*debugWrapper)
		dw.backend = backend

		dw.Draw()

		if backend.DrawFPSCalls != 0 {
			t.Fatalf("expected DrawFPS not to be called when debug disabled, got %d", backend.DrawFPSCalls)
		}
		if len(dw.debug.draws) != 0 {
			t.Fatalf("expected debug draws to be cleared when disabled, got %d", len(dw.debug.draws))
		}
	})

	t.Run("Draw skips FPS when hidden", func(t *testing.T) {
		restore := withDebugInstance(t)
		defer restore()

		backend := &mocks.DebugBackendMock{}

		game := &fakeGame{}
		cfg := DebugConfig{
			StartEnabled: true,
			ToggleKey:    rl.KeyF1,
			ShowFPS:      false,
		}
		wrapped := DebugMiddlewareWithConfig(cfg)(game)

		dw := wrapped.(*debugWrapper)
		dw.backend = backend

		dw.Draw()

		if backend.DrawFPSCalls != 0 {
			t.Fatalf("expected DrawFPS not to be called when ShowFPS is false, got %d", backend.DrawFPSCalls)
		}
	})
}
