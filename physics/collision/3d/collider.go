package col3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ShapeKind uint8

const (
	// Axis-aligned box collider
	ShapeBox ShapeKind = iota
	// Upright cylinder aligned to Y axis
	ShapeCylinderY
	// Finite axis-aligned plane rectangle, see PlaneAxis type
	ShapePlaneRect
	// Point in 3D space
	ShapePoint
)

// Pure collider interface
type Collider interface {
	Kind() ShapeKind
	Collide(Collider) Contact
	BoundingBox() rl.BoundingBox
}

// Helper function to dispatch collision check from a to b
func Collide(a, b Collider) Contact {
	return a.Collide(b)
}

// Applies minimum translation vector to object position.
func ResolveByMTV(getPosition func() rl.Vector3, setPosition func(rl.Vector3), hit Contact) {
	if !hit.Hit || hit.Penetration <= 0 {
		return
	}
	setPosition(rl.Vector3Add(getPosition(), rl.Vector3Scale(hit.Normal, hit.Penetration)))
}
