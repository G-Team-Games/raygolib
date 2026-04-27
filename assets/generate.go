package rga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

const (
	defaultOutputFile   = "zz_generated_assets.go"
	defaultFlatTypeName = "AssetName"
)

// Config controls asset generation and doubles as the JSON config file format.
// All fields are JSON-serialisable; runtime-only fields carry json:"-" tags and
// are simply ignored when reading or writing a config file.
//
// Zero-value string fields fall back to sensible hardcoded defaults when
// Generate is called. Boolean fields default to false — call DefaultConfig()
// when you want SingleRoot and Recursive pre-set to true.
type Config struct {
	Root         string       `json:"root,omitempty"`           // root directory to scan; default "assets"
	Output       string       `json:"output,omitempty"`         // output file path; default zz_generated_assets.go
	Package      string       `json:"package,omitempty"`        // Go package name; default "assets"
	SingleRoot   bool         `json:"single_root,omitempty"`    // scan root dir directly instead of per-kind subdirs
	Recursive    bool         `json:"recursive,omitempty"`      // recurse into subdirectories
	Glob         string       `json:"glob,omitempty"`           // pipe-separated filter applied globally to all kinds
	Kinds        []KindConfig `json:"kinds,omitempty"`          // per-kind config; omit to use all built-in kinds
	Naming       string       `json:"naming,omitempty"`         // pascal (default), camel, snake, upper_snake
	Prefix       string       `json:"prefix,omitempty"`         // prepended to every constant name
	StripPrefix  string       `json:"strip_prefix,omitempty"`   // path prefix to strip from emitted asset paths
	Includes     []string     `json:"include,omitempty"`        // filename globs — only matching files are included
	Excludes     []string     `json:"exclude,omitempty"`        // filename globs — matching files are excluded
	FlatMode     bool         `json:"flat_mode,omitempty"`      // collapse all assets into one type
	FlatTypeName string       `json:"flat_type_name,omitempty"` // type name for flat mode; default AssetName
	Verbose      bool         `json:"verbose,omitempty"`        // print scan details to stderr

	// Runtime only fields — ignored in JSON config files
	ConfigFile   string    `json:"-"` // path to a JSON config file
	TemplateFile string    `json:"-"` // custom Go template file; uses DefaultTemplate when empty
	DryRun       bool      `json:"-"` // write to DryRunWriter/stdout instead of Output file
	DryRunWriter io.Writer `json:"-"` // destination for dry-run output; nil means stdout
}

// DefaultConfig returns a Config with the most common defaults pre-set.
func DefaultConfig() Config {
	return Config{
		Root:       "assets",
		Output:     defaultOutputFile,
		Package:    "assets",
		Naming:     "pascal",
		SingleRoot: true,
		Recursive:  true,
	}
}

// KindConfig configures how one kind of asset is scanned and what code is generated.
// Dir, Extensions/Glob, and Type can all be overridden; everything else has a
// sensible built-in default derived from the Kind constant.
type KindConfig struct {
	Kind       AssetKind     `json:"kind"`
	Dir        string   `json:"dir,omitempty"`        // subdirectory to scan; defaults to Kind.DefaultDir()
	Extensions []string `json:"extensions,omitempty"` // file extensions, e.g. [".png",".jpg"] — alternative to Glob
	Glob       string   `json:"glob,omitempty"`       // pipe-separated filename filter, e.g. "*.png|*.jpg"
	Include    []string `json:"include,omitempty"`    // additional per-kind include patterns
	Exclude    []string `json:"exclude,omitempty"`    // additional per-kind exclude patterns
	Type       string   `json:"type,omitempty"`       // generated Go type name, e.g. "TextureName"
	Plural     string   `json:"plural,omitempty"`     // plural label used in generated comments
	ScanRoot   bool     `json:"scan_root,omitempty"`  // scan root dir directly; set automatically by Config.SingleRoot
}

var defaultKindsSlice = []AssetKind{
	KindModel,
	KindTexture,
	KindImage,
	KindSound,
	KindMusic,
	KindFont,
	KindShader,
}

// DefaultKinds derives the default KindConfig slice from DefaultKindDirs so
// there is one canonical place to update when adding new built-in kinds.
func DefaultKinds() []KindConfig {
	kinds := make([]KindConfig, len(defaultKindsSlice))
	for i, k := range defaultKindsSlice {
		kinds[i] = KindConfig{
			Kind:   k,
			Type:   k.TypeName(),
			Plural: k.Plural(),
			Dir:    k.DefaultDir(),
		}
	}
	return kinds
}

// Asset is one discovered file with its generated constant name and relative path.
type Asset struct {
	ConstName string
	Path      string
}

// TypeName is a group of assets that share a generated Go type.
type TypeName struct {
	Name     string
	Kind     AssetKind
	TypeName string
	Dir      string
	Assets   []Asset
}

// TemplateData is passed to the Go template when rendering the output file.
type TemplateData struct {
	Package string
	Types   []TypeName
}

// Generate scans the asset tree described by cfg and writes typed Go constants
// to cfg.Output (or stdout when cfg.DryRun is true).
func Generate(cfg Config) error {
	// Load and merge the JSON config file first; programmatic values take priority.
	if cfg.ConfigFile != "" {
		fileCfg, err := loadFileConfig(cfg.ConfigFile)
		if err != nil {
			return err
		}
		cfg = mergeWithFile(cfg, fileCfg)
	}
	cfg = applyStringDefaults(cfg)

	if err := validateNaming(cfg.Naming); err != nil {
		return err
	}

	categories, err := buildCategories(cfg)
	if err != nil {
		return err
	}

	out, err := renderTemplate(cfg.TemplateFile, TemplateData{Package: cfg.Package, Types: categories})
	if err != nil {
		return err
	}

	if cfg.DryRun {
		w := cfg.DryRunWriter
		if w == nil {
			w = os.Stdout
		}
		_, err = w.Write(out)
		return err
	}

	if err := os.WriteFile(cfg.Output, out, 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

// GenerateDryRun renders the generated code to w without writing any files.
func GenerateDryRun(cfg Config, w io.Writer) error {
	cfg.DryRun = true
	cfg.DryRunWriter = w
	return Generate(cfg)
}

// Dump prints a human-readable summary of cfg to w. Useful for -v / debug output.
func (cfg Config) Dump(w io.Writer) {
	fmt.Fprintf(w, "  Root:       %s\n", cfg.Root)
	fmt.Fprintf(w, "  Output:     %s\n", cfg.Output)
	fmt.Fprintf(w, "  Package:    %s\n", cfg.Package)
	fmt.Fprintf(w, "  Naming:     %s\n", cfg.Naming)
	fmt.Fprintf(w, "  SingleRoot: %t\n", cfg.SingleRoot)
	fmt.Fprintf(w, "  Recursive:  %t\n", cfg.Recursive)
	fmt.Fprintf(w, "  FlatMode:   %t\n", cfg.FlatMode)
	fmt.Fprintf(w, "  DryRun:     %t\n", cfg.DryRun)
}

// ---- config resolution -----------------------------------------------------

func loadFileConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// mergeWithFile returns cfg with any zero-value fields filled from file.
// Programmatic values always win; file values fill in what is missing.
// Booleans use OR semantics: true in either source wins.
func mergeWithFile(cfg, file Config) Config {
	if cfg.Root == "" {
		cfg.Root = file.Root
	}
	if cfg.Output == "" {
		cfg.Output = file.Output
	}
	if cfg.Package == "" {
		cfg.Package = file.Package
	}
	if cfg.Naming == "" {
		cfg.Naming = file.Naming
	}
	if cfg.Prefix == "" {
		cfg.Prefix = file.Prefix
	}
	if cfg.StripPrefix == "" {
		cfg.StripPrefix = file.StripPrefix
	}
	if cfg.Glob == "" {
		cfg.Glob = file.Glob
	}
	if cfg.FlatTypeName == "" {
		cfg.FlatTypeName = file.FlatTypeName
	}
	if len(cfg.Kinds) == 0 {
		cfg.Kinds = file.Kinds
	}
	if len(cfg.Includes) == 0 {
		cfg.Includes = file.Includes
	}
	if len(cfg.Excludes) == 0 {
		cfg.Excludes = file.Excludes
	}
	cfg.SingleRoot = cfg.SingleRoot || file.SingleRoot
	cfg.Recursive = cfg.Recursive || file.Recursive
	cfg.FlatMode = cfg.FlatMode || file.FlatMode
	cfg.Verbose = cfg.Verbose || file.Verbose
	return cfg
}

// applyStringDefaults fills remaining zero-value string fields with hardcoded
// defaults. Boolean fields are intentionally not touched here — use DefaultConfig().
func applyStringDefaults(cfg Config) Config {
	if cfg.Root == "" {
		cfg.Root = "assets"
	}
	if cfg.Output == "" {
		cfg.Output = defaultOutputFile
	}
	if cfg.Package == "" {
		cfg.Package = "assets"
	}
	if cfg.Naming == "" {
		cfg.Naming = "pascal"
	}
	if cfg.FlatTypeName == "" {
		cfg.FlatTypeName = defaultFlatTypeName
	}
	return cfg
}

// ---- scanning --------------------------------------------------------------

// buildCategories resolves kinds, scans the file tree, and returns categories
// sorted by directory name.
func buildCategories(cfg Config) ([]TypeName, error) {
	kinds := cfg.Kinds
	if len(kinds) == 0 {
		kinds = DefaultKinds()
	}
	kinds = normalizeKinds(kinds, cfg)

	if cfg.Verbose {
		warnUnknownKinds(kinds)
	}

	namer := newNamer(cfg.Naming, cfg.Prefix)
	categories, err := collectCategories(cfg, kinds, namer)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("no assets found under %q", cfg.Root)
	}

	if cfg.FlatMode {
		categories = flattenCategories(categories, cfg.FlatTypeName)
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Dir < categories[j].Dir
	})
	return categories, nil
}

// normalizeKinds fills any omitted KindConfig fields with their defaults.
func normalizeKinds(kinds []KindConfig, cfg Config) []KindConfig {
	for i := range kinds {
		k := &kinds[i]
		if k.Type == "" {
			k.Type = k.Kind.TypeName()
		}
		if k.Plural == "" {
			k.Plural = k.Kind.Plural()
		}
		// Build a glob from extensions when no explicit glob is set.
		if k.Glob == "" {
			exts := k.Extensions
			if len(exts) == 0 {
				exts = k.Kind.DefaultExtensions()
				k.Extensions = exts
			}
			k.Glob = extensionsToGlob(exts)
		}
		// Fall back to the global glob if the kind still has none.
		if k.Glob == "" {
			k.Glob = cfg.Glob
		}
		// SingleRoot overrides the per-kind directory.
		if cfg.SingleRoot {
			k.ScanRoot = true
			k.Dir = ""
		}
	}
	return kinds
}

func warnUnknownKinds(kinds []KindConfig) {
	for _, k := range kinds {
		if !k.Kind.IsKnown() {
			fmt.Fprintf(os.Stderr, "assetgen: custom kind %q — no built-in defaults, using configured dir/extensions\n", k.Kind)
		}
	}
}

// collectCategories scans the directory tree for each kind and returns one
// Category per kind that found at least one asset.
func collectCategories(cfg Config, kinds []KindConfig, namer *namer) ([]TypeName, error) {
	var categories []TypeName
	// scanSeen deduplicates identical (kind, dirPath, glob) combinations.
	scanSeen := map[string]bool{}
	// knownNames is shared across all kinds to catch cross-kind name collisions.
	knownNames := map[string]string{}

	for _, kind := range kinds {
		dir := kind.Dir
		if dir == "" && !kind.ScanRoot {
			dir = kind.Kind.DefaultDir()
		}
		dirPath := filepath.Join(cfg.Root, dir)

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			continue // kind simply not used in this project
		}

		scanKey := string(kind.Kind) + "|" + dirPath + "|" + kind.Glob
		if scanSeen[scanKey] {
			continue
		}
		scanSeen[scanKey] = true

		// Merge global and per-kind include/exclude lists.
		inc := append(append([]string(nil), cfg.Includes...), kind.Include...)
		exc := append(append([]string(nil), cfg.Excludes...), kind.Exclude...)

		assets, err := scanDir(dirPath, kind.Glob, inc, exc, cfg.StripPrefix, namer, cfg.Recursive, knownNames, cfg.Verbose)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", dirPath, err)
		}
		if len(assets) == 0 {
			continue
		}
		categories = append(categories, TypeName{
			Name:     kind.Plural,
			Kind:     kind.Kind,
			TypeName: kind.Type,
			Dir:      dir,
			Assets:   assets,
		})
	}
	return categories, nil
}

// flattenCategories merges all assets into a single category under typeName.
func flattenCategories(categories []TypeName, typeName string) []TypeName {
	if typeName == "" {
		typeName = defaultFlatTypeName
	}
	var all []Asset
	for _, c := range categories {
		all = append(all, c.Assets...)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ConstName < all[j].ConstName })
	return []TypeName{{
		Name:     "assets",
		Kind:     AssetKind("asset"),
		TypeName: typeName,
		Assets:   all,
	}}
}

// scanDir walks dirPath and returns assets that pass the glob, include, and
// exclude filters. knownNames is used to detect constant name collisions across
// the whole generation run.
func scanDir(dirPath, glob string, includes, excludes []string, stripPrefix string, namer *namer, recursive bool, knownNames map[string]string, verbose bool) ([]Asset, error) {
	var assets []Asset

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if skip, err := shouldSkipDir(d, dirPath, recursive); skip {
			return err
		}
		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if skipForHidden(name, verbose, path) {
			return nil
		}
		if skipForGlob(glob, name, verbose, path) {
			return nil
		}
		if skipForFilters(path, includes, excludes, verbose, path) {
			return nil
		}

		assetPath := transformPath(path, stripPrefix)
		constName := namer.name(name)
		if prev, exists := knownNames[constName]; exists {
			return fmt.Errorf("constant name collision %q: files %q and %q produce the same name", constName, prev, path)
		}
		knownNames[constName] = path

		if verbose {
			fmt.Fprintf(os.Stderr, "assetgen include: %s -> %s\n", path, constName)
		}
		assets = append(assets, Asset{ConstName: constName, Path: assetPath})
		return nil
	})

	return assets, err
}

func shouldSkipDir(d os.DirEntry, dirPath string, recursive bool) (skip bool, err error) {
	if !d.IsDir() {
		return false, nil
	}
	if !recursive && d.Name() != filepath.Base(dirPath) {
		return true, filepath.SkipDir
	}
	return true, nil
}

func skipForHidden(name string, verbose bool, path string) bool {
	if strings.HasPrefix(name, ".") {
		if verbose {
			fmt.Fprintf(os.Stderr, "assetgen skip (hidden): %s\n", path)
		}
		return true
	}
	return false
}

func skipForGlob(glob, name string, verbose bool, path string) bool {
	if glob != "" && !matchGlob(glob, name) {
		if verbose {
			fmt.Fprintf(os.Stderr, "assetgen skip (glob): %s\n", path)
		}
		return true
	}
	return false
}

func skipForFilters(path string, includes, excludes []string, verbose bool, fullPath string) bool {
	if !matchInclExcl(path, includes, excludes) {
		if verbose {
			fmt.Fprintf(os.Stderr, "assetgen skip (filter): %s\n", fullPath)
		}
		return true
	}
	return false
}

func transformPath(path, stripPrefix string) string {
	assetPath := path
	if stripPrefix != "" {
		if trimmed, ok := strings.CutPrefix(path, stripPrefix); ok {
			assetPath = trimmed
		}
	}
	return "./" + strings.TrimPrefix(filepath.ToSlash(assetPath), "/")
}

// ---- filtering helpers -----------------------------------------------------

// matchGlob reports whether name matches any pattern in the pipe-separated glob string.
func matchGlob(glob, name string) bool {
	for _, pattern := range strings.Split(glob, "|") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		if ok, _ := filepath.Match(pattern, name); ok {
			return true
		}
	}
	return false
}

// matchInclExcl returns false if path's filename matches any exclude pattern.
// If includes are configured, at least one must match. Otherwise true.
func matchInclExcl(path string, includes, excludes []string) bool {
	base := filepath.Base(path)
	for _, pat := range excludes {
		if pat = strings.TrimSpace(pat); pat == "" {
			continue
		}
		if ok, _ := filepath.Match(pat, base); ok {
			return false
		}
	}
	if len(includes) == 0 {
		return true
	}
	for _, pat := range includes {
		if pat = strings.TrimSpace(pat); pat == "" {
			continue
		}
		if ok, _ := filepath.Match(pat, base); ok {
			return true
		}
	}
	return false
}

// ---- naming ----------------------------------------------------------------

type namer struct {
	style  string
	prefix string
}

func newNamer(style, prefix string) *namer {
	return &namer{style: style, prefix: prefix}
}

func (n *namer) name(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.FieldsFunc(base, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var words []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		switch n.style {
		case "pascal":
			words = append(words, strings.ToUpper(string(p[0]))+strings.ToLower(p[1:]))
		case "camel":
			if len(words) == 0 {
				words = append(words, strings.ToLower(p))
			} else {
				words = append(words, strings.ToUpper(string(p[0]))+strings.ToLower(p[1:]))
			}
		case "snake":
			words = append(words, strings.ToLower(p))
		case "upper_snake":
			words = append(words, strings.ToUpper(p))
		}
	}

	sep := ""
	if n.style == "snake" || n.style == "upper_snake" {
		sep = "_"
	}
	joined := strings.Join(words, sep)

	// Prepend underscore if the identifier would start with a digit.
	if len(joined) > 0 {
		if first := []rune(joined)[0]; !unicode.IsLetter(first) && first != '_' {
			joined = "_" + joined
		}
	}
	if n.prefix != "" {
		joined = n.prefix + "_" + joined
	}
	return joined
}

func validateNaming(style string) error {
	switch style {
	case "pascal", "camel", "snake", "upper_snake":
		return nil
	default:
		return fmt.Errorf("invalid naming style %q: want pascal, camel, snake, or upper_snake", style)
	}
}

// ---- template --------------------------------------------------------------

func renderTemplate(templateFile string, data TemplateData) ([]byte, error) {
	src := DefaultTemplate
	if templateFile != "" {
		b, err := os.ReadFile(templateFile)
		if err != nil {
			return nil, fmt.Errorf("read template: %w", err)
		}
		src = string(b)
	}

	tmpl, err := template.New("assets").Parse(src)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}

// DefaultTemplate is the Go template used when no TemplateFile is configured.
// It receives TemplateData and emits one named string type per Category.
const DefaultTemplate = `// Code generated by raygolib/assetgen. DO NOT EDIT.

package {{.Package}}
{{range $cat := .Types}}
// {{$cat.Name}}
type {{$cat.TypeName}} string

const (
{{- range $cat.Assets}}
	{{.ConstName}} {{$cat.TypeName}} = "{{.Path}}"
{{- end}}
)
{{end}}`

// ---- small utilities -------------------------------------------------------

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func extensionsToGlob(exts []string) string {
	parts := make([]string, 0, len(exts))
	for _, ext := range exts {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		parts = append(parts, "*"+strings.ToLower(ext))
	}
	return strings.Join(parts, "|")
}

func isVowel(r byte) bool {
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'A' || r == 'E' || r == 'I' || r == 'O' || r == 'U'
}
