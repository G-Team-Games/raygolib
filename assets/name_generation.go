package assets

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed template.go.tmpl
var templateFS embed.FS

type AssetCategory struct {
	Name      string // Models, Textures, etc.
	Directory string // models, graphics, etc.
	Assets    []Asset
}

type Asset struct {
	ConstName string // MODEL_PLAYER
	Filename  string // player.glb
}

type TemplateData struct {
	Categories []AssetCategory
}

func main() {
	if err := generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated assets_gen.go successfully")
}

// When dev wants to run generate from client or from root
func findAssetsPath() string {
	possiblePaths := []string{
		"../../assets",
		"../assets",
		"client/assets",
		"assets",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	return ""
}

func generate() error {
	assetsPath := findAssetsPath()
	if assetsPath == "" {
		return fmt.Errorf("could not find assets directory")
	}

	categories := []struct {
		name string
		dir  string
	}{
		{"Models", "models"},
		{"Textures", "textures"},
		{"Images", "images"},
		{"Audio", "audio"},
		{"Fonts", "fonts"},
		{"Shaders", "shaders"},
	}

	var assetCategories []AssetCategory

	for _, cat := range categories {
		dirPath := filepath.Join(assetsPath, cat.dir)

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("Directory not found: %s (skipping)\n", dirPath)
			continue
		}

		assets, err := scanDirectory(dirPath, cat.name)
		if err != nil {
			return fmt.Errorf("failed to scan %s: %w", cat.dir, err)
		}

		if len(assets) > 0 {
			assetCategories = append(assetCategories, AssetCategory{
				Name:      cat.name,
				Directory: cat.dir,
				Assets:    assets,
			})
			fmt.Printf("Found %d %s\n", len(assets), cat.name)
		}
	}

	tmplContent, err := templateFS.ReadFile("template.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("assets").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := TemplateData{
		Categories: assetCategories,
	}

	outputPath := filepath.Join(assetsPath, "assets_gen.go")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func scanDirectory(dirPath, categoryName string) ([]Asset, error) {
	var assets []Asset

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		// Skip hidden files
		if strings.HasPrefix(filename, ".") {
			continue
		}

		constName := generateConstName(categoryName, filename)

		assets = append(assets, Asset{
			ConstName: constName,
			Filename:  filename,
		})
	}

	return assets, nil
}

func generateConstName(category, filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// create parts from string, split by interpunction
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '.'
	})

	var result strings.Builder
	result.WriteString(strings.TrimSuffix(category, "s")) // Remove plural: Models -> Model

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}

	if filepath.Ext(filename) == ".vs" || filepath.Ext(filename) == ".fs" {
		result.WriteString(strings.ToTitle(filepath.Ext(filename)[1:]))
	}

	return result.String()
}
