package col3d

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CylinderCollider struct {
	Position rl.Vector3
	Radius   float32
	Height   float32
}

// NewCylinderCollider creates upright cylinder collider.
func NewCylinderCollider(position rl.Vector3, radius float32, height float32) *CylinderCollider {
	return &CylinderCollider{
		Position: position,
		Radius:   radius,
		Height:   height,
	}
}

// Kind returns cylinder collider type.
func (c *CylinderCollider) Kind() ShapeKind {
	return ShapeCylinderY
}

// Collide returns contact between cylinder and another collider.
func (c *CylinderCollider) Collide(other Collider) Contact {
	switch o := other.(type) {
	case *CylinderCollider:
		return cylinderVsCylinderContact(c, o)
	case *BoxCollider:
		return cylinderVsBoxContact(c, o)
	case *PlaneCollider:
		return cylinderVsPlaneContact(c, o)
	default:
		return Contact{}
	}
}

// BoundingBox returns enclosing AABB for broad-phase checks.
func (c *CylinderCollider) BoundingBox() rl.BoundingBox {
	min := rl.NewVector3(c.Position.X-c.Radius, c.Position.Y, c.Position.Z-c.Radius)
	max := rl.NewVector3(c.Position.X+c.Radius, c.Position.Y+c.Height, c.Position.Z+c.Radius)
	return rl.NewBoundingBox(min, max)
}

// Center returns geometric center of cylinder volume.
func (c *CylinderCollider) Center() rl.Vector3 {
	return rl.NewVector3(c.Position.X, c.Position.Y+c.Height*0.5, c.Position.Z)
}

// GetSides returns two tangent offsets on cylinder edge from point direction.
func (c *CylinderCollider) GetSides(position rl.Vector2) (rl.Vector2, rl.Vector2) {
	cylinderPosition := rl.NewVector2(c.Position.X, c.Position.Z)
	difference := rl.Vector2Subtract(cylinderPosition, position)
	normalized := safeNormalize2(difference)
	return rl.Vector2Scale(rl.NewVector2(-normalized.Y, normalized.X), c.Radius), rl.Vector2Scale(rl.NewVector2(normalized.Y, -normalized.X), c.Radius)
}

// GetPosition returns cylinder base-center position.
func (c *CylinderCollider) GetPosition() rl.Vector3 {
	return c.Position
}

// SetPosition sets cylinder base-center position.
func (c *CylinderCollider) SetPosition(vec rl.Vector3) {
	c.Position = vec
}

func (c *CylinderCollider) DistanceTo(other Collider) float32 {
	switch o := other.(type) {
	case *CylinderCollider:
		return cylinderVsCylinderDistance(*c, *o)
	case *BoxCollider:
		return cylinderVsBoxDistance(*c, *o)
	case *PointCollider:
		return cylinderVsPointDistance(*c, *o)
	case *PlaneCollider:
		return cylinderVsPlaneDistance(*c, *o)
	default:
		return 0
	}
}

func cylinderVsCylinderDistance(a, b CylinderCollider) float32 {
	dx := a.Center().X - b.Center().X
	dz := a.Center().Z - b.Center().Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)
	distY1 := a.Position.Y - (b.Position.Y + b.Height)
	distY2 := b.Position.Y - (a.Position.Y + a.Height)
	radiusSum := a.Radius + b.Radius

	sideDist := distXZ - radiusSum
	botDist := -distY1
	topDist := -distY2

	distSq := sideDist*sideDist + botDist*botDist + topDist*topDist
	return math32.Sqrt(distSq)
}

func cylinderVsBoxDistance(cyl CylinderCollider, box BoxCollider) float32 {
	closestX := math32.Max(cyl.Position.X-cyl.Radius, math32.Min(cyl.Position.X+cyl.Radius, box.Position.X))
	closestY := math32.Max(cyl.Position.Y, math32.Min(cyl.Position.Y+cyl.Height, box.Position.Y))
	closestZ := math32.Max(cyl.Position.Z-cyl.Radius, math32.Min(cyl.Position.Z+cyl.Radius, box.Position.Z))

	dx := box.Center().X - closestX
	dy := box.Center().Y - closestY
	dz := box.Center().Z - closestZ

	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

func cylinderVsPointDistance(cyl CylinderCollider, pt PointCollider) float32 {
	closestX := math32.Max(cyl.Position.X-cyl.Radius, math32.Min(cyl.Position.X+cyl.Radius, pt.Position.X))
	closestZ := math32.Max(cyl.Position.Z-cyl.Radius, math32.Min(cyl.Position.Z+cyl.Radius, pt.Position.Z))
	closestY := math32.Max(cyl.Position.Y, math32.Min(cyl.Position.Y+cyl.Height, pt.Position.Y))

	dx := pt.Position.X - closestX
	dy := pt.Position.Y - closestY
	dz := pt.Position.Z - closestZ

	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

func cylinderVsPlaneDistance(cyl CylinderCollider, plane PlaneCollider) float32 {
	closest := closestPointOnPlane(plane, cyl.Center())
	dx := cyl.Center().X - closest.X
	dy := cyl.Center().Y - closest.Y
	dz := cyl.Center().Z - closest.Z
	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

func closestPointOnPlane(plane PlaneCollider, point rl.Vector3) rl.Vector3 {
	switch plane.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		return rl.NewVector3(plane.Position.X, math32.Max(plane.Position.Y, math32.Min(point.Y, plane.Position.Y+plane.Height)), math32.Max(plane.Position.Z, math32.Min(point.Z, plane.Position.Z+plane.Width)))
	case PlaneAxisZPos, PlaneAxisZNeg:
		return rl.NewVector3(math32.Max(plane.Position.X, math32.Min(point.X, plane.Position.X+plane.Width)), math32.Max(plane.Position.Y, math32.Min(point.Y, plane.Position.Y+plane.Height)), plane.Position.Z)
	default:
		return rl.NewVector3(math32.Max(plane.Position.X, math32.Min(point.X, plane.Position.X+plane.Width)), plane.Position.Y, math32.Max(plane.Position.Z, math32.Min(point.Z, plane.Position.Z+plane.Height)))
	}
}



var _ Collider = (*CylinderCollider)(nil)
