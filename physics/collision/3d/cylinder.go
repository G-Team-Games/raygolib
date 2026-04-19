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
		return unsupportedDistance()
	}
}

func cylinderVsCylinderDistance(a, b CylinderCollider) float32 {
	dx := a.Position.X - b.Position.X
	dz := a.Position.Z - b.Position.Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)
	horizontalGap := math32.Max(0, distXZ-(a.Radius+b.Radius))
	verticalGap := intervalGap(a.Position.Y, a.Position.Y+a.Height, b.Position.Y, b.Position.Y+b.Height)

	return combineOrthogonalGaps(horizontalGap, verticalGap)
}

func cylinderVsBoxDistance(cyl CylinderCollider, box BoxCollider) float32 {
	boxMax := box.Max()

	horizontalGap := circleRectGapXZ(
		cyl.Position.X,
		cyl.Position.Z,
		cyl.Radius,
		box.Position.X,
		boxMax.X,
		box.Position.Z,
		boxMax.Z,
	)
	verticalGap := intervalGap(cyl.Position.Y, cyl.Position.Y+cyl.Height, box.Position.Y, boxMax.Y)

	return combineOrthogonalGaps(horizontalGap, verticalGap)
}

func cylinderVsPointDistance(cyl CylinderCollider, pt PointCollider) float32 {
	dx := pt.Position.X - cyl.Position.X
	dz := pt.Position.Z - cyl.Position.Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)
	horizontalGap := math32.Max(0, distXZ-cyl.Radius)
	verticalGap := intervalGap(pt.Position.Y, pt.Position.Y, cyl.Position.Y, cyl.Position.Y+cyl.Height)

	return combineOrthogonalGaps(horizontalGap, verticalGap)
}

func cylinderVsPlaneDistance(cyl CylinderCollider, plane PlaneCollider) float32 {
	var gapX, gapY, gapZ float32

	switch plane.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		gapX = intervalGap(cyl.Position.X-cyl.Radius, cyl.Position.X+cyl.Radius, plane.Position.X, plane.Position.X)
		gapY = intervalGap(cyl.Position.Y, cyl.Position.Y+cyl.Height, plane.Position.Y, plane.Position.Y+plane.Height)
		gapZ = intervalGap(cyl.Position.Z-cyl.Radius, cyl.Position.Z+cyl.Radius, plane.Position.Z, plane.Position.Z+plane.Width)
	case PlaneAxisYPos, PlaneAxisYNeg:
		gapX = intervalGap(cyl.Position.X-cyl.Radius, cyl.Position.X+cyl.Radius, plane.Position.X, plane.Position.X+plane.Width)
		gapY = intervalGap(cyl.Position.Y, cyl.Position.Y+cyl.Height, plane.Position.Y, plane.Position.Y)
		gapZ = intervalGap(cyl.Position.Z-cyl.Radius, cyl.Position.Z+cyl.Radius, plane.Position.Z, plane.Position.Z+plane.Height)
	case PlaneAxisZPos, PlaneAxisZNeg:
		gapX = intervalGap(cyl.Position.X-cyl.Radius, cyl.Position.X+cyl.Radius, plane.Position.X, plane.Position.X+plane.Width)
		gapY = intervalGap(cyl.Position.Y, cyl.Position.Y+cyl.Height, plane.Position.Y, plane.Position.Y+plane.Height)
		gapZ = intervalGap(cyl.Position.Z-cyl.Radius, cyl.Position.Z+cyl.Radius, plane.Position.Z, plane.Position.Z)
	}

	return math32.Sqrt(gapX*gapX + gapY*gapY + gapZ*gapZ)
}

var _ Collider = (*CylinderCollider)(nil)
