package main

import (
	"fmt"

	rgl "github.com/G-Team-Games/raygolib"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/G-Team-Games/raygolib/examples/assetgen/assets"
)

type Game struct {
	playerTex rl.Texture2D
	enemyTex  rl.Texture2D
	clickSfx  rl.Sound
}

func (g *Game) Init() error {
	rl.InitAudioDevice()
	g.playerTex = rl.LoadTexture(string(assets.Player))
	g.enemyTex = rl.LoadTexture(string(assets.Enemy))
	g.clickSfx = rl.LoadSound(string(assets.Click))
	return nil
}

func (g *Game) Update(dt float32) error {
	if rl.IsKeyPressed(rl.KeySpace) {
		rl.PlaySound(g.clickSfx)
	}
	return nil
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTexture(g.playerTex, 100, 200, rl.White)
	rl.DrawTexture(g.enemyTex, 500, 200, rl.White)
	rl.DrawText("Asset Name Gen Demo", 200, 60, 24, rl.Black)
	rl.DrawText("SPACE = play click sound", 250, 400, 20, rl.DarkGray)
}

func (g *Game) Close() error {
	rl.UnloadTexture(g.playerTex)
	rl.UnloadTexture(g.enemyTex)
	rl.UnloadSound(g.clickSfx)
	rl.CloseAudioDevice()
	return nil
}

func main() {
	game := rgl.InitGameWithConfig(&Game{}, &rgl.InitGameConfig{
		ScreenWidth:  800,
		ScreenHeight: 450,
		WindowTitle:  "Asset Name Gen Demo",
		TargetFPS:    60,
	})

	if err := game.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}