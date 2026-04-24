package assetgen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	tmpDir := t.TempDir()

	texturesDir := filepath.Join(tmpDir, "textures")
	if err := os.MkdirAll(texturesDir, 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{
		filepath.Join(texturesDir, "player.png"),
		filepath.Join(texturesDir, "enemy.png"),
	}
	for _, f := range files {
		if err := os.WriteFile(f, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	outputPath := filepath.Join(tmpDir, "generated.go")

	cfg := Config{
		Root:    tmpDir,
		Output:  outputPath,
		Package: "gameassets",
		Kinds:   "texture",
		Naming:  "pascal",
	}

	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Read output failed: %v", err)
	}

	content := string(data)

	if !containsAny(content, "package gameassets") {
		t.Error("missing package declaration")
	}
	if !containsAny(content, "type TextureName string") {
		t.Error("missing type declaration")
	}
	if !containsAny(content, "Player") {
		t.Error("missing Player constant")
	}
	if !containsAny(content, "Enemy") {
		t.Error("missing Enemy constant")
	}
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if len(sub) > 0 && contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr)))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
