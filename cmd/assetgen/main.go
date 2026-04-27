package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	rgl "github.com/G-Team-Games/raygolib/assets"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, rgl.Generate))
}

func run(args []string, stdout, stderr io.Writer, generate func(rgl.Config) error) int {
	cfg := rgl.DefaultConfig()
	if generate == nil {
		generate = rgl.Generate
	}

	fs := flag.NewFlagSet("assetgen", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var singleRoot bool
	var glob string
	var recursive bool
	var dryRun bool
	var verbose bool
	var flatMode bool
	var flatTypeName string
	var kinds string

	fs.StringVar(&cfg.Root, "root", "", "Root directory to scan for assets")
	fs.StringVar(&cfg.Output, "out", "", "Output file path")
	fs.StringVar(&cfg.Package, "pkg", "", "Go package name for output")
	fs.StringVar(&cfg.ConfigFile, "config", "", "JSON config file (overrides flags)")
	fs.StringVar(&cfg.Naming, "naming", "pascal", "Naming style: pascal, camel, snake, upper_snake")
	fs.StringVar(&cfg.Prefix, "prefix", "", "Constant prefix (optional)")
	fs.StringVar(&kinds, "kinds", "", "Comma-separated kinds to include")
	fs.StringVar(&cfg.StripPrefix, "strip-prefix", "", "Path prefix to strip from keys")
	fs.StringVar(&cfg.TemplateFile, "template", "", "Custom Go template file")
	fs.BoolVar(&singleRoot, "single-root", cfg.SingleRoot, "Scan root directory for all assets (no subdirs)")
	fs.BoolVar(&recursive, "recursive", cfg.Recursive, "Recursively scan subdirectories")
	fs.StringVar(&glob, "glob", "", "File pattern filter (e.g. '*.ttf|*.otf')")
	fs.BoolVar(&dryRun, "dry-run", false, "Print generated output to stdout instead of writing file")
	fs.BoolVar(&verbose, "v", false, "Print verbose debug output")
	fs.BoolVar(&flatMode, "flat-mode", false, "Generate one asset category and type")
	fs.StringVar(&flatTypeName, "flat-type", "", "Type name used when -flat-mode is enabled")

	if err := fs.Parse(args); err != nil {
		return 2
	}

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
		fmt.Fprintf(stderr, "Config:\n")
		cfg.Dump(stderr)
	}

	if err := generate(cfg); err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}
	if dryRun {
		fmt.Fprintln(stdout, "Dry run complete")
		return 0
	}
	fmt.Fprintln(stdout, "Generated successfully")
	return 0
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
