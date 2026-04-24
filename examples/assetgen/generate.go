//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/G-Team-Games/raygolib/assets/assetgen"
)

func main() {
	cfg := assetgen.DefaultConfig()
	cfg.Root = "./assets"
	cfg.Package = "assets"
	cfg.Output = "assets/generated.go"

	cfg.KindsConfig = []assetgen.KindConfig{
		{Kind: assetgen.KindTexture, Dir: "textures"},
		{Kind: assetgen.KindSound, Dir: "audio"},
	}

	if err := assetgen.Generate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated successfully")
}