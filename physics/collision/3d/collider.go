// Package col3d provides simple 3D collision and raycast primitives.
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

// Collider defines shape operations used by the narrow phase.
//
// Package-wide contracts:
//   - For a.Collide(b), returned contact normal points from b toward a
//     (direction to move a out of b).
//   - For a hit contact, penetration is always >= 0.
//   - Touching counts as hit with penetration == 0.
//   - DistanceTo returns the minimum Euclidean separation between volumes,
//     and returns 0 when touching or overlapping.
//   - Unsupported DistanceTo pairs return +Inf (never a silent zero).
type Collider interface {
	Kind() ShapeKind
	Collide(Collider) Contact
	DistanceTo(Collider) float32
	BoundingBox() rl.BoundingBox
}

// Collide dispatches collision using the central registry 
// and returns the contact between two colliders
func Collide(a, b Collider) Contact {
	if handler, ok := collisionRegistry[a.Kind()][b.Kind()]; ok {
		return handler(a, b)
	}
	return Contact{}
}

// Distance returns distance between two colliders using the central registry.
func Distance(a, b Collider) float32 {
	if handler, ok := distanceRegistry[a.Kind()][b.Kind()]; ok {
		return handler(a, b)
	}
	return infiniteDistance()
}

// ResolveByMTV applies the contact minimum translation vector to a position.
//
// It assumes hit.Normal and hit.Penetration follow the contact contract:
// normal points in the direction that moves the resolved object out of overlap,
// and penetration is non-negative.
func ResolveByMTV(getPosition func() rl.Vector3, setPosition func(rl.Vector3), hit Contact) {
	if !hit.Hit || hit.Penetration <= 0 {
		return
	}
	setPosition(rl.Vector3Add(getPosition(), rl.Vector3Scale(hit.Normal, hit.Penetration)))
}
