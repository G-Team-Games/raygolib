package col3d

import (
	"github.com/chewxy/math32"
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
	switch o := other.(type) {
	case *BoxCollider:
		return boxVsBoxContact(b, o)
	case *CylinderCollider:
		contact := cylinderVsBoxContact(o, b)
		contact.Normal = rl.Vector3Negate(contact.Normal)
		return contact
	case *PointCollider:
		return boxVsPointContact(b, o)
	default:
		return Contact{}
	}
}

// TODO: Add custom collide method that accepts the collision solver

func (b *BoxCollider) DistanceTo(other Collider) float32 {
	switch c := other.(type) {
	case *BoxCollider:
		return boxVsBoxDistance(*b, *c)
	case *CylinderCollider:
		return boxVsCylinderDistance(*b, *c)
	case *PointCollider:
		return boxVsPointDistance(*b, *c)
	case *PlaneCollider:
		return boxVsPlaneDistance(*b, *c)
	default:
		return 0
	}
}

func boxVsBoxDistance(a, b BoxCollider) float32 {
	closestX := math32.Max(a.Position.X, math32.Min(b.Position.X, a.Position.X+a.Size.X))
	closestY := math32.Max(a.Position.Y, math32.Min(b.Position.Y, a.Position.Y+a.Size.Y))
	closestZ := math32.Max(a.Position.Z, math32.Min(b.Position.Z, a.Position.Z+a.Size.Z))

	dx := b.Center().X - closestX
	dy := b.Center().Y - closestY
	dz := b.Center().Z - closestZ

	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

func boxVsCylinderDistance(box BoxCollider, cyl CylinderCollider) float32 {
	closestX := math32.Max(cyl.Position.X-cyl.Radius, math32.Min(cyl.Position.X+cyl.Radius, box.Position.X))
	closestY := math32.Max(cyl.Position.Y, math32.Min(cyl.Position.Y+cyl.Height, box.Position.Y))
	closestZ := math32.Max(cyl.Position.Z-cyl.Radius, math32.Min(cyl.Position.Z+cyl.Radius, box.Position.Z))

	dx := box.Center().X - closestX
	dy := box.Center().Y - closestY
	dz := box.Center().Z - closestZ

	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

func boxVsPointDistance(box BoxCollider, pt PointCollider) float32 {
	closestX := math32.Max(box.Position.X, math32.Min(pt.Position.X, box.Position.X+box.Size.X))
	closestY := math32.Max(box.Position.Y, math32.Min(pt.Position.Y, box.Position.Y+box.Size.Y))
	closestZ := math32.Max(box.Position.Z, math32.Min(pt.Position.Z, box.Position.Z+box.Size.Z))

	dx := pt.Position.X - closestX
	dy := pt.Position.Y - closestY
	dz := pt.Position.Z - closestZ

	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

func boxVsPlaneDistance(box BoxCollider, plane PlaneCollider) float32 {
	closest := closestPointOnPlane(plane, box.Center())
	dx := box.Center().X - closest.X
	dy := box.Center().Y - closest.Y
	dz := box.Center().Z - closest.Z
	distSq := dx*dx + dy*dy + dz*dz
	return math32.Sqrt(distSq)
}

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
