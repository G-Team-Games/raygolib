package rgcol3d

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	defaultPointRadius    float32 = 0.08
	defaultCylinderSlices int32   = 24
	defaultSphereRings    int32   = 8
	defaultSphereSlices   int32   = 8
	defaultPlaneThickness float32 = 0.01
)

func DrawCollider(c Collider, color rl.Color) {
	drawCollider(c, color)
}

func DrawColliders(colliders []Collider, color rl.Color) {
	for _, c := range colliders {
		DrawCollider(c, color)
	}
}

func DrawColliderWires(c Collider, color rl.Color) {
	drawColliderWires(c, color)
}

func DrawCollidersWires(colliders []Collider, color rl.Color) {
	for _, c := range colliders {
		DrawColliderWires(c, color)
	}
}

func drawCollider(c Collider, color rl.Color) {
	switch shape := c.(type) {
	case *BoxCollider:
		rl.DrawCubeV(shape.Center(), shape.Size, color)
	case *CylinderCollider:
		rl.DrawCylinder(shape.Position, shape.Radius, shape.Radius, shape.Height, defaultCylinderSlices, color)
	case *PlaneCollider:
		center, size := planeDrawBox(shape)
		rl.DrawCubeV(center, size, color)
	case *PointCollider:
		rl.DrawSphere(shape.Position, defaultPointRadius, color)
	}
}

func drawColliderWires(c Collider, color rl.Color) {
	switch shape := c.(type) {
	case *BoxCollider:
		rl.DrawCubeWiresV(shape.Center(), shape.Size, color)
	case *CylinderCollider:
		rl.DrawCylinderWires(shape.Position, shape.Radius, shape.Radius, shape.Height, defaultCylinderSlices, color)
	case *PlaneCollider:
		center, size := planeDrawBox(shape)
		rl.DrawCubeWiresV(center, size, color)
	case *PointCollider:
		rl.DrawSphereWires(shape.Position, defaultPointRadius, defaultSphereRings, defaultSphereSlices, color)
	}
}

func planeDrawBox(plane *PlaneCollider) (rl.Vector3, rl.Vector3) {
	box := plane.BoundingBox()
	size := rl.Vector3Subtract(box.Max, box.Min)

	if size.X == 0 {
		size.X = defaultPlaneThickness
	}
	if size.Y == 0 {
		size.Y = defaultPlaneThickness
	}
	if size.Z == 0 {
		size.Z = defaultPlaneThickness
	}

	center := rl.Vector3Add(box.Min, rl.Vector3Scale(size, 0.5))
	return center, size
}
