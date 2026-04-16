package main

import (
	"fmt"

	rgl "github.com/G-Team-Games/raygolib"
	"github.com/G-Team-Games/raygolib/physics/collision"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	colliderMoveSpeed float32 = 4.0
	cameraMoveSpeed   float32 = 8.0
	cameraTurnSpeed   float32 = 0.0035
)

type Game struct {
	camera  rl.Camera3D
	active  int
	box     *collision.BoxCollider
	cyl     *collision.CylinderCollider
	contact collision.Contact
	states  []colliderState
}

type colliderState struct {
	name  string
	color rl.Color
}

func NewGame() *Game {
	return &Game{
		camera: rl.Camera3D{
			Position:   rl.NewVector3(9, 9, 9),
			Target:     rl.NewVector3(0, 1.5, 0),
			Up:         rl.NewVector3(0, 1, 0),
			Fovy:       45,
			Projection: rl.CameraPerspective,
		},
		box: collision.NewBoxCollider(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
		cyl: collision.NewCylinderCollider(rl.NewVector3(3, 0, 3), 1.0, 2.0),
		states: []colliderState{
			{name: "Box", color: rl.Blue},
			{name: "Cylinder", color: rl.Green},
		},
	}
}

func (g *Game) Update(dt float32) error {
	g.moveCamera(dt)

	if rl.IsKeyPressed(rl.KeyTab) {
		g.active = (g.active + 1) % len(g.states)
	}

	g.contact = collision.Contact{}
	g.moveActive(dt)

	g.contact = g.box.Collide(g.cyl)
	collision.ResolveByMTV(g.box.GetPosition, g.box.SetPosition, g.contact)

	return nil
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(g.camera)
	rl.DrawGrid(20, 1)
	g.drawBox()
	g.drawCylinder()
	g.drawPositionPoints()
	g.drawContact()
	rl.EndMode3D()
	g.drawUI()
}

func (g *Game) moveCamera(dt float32) {
	if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
		delta := rl.GetMouseDelta()
		rl.CameraYaw(&g.camera, -delta.X*cameraTurnSpeed, 0)
		rl.CameraPitch(&g.camera, -delta.Y*cameraTurnSpeed, 1, 0, 0)
	}

	step := cameraMoveSpeed * dt
	move := rl.NewVector3(0, 0, 0)

	forward := rl.Vector3Subtract(g.camera.Target, g.camera.Position)
	forward.Y = 0
	if rl.Vector3Length(forward) > 0 {
		forward = rl.Vector3Normalize(forward)
	} else {
		forward = rl.NewVector3(0, 0, -1)
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

func (g *Game) moveActive(dt float32) {
	move := rl.NewVector3(0, 0, 0)
	step := colliderMoveSpeed * dt

	if rl.IsKeyDown(rl.KeyLeft) {
		move.X -= step
	}
	if rl.IsKeyDown(rl.KeyRight) {
		move.X += step
	}
	if rl.IsKeyDown(rl.KeyUp) {
		move.Z -= step
	}
	if rl.IsKeyDown(rl.KeyDown) {
		move.Z += step
	}
	if rl.IsKeyDown(rl.KeySpace) {
		move.Y += step
	}
	if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
		move.Y -= step
	}

	if move == (rl.Vector3{}) {
		return
	}

	switch g.active {
	case 0:
		g.box.SetPosition(rl.Vector3Add(g.box.GetPosition(), move))
	case 1:
		g.cyl.SetPosition(rl.Vector3Add(g.cyl.GetPosition(), move))
	}

}

func (g *Game) drawBox() {
	fill := rl.Fade(g.states[0].color, 0.30)
	if g.active == 0 {
		fill = rl.Fade(rl.Orange, 0.40)
	}
	rl.DrawCubeV(g.box.Center(), g.box.Size, fill)
	rl.DrawCubeWiresV(g.box.Center(), g.box.Size, g.states[0].color)
}

func (g *Game) drawCylinder() {
	fill := rl.Fade(g.states[1].color, 0.30)
	if g.active == 1 {
		fill = rl.Fade(rl.Orange, 0.40)
	}
	rl.DrawCylinder(g.cyl.Position, g.cyl.Radius, g.cyl.Radius, g.cyl.Height, 24, fill)
	rl.DrawCylinderWires(g.cyl.Position, g.cyl.Radius, g.cyl.Radius, g.cyl.Height, 24, g.states[1].color)
}

func (g *Game) drawContact() {
	if !g.contact.Hit {
		return
	}
	rl.DrawSphere(g.contact.Point, 0.12, rl.Red)
	end := rl.Vector3Add(g.contact.Point, rl.Vector3Scale(g.contact.Normal, 1.0))
	rl.DrawLine3D(g.contact.Point, end, rl.Maroon)
}

func (g *Game) drawPositionPoints() {
	boxPos := g.box.GetPosition()
	cylPos := g.cyl.GetPosition()

	rl.DrawSphere(boxPos, 0.10, rl.DarkBlue)
	rl.DrawSphere(cylPos, 0.10, rl.DarkGreen)

	rl.DrawLine3D(boxPos, rl.Vector3Add(boxPos, rl.NewVector3(0, 0.45, 0)), rl.DarkBlue)
	rl.DrawLine3D(cylPos, rl.Vector3Add(cylPos, rl.NewVector3(0, 0.45, 0)), rl.DarkGreen)
}

func (g *Game) drawUI() {
	active := g.states[g.active].name
	boxPos := g.box.GetPosition()
	cylPos := g.cyl.GetPosition()
	camPos := g.camera.Position

	rl.DrawText("Collision 3D demo", 16, 12, 24, rl.Black)
	rl.DrawText("TAB: switch collider", 16, 42, 20, rl.DarkGray)
	rl.DrawText("Camera: WASD + Q/E | MMB drag rotate", 16, 66, 20, rl.DarkGray)
	rl.DrawText("Collider: Arrows + Space/LShift (plane static)", 16, 90, 20, rl.DarkGray)
	rl.DrawText("Active: "+active, 16, 120, 22, rl.Black)
	rl.DrawText(fmt.Sprintf("Camera pos: (%.2f, %.2f, %.2f)", camPos.X, camPos.Y, camPos.Z), 16, 148, 20, rl.Brown)
	rl.DrawText(fmt.Sprintf("Box pos: (%.2f, %.2f, %.2f)", boxPos.X, boxPos.Y, boxPos.Z), 16, 172, 20, rl.DarkBlue)
	rl.DrawText(fmt.Sprintf("Cylinder pos: (%.2f, %.2f, %.2f)", cylPos.X, cylPos.Y, cylPos.Z), 16, 196, 20, rl.DarkGreen)

	status := "No collision"
	if g.contact.Hit {
		status = "Collision detected"
	}
	rl.DrawText(status, 16, 248, 22, rl.Maroon)
}

func main() {
	game := rgl.InitGameWithConfig(NewGame(), &rgl.InitGameConfig{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		WindowTitle:  "raygolib example - collision-3d",
		TargetFPS:    60,
	})

	if err := game.Run(); err != nil {
		panic(err)
	}
}
