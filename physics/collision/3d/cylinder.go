package col3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CylinderCollider struct {
	Position rl.Vector3
	Radius   float32
	Height   float32
}

// NewCylinderColliderV creates upright cylinder collider.
func NewCylinderColliderV(position rl.Vector3, radius float32, height float32) *CylinderCollider {
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
	case *PointCollider:
		contact := cylinderVsPointContact(c, o)
		contact.Normal = rl.Vector3Negate(contact.Normal)
		return contact
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

// DistanceTo returns distance between cylinder and another collider.
func (c *CylinderCollider) DistanceTo(other Collider) float32 {
	switch o := other.(type) {
	case *CylinderCollider:
		return cylinderVsCylinderDistance(c, o)
	case *BoxCollider:
		return cylinderVsBoxDistance(c, o)
	case *PointCollider:
		return cylinderVsPointDistance(c, o)
	case *PlaneCollider:
		return cylinderVsPlaneDistance(c, o)
	default:
		return infiniteDistance()
	}
}

var _ Collider = (*CylinderCollider)(nil)
