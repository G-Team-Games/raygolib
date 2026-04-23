package main

import (
	"fmt"

	rgl "github.com/G-Team-Games/raygolib"
	col3d "github.com/G-Team-Games/raygolib/physics/collision/3d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	playerSpeed    float32 = 8.0
	projectileVel  float32 = 15.0
	gravity        float32 = -20.0
	cameraSpeed    float32 = 8.0
	cameraTurn     float32 = 0.0035
)

type Game struct {
	camera rl.Camera3D

	playerBox *col3d.BoxCollider
	playerVel rl.Vector3

	projectile *col3d.PointCollider
	projVel    rl.Vector3
	fired      bool

	walls []*col3d.BoxCollider
	floor  *col3d.BoxCollider
}

func NewGame() *Game {
	floor := col3d.NewBoxColliderV(rl.NewVector3(0, -0.5, 0), rl.NewVector3(20, 1, 20))
	return &Game{
		camera: rl.Camera3D{
			Position:   rl.NewVector3(0, 15, 15),
			Target:     rl.NewVector3(0, 0, 0),
			Up:         rl.NewVector3(0, 1, 0),
			Fovy:       60,
			Projection: rl.CameraPerspective,
		},
		playerBox: col3d.NewBoxColliderV(rl.NewVector3(0, 0.5, 0), rl.NewVector3(1, 1, 1)),
		playerVel: rl.NewVector3(0, 0, 0),
		projectile: col3d.NewPointXYZ(0, 0.5, -3),
		projVel:    rl.NewVector3(0, 0, 0),
		floor: floor,
		walls: []*col3d.BoxCollider{
			col3d.NewBoxColliderV(rl.NewVector3(-8, 0, 0), rl.NewVector3(1, 3, 10)),
			col3d.NewBoxColliderV(rl.NewVector3(8, 0, 0), rl.NewVector3(1, 3, 10)),
			col3d.NewBoxColliderV(rl.NewVector3(0, 0, -8), rl.NewVector3(10, 3, 1)),
			col3d.NewBoxColliderV(rl.NewVector3(0, 0, 8), rl.NewVector3(10, 3, 1)),
			col3d.NewBoxColliderV(rl.NewVector3(-3, 0, -3), rl.NewVector3(1, 3, 1)),
			col3d.NewBoxColliderV(rl.NewVector3(3, 0, 3), rl.NewVector3(1, 3, 1)),
		},
	}
}

func (g *Game) Init() error {
	return nil
}

func (g *Game) Close() error {
	return nil
}

func (g *Game) Update(dt float32) error {
	g.moveCamera(dt)
	g.movePlayer(dt)
	g.updateProjectile(dt)
	return nil
}

func (g *Game) moveCamera(dt float32) {
	if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
		delta := rl.GetMouseDelta()
		rl.CameraYaw(&g.camera, -delta.X*cameraTurn, 0)
		rl.CameraPitch(&g.camera, -delta.Y*cameraTurn, 1, 0, 0)
	}

	step := cameraSpeed * dt
	move := rl.NewVector3(0, 0, 0)

	forward := rl.Vector3Subtract(g.camera.Target, g.camera.Position)
	forward.Y = 0
	if rl.Vector3Length(forward) > 0 {
		forward = rl.Vector3Normalize(forward)
	}
	right := rl.Vector3Normalize(rl.Vector3CrossProduct(forward, g.camera.Up))

	if rl.IsKeyDown(rl.KeyW) {
		move = rl.Vector3Add(move, rl.Vector3Scale(forward, step))
	}
	if rl.IsKeyDown(rl.KeyS) {
		move = rl.Vector3Add(move, rl.Vector3Scale(forward, -step))
	}
	if rl.IsKeyDown(rl.KeyA) {
		move = rl.Vector3Add(move, rl.Vector3Scale(right, -step))
	}
	if rl.IsKeyDown(rl.KeyD) {
		move = rl.Vector3Add(move, rl.Vector3Scale(right, step))
	}
	if rl.IsKeyDown(rl.KeyQ) {
		move.Y += step
	}
	if rl.IsKeyDown(rl.KeyE) {
		move.Y -= step
	}

	if move == (rl.Vector3{}) {
		return
	}

	g.camera.Position = rl.Vector3Add(g.camera.Position, move)
	g.camera.Target = rl.Vector3Add(g.camera.Target, move)
}

func (g *Game) movePlayer(dt float32) {
	input := rl.NewVector3(0, 0, 0)
	speed := playerSpeed * dt

	if rl.IsKeyDown(rl.KeyUp) {
		input.Z -= speed
	}
	if rl.IsKeyDown(rl.KeyDown) {
		input.Z += speed
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		input.X -= speed
	}
	if rl.IsKeyDown(rl.KeyRight) {
		input.X += speed
	}
	if rl.IsKeyDown(rl.KeySpace) {
		input.Y += speed
	}
	if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
		input.Y -= speed
	}

	if input != (rl.Vector3{}) {
		g.playerVel = input
	} else {
		g.playerVel = rl.NewVector3(0, 0, 0)
	}

	others := make([]col3d.Collider, len(g.walls)+1)
	others[0] = g.floor
	for i, w := range g.walls {
		others[i+1] = w
	}

	if g.playerVel != (rl.Vector3{}) {
		g.playerBox.SetPosition(rl.Vector3Add(g.playerBox.GetPosition(), g.playerVel))
		col3d.ResolveMultiMTV(g.playerBox, others, 4)
	}
}

func (g *Game) updateProjectile(dt float32) {
	if rl.IsKeyPressed(rl.KeyF) && !g.fired {
		g.fired = true
		dir := rl.Vector3Normalize(rl.Vector3Subtract(g.camera.Target, g.camera.Position))
		g.projVel = rl.Vector3Scale(dir, projectileVel)
		g.projectile.SetPosition(g.playerBox.Center())
	}

	if g.fired {
		g.projVel.Y += gravity * dt

		others := make([]col3d.Collider, len(g.walls)+2)
		others[0] = g.floor
		others[1] = g.playerBox
		for i, w := range g.walls {
			others[i+2] = w
		}

		g.projectile.SetPosition(rl.Vector3Add(g.projectile.GetPosition(), rl.Vector3Scale(g.projVel, dt)))
		col3d.ResolveMultiMTV(g.projectile, others, 4)

		if g.projectile.GetPosition().Y < -10 {
			g.fired = false
		}
	}
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(g.camera)
	rl.DrawGrid(20, 1)

	rl.DrawCubeV(g.floor.Center(), g.floor.Size, rl.Fade(rl.LightGray, 0.5))
	rl.DrawCubeWiresV(g.floor.Center(), g.floor.Size, rl.Gray)

	rl.DrawCubeV(g.playerBox.Center(), g.playerBox.Size, rl.Fade(rl.Blue, 0.5))
	rl.DrawCubeWiresV(g.playerBox.Center(), g.playerBox.Size, rl.Blue)

	for _, w := range g.walls {
		rl.DrawCubeV(w.Center(), w.Size, rl.Fade(rl.DarkGray, 0.3))
		rl.DrawCubeWiresV(w.Center(), w.Size, rl.DarkGray)
	}

	if g.fired {
		rl.DrawSphere(g.projectile.GetPosition(), 0.3, rl.Red)
	}

	rl.EndMode3D()
	g.drawUI()
}

func (g *Game) drawUI() {
	rl.DrawText("Collision Resolvers Demo", 16, 12, 24, rl.Black)
	rl.DrawText("WASD + Space/LShift: move player", 16, 42, 20, rl.DarkGray)
	rl.DrawText("F: fire projectile (physics enabled)", 16, 66, 20, rl.DarkGray)
	rl.DrawText("Camera: WASD + Q/E | MMB drag", 16, 90, 20, rl.DarkGray)

	if g.fired {
		rl.DrawText("Projectile: ACTIVE", 16, 142, 20, rl.Red)
	} else {
		rl.DrawText("Projectile: READY (press F)", 16, 142, 20, rl.Green)
	}
}

func main() {
	game := rgl.InitGameWithConfig(NewGame(), &rgl.InitGameConfig{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		WindowTitle:  "raygolib example - collision-resolvers",
		TargetFPS:    60,
	})

	if err := game.Run(); err != nil {
		fmt.Println(err)
	}
}