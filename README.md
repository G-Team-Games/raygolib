# Raygolib

Zombie-simple game library for Go, built on top of raylib. It provides a simple game loop and structure to build games with raylib in Go.

Note that this library has [raylib-go](https://github.com/gen2brain/raylib-go) ([raylib](https://github.com/raysan5/raylib) binding for Go) as a _dependency_. You can use raylib features like you would in any other raylib project. It it also inspired by [Ebitengine](https://github.com/hajimehoshi/ebiten) and its _dead simplicity_.

> [!NOTE]
> This library is during active development.
> There is no release yet, but if you want to contribute scroll down to know more.

## Quick usage

```go
package main

import (
	"fmt"

	rgl "github.com/G-Team-Games/raygolib"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct{}

// Init for game initialization
func (g *Game) Init() error              { return nil }
// Update for game logic, called every frame with delta time
func (g *Game) Update(dt float32) error  { return nil }
// Draw for rendering/drawing, called every frame after Update
func (g *Game) Draw() {
	rl.ClearBackground(rl.White)
	rl.DrawText("Hello from raygolib", 20, 20, 24, rl.Black)
}
// Close for cleanup when game ends
func (g *Game) Close() error             { return nil }

func main() {
	if err := rgl.InitGame(&Game{}).Run(); err != nil {
		fmt.Println("run error:", err)
	}
}
```

## Examples
Examples live in `examples/`:
- `examples/basic/simple_screen_init`
- `examples/basic/debug_mode`
- `examples/assetgen`
- `examples/asset_demo`
- `examples/collision-3d`

Run any example from its own directory:

```bash
cd examples/basic/simple_screen_init
go run .
```

### Want to contribute?
- You can always open an issue
- If you want to add features or fix bugs, fork the repo and open a PR to this repo with your changes.

**Simple rules to follow while contributing**
- Keep changes focused. Small PRs are easier to review
- Add or update tests when behavior changes
- Run checks before opening PR:
  - `make test`
  - `make coverage`
- Update docs when you add features, flags, or public API changes.
- Follow existing style in code and commit messages.
