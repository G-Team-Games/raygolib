package raygolib

import (
	"errors"
	"reflect"
	"testing"

	"github.com/G-Team-Games/raygolib/internal/testutils/mocks"
)

type fakeGameForGameTests struct {
	updateDts   []float32
	updateErr   error
	updateCalls int
	drawCalls   int
}

func (f *fakeGameForGameTests) Init() error {
	return nil
}

func (f *fakeGameForGameTests) Close() error {
	return nil
}

func (f *fakeGameForGameTests) Update(dt float32) error {
	f.updateCalls++
	f.updateDts = append(f.updateDts, dt)
	return f.updateErr
}

func (f *fakeGameForGameTests) Draw() {
	f.drawCalls++
}

type orderMiddleware struct {
	name string
	log  *[]string
	next Game
}

func (m *orderMiddleware) Init() error {
	*m.log = append(*m.log, m.name+"-init")
	return m.next.Init()
}

func (m *orderMiddleware) Close() error {
	*m.log = append(*m.log, m.name+"-close")
	return m.next.Close()
}

func (m *orderMiddleware) Update(dt float32) error {
	*m.log = append(*m.log, m.name+"-update")
	return m.next.Update(dt)
}

func (m *orderMiddleware) Draw() {
	*m.log = append(*m.log, m.name+"-draw")
	m.next.Draw()
}

type orderGame struct {
	log *[]string
}

func (g *orderGame) Init() error {
	*g.log = append(*g.log, "game-init")
	return nil
}

func (g *orderGame) Close() error {
	*g.log = append(*g.log, "game-close")
	return nil
}

func (g *orderGame) Update(dt float32) error {
	*g.log = append(*g.log, "game-update")
	return nil
}

func (g *orderGame) Draw() {
	*g.log = append(*g.log, "game-draw")
}

func makeOrderMiddleware(name string, log *[]string) Middleware {
	return func(next Game) Game {
		return &orderMiddleware{
			name: name,
			log:  log,
			next: next,
		}
	}
}

type windowCloseSeq struct {
	results []bool
	index   int
}

func (s *windowCloseSeq) Next() bool {
	if s.index >= len(s.results) {
		return true
	}
	result := s.results[s.index]
	s.index++
	return result
}

type frameTimeSeq struct {
	results []float32
	index   int
}

func (s *frameTimeSeq) Next() float32 {
	if s.index >= len(s.results) {
		return 0
	}
	result := s.results[s.index]
	s.index++
	return result
}

func TestInitGameSuite(t *testing.T) {
	t.Run("uses default config", func(t *testing.T) {
		game := &fakeGameForGameTests{}
		gi := InitGame(game)

		if gi.game != game {
			t.Fatalf("expected game to be stored in initializer")
		}
		if gi.cfg == nil {
			t.Fatalf("expected config to be initialized")
		}
		if gi.cfg.ScreenWidth != 800 || gi.cfg.ScreenHeight != 600 || gi.cfg.WindowTitle != "Game" || gi.cfg.TargetFPS != 60 {
			t.Fatalf("unexpected default config: %+v", gi.cfg)
		}
	})

	t.Run("merges custom config", func(t *testing.T) {
		game := &fakeGameForGameTests{}
		cfg := &InitGameConfig{
			ScreenWidth:  1024,
			WindowTitle:  "Custom Title",
			TargetFPS:    120,
			ScreenHeight: 0,
		}
		gi := InitGameWithConfig(game, cfg)

		if gi.cfg.ScreenWidth != 1024 {
			t.Fatalf("expected ScreenWidth=1024, got %d", gi.cfg.ScreenWidth)
		}
		if gi.cfg.ScreenHeight != 600 {
			t.Fatalf("expected ScreenHeight to default to 600, got %d", gi.cfg.ScreenHeight)
		}
		if gi.cfg.WindowTitle != "Custom Title" {
			t.Fatalf("expected WindowTitle to be overridden, got %q", gi.cfg.WindowTitle)
		}
		if gi.cfg.TargetFPS != 120 {
			t.Fatalf("expected TargetFPS=120, got %d", gi.cfg.TargetFPS)
		}
	})
}

func TestMergeInitGameConfigsSuite(t *testing.T) {
	t.Run("overrides non-zero values", func(t *testing.T) {
		base := defaultInitGameConfig()
		other := &InitGameConfig{
			ScreenWidth:  1280,
			ScreenHeight: 720,
			WindowTitle:  "Merged",
			TargetFPS:    144,
		}

		merged := mergeInitGameConfigs(base, other)

		if merged.ScreenWidth != 1280 || merged.ScreenHeight != 720 || merged.WindowTitle != "Merged" || merged.TargetFPS != 144 {
			t.Fatalf("merge result unexpected: %+v", merged)
		}
	})
}

func TestWithMiddlewareSuite(t *testing.T) {
	t.Run("applies in order", func(t *testing.T) {
		windowSeq := &windowCloseSeq{results: []bool{false, true}}
		frameSeq := &frameTimeSeq{results: []float32{0.016}}

		backend := &mocks.InitBackendMock{
			WindowShouldCloseFunc: windowSeq.Next,
			GetFrameTimeFunc:      frameSeq.Next,
		}

		log := []string{}
		game := &orderGame{log: &log}

		gi := InitGame(game).
			WithMiddleware(makeOrderMiddleware("m1", &log)).
			WithMiddleware(makeOrderMiddleware("m2", &log))
		gi.backend = backend

		if err := gi.Run(); err != nil {
			t.Fatalf("unexpected Run error: %v", err)
		}

		expected := []string{
			"m2-init",
			"m1-init",
			"game-init",
			"m2-update",
			"m1-update",
			"game-update",
			"m2-draw",
			"m1-draw",
			"game-draw",
			"m2-close",
			"m1-close",
			"game-close",
		}

		if !reflect.DeepEqual(log, expected) {
			t.Fatalf("unexpected middleware order. got=%v want=%v", log, expected)
		}
	})
}

func TestRunSuite(t *testing.T) {
	t.Run("calls game loop and backend", func(t *testing.T) {
		windowSeq := &windowCloseSeq{results: []bool{false, false, true}}
		frameSeq := &frameTimeSeq{results: []float32{0.1, 0.2}}

		var initWidth, initHeight int32
		var initTitle string
		var targetFPS int32

		backend := &mocks.InitBackendMock{
			InitWindowFunc: func(width, height int32, title string) {
				initWidth = width
				initHeight = height
				initTitle = title
			},
			SetTargetFPSFunc: func(fps int32) {
				targetFPS = fps
			},
			WindowShouldCloseFunc: windowSeq.Next,
			GetFrameTimeFunc:      frameSeq.Next,
		}

		game := &fakeGameForGameTests{}
		cfg := &InitGameConfig{
			ScreenWidth:  640,
			ScreenHeight: 480,
			WindowTitle:  "Test Game",
			TargetFPS:    30,
		}

		gi := InitGameWithConfig(game, cfg)
		gi.backend = backend

		if err := gi.Run(); err != nil {
			t.Fatalf("unexpected Run error: %v", err)
		}

		if backend.InitWindowCalls != 1 {
			t.Fatalf("expected InitWindow to be called once, got %d", backend.InitWindowCalls)
		}
		if initWidth != 640 || initHeight != 480 || initTitle != "Test Game" {
			t.Fatalf("unexpected InitWindow args: %dx%d %q", initWidth, initHeight, initTitle)
		}
		if backend.SetTargetFPSCalls != 1 || targetFPS != 30 {
			t.Fatalf("expected SetTargetFPS(30), got calls=%d fps=%d", backend.SetTargetFPSCalls, targetFPS)
		}
		if backend.BeginDrawingCalls != 2 || backend.EndDrawingCalls != 2 {
			t.Fatalf("expected Begin/EndDrawing to be called twice, got begin=%d end=%d", backend.BeginDrawingCalls, backend.EndDrawingCalls)
		}
		if backend.CloseWindowCalls != 1 {
			t.Fatalf("expected CloseWindow to be called once, got %d", backend.CloseWindowCalls)
		}

		if game.updateCalls != 2 || game.drawCalls != 2 {
			t.Fatalf("expected 2 Update and 2 Draw calls, got update=%d draw=%d", game.updateCalls, game.drawCalls)
		}
		if len(game.updateDts) != 2 || game.updateDts[0] != 0.1 || game.updateDts[1] != 0.2 {
			t.Fatalf("unexpected Update dt values: %v", game.updateDts)
		}
	})

	t.Run("returns update error and closes window", func(t *testing.T) {
		windowSeq := &windowCloseSeq{results: []bool{false}}
		frameSeq := &frameTimeSeq{results: []float32{0.25}}

		backend := &mocks.InitBackendMock{
			WindowShouldCloseFunc: windowSeq.Next,
			GetFrameTimeFunc:      frameSeq.Next,
		}

		expectedErr := errors.New("update failed")
		game := &fakeGameForGameTests{updateErr: expectedErr}

		gi := InitGame(game)
		gi.backend = backend

		err := gi.Run()
		if err == nil || err.Error() != expectedErr.Error() {
			t.Fatalf("expected Run to return update error, got %v", err)
		}

		if game.updateCalls != 1 {
			t.Fatalf("expected Update to be called once, got %d", game.updateCalls)
		}
		if game.drawCalls != 0 {
			t.Fatalf("expected Draw not to be called after update error, got %d", game.drawCalls)
		}
		if backend.BeginDrawingCalls != 0 || backend.EndDrawingCalls != 0 {
			t.Fatalf("expected no Begin/EndDrawing after update error, got begin=%d end=%d", backend.BeginDrawingCalls, backend.EndDrawingCalls)
		}
		if backend.CloseWindowCalls != 1 {
			t.Fatalf("expected CloseWindow to be called once, got %d", backend.CloseWindowCalls)
		}
	})
}
