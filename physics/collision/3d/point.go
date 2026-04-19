package col3d

import rl "github.com/gen2brain/raylib-go/raylib"

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
	case *PlaneCollider:
		contact := pointVsPlaneContact(p, o)
		contact.Normal = rl.Vector3Negate(contact.Normal)
		return contact
	case *PointCollider:
		return pointVsPointContact(p, o)
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
	case *PlaneCollider:
		return pointVsPlaneDistance(p, o)
	case *PointCollider:
		return pointVsPointDistance(p, o)
	default:
		return infiniteDistance()
	}
}

func (p *PointCollider) GetPosition() rl.Vector3 {
	return p.Position
}

func (p *PointCollider) SetPosition(vec rl.Vector3) {
	p.Position = vec
}

var _ Collider = (*PointCollider)(nil)

