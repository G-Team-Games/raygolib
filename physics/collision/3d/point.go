package col3d

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PointCollider struct {
	Position rl.Vector3
}

func NewPointV(position rl.Vector3) *PointCollider {
	return &PointCollider{Position: position}
}

func NewPointXYZ(x, y, z float32) *PointCollider {
	return &PointCollider{Position: rl.NewVector3(x, y, z)}
}

func (p *PointCollider) BoundingBox() rl.BoundingBox {
	return rl.NewBoundingBox(p.Position, p.Position)
}

func (p *PointCollider) Collide(other Collider) Contact {
	switch o := other.(type) {
	case *BoxCollider:
		return boxVsPointContact(o, p)
	case *CylinderCollider:
		return cylinderVsPointContact(o, p)
	default:
		return Contact{}
	}
}

func (p *PointCollider) Kind() ShapeKind {
	return ShapePoint
}

func (p *PointCollider) DistanceTo(other Collider) float32 {
	switch o := other.(type) {
	case *BoxCollider:
		return pointVsBoxDistance(p, o)
	case *CylinderCollider:
		return pointVsCylinderDistance(p, o)
	default:
		return unsupportedDistance()
	}
}

func pointVsBoxDistance(pt *PointCollider, box *BoxCollider) float32 {
	closestX := math32.Max(box.Position.X, math32.Min(pt.Position.X, box.Position.X+box.Size.X))
	closestY := math32.Max(box.Position.Y, math32.Min(pt.Position.Y, box.Position.Y+box.Size.Y))
	closestZ := math32.Max(box.Position.Z, math32.Min(pt.Position.Z, box.Position.Z+box.Size.Z))

	dx := pt.Position.X - closestX
	dy := pt.Position.Y - closestY
	dz := pt.Position.Z - closestZ

	return math32.Sqrt(dx*dx + dy*dy + dz*dz)
}

func pointVsCylinderDistance(pt *PointCollider, cyl *CylinderCollider) float32 {
	dx := pt.Position.X - cyl.Position.X
	dz := pt.Position.Z - cyl.Position.Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)
	horizontalGap := math32.Max(0, distXZ-cyl.Radius)
	verticalGap := intervalGap(pt.Position.Y, pt.Position.Y, cyl.Position.Y, cyl.Position.Y+cyl.Height)

	return combineOrthogonalGaps(horizontalGap, verticalGap)
}

func (p *PointCollider) GetPosition() rl.Vector3 {
	return p.Position
}

func (p *PointCollider) SetPosition(vec rl.Vector3) {
	p.Position = vec
}

var _ Collider = (*PointCollider)(nil)
