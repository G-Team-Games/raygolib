package assets

import (
	"path/filepath"
	"testing"
)

type fakeResolver struct {
	paths map[Kind]string
}

func (r *fakeResolver) Resolve(kind Kind, key string) (string, error) {
	base := r.paths[kind]
	if base == "" {
		base = "assets"
	}
	return filepath.Join(base, key), nil
}

func (r *fakeResolver) KeysForPath(path string) []AssetRef {
	return nil
}

func (r *fakeResolver) WatchRoots() []string {
	return nil
}

func TestAssetManagerUsesProvidedResolver(t *testing.T) {
	r := &fakeResolver{paths: map[Kind]string{KindTexture: "my_textures"}}
	am := NewAssetManagerWithResolver(r)

	path, err := am.resolvePath(KindTexture, "player.png")
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}

	want := filepath.Join("my_textures", "player.png")
	if path != want {
		t.Fatalf("unexpected path. got=%q want=%q", path, want)
	}
}

func TestNewAssetManagerHasDefaultResolver(t *testing.T) {
	am := NewAssetManager()

	if am.resolver == nil {
		t.Fatalf("expected default resolver")
	}

	path, err := am.resolvePath(KindTexture, "a.png")
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}

	want := filepath.Join("assets", "a.png")
	if path != want {
		t.Fatalf("unexpected path. got=%q want=%q", path, want)
	}
}
