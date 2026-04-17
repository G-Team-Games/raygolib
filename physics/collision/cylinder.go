package collision

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
		return cylinderVsCylinderContact(*c, *o)
	case *BoxCollider:
		return cylinderVsBoxContact(*c, *o)
	case *PlaneCollider:
		return cylinderVsPlaneContact(*c, *o)
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



// cylinderVsPlaneContact computes contact between cylinder and finite plane.
func cylinderVsPlaneContact(cylinder CylinderCollider, plane PlaneCollider) Contact {
	switch plane.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg, PlaneAxisZPos, PlaneAxisZNeg:
		var difference rl.Vector2
		if plane.Axis == PlaneAxisZPos || plane.Axis == PlaneAxisZNeg {
			difference = rl.Vector2Subtract(
				rl.NewVector2(cylinder.Position.X, cylinder.Position.Z),
				rl.NewVector2(
					math32.Min(plane.Position.X+plane.Width, math32.Max(plane.Position.X, cylinder.Position.X)),
					plane.Position.Z,
				),
			)
		} else {
			difference = rl.Vector2Subtract(
				rl.NewVector2(cylinder.Position.X, cylinder.Position.Z),
				rl.NewVector2(
					plane.Position.X,
					math32.Min(plane.Position.Z+plane.Width, math32.Max(plane.Position.Z, cylinder.Position.Z)),
				),
			)
		}

		distanceXZ := rl.Vector2Length(difference)
		penetrationXZ := cylinder.Radius - distanceXZ
		distanceY1 := cylinder.Position.Y - (plane.Position.Y + plane.Height)
		distanceY2 := plane.Position.Y - (cylinder.Position.Y + cylinder.Height)

		if penetrationXZ < 0 || distanceY1 > 0 || distanceY2 > 0 {
			return Contact{}
		}

		normalXZ := safeNormalize2(difference)
		if rl.Vector2Length(normalXZ) == 0 {
			normal3 := plane.Axis.Normal()
			normalXZ = rl.NewVector2(-normal3.X, -normal3.Z)
		}
		normal := rl.NewVector3(normalXZ.X, 0, normalXZ.Y)
		return Contact{Hit: true, Normal: normal, Penetration: penetrationXZ}

	case PlaneAxisYPos, PlaneAxisYNeg:
		difference := rl.Vector2Subtract(
			rl.NewVector2(cylinder.Position.X, cylinder.Position.Z),
			rl.NewVector2(
				math32.Min(plane.Position.X+plane.Width, math32.Max(plane.Position.X, cylinder.Position.X)),
				math32.Min(plane.Position.Z+plane.Height, math32.Max(plane.Position.Z, cylinder.Position.Z)),
			),
		)
		distanceXZ := rl.Vector2Length(difference)
		penetrationXZ := cylinder.Radius - distanceXZ
		distanceY1 := cylinder.Position.Y - plane.Position.Y
		distanceY2 := plane.Position.Y - (cylinder.Position.Y + cylinder.Height)

		if penetrationXZ < 0 || distanceY1 > 0 || distanceY2 > 0 {
			return Contact{}
		}

		if plane.Axis == PlaneAxisYNeg {
			return Contact{
				Hit:         true,
			
				Normal:      rl.NewVector3(0, 1, 0),
				Penetration: -distanceY2,
			}
		}

		return Contact{
			Hit:         true,
			Normal:      rl.NewVector3(0, -1, 0),
			Penetration: -distanceY1,
		}
	}

	return Contact{}
}


var _ Collider = (*CylinderCollider)(nil)
