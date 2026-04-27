package rga

import "testing"

func TestKindDefaultDirKnownAndCustom(t *testing.T) {
	tests := []struct {
		kind AssetKind
		want string
	}{
		{KindModel, "models"},
		{KindTexture, "textures"},
		{KindImage, "images"},
		{KindSound, "audio"},
		{KindMusic, "audio"},
		{KindFont, "fonts"},
		{KindShader, "shaders"},
		{AssetKind("sprite"), "sprites"},
	}

	for _, tc := range tests {
		if got := tc.kind.DefaultDir(); got != tc.want {
			t.Fatalf("DefaultDir(%q) = %q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestKindPluralRules(t *testing.T) {
	tests := []struct {
		kind AssetKind
		want string
	}{
		{KindShader, "shaders"},
		{AssetKind("brush"), "brushes"},
		{AssetKind("enemy"), "enemies"},
		{AssetKind("toy"), "toys"},
	}

	for _, tc := range tests {
		if got := tc.kind.Plural(); got != tc.want {
			t.Fatalf("Plural(%q) = %q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestKindIsKnown(t *testing.T) {
	if !KindTexture.IsKnown() {
		t.Fatal("expected built-in kind to be known")
	}
	if AssetKind("custom_kind").IsKnown() {
		t.Fatal("expected custom kind to be unknown")
	}
}

func TestKindDefaultExtensions(t *testing.T) {
	if got := KindTexture.DefaultExtensions(); len(got) == 0 {
		t.Fatal("expected texture default extensions")
	}
	if got := AssetKind("custom_kind").DefaultExtensions(); got != nil {
		t.Fatal("expected nil extensions for unknown kind")
	}
}
