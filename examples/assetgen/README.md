# Asset Name Generation

Generate typed asset names for compile-time autocomplete and refactor safety.

## Directory Layouts

### Standard - Subdirs Per Kind

```
assets/
  textures/player.png
  textures/enemy.png
  audio/click.ogg
```

### Single Root - All In One Dir

```
assets/
  player.png
  enemy.png
  click.ogg
```

### Mixed - Subdirs + Root

```
assets/
  textures/player.png
  textures/enemy.png
  audio/click.ogg
  font.ttf
  shader.glsl
```

---

## Standard Layout (Subdirs)

### CLI

```bash
go run github.com/G-Team-Games/raygolib/cmd/assetgen \
  -root ./assets \
  -pkg assets \
  -out assets/generated.go
```

### Programmatic

```go
cfg := assetgen.DefaultConfig()
cfg.Root = "./assets"
cfg.Package = "assets"
cfg.Output = "assets/generated.go"
cfg.Kinds = "texture,sound"
```

### Config

```json
{
  "roots": ["assets"],
  "kinds": [
    {"kind": "texture", "dir": "textures"},
    {"kind": "sound", "dir": "audio"}
  ],
  "package": "assets"
}
```

---

## Single Root Layout

### CLI

```bash
go run github.com/G-Team-Games/raygolib/cmd/assetgen \
  -root ./assets \
  -pkg assets \
  -out assets/generated.go \
  -single-root \
  -glob "*.png|*.ogg"
```

### Programmatic

```go
cfg := assetgen.DefaultConfig()
cfg.Root = "./assets"
cfg.Package = "assets"
cfg.Output = "assets/generated.go"
cfg.Kinds = "texture,sound"
cfg.SingleRoot = true
cfg.Glob = "*.png|*.ogg"
```

### Config

```json
{
  "roots": ["assets"],
  "kinds": [
    {"kind": "texture", "scan_root": true, "glob": "*.png"},
    {"kind": "sound", "scan_root": true, "glob": "*.ogg"}
  ],
  "package": "assets"
}
```

---

## Mixed Layout (Subdirs + Root)

### CLI

Config file required for mixed layout.

```json
{
  "roots": ["assets"],
  "kinds": [
    {"kind": "texture", "dir": "textures"},
    {"kind": "sound", "dir": "audio"},
    {"kind": "font", "dir": ".", "glob": "*.ttf"},
    {"kind": "shader", "dir": ".", "glob": "*.glsl"}
  ],
  "package": "assets"
}
```

```bash
go run github.com/G-Team-Games/raygolib/cmd/assetgen -config config.json
```

### Programmatic

```go
cfg := assetgen.DefaultConfig()
cfg.Root = "./assets"
cfg.Package = "assets"
cfg.Output = "assets/generated.go"

cfg.KindsConfig = []assetgen.KindConfig{
    {Kind: assetgen.KindTexture, Dir: "textures"},
    {Kind: assetgen.KindSound, Dir: "audio"},
    {Kind: assetgen.KindFont, Dir: ".", Glob: "*.ttf", Type: "FontName"},
    {Kind: assetgen.KindShader, Dir: ".", Glob: "*.glsl", Type: "ShaderName"},
}
```

### Config

```json
{
  "roots": ["assets"],
  "kinds": [
    {"kind": "texture", "dir": "textures"},
    {"kind": "sound", "dir": "audio"},
    {"kind": "font", "dir": ".", "glob": "*.ttf"},
    {"kind": "shader", "dir": ".", "glob": "*.glsl"}
  ],
  "package": "assets"
}
```

---

## CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-root` | Root directory | `assets` |
| `-out` | Output file | `zz_generated_assets.go` |
| `-pkg` | Package name | `gameassets` |
| `-config` | JSON config file | - |
| `-naming` | `pascal`, `camel`, `snake`, `upper_snake` | `pascal` |
| `-kinds` | Comma-separated kinds | all |
| `-single-root` | Scan root instead of subdirs | false |
| `-glob` | File pattern filter | - |

## KindConfig Fields

| Field | Type | Description |
|-------|------|-------------|
| `kind` | string | `texture`, `sound`, `font`, etc. |
| `dir` | string | Subdirectory (empty/`.` = root) |
| `glob` | string | Pattern filter (`*.png|*.jpg`) |
| `type` | string | Generated type name |
| `scan_root` | bool | Also scan root |

## Using in Code

```go
import "mypkg/assets"

playerTex := rl.LoadTexture(string(assets.Player))
clickSfx := rl.LoadSound(string(assets.Click))
```

## Output

```go
package assets

type TextureName string
const Player TextureName = "player.png"

type SoundName string
const Click SoundName = "click.ogg"
```