package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	rgl "github.com/G-Team-Games/raygolib/assets"
)

func main() {
	cfg := rgl.DefaultConfig()

	var singleRoot bool
	var glob string
	var recursive bool
	var dryRun bool
	var verbose bool
	var flatMode bool
	var flatTypeName string
	var kinds string

	flag.StringVar(&cfg.Root, "root", "", "Root directory to scan for assets")
	flag.StringVar(&cfg.Output, "out", "", "Output file path")
	flag.StringVar(&cfg.Package, "pkg", "", "Go package name for output")
	flag.StringVar(&cfg.ConfigFile, "config", "", "JSON config file (overrides flags)")
	flag.StringVar(&cfg.Naming, "naming", "pascal", "Naming style: pascal, camel, snake, upper_snake")
	flag.StringVar(&cfg.Prefix, "prefix", "", "Constant prefix (optional)")
	flag.StringVar(&kinds, "kinds", "", "Comma-separated kinds to include")
	flag.StringVar(&cfg.StripPrefix, "strip-prefix", "", "Path prefix to strip from keys")
	flag.StringVar(&cfg.TemplateFile, "template", "", "Custom Go template file")
	flag.BoolVar(&singleRoot, "single-root", cfg.SingleRoot, "Scan root directory for all assets (no subdirs)")
	flag.BoolVar(&recursive, "recursive", cfg.Recursive, "Recursively scan subdirectories")
	flag.StringVar(&glob, "glob", "", "File pattern filter (e.g. '*.ttf|*.otf')")
	flag.BoolVar(&dryRun, "dry-run", false, "Print generated output to stdout instead of writing file")
	flag.BoolVar(&verbose, "v", false, "Print verbose debug output")
	flag.BoolVar(&flatMode, "flat-mode", false, "Generate one asset category and type")
	flag.StringVar(&flatTypeName, "flat-type", "", "Type name used when -flat-mode is enabled")

	flag.Parse()

	cfg.SingleRoot = singleRoot
	cfg.Recursive = recursive
	if glob != "" {
		cfg.Glob = glob
	}
	cfg.Verbose = verbose
	cfg.DryRun = dryRun
	cfg.FlatMode = flatMode
	if flatTypeName != "" {
		cfg.FlatTypeName = flatTypeName
	}

	if kinds != "" {
		cfg.Kinds = parseStringKinds(kinds)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Config:\n")
		cfg.Dump(os.Stderr)
	}

	if err := rgl.Generate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if dryRun {
		fmt.Println("Dry run complete")
		return
	}
	fmt.Println("Generated successfully")
}

func parseStringKinds(kindsStr string) []rgl.KindConfig {
	if kindsStr == "" {
		return nil
	}
	enabled := strings.Split(kindsStr, ",")
	var result []rgl.KindConfig
	for _, k := range enabled {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		kind := rgl.AssetKind(k)
		result = append(result, rgl.KindConfig{
			Kind:       kind,
			Dir:        kind.DefaultDir(),
			Type:       strings.ToUpper(string(k[0])) + k[1:] + "Name", // capitalize first letter 
			Plural:     kind.Plural(),
			Extensions: kind.DefaultExtensions(),
		})
	}
	return result
}
