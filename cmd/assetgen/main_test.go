package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	rga "github.com/G-Team-Games/raygolib/assets"
)

func TestParseStringKinds(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		if got := parseStringKinds(""); got != nil {
			t.Fatalf("expected nil for empty input, got %#v", got)
		}
	})

	t.Run("trims and skips empties", func(t *testing.T) {
		got := parseStringKinds(" texture, ,shader ,, custom ")
		if len(got) != 3 {
			t.Fatalf("expected 3 kinds, got %d", len(got))
		}

		if got[0].Kind != rga.KindTexture || got[0].Dir != "textures" || got[0].Type != "TextureName" || got[0].Plural != "textures" {
			t.Fatalf("unexpected first kind config: %+v", got[0])
		}
		if got[1].Kind != rga.KindShader || got[1].Dir != "shaders" || got[1].Type != "ShaderName" || got[1].Plural != "shaders" {
			t.Fatalf("unexpected second kind config: %+v", got[1])
		}
		if got[2].Kind != rga.AssetKind("custom") || got[2].Dir != "customs" || got[2].Type != "CustomName" || got[2].Plural != "customs" {
			t.Fatalf("unexpected custom kind config: %+v", got[2])
		}
	})

	t.Run("single letter kind", func(t *testing.T) {
		got := parseStringKinds("x")
		if len(got) != 1 {
			t.Fatalf("expected 1 kind, got %d", len(got))
		}
		if got[0].Type != "XName" {
			t.Fatalf("expected type XName, got %q", got[0].Type)
		}
	})
}

func TestRun_FlagMappingAndSuccessMessage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var gotCfg rga.Config
	genCalled := 0
	generate := func(cfg rga.Config) error {
		genCalled++
		gotCfg = cfg
		return nil
	}

	code := run([]string{
		"-root", "assets_dir",
		"-out", "out.go",
		"-pkg", "mypkg",
		"-config", "assetgen.json",
		"-naming", "snake",
		"-prefix", "pre",
		"-kinds", "texture,shader",
		"-strip-prefix", "assets/",
		"-template", "tmpl.tmpl",
		"-single-root=false",
		"-recursive=false",
		"-glob", "*.png|*.jpg",
		"-v",
		"-flat-mode",
		"-flat-type", "AnyAsset",
	}, &stdout, &stderr, generate)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q", code, stderr.String())
	}
	if genCalled != 1 {
		t.Fatalf("expected generate called once, got %d", genCalled)
	}

	if gotCfg.Root != "assets_dir" || gotCfg.Output != "out.go" || gotCfg.Package != "mypkg" || gotCfg.ConfigFile != "assetgen.json" {
		t.Fatalf("unexpected root/out/pkg/config mapping: %+v", gotCfg)
	}
	if gotCfg.Naming != "snake" || gotCfg.Prefix != "pre" || gotCfg.StripPrefix != "assets/" || gotCfg.TemplateFile != "tmpl.tmpl" {
		t.Fatalf("unexpected naming/prefix/strip/template mapping: %+v", gotCfg)
	}
	if gotCfg.SingleRoot || gotCfg.Recursive {
		t.Fatalf("expected single-root=false and recursive=false, got SingleRoot=%v Recursive=%v", gotCfg.SingleRoot, gotCfg.Recursive)
	}
	if gotCfg.Glob != "*.png|*.jpg" || !gotCfg.Verbose || !gotCfg.FlatMode || gotCfg.FlatTypeName != "AnyAsset" {
		t.Fatalf("unexpected glob/verbose/flat mapping: %+v", gotCfg)
	}
	if gotCfg.DryRun {
		t.Fatalf("expected dry-run false by default")
	}

	wantKinds := []rga.AssetKind{rga.KindTexture, rga.KindShader}
	if len(gotCfg.Kinds) != len(wantKinds) {
		t.Fatalf("expected %d kinds, got %d", len(wantKinds), len(gotCfg.Kinds))
	}
	for i, want := range wantKinds {
		if gotCfg.Kinds[i].Kind != want {
			t.Fatalf("unexpected kind at %d: got %q want %q", i, gotCfg.Kinds[i].Kind, want)
		}
	}

	if !strings.Contains(stderr.String(), "Config:\n") {
		t.Fatalf("expected verbose config dump in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Generated successfully") {
		t.Fatalf("expected success message on stdout, got %q", stdout.String())
	}
}

func TestRun_DryRunMessage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var gotCfg rga.Config
	code := run([]string{"-dry-run"}, &stdout, &stderr, func(cfg rga.Config) error {
		gotCfg = cfg
		return nil
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !gotCfg.DryRun {
		t.Fatalf("expected dry-run flag mapped to config")
	}
	if strings.Contains(stdout.String(), "Generated successfully") {
		t.Fatalf("did not expect non-dry-run message, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Dry run complete") {
		t.Fatalf("expected dry-run message, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr on dry-run success, got %q", stderr.String())
	}
}

func TestRun_GenerateError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(nil, &stdout, &stderr, func(cfg rga.Config) error {
		return errors.New("boom")
	})

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout on error, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Error: boom") {
		t.Fatalf("expected stderr error output, got %q", stderr.String())
	}
}

func TestRun_InvalidFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-unknown-flag"}, &stdout, &stderr, func(cfg rga.Config) error { return nil })

	if code != 2 {
		t.Fatalf("expected exit code 2 for parse error, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout for parse error, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("expected parse error in stderr, got %q", stderr.String())
	}
}

func TestRun_DefaultGenerateWhenNil(t *testing.T) {
	tmpDir := t.TempDir()
	asset := filepath.Join(tmpDir, "player.png")
	if err := os.WriteFile(asset, []byte(""), 0o644); err != nil {
		t.Fatalf("write asset file: %v", err)
	}

	outPath := filepath.Join(tmpDir, "zz_generated_assets.go")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{
		"-root", tmpDir,
		"-out", outPath,
		"-pkg", "assetspkg",
		"-kinds", "texture",
		"-single-root",
	}, &stdout, &stderr, nil)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Generated successfully") {
		t.Fatalf("expected success output, got %q", stdout.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("expected generated file at %q: %v", outPath, err)
	}
	content := string(data)
	if !strings.Contains(content, "package assetspkg") {
		t.Fatalf("expected generated package name, got %q", content)
	}
	if !strings.Contains(content, "type TextureName string") {
		t.Fatalf("expected generated texture type, got %q", content)
	}
}

func TestRun_KindsParsingOrderPreserved(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var gotKinds []rga.KindConfig
	code := run([]string{"-kinds", "shader,texture,custom"}, &stdout, &stderr, func(cfg rga.Config) error {
		gotKinds = cfg.Kinds
		return nil
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	want := []rga.AssetKind{rga.KindShader, rga.KindTexture, rga.AssetKind("custom")}
	got := make([]rga.AssetKind, 0, len(gotKinds))
	for _, k := range gotKinds {
		got = append(got, k.Kind)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected kinds order: got %v want %v", got, want)
	}
}
