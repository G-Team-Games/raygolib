package col3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlaneCollider struct {
	Position rl.Vector3
	Axis     PlaneAxis
	Width    float32
	Height   float32
}

// NewPlaneCollider creates finite axis-aligned plane rectangle.
func NewPlaneCollider(position rl.Vector3, width float32, height float32, axis PlaneAxis) *PlaneCollider {
	return &PlaneCollider{
		Position: position,
		Width:    width,
		Height:   height,
		Axis:     axis,
	}
}

// Kind returns plane-rectangle collider type.
func (p *PlaneCollider) Kind() ShapeKind {
	return ShapePlaneRect
}

// Collide returns contact between plane and another collider.
func (p *PlaneCollider) Collide(other Collider) Contact {
	switch c := other.(type) {
	case *CylinderCollider:
		contact := cylinderVsPlaneContact(c, p)
		contact.Normal = rl.Vector3Negate(contact.Normal)
		return contact
	case *PlaneCollider:
		return Contact{}
	default:
		return Contact{}
	}
}

// BoundingBox returns enclosing AABB of plane rectangle.
func (p *PlaneCollider) BoundingBox() rl.BoundingBox {
	min := p.Position
	max := p.Position
	switch p.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		max.Y += p.Height
		max.Z += p.Width
	case PlaneAxisYPos, PlaneAxisYNeg:
		max.X += p.Width
		max.Z += p.Height
	case PlaneAxisZPos, PlaneAxisZNeg:
		max.X += p.Width
		max.Y += p.Height
	}
	return rl.NewBoundingBox(min, max)
}

// Center returns geometric center of finite plane rectangle.
func (p *PlaneCollider) Center() rl.Vector3 {
	switch p.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		return rl.NewVector3(p.Position.X, p.Position.Y+p.Height*0.5, p.Position.Z+p.Width*0.5)
	case PlaneAxisYPos, PlaneAxisYNeg:
		return rl.NewVector3(p.Position.X+p.Width*0.5, p.Position.Y, p.Position.Z+p.Height*0.5)
	default:
		return rl.NewVector3(p.Position.X+p.Width*0.5, p.Position.Y+p.Height*0.5, p.Position.Z)
	}
}

// GetPosition returns plane position anchor.
func (p *PlaneCollider) GetPosition() rl.Vector3 {
	return p.Position
}

// SetPosition sets plane position anchor.
func (p *PlaneCollider) SetPosition(pos rl.Vector3) {
	p.Position = pos
}

func (p *PlaneCollider) DistanceTo(other Collider) float32 {
	switch c := other.(type) {
	case *CylinderCollider:
		return cylinderVsPlaneDistance(*c, *p)
	case *BoxCollider:
		return boxVsPlaneDistance(*c, *p)
	default:
		return infiniteDistance()
	}
}

var _ Collider = (*PlaneCollider)(nil)
