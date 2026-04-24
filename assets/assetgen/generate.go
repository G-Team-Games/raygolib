package assetgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// Config controls asset generation behavior.
type Config struct {
	Root         string       // Root directory containing asset subdirectories
	Output       string       // Output file path
	Package      string       // Go package name for generated file
	ConfigFile   string       // JSON config file (overrides flags)
	Naming       string       // Naming style: pascal, camel, snake, upper_snake
	Prefix       string       // Constant prefix (optional)
	Kinds        string       // Comma-separated kinds to include
	KindsConfig  []KindConfig // Per-kind configuration (overrides Kinds string)
	Include      string       // Comma-separated glob patterns to include
	Exclude      string       // Comma-separated glob patterns to exclude
	StripPrefix  string       // Path prefix to strip from keys
	TemplateFile string       // Custom Go template file path
	SingleRoot   bool         // Scan root directory instead of subdirs for each kind
	Glob         string       // File pattern filter (pipe-separated, e.g. "*.ttf|*.otf")
}

// DefaultConfig returns sensible defaults for quick setup.
func DefaultConfig() Config {
	return Config{
		Naming:  "pascal",
		Root:    "assets",
		Output:  "zz_generated_assets.go",
		Package: "gameassets",
	}
}

// Kind represents the type of asset (texture, sound, model, etc.).
type Kind string

// Asset kind constants.
const (
	KindModel   Kind = "model"
	KindTexture Kind = "texture"
	KindImage   Kind = "image"
	KindSound   Kind = "sound"
	KindMusic   Kind = "music"
	KindFont    Kind = "font"
	KindShader  Kind = "shader"
)

// KindDir maps a Kind to its default directory name and generated type name.
type KindDir struct {
	Kind     Kind
	Dir      string
	TypeName string
}

// DefaultKindDirs provides conventional mappings for common asset layouts.
var DefaultKindDirs = []KindDir{
	{KindModel, "models", "ModelName"},
	{KindTexture, "textures", "TextureName"},
	{KindImage, "images", "ImageName"},
	{KindSound, "audio", "SoundName"},
	{KindMusic, "audio", "MusicName"},
	{KindFont, "fonts", "FontName"},
	{KindShader, "shaders", "ShaderName"},
}

type KindConfig struct {
	Kind     Kind
	Glob     string   `json:"glob,omitempty"`      // File pattern filter (e.g. "*.ttf|*.otf")
	Include  []string `json:"include,omitempty"`   // Include globs
	Exclude  []string `json:"exclude,omitempty"`   // Exclude globs
	Dir      string   `json:"dir,omitempty"`       // Subdirectory (empty = scan root of root)
	ScanRoot bool     `json:"scan_root,omitempty"` // Also scan root directory for this kind
	Type     string   `json:"type,omitempty"`
	Priority int      `json:"priority,omitempty"`
}

type FileConfig struct {
	Roots       []string              `json:"roots,omitempty"`
	Kinds       []KindConfig          `json:"kinds,omitempty"`
	Naming      string                `json:"naming,omitempty"`
	Prefix      string                `json:"prefix,omitempty"`
	StripPrefix string                `json:"strip_prefix,omitempty"`
	Include     []string              `json:"include,omitempty"`
	Exclude     []string              `json:"exclude,omitempty"`
	Output      string                `json:"output,omitempty"`
	Package     string                `json:"package,omitempty"`
	Types       map[string]KindConfig `json:"types,omitempty"`
	SingleRoot  bool                  `json:"single_root,omitempty"`
	Glob        string                `json:"glob,omitempty"`
}

type Asset struct {
	ConstName string
	Filename  string
	Key       string
}

type Category struct {
	Name     string
	Kind     Kind
	TypeName string
	Dir      string
	Assets   []Asset
}

type TemplateData struct {
	Package string
	Types   []Category
}

// Generate creates asset name constants based on Config settings.
// It scans the root directory for asset files matching the configured kinds,
// applies naming rules, and writes typed constants to the output file.
func Generate(cfg Config) error {
	if cfg.Package == "" {
		cfg.Package = "gameassets"
	}
	if cfg.Root == "" {
		cfg.Root = "assets"
	}

	var fileCfg FileConfig
	if cfg.ConfigFile != "" {
		data, err := os.ReadFile(cfg.ConfigFile)
		if err != nil {
			return fmt.Errorf("read config: %w", err)
		}
		if err := jsonLoad(&fileCfg, data); err != nil {
			return fmt.Errorf("parse config file: %w", err)
		}
	}

	roots := fileCfg.Roots
	if len(roots) == 0 {
		roots = []string{cfg.Root}
	}

	kinds := parseKinds(cfg.Kinds, cfg.KindsConfig, fileCfg.Kinds)
	includes := parseGlobs(cfg.Include, fileCfg.Include)
	excludes := parseGlobs(cfg.Exclude, fileCfg.Exclude)

	if cfg.Output == "" && fileCfg.Output != "" {
		cfg.Output = fileCfg.Output
	}
	if cfg.Output == "" {
		cfg.Output = "zz_generated_assets.go"
	}

	if cfg.Package == "gameassets" && fileCfg.Package != "" {
		cfg.Package = fileCfg.Package
	}
	if cfg.Prefix == "" {
		cfg.Prefix = fileCfg.Prefix
	}
	if cfg.StripPrefix == "" {
		cfg.StripPrefix = fileCfg.StripPrefix
	}
	if cfg.Naming == "" || cfg.Naming == "pascal" {
		if fileCfg.Naming != "" {
			cfg.Naming = fileCfg.Naming
		}
	}

	if cfg.SingleRoot {
		for i := range kinds {
			kinds[i].ScanRoot = true
			kinds[i].Dir = ""
		}
	}
	if cfg.Glob != "" {
		for i := range kinds {
			if kinds[i].Glob == "" {
				kinds[i].Glob = cfg.Glob
			}
		}
	}

	namer := newNamer(cfg.Naming, cfg.Prefix)

	var categories []Category
	seenDirs := map[string]bool{}

	for _, kind := range kinds {
		typeName := kind.Type
		if typeName == "" {
			typeName = capitalize(kind.Kind.String()) + "Name"
		}

		dir := kind.Dir
		if dir == "" && !kind.ScanRoot {
			dir = kind.Kind.DefaultDir()
		}

		for _, root := range roots {
			dirPath := filepath.Join(root, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				if !kind.ScanRoot || dir != "" {
					continue
				}
			}
			seenDirs[dirPath] = true

			assets, err := scanDir(dirPath, kind, includes, excludes, cfg.StripPrefix, namer)
			if err != nil {
				return fmt.Errorf("scan %s: %w", dirPath, err)
			}
			if len(assets) > 0 {
				categories = append(categories, Category{
					Name:     kind.Kind.Plural(),
					Kind:     kind.Kind,
					TypeName: typeName,
					Dir:      dir,
					Assets:   assets,
				})
			}
		}
	}

	if len(categories) == 0 {
		return fmt.Errorf("no assets found")
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Dir < categories[j].Dir
	})

	data := TemplateData{
		Package: cfg.Package,
		Types:   categories,
	}

	var buf bytes.Buffer

	tmplContent, err := os.ReadFile(cfg.TemplateFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read template: %w", err)
	}

	if len(tmplContent) == 0 {
		tmplContent = []byte(DefaultTemplate)
	}

	tmpl, err := template.New("assets").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	if err := os.WriteFile(cfg.Output, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func scanDir(dirPath string, kind KindConfig, includes, excludes []string, stripPrefix string, namer *namer) ([]Asset, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var assets []Asset
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if strings.HasPrefix(filename, ".") {
			continue
		}

		if kind.Glob != "" {
			matched := false
			for g := range strings.SplitSeq(kind.Glob, "|") {
				if ok, _ := filepath.Match(strings.TrimSpace(g), filename); ok {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		relPath := filepath.Join(dirPath, filename)
		if !matchGlobs(relPath, includes, excludes) {
			continue
		}

		key := filename
		if stripPrefix != "" {
			if trimmed, ok := strings.CutPrefix(key, stripPrefix); ok {
				key = trimmed
			}
		}

		constName := namer.name(kind.Kind, key)

		assets = append(assets, Asset{
			ConstName: constName,
			Filename:  filename,
			Key:       key,
		})
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ConstName < assets[j].ConstName
	})

	return assets, nil
}

func (k Kind) String() string { return string(k) }

func (k Kind) DefaultDir() string {
	switch k {
	case KindModel:
		return "models"
	case KindTexture:
		return "textures"
	case KindImage:
		return "images"
	case KindSound, KindMusic:
		return "audio"
	case KindFont:
		return "fonts"
	case KindShader:
		return "shaders"
	default:
		return string(k) + "s"
	}
}

func (k Kind) Plural() string {
	s := string(k)
	if strings.HasSuffix(s, "h") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(s[len(s)-2]) {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

func isVowel(r byte) bool {
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u'
}

type namer struct {
	style  string
	prefix string
}

func newNamer(style, prefix string) *namer {
	if style == "" {
		style = "pascal"
	}
	return &namer{style: style, prefix: prefix}
}

func (n *namer) name(kind Kind, filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.FieldsFunc(base, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '.'
	})

	var parts2 []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		switch n.style {
		case "pascal":
			if len(p) > 0 {
				parts2 = append(parts2, strings.ToUpper(string(p[0]))+strings.ToLower(p[1:]))
			}
		case "camel":
			if len(parts2) == 0 {
				parts2 = append(parts2, strings.ToLower(p))
			} else {
				if len(p) > 0 {
					parts2 = append(parts2, strings.ToUpper(string(p[0]))+strings.ToLower(p[1:]))
				}
			}
		case "snake":
			parts2 = append(parts2, strings.ToLower(p))
		case "upper_snake":
			parts2 = append(parts2, strings.ToUpper(p))
		default:
			if len(p) > 0 {
				parts2 = append(parts2, strings.ToUpper(string(p[0]))+strings.ToLower(p[1:]))
			}
		}
	}

	var joined string
	switch n.style {
	case "snake":
		joined = strings.Join(parts2, "_")
	case "upper_snake":
		joined = strings.Join(parts2, "_")
	default:
		joined = strings.Join(parts2, "")
	}

	if n.prefix != "" {
		joined = n.prefix + "_" + joined
	}

	return joined
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func parseKinds(kindsStr string, programmaticKinds []KindConfig, fileKinds []KindConfig) []KindConfig {
	if len(programmaticKinds) > 0 {
		return programmaticKinds
	}

	if kindsStr != "" {
		enabled := strings.Split(kindsStr, ",")
		var result []KindConfig
		for _, k := range enabled {
			k = strings.TrimSpace(k)
			result = append(result, KindConfig{
				Kind: Kind(k),
				Dir:  Kind(k).DefaultDir(),
				Type: capitalize(k) + "Name",
			})
		}
		return result
	}

	if len(fileKinds) > 0 {
		return fileKinds
	}

	return defaultKinds()
}

func defaultKinds() []KindConfig {
	return []KindConfig{
		{Kind: KindModel, Dir: "models", Type: "ModelName", Priority: 1},
		{Kind: KindTexture, Dir: "textures", Type: "TextureName", Priority: 2},
		{Kind: KindImage, Dir: "images", Type: "ImageName", Priority: 3},
		{Kind: KindSound, Dir: "audio", Type: "SoundName", Priority: 4},
		{Kind: KindMusic, Dir: "audio", Type: "MusicName", Priority: 5},
		{Kind: KindFont, Dir: "fonts", Type: "FontName", Priority: 6},
		{Kind: KindShader, Dir: "shaders", Type: "ShaderName", Priority: 7},
	}
}

func parseGlobs(str string, fileStr []string) []string {
	if str != "" {
		return strings.Split(str, ",")
	}
	return fileStr
}

func matchGlobs(path string, includes, excludes []string) bool {
	if len(includes) == 0 && len(excludes) == 0 {
		return true
	}

	for _, excl := range excludes {
		excl = strings.TrimSpace(excl)
		if excl == "" {
			continue
		}
		matched, _ := filepath.Match(excl, filepath.Base(path))
		if matched {
			return false
		}
	}

	if len(includes) == 0 {
		return true
	}

	for _, incl := range includes {
		incl = strings.TrimSpace(incl)
		if incl == "" {
			continue
		}
		matched, _ := filepath.Match(incl, filepath.Base(path))
		if matched {
			return true
		}
	}

	return len(includes) == 0
}

func jsonLoad(v any, data []byte) error {
	return json.Unmarshal(data, v)
}

const DefaultTemplate = `// Code generated by raygolib/assetgen. DO NOT EDIT.

package {{.Package}}

{{range $cat := .Types}}
// {{$cat.Name}} from {{$cat.Dir}}/
type {{$cat.TypeName}} string

const (
{{- range $cat.Assets}}
	{{.ConstName}} {{$cat.TypeName}} = "{{.Filename}}"
{{- end}}
)
{{end}}
`
