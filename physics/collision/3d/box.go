package col3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BoxCollider struct {
	Position rl.Vector3
	Size     rl.Vector3
}

// Creates box collider from min corner and size as rl.Vector3
func NewBoxColliderV(position rl.Vector3, size rl.Vector3) *BoxCollider {
	return &BoxCollider{Position: position, Size: size}
}

// Creates box collider from scalar sizes
func NewBoxColliderXYZ(position rl.Vector3, sizeX, sizeY, sizeZ float32) *BoxCollider {
	return &BoxCollider{Position: position, Size: rl.NewVector3(sizeX, sizeY, sizeZ)}
}

func (b *BoxCollider) Kind() ShapeKind {
	return ShapeBox
}

// Returns contact between box and another collider
func (b *BoxCollider) Collide(other Collider) Contact {
	switch c := other.(type) {
	case *BoxCollider:
		return boxVsBoxContact(*b, *c)
	case *CylinderCollider:
		contact := cylinderVsBoxContact(*c, *b)
		contact.Normal = rl.Vector3Negate(contact.Normal)
		return contact
	default:
		return Contact{}
	}
}

// TODO: Add custom collide method that accepts the collision solver

// Returns raylib AABB representation of box
func (b *BoxCollider) BoundingBox() rl.BoundingBox {
	return rl.NewBoundingBox(b.Position, rl.Vector3Add(b.Position, b.Size))
}

// Geometric center of box
func (b *BoxCollider) Center() rl.Vector3 {
	return rl.Vector3Add(b.Position, rl.Vector3Scale(b.Size, 0.5))
}

// Minimum corner of box
func (b *BoxCollider) Min() rl.Vector3 {
	return b.Position
}

// Maximum corner of box
func (b *BoxCollider) Max() rl.Vector3 {
	return rl.Vector3Add(b.Position, b.Size)
}

// Current box position
func (b *BoxCollider) GetPosition() rl.Vector3 {
	return b.Position
}

// Sets box position
func (b *BoxCollider) SetPosition(vec rl.Vector3) {
	b.Position = vec
}

var _ Collider = (*BoxCollider)(nil)
