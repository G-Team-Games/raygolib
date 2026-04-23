package assets

import (
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Resolver == nil {
		t.Fatalf("expected default resolver")
	}
	if cfg.HotReload {
		t.Fatalf("expected hot reload disabled by default")
	}
	if cfg.WatchDebounce <= 0 {
		t.Fatalf("expected positive debounce")
	}
	if cfg.ThreadPolicy != ThreadPolicyStrict {
		t.Fatalf("expected strict thread policy default, got %v", cfg.ThreadPolicy)
	}
}

func TestNewManagerOptions(t *testing.T) {
	m, err := NewManager(
		WithSingleRoot("game_assets"),
		WithHotReload(true),
		WithWatchDebounce(250*time.Millisecond),
		WithThreadPolicy(ThreadPolicyQueueOnly),
		WithRule(KindRule{Kind: KindTexture, Include: []string{"**/*.png"}, Priority: 10}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := m.Config()
	if !cfg.HotReload {
		t.Fatalf("expected hot reload enabled")
	}
	if cfg.WatchDebounce != 250*time.Millisecond {
		t.Fatalf("unexpected debounce: %v", cfg.WatchDebounce)
	}
	if cfg.ThreadPolicy != ThreadPolicyQueueOnly {
		t.Fatalf("unexpected thread policy: %v", cfg.ThreadPolicy)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Kind != KindTexture {
		t.Fatalf("unexpected rules: %+v", cfg.Rules)
	}

	path, err := cfg.Resolver.Resolve(KindTexture, "player.png")
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	expected := filepath.Join("game_assets", "player.png")
	if path != expected {
		t.Fatalf("unexpected path. got=%q want=%q", path, expected)
	}
}

func TestOptionValidation(t *testing.T) {
	if _, err := NewManager(WithWatchDebounce(0)); err == nil {
		t.Fatalf("expected debounce validation error")
	}

	if _, err := NewManager(WithResolver(nil)); err == nil {
		t.Fatalf("expected nil resolver error")
	}

	if _, err := NewManager(WithThreadPolicy(ThreadPolicy(99))); err == nil {
		t.Fatalf("expected invalid thread policy error")
	}
}

func TestFixedDirsResolver(t *testing.T) {
	r := NewFixedDirsResolver("assets")

	path, err := r.Resolve(KindTexture, "ui/button.png")
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	want := filepath.Join("assets", "textures", "ui/button.png")
	if path != want {
		t.Fatalf("unexpected texture path. got=%q want=%q", path, want)
	}

	refs := r.KeysForPath(filepath.Join("assets", "textures", "ui/button.png"))
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Kind != KindTexture || refs[0].Key != "ui/button.png" {
		t.Fatalf("unexpected ref: %+v", refs[0])
	}

	roots := r.WatchRoots()
	if len(roots) == 0 {
		t.Fatalf("expected watch roots")
	}
}

func TestSingleRootResolver(t *testing.T) {
	r := NewSingleRootResolver("assets")

	path, err := r.Resolve(KindSound, "audio/shot.wav")
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	want := filepath.Join("assets", "audio/shot.wav")
	if path != want {
		t.Fatalf("unexpected path. got=%q want=%q", path, want)
	}

	refs := r.KeysForPath(filepath.Join("assets", "audio/shot.wav"))
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Key != "audio/shot.wav" {
		t.Fatalf("unexpected key: %q", refs[0].Key)
	}

	roots := r.WatchRoots()
	if len(roots) != 1 || roots[0] != "assets" {
		t.Fatalf("unexpected roots: %v", roots)
	}
}
