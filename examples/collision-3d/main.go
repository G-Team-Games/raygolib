package main

import (
	"fmt"

	rgl "github.com/G-Team-Games/raygolib"
	col3d "github.com/G-Team-Games/raygolib/physics/collision/3d"
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
	box     *col3d.BoxCollider
	cyl     *col3d.CylinderCollider
	point   *col3d.PointCollider
	plane   *col3d.PlaneCollider
	contact col3d.Contact
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
		box:   col3d.NewBoxColliderV(rl.NewVector3(-6, 0, 0), rl.NewVector3(2, 2, 2)),
		cyl:   col3d.NewCylinderColliderV(rl.NewVector3(-2, 0, 1), 1.0, 2.0),
		point: col3d.NewPointXYZ(2, 1, 1),
		plane: col3d.NewPlaneCollider(rl.NewVector3(6, 0, 0), 4, 4, col3d.PlaneAxisXPos),
		states: []colliderState{
			{name: "Box", color: rl.Blue},
			{name: "Cylinder", color: rl.Green},
			{name: "Point", color: rl.Red},
			{name: "Plane", color: rl.Purple},
		},
	}
}

func (g *Game) colliders() []col3d.Collider {
	return []col3d.Collider{g.box, g.cyl, g.point, g.plane}
}

func (g *Game) activeCollider() col3d.SpatialCollider {
	return g.colliders()[g.active].(col3d.SpatialCollider)
}

func (g *Game) Init() error {
	return nil
}

func (g *Game) Close() error {
	return nil
}

func (g *Game) Update(dt float32) error {
	g.moveCamera(dt)

	if rl.IsKeyPressed(rl.KeyTab) {
		g.active = (g.active + 1) % len(g.states)
	}

	g.moveActive(dt)
	g.contact = col3d.Contact{}
	activeCol := g.activeCollider()

	col3d.ResolveMultiMTV(activeCol, g.colliders(), 4)

	return nil
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(g.camera)
	rl.DrawGrid(20, 1)
	for i, collider := range g.colliders() {
		fill := rl.Fade(g.states[i].color, 0.30)
		col3d.DrawCollider(collider, fill)
		col3d.DrawColliderWires(collider, g.states[i].color)
	}

	col3d.DrawCollider(g.activeCollider(), rl.Orange)
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
	case 2:
		g.point.SetPosition(rl.Vector3Add(g.point.GetPosition(), move))
	case 3:
		g.plane.SetPosition(rl.Vector3Add(g.plane.GetPosition(), move))
	}

}

func (g *Game) drawUI() {
	active := g.states[g.active].name
	camPos := g.camera.Position
	activeCol := g.activeCollider()

	rl.DrawText("Collision 3D demo", 16, 12, 24, rl.Black)
	rl.DrawText("TAB: switch collider", 16, 42, 20, rl.DarkGray)
	rl.DrawText("Camera: WASD + Q/E | MMB drag rotate", 16, 66, 20, rl.DarkGray)
	rl.DrawText("Collider: Arrows + Space/LShift", 16, 90, 20, rl.DarkGray)
	rl.DrawText("Active: "+active, 16, 120, 22, rl.Black)
	rl.DrawText(fmt.Sprintf("Camera pos: (%.2f, %.2f, %.2f)", camPos.X, camPos.Y, camPos.Z), 16, 148, 20, rl.Brown)

	// Display position and distance stats
	y := int32(172)
	for i, col := range g.colliders() {
		var pos rl.Vector3
		switch c := col.(type) {
		case *col3d.BoxCollider:
			pos = c.GetPosition()
		case *col3d.CylinderCollider:
			pos = c.GetPosition()
		case *col3d.PointCollider:
			pos = c.GetPosition()
		case *col3d.PlaneCollider:
			pos = c.GetPosition()
		}
		state := g.states[i]

		if i == g.active {
			rl.DrawText(fmt.Sprintf("%s pos: (%.2f, %.2f, %.2f) [ACTIVE]", state.name, pos.X, pos.Y, pos.Z), 16, y, 20, state.color)
		} else {
			dist := activeCol.DistanceTo(col)
			rl.DrawText(fmt.Sprintf("%s pos: (%.2f, %.2f, %.2f) | Dist: %.2f", state.name, pos.X, pos.Y, pos.Z, dist), 16, y, 20, state.color)
		}
		y += 24
	}

	status := "No collision"
	if g.contact.Hit {
		status = "Collision detected! Active shifted."
	}
	rl.DrawText(status, 16, y+12, 22, rl.Maroon)
}

func main() {
	game := rgl.InitGameWithConfig(NewGame(), &rgl.InitGameConfig{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		WindowTitle:  "3D collision demo",
		TargetFPS:    60,
	})

	if err := game.Run(); err != nil {
		panic(err)
	}
}
