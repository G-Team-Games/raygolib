package rgl

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateBasic(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "player.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "enemy.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(tmpDir, "generated.go")
	cfg := Config{
		Root:       tmpDir,
		Output:     out,
		Package:    "gameassets",
		SingleRoot: true,
		Recursive:  true,
	}

	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Read output failed: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "package gameassets") {
		t.Fatal("missing package declaration")
	}
	if !strings.Contains(content, "type TextureName string") {
		t.Fatal("missing texture type declaration")
	}
	if !strings.Contains(content, "Player") || !strings.Contains(content, "Enemy") {
		t.Fatal("missing generated texture constants")
	}
}

func TestGenerateWithExtensions(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "font.ttf"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "shader.glsl"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "data.txt"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(tmpDir, "generated.go")
	cfg := Config{
		Root:    tmpDir,
		Output:  out,
		Package: "gameassets",
		Kinds: []KindConfig{
			{Kind: KindFont, ScanRoot: true, Extensions: []string{".ttf"}},
			{Kind: KindShader, ScanRoot: true, Extensions: []string{".glsl"}},
		},
	}

	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Read output failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Font") || !strings.Contains(content, "Shader") {
		t.Fatal("missing generated constants for extension-mapped kinds")
	}
	if strings.Contains(content, "data.txt") {
		t.Fatal("unexpected non-matching file included")
	}
}

func TestFlatMode(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "player.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.vert"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(tmpDir, "generated.go")
	cfg := Config{
		Root:         tmpDir,
		Output:       out,
		Package:      "gameassets",
		SingleRoot:   true,
		FlatMode:     true,
		FlatTypeName: "AnyAsset",
		Kinds: []KindConfig{
			{Kind: KindTexture, ScanRoot: true, Extensions: []string{".png"}},
			{Kind: KindShader, ScanRoot: true, Extensions: []string{".vert"}},
		},
	}

	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Read output failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "type AnyAsset string") {
		t.Fatal("missing flat mode type declaration")
	}
	if strings.Count(content, "type ") != 1 {
		t.Fatal("expected a single generated type in flat mode")
	}
}

func TestDryRunWritesToWriter(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "player.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cfg := Config{
		Root:         tmpDir,
		Package:      "gameassets",
		SingleRoot:   true,
		DryRun:       true,
		DryRunWriter: &buf,
	}

	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if !strings.Contains(buf.String(), "Player") {
		t.Fatal("dry-run did not write generated output")
	}
}

func TestNameCollisionIncludesBothPaths(t *testing.T) {
	tmpDir := t.TempDir()
	a := filepath.Join(tmpDir, "a")
	b := filepath.Join(tmpDir, "b")
	if err := os.MkdirAll(a, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(a, "player.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "player.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Root:         tmpDir,
		Package:      "gameassets",
		SingleRoot:   true,
		Recursive:    true,
		DryRun:       true,
		DryRunWriter: &bytes.Buffer{},
	}

	err := Generate(cfg)
	if err == nil {
		t.Fatal("expected name collision error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "constant name collision") || !strings.Contains(msg, "a/player.png") || !strings.Contains(msg, "b/player.png") {
		t.Fatalf("unexpected collision error: %v", err)
	}
}

func TestGenerateDryRunAndDump(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "x.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := GenerateDryRun(Config{Root: tmpDir, SingleRoot: true, Package: "pkgx"}, &out); err != nil {
		t.Fatalf("GenerateDryRun failed: %v", err)
	}
	if !strings.Contains(out.String(), "package pkgx") {
		t.Fatalf("unexpected dry run output: %q", out.String())
	}

	var dump bytes.Buffer
	Config{Root: "r", Output: "o", Package: "p", Naming: "pascal", SingleRoot: true, Recursive: true, FlatMode: true, DryRun: true}.Dump(&dump)
	if !strings.Contains(dump.String(), "Root:") || !strings.Contains(dump.String(), "DryRun:") {
		t.Fatalf("unexpected dump output: %q", dump.String())
	}
}

func TestGenerateRejectsInvalidNaming(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "x.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(Config{
		Root:       tmpDir,
		Naming:     "invalid_style",
		SingleRoot: true,
		DryRun:     true,
	})
	if err == nil || !strings.Contains(err.Error(), "invalid naming style") {
		t.Fatalf("expected invalid naming error, got: %v", err)
	}
}

func TestGenerateConfigFileParseError(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "assetgen.json")
	if err := os.WriteFile(cfgFile, []byte("{invalid_json"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(Config{ConfigFile: cfgFile})
	if err == nil || !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got: %v", err)
	}
}

func TestGenerateTemplateFileReadError(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "x.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(Config{
		Root:         tmpDir,
		SingleRoot:   true,
		TemplateFile: filepath.Join(tmpDir, "missing.tmpl"),
		DryRun:       true,
	})
	if err == nil || !strings.Contains(err.Error(), "read template") {
		t.Fatalf("expected read template error, got: %v", err)
	}
}

func TestGenerateTemplateParseError(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "x.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	tmpl := filepath.Join(tmpDir, "bad.tmpl")
	if err := os.WriteFile(tmpl, []byte("{{ if"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(Config{
		Root:         tmpDir,
		SingleRoot:   true,
		TemplateFile: tmpl,
		DryRun:       true,
	})
	if err == nil || !strings.Contains(err.Error(), "parse template") {
		t.Fatalf("expected parse template error, got: %v", err)
	}
}



func TestHelperFunctions(t *testing.T) {
	if got := extensionsToGlob([]string{"png", ".JPG"}); got != "*.png|*.jpg" {
		t.Fatalf("unexpected extensionsToGlob output: %q", got)
	}
	if got := capitalize("hello"); got != "Hello" {
		t.Fatalf("unexpected capitalize output: %q", got)
	}
	if got := capitalize("h"); got != "H" {
		t.Fatalf("unexpected capitalize single char: %q", got)
	}
	if got := isVowel('a'); !got {
		t.Fatalf("unexpected isVowel for 'a': %v", got)
	}
	if got := isVowel('E'); !got {
		t.Fatalf("unexpected isVowel for 'E': %v", got)
	}
	if got := isVowel('x'); got {
		t.Fatalf("unexpected isVowel for 'x': %v", got)
	}
}

func TestNamerAllStylesAndPrefix(t *testing.T) {
	tests := []struct {
		style string
		name  string
		want  string
	}{
		{"pascal", "my-file.png", "MyFile"},
		{"camel", "my-file.png", "myFile"},
		{"snake", "my-file.png", "my_file"},
		{"upper_snake", "my-file.png", "MY_FILE"},
	}

	for _, tc := range tests {
		n := newNamer(tc.style, "")
		if got := n.name(tc.name); got != tc.want {
			t.Fatalf("name(%q, %q) = %q, want %q", tc.style, tc.name, got, tc.want)
		}
	}

	if got := newNamer("pascal", "pre").name("1-file.png"); got != "pre__1File" {
		t.Fatalf("unexpected prefixed/sanitized name: %q", got)
	}

	if got := capitalize("abc"); got != "Abc" {
		t.Fatalf("unexpected capitalize: %q", got)
	}
	if got := capitalize(""); got != "" {
		t.Fatalf("unexpected capitalize for empty string: %q", got)
	}
}

func TestReadFileConfigBranches(t *testing.T) {
	if _, err := loadFileConfig(""); err == nil {
		t.Fatal("expected loadFileConfig empty path error")
	}

	if _, err := loadFileConfig("/definitely/missing/config.json"); err == nil {
		t.Fatal("expected loadFileConfig missing path error")
	}
}

func TestNormalizeKindsBranches(t *testing.T) {
	kinds := []KindConfig{{Kind: KindTexture}, {Kind: AssetKind("custom")}}
	got := normalizeKinds(kinds, Config{SingleRoot: false, Glob: "*.forced"})

	if got[0].Type == "" || got[0].Plural == "" {
		t.Fatal("expected auto-filled type/plural")
	}
	if got[0].Glob == "" {
		t.Fatal("expected known kind to receive default-extension glob")
	}
	if got[1].Glob != "*.forced" {
		t.Fatalf("expected cfg glob fallback for custom kind, got %q", got[1].Glob)
	}

	got2 := normalizeKinds([]KindConfig{{Kind: KindTexture, Dir: "textures"}}, Config{SingleRoot: true})
	if !got2[0].ScanRoot || got2[0].Dir != "" {
		t.Fatal("expected single-root normalization to enable ScanRoot and clear Dir")
	}
}



func TestWarnUnknownKindsVerbose(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	orig := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	warnUnknownKinds([]KindConfig{{Kind: AssetKind("mykind")}})
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "custom kind") {
		t.Fatalf("expected warning output, got: %q", buf.String())
	}
}

func TestLoadFileConfigReturnsProperConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "assetgen.json")

	jsonCfg := `{
		"root": "custom_root",
		"output": "custom_out.go",
		"package": "custompkg",
		"naming": "snake",
		"single_root": true,
		"recursive": true,
		"flat_mode": true,
		"flat_type_name": "CustomAsset",
		"verbose": true
	}`
	if err := os.WriteFile(cfgFile, []byte(jsonCfg), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadFileConfig(cfgFile)
	if err != nil {
		t.Fatalf("loadFileConfig failed: %v", err)
	}

	if cfg.Root != "custom_root" {
		t.Fatalf("Root: got %q, want %q", cfg.Root, "custom_root")
	}
	if cfg.Output != "custom_out.go" {
		t.Fatalf("Output: got %q, want %q", cfg.Output, "custom_out.go")
	}
	if cfg.Package != "custompkg" {
		t.Fatalf("Package: got %q, want %q", cfg.Package, "custompkg")
	}
	if cfg.Naming != "snake" {
		t.Fatalf("Naming: got %q, want %q", cfg.Naming, "snake")
	}
	if !cfg.SingleRoot {
		t.Fatal("SingleRoot: expected true")
	}
	if !cfg.Recursive {
		t.Fatal("Recursive: expected true")
	}
	if !cfg.FlatMode {
		t.Fatal("FlatMode: expected true")
	}
	if cfg.FlatTypeName != "CustomAsset" {
		t.Fatalf("FlatTypeName: got %q, want %q", cfg.FlatTypeName, "CustomAsset")
	}
	if !cfg.Verbose {
		t.Fatal("Verbose: expected true")
	}
}

func TestLoadFileConfigWithKinds(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "assetgen.json")

	jsonCfg := `{
		"kinds": [
			{"kind": "texture", "dir": "textures", "extensions": [".png", ".jpg"]},
			{"kind": "shader", "dir": "shaders", "glob": "*.glsl"}
		]
	}`
	if err := os.WriteFile(cfgFile, []byte(jsonCfg), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadFileConfig(cfgFile)
	if err != nil {
		t.Fatalf("loadFileConfig failed: %v", err)
	}

	if len(cfg.Kinds) != 2 {
		t.Fatalf("Kinds count: got %d, want 2", len(cfg.Kinds))
	}
	if cfg.Kinds[0].Kind != KindTexture || cfg.Kinds[0].Dir != "textures" {
		t.Fatalf("Kinds[0]: got %+v", cfg.Kinds[0])
	}
	if cfg.Kinds[1].Kind != KindShader || cfg.Kinds[1].Glob != "*.glsl" {
		t.Fatalf("Kinds[1]: got %+v", cfg.Kinds[1])
	}
}

func TestMergeWithFile(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		fileCfg  Config
		check    func(t *testing.T, got Config)
	}{
		{
			name:    "programmatic string fields win",
			cfg:     Config{Root: "prog_root", Package: "progpkg"},
			fileCfg: Config{Root: "file_root", Package: "filepkg", Output: "file_out"},
			check: func(t *testing.T, got Config) {
				if got.Root != "prog_root" {
					t.Errorf("Root: got %q, want %q", got.Root, "prog_root")
				}
				if got.Package != "progpkg" {
					t.Errorf("Package: got %q, want %q", got.Package, "progpkg")
				}
				if got.Output != "file_out" {
					t.Errorf("Output: got %q, want %q", got.Output, "file_out")
				}
			},
		},
		{
			name:    "file fills zero-value strings",
			cfg:     Config{Root: "", Package: "", Naming: ""},
			fileCfg: Config{Root: "file_root", Package: "filepkg", Naming: "camel", Prefix: "pre"},
			check: func(t *testing.T, got Config) {
				if got.Root != "file_root" {
					t.Errorf("Root: got %q, want %q", got.Root, "file_root")
				}
				if got.Package != "filepkg" {
					t.Errorf("Package: got %q, want %q", got.Package, "filepkg")
				}
				if got.Naming != "camel" {
					t.Errorf("Naming: got %q, want %q", got.Naming, "camel")
				}
				if got.Prefix != "pre" {
					t.Errorf("Prefix: got %q, want %q", got.Prefix, "pre")
				}
			},
		},
		{
			name:    "boolean OR semantics",
			cfg:     Config{SingleRoot: true, Recursive: false},
			fileCfg: Config{SingleRoot: false, Recursive: true, FlatMode: true},
			check: func(t *testing.T, got Config) {
				if !got.SingleRoot {
					t.Error("SingleRoot: expected true (cfg OR file)")
				}
				if !got.Recursive {
					t.Error("Recursive: expected true (file OR cfg)")
				}
				if !got.FlatMode {
					t.Error("FlatMode: expected true (file)")
				}
			},
		},
		{
			name:    "programmatic bools preserve when file is false",
			cfg:     Config{Verbose: true},
			fileCfg: Config{Verbose: false},
			check: func(t *testing.T, got Config) {
				if !got.Verbose {
					t.Error("Verbose: expected true (cfg)")
				}
			},
		},
		{
			name:    "kinds from file when cfg has no kinds",
			cfg:     Config{Kinds: nil},
			fileCfg: Config{Kinds: []KindConfig{{Kind: KindTexture, Dir: "textures"}}},
			check: func(t *testing.T, got Config) {
				if len(got.Kinds) != 1 {
					t.Errorf("Kinds: got %d, want 1", len(got.Kinds))
				}
				if got.Kinds[0].Kind != KindTexture {
					t.Errorf("Kinds[0].Kind: got %v, want KindTexture", got.Kinds[0].Kind)
				}
			},
		},
		{
			name:    "cfg kinds win over file kinds",
			cfg:     Config{Kinds: []KindConfig{{Kind: KindShader}}},
			fileCfg: Config{Kinds: []KindConfig{{Kind: KindTexture}}},
			check: func(t *testing.T, got Config) {
				if len(got.Kinds) != 1 {
					t.Errorf("Kinds: got %d, want 1", len(got.Kinds))
				}
				if got.Kinds[0].Kind != KindShader {
					t.Errorf("Kinds[0].Kind: got %v, want KindShader", got.Kinds[0].Kind)
				}
			},
		},
		{
			name:    "includes/excludes merge",
			cfg:     Config{Includes: []string{"*.png"}, Excludes: []string{".git"}},
			fileCfg: Config{Includes: []string{"*.jpg"}, Excludes: []string{"*.tmp"}},
			check: func(t *testing.T, got Config) {
				if len(got.Includes) != 1 || got.Includes[0] != "*.png" {
					t.Errorf("Includes: got %v, want [*.png]", got.Includes)
				}
				if len(got.Excludes) != 1 || got.Excludes[0] != ".git" {
					t.Errorf("Excludes: got %v, want [.git]", got.Excludes)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeWithFile(tc.cfg, tc.fileCfg)
			tc.check(t, got)
		})
	}
}

func TestGenerateMergesConfigWithFileConfig(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "player.png"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	cfgFile := filepath.Join(tmpDir, "assetgen.json")
	jsonCfg := `{
		"root": "` + tmpDir + `",
		"output": "` + filepath.Join(tmpDir, "out1.go") + `",
		"package": "filepkg",
		"naming": "camel",
		"single_root": true
	}`
	if err := os.WriteFile(cfgFile, []byte(jsonCfg), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(Config{
		ConfigFile: cfgFile,
		Package:    "progpkg",
		DryRun:     true,
	})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var buf bytes.Buffer
	err = Generate(Config{
		ConfigFile: cfgFile,
		Package:    "progpkg",
		DryRun:     true,
		DryRunWriter: &buf,
	})
	if err != nil {
		t.Fatalf("Generate with dry-run failed: %v", err)
	}
	if !strings.Contains(buf.String(), "package progpkg") {
		t.Fatalf("programmatic Package should win, got: %s", buf.String())
	}
	if strings.Contains(buf.String(), "filepkg") {
		t.Fatal("file Package should not appear when programmatic is set")
	}
}

func TestGenerateUsesFileConfigWhenProgrammaticIsZero(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "font.ttf"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	cfgFile := filepath.Join(tmpDir, "assetgen.json")
	jsonCfg := `{
		"root": "` + tmpDir + `",
		"package": "fileconfigpkg",
		"naming": "snake",
		"single_root": true,
		"kinds": [
			{"kind": "font", "extensions": [".ttf"]}
		]
	}`
	if err := os.WriteFile(cfgFile, []byte(jsonCfg), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := Generate(Config{
		ConfigFile:   cfgFile,
		DryRun:       true,
		DryRunWriter: &buf,
	})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if !strings.Contains(buf.String(), "package fileconfigpkg") {
		t.Fatalf("expected file Package, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "player_font") && !strings.Contains(buf.String(), "Font") {
		t.Fatalf("expected font asset, got: %s", buf.String())
	}
}

func TestApplyStringDefaults(t *testing.T) {
	cfg := applyStringDefaults(Config{})
	if cfg.Root != "assets" {
		t.Fatalf("Root: got %q, want %q", cfg.Root, "assets")
	}
	if cfg.Output != defaultOutputFile {
		t.Fatalf("Output: got %q, want %q", cfg.Output, defaultOutputFile)
	}
	if cfg.Package != "assets" {
		t.Fatalf("Package: got %q, want %q", cfg.Package, "assets")
	}
	if cfg.Naming != "pascal" {
		t.Fatalf("Naming: got %q, want %q", cfg.Naming, "pascal")
	}
	if cfg.FlatTypeName != "AssetName" {
		t.Fatalf("FlatTypeName: got %q, want %q", cfg.FlatTypeName, "AssetName")
	}

	cfg = applyStringDefaults(Config{Root: "custom", Output: "out.go", Package: "mypkg", Naming: "snake", FlatTypeName: "MyAsset"})
	if cfg.Root != "custom" {
		t.Fatalf("Root: got %q, want %q", cfg.Root, "custom")
	}
	if cfg.Output != "out.go" {
		t.Fatalf("Output: got %q, want %q", cfg.Output, "out.go")
	}
	if cfg.Package != "mypkg" {
		t.Fatalf("Package: got %q, want %q", cfg.Package, "mypkg")
	}
	if cfg.Naming != "snake" {
		t.Fatalf("Naming: got %q, want %q", cfg.Naming, "snake")
	}
	if cfg.FlatTypeName != "MyAsset" {
		t.Fatalf("FlatTypeName: got %q, want %q", cfg.FlatTypeName, "MyAsset")
	}
}

func TestMatchInclExcl(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		includes []string
		excludes []string
		want     bool
	}{
		{
			name:     "no filters includes all",
			path:     "path/to/asset.png",
			includes: nil,
			excludes: nil,
			want:     true,
		},
		{
			name:     "excluded by pattern",
			path:     "path/to/.gitignore",
			includes: nil,
			excludes: []string{".git*"},
			want:     false,
		},
		{
			name:     "excluded by exact match",
			path:     "path/to/data.tmp",
			includes: nil,
			excludes: []string{"*.tmp", "*.log"},
			want:     false,
		},
		{
			name:     "excludes trimmed whitespace",
			path:     "path/to/file.cache",
			includes: nil,
			excludes: []string{"  ", "*.cache", ""},
			want:     false,
		},
		{
			name:     "included when no includes",
			path:     "path/to/asset.png",
			includes: nil,
			excludes: []string{".git"},
			want:     true,
		},
		{
			name:     "included by matching include",
			path:     "path/to/asset.png",
			includes: []string{"*.png", "*.jpg"},
			excludes: nil,
			want:     true,
		},
		{
			name:     "excluded by matching exclude despite include",
			path:     "path/to/thumb.png",
			includes: []string{"*.png"},
			excludes: []string{"thumb*"},
			want:     false,
		},
		{
			name:     "not included when no include matches",
			path:     "path/to/asset.gif",
			includes: []string{"*.png", "*.jpg"},
			excludes: nil,
			want:     false,
		},
		{
			name:     "includes trimmed whitespace",
			path:     "path/to/asset.png",
			includes: []string{"  ", "*.png", ""},
			excludes: nil,
			want:     true,
		},
		{
			name:     "not included when includes empty strings only",
			path:     "path/to/asset.png",
			includes: []string{"", "  "},
			excludes: nil,
			want:     false,
		},
		{
			name:     "excludes empty strings only still includes",
			path:     "path/to/asset.png",
			includes: nil,
			excludes: []string{"", "  "},
			want:     true,
		},
		{
			name:     "multiple excludes first match wins",
			path:     "path/to/file.log",
			includes: nil,
			excludes: []string{"*.tmp", "*.log", "*.cache"},
			want:     false,
		},
		{
			name:     "multiple includes first match wins",
			path:     "path/to/file.jpg",
			includes: []string{"*.png", "*.jpg"},
			excludes: nil,
			want:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := matchInclExcl(tc.path, tc.includes, tc.excludes)
			if got != tc.want {
				t.Errorf("matchInclExcl(%q, %v, %v) = %v, want %v",
					tc.path, tc.includes, tc.excludes, got, tc.want)
			}
		})
	}
}

func TestScanDir(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		glob      string
		recursive bool
		check     func(t *testing.T, assets []Asset)
	}{
		{
			name:  "finds all files in root",
			files: []string{"a.png", "b.png"},
			check: func(t *testing.T, assets []Asset) {
				if len(assets) != 2 {
					t.Fatalf("got %d assets, want 2", len(assets))
				}
			},
		},
		{
			name:  "non-recursive skips subdirs",
			files: []string{"a.png", "sub/b.png"},
			check: func(t *testing.T, assets []Asset) {
				if len(assets) != 1 {
					t.Fatalf("got %d assets, want 1", len(assets))
				}
			},
		},
		{
			name:      "recursive finds subdir files",
			files:     []string{"x.png", "sub/y.png", "sub/z.png"},
			recursive: true,
			check: func(t *testing.T, assets []Asset) {
				if len(assets) != 3 {
					t.Fatalf("got %d assets, want 3", len(assets))
				}
			},
		},
		{
			name:  "glob filters files",
			files: []string{"a.png", "b.jpg", "c.png"},
			glob: "*.png",
			check: func(t *testing.T, assets []Asset) {
				if len(assets) != 2 {
					t.Fatalf("got %d assets, want 2", len(assets))
				}
			},
		},
		{
			name:  "excludes hidden files",
			files: []string{"a.png", ".hidden", "b.png"},
			check: func(t *testing.T, assets []Asset) {
				if len(assets) != 2 {
					t.Fatalf("got %d assets, want 2", len(assets))
				}
			},
		},
		{
			name:  "collision detected",
			files: []string{"a-file.png", "sub/a-file.png"},
			check: func(t *testing.T, assets []Asset) {
				if len(assets) != 1 {
					t.Fatalf("got %d assets, want 1", len(assets))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createFile := func(rel string) {
				full := filepath.Join(tmpDir, rel)
				if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(full, []byte("data"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			for _, f := range tc.files {
				createFile(f)
			}
			namer := newNamer("pascal", "")
			knownNames := make(map[string]string)
			assets, err := scanDir(tmpDir, tc.glob, nil, nil, "", namer, tc.recursive, knownNames, false)
			if err != nil {
				t.Fatalf("scanDir error: %v", err)
			}
			tc.check(t, assets)
		})
	}
}

func TestScanDirWithFilters(t *testing.T) {
	tmpDir := t.TempDir()
	createFile := func(rel string) {
		full := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	createFile("a.png")
	createFile("b.jpg")
	createFile("c.tmp")
	createFile("data.log")

	tests := []struct {
		name     string
		includes []string
		excludes []string
		want     int
	}{
		{
			name:     "include png only",
			includes: []string{"*.png"},
			want:     1,
		},
		{
			name:     "exclude tmp and log",
			excludes: []string{"*.tmp", "*.log"},
			want:     2,
		},
		{
			name:     "include jpg and png, exclude log",
			includes: []string{"*.jpg", "*.png"},
			excludes: []string{"*.log"},
			want:     2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			namer := newNamer("pascal", "")
			knownNames := make(map[string]string)
			assets, err := scanDir(tmpDir, "", tc.includes, tc.excludes, "", namer, true, knownNames, false)
			if err != nil {
				t.Fatalf("scanDir error: %v", err)
			}
			if len(assets) != tc.want {
				t.Errorf("got %d assets, want %d", len(assets), tc.want)
			}
		})
	}
}

func TestScanDirStripPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	createFile := func(rel string) {
		full := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	createFile("textures/a.png")
	createFile("textures/b.png")

	namer := newNamer("pascal", "")
	knownNames := make(map[string]string)
	assets, err := scanDir(tmpDir, "", nil, nil, tmpDir+"/textures/", namer, true, knownNames, false)
	if err != nil {
		t.Fatalf("scanDir error: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("got %d assets, want 2", len(assets))
	}
	if assets[0].Path == tmpDir || strings.Contains(assets[0].Path, "textures") {
		t.Errorf("Path should not contain prefix, got: %s", assets[0].Path)
	}
}

func TestScanDirCollision(t *testing.T) {
	tmpDir := t.TempDir()
	createFile := func(rel string) {
		full := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	createFile("a-file.png")
	createFile("sub/a-file.png")

	namer := newNamer("pascal", "")
	knownNames := make(map[string]string)
	_, err := scanDir(tmpDir, "", nil, nil, "", namer, true, knownNames, false)
	if err == nil {
		t.Fatal("expected collision error")
	}
	if !strings.Contains(err.Error(), "constant name collision") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransformPath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		stripPrefix string
		want        string
	}{
		{
			name: "absolute path",
			path: "/absolute/path/a.png",
			want: "./absolute/path/a.png",
		},
		{
			name: "relative path",
			path: "path/to/a.png",
			want: "./path/to/a.png",
		},
		{
			name: "with strip prefix",
			path: "/project/assets/textures/a.png",
			stripPrefix: "/project/assets/",
			want:        "./textures/a.png",
		},
		{
			name: "no strip match",
			path: "/other/path/a.png",
			stripPrefix: "/project/assets/",
			want:        "./other/path/a.png",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := transformPath(tc.path, tc.stripPrefix)
			if got != tc.want {
				t.Errorf("transformPath(%q, %q) = %q, want %q", tc.path, tc.stripPrefix, got, tc.want)
			}
		})
	}
}

func TestSkipForHidden(t *testing.T) {
	if !skipForHidden(".hidden", false, "/path/.hidden") {
		t.Error("expected hidden file to be skipped")
	}
	if !skipForHidden(".gitignore", false, "/path/.gitignore") {
		t.Error("expected .gitignore to be skipped")
	}
	if skipForHidden("visible.png", false, "/path/visible.png") {
		t.Error("expected visible file not to be skipped")
	}
}

func TestSkipForGlob(t *testing.T) {
	if !skipForGlob("*.png", "file.jpg", false, "/path/file.jpg") {
		t.Error("expected non-matching file to be skipped")
	}
	if skipForGlob("*.png", "file.png", false, "/path/file.png") {
		t.Error("expected matching file not to be skipped")
	}
	if skipForGlob("", "file.png", false, "/path/file.png") {
		t.Error("expected empty glob to not skip")
	}
}
