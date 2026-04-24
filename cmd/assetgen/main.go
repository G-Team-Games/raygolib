package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/G-Team-Games/raygolib/assets/assetgen"
)

func main() {
	cfg := assetgen.DefaultConfig()

	flag.StringVar(&cfg.Root, "root", "", "Root directory to scan for assets")
	flag.StringVar(&cfg.Output, "out", "", "Output file path")
	flag.StringVar(&cfg.Package, "pkg", "", "Go package name for output")
	flag.StringVar(&cfg.ConfigFile, "config", "", "YAML/JSON config file (overrides flags)")
	flag.StringVar(&cfg.Naming, "naming", "pascal", "Naming style: pascal, camel, snake, upper_snake")
	flag.StringVar(&cfg.Prefix, "prefix", "", "Constant prefix (optional)")
	flag.StringVar(&cfg.Kinds, "kinds", "", "Comma-separated kinds to include")
	flag.StringVar(&cfg.Include, "include", "", "Comma-separated glob patterns to include")
	flag.StringVar(&cfg.Exclude, "exclude", "", "Comma-separated glob patterns to exclude")
	flag.StringVar(&cfg.StripPrefix, "strip-prefix", "", "Path prefix to strip from keys")
	flag.StringVar(&cfg.TemplateFile, "template", "", "Custom Go template file")

	flag.Parse()

	if err := assetgen.Generate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated successfully")
}