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

func (g *Game) activeCollider() col3d.Collider {
	return g.colliders()[g.active]
}

func (g *Game) Update(dt float32) error {
	g.moveCamera(dt)

	if rl.IsKeyPressed(rl.KeyTab) {
		g.active = (g.active + 1) % len(g.states)
	}

	g.moveActive(dt)

	g.contact = col3d.Contact{}
	
	activeCol := g.activeCollider()

	// Multi-iteration collision solver to handle overlapping push-backs
	// If active is pushed out of Box into Cylinder, next iter pushes out of Cylinder
	for range 4 {
		hitThisIter := false
		
		for i, other := range g.colliders() {
			if i == g.active {
				continue // Skip self
			}

			// Compute hit from active towards other
			contact := activeCol.Collide(other)
			
			if contact.Hit {
				g.contact = contact
				hitThisIter = true
				
				// type assertion since Collider interface doesn't define Get/SetPosition directly
				var getPos func() rl.Vector3
				var setPos func(rl.Vector3)
				
				switch c := activeCol.(type) {
				case *col3d.BoxCollider:
					getPos, setPos = c.GetPosition, c.SetPosition
				case *col3d.CylinderCollider:
					getPos, setPos = c.GetPosition, c.SetPosition
				case *col3d.PointCollider:
					getPos, setPos = c.GetPosition, c.SetPosition
				case *col3d.PlaneCollider:
					getPos, setPos = c.GetPosition, c.SetPosition
				}

				// Resolve collision using MTV (Minimum Translation Vector)
				col3d.ResolveByMTV(getPos, setPos, contact)
			}
		}
		
		// If no collisions found in this pass, we are fully resolved
		if !hitThisIter {
			break
		}
	}

	return nil
}

func (g *Game) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(g.camera)
	rl.DrawGrid(20, 1)
	g.drawBox()
	g.drawCylinder()
	g.drawPlane()
	g.drawPositionPoints()
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

func (g *Game) drawPlane() {
	fill := rl.Fade(g.states[3].color, 0.30)
	if g.active == 3 {
		fill = rl.Fade(rl.Orange, 0.40)
	}

	// Plane is a 2D quad in 3D space
	// We extract its BoundingBox corners for drawing since DrawPlane assumes ground Y=0
	box := g.plane.BoundingBox()
	size := rl.Vector3Subtract(box.Max, box.Min)
	
	// Ensure minimum thickness so rl.DrawCubeV renders something visible (as plane has 0 depth)
	if size.X == 0 { size.X = 0.01 }
	if size.Y == 0 { size.Y = 0.01 }
	if size.Z == 0 { size.Z = 0.01 }
	
	center := rl.Vector3Add(box.Min, rl.Vector3Scale(size, 0.5))
	
	rl.DrawCubeV(center, size, fill)
	rl.DrawCubeWiresV(center, size, g.states[3].color)
}

func (g *Game) drawPositionPoints() {
	boxPos := g.box.GetPosition()
	cylPos := g.cyl.GetPosition()
	ptPos := g.point.GetPosition()
	planePos := g.plane.GetPosition()

	rl.DrawSphere(boxPos, 0.10, rl.DarkBlue)
	rl.DrawSphere(ptPos, 0.10, rl.Red)
	rl.DrawSphere(cylPos, 0.10, rl.DarkGreen)
	rl.DrawSphere(planePos, 0.10, rl.Purple)

	rl.DrawLine3D(boxPos, rl.Vector3Add(boxPos, rl.NewVector3(0, 0.45, 0)), rl.DarkBlue)
	rl.DrawLine3D(cylPos, rl.Vector3Add(cylPos, rl.NewVector3(0, 0.45, 0)), rl.DarkGreen)
	rl.DrawLine3D(ptPos, rl.Vector3Add(ptPos, rl.NewVector3(0, 0.45, 0)), rl.Red)
	rl.DrawLine3D(planePos, rl.Vector3Add(planePos, rl.NewVector3(0, 0.45, 0)), rl.Purple)
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
			// Calculate minimal distance between volumes.
			// Distance is 0 when touching or overlapping.
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
		WindowTitle:  "raygolib example - collision-3d",
		TargetFPS:    60,
	})

	if err := game.Run(); err != nil {
		panic(err)
	}
}
