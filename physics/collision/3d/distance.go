package col3d

import "github.com/chewxy/math32"

func boxVsBoxDistance(a, b *BoxCollider) float32 {
	aMax := a.Max()
	bMax := b.Max()

	gapX := intervalGap(a.Position.X, aMax.X, b.Position.X, bMax.X)
	gapY := intervalGap(a.Position.Y, aMax.Y, b.Position.Y, bMax.Y)
	gapZ := intervalGap(a.Position.Z, aMax.Z, b.Position.Z, bMax.Z)

	return math32.Sqrt(gapX*gapX + gapY*gapY + gapZ*gapZ)
}

func boxVsCylinderDistance(box *BoxCollider, cyl *CylinderCollider) float32 {
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
	verticalGap := intervalGap(box.Position.Y, boxMax.Y, cyl.Position.Y, cyl.Position.Y+cyl.Height)

	return math32.Sqrt(horizontalGap*horizontalGap + verticalGap*verticalGap)
}

func boxVsPointDistance(box *BoxCollider, pt *PointCollider) float32 {
	boxMax := box.Max()
	return pointAABBDistance3D(
		pt.Position.X,
		pt.Position.Y,
		pt.Position.Z,
		box.Position.X,
		boxMax.X,
		box.Position.Y,
		boxMax.Y,
		box.Position.Z,
		boxMax.Z,
	)
}

func boxVsPlaneDistance(box *BoxCollider, plane *PlaneCollider) float32 {
	boxMax := box.Max()
	return aabbDistanceToPlaneRect(
		box.Position.X,
		boxMax.X,
		box.Position.Y,
		boxMax.Y,
		box.Position.Z,
		boxMax.Z,
		plane,
	)
}


func cylinderVsCylinderDistance(a, b *CylinderCollider) float32 {
	dx := a.Position.X - b.Position.X
	dz := a.Position.Z - b.Position.Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)
	horizontalGap := math32.Max(0, distXZ-(a.Radius+b.Radius))
	verticalGap := intervalGap(a.Position.Y, a.Position.Y+a.Height, b.Position.Y, b.Position.Y+b.Height)

	return math32.Sqrt(horizontalGap*horizontalGap + verticalGap*verticalGap)
}

func cylinderVsBoxDistance(cyl *CylinderCollider, box *BoxCollider) float32 {
	return boxVsCylinderDistance(box, cyl)
}

func cylinderVsPointDistance(cyl *CylinderCollider, pt *PointCollider) float32 {
	return pointVsCylinderDistance(pt, cyl)
}

func cylinderVsPlaneDistance(cyl *CylinderCollider, plane *PlaneCollider) float32 {
	return aabbDistanceToPlaneRect(
		cyl.Position.X-cyl.Radius,
		cyl.Position.X+cyl.Radius,
		cyl.Position.Y,
		cyl.Position.Y+cyl.Height,
		cyl.Position.Z-cyl.Radius,
		cyl.Position.Z+cyl.Radius,
		plane,
	)
}


func pointVsBoxDistance(pt *PointCollider, box *BoxCollider) float32 {
	boxMax := box.Max()
	return pointAABBDistance3D(
		pt.Position.X,
		pt.Position.Y,
		pt.Position.Z,
		box.Position.X,
		boxMax.X,
		box.Position.Y,
		boxMax.Y,
		box.Position.Z,
		boxMax.Z,
	)
}

func pointVsCylinderDistance(pt *PointCollider, cyl *CylinderCollider) float32 {
	return pointCylinderDistance(
		pt.Position.X,
		pt.Position.Y,
		pt.Position.Z,
		cyl.Position.X,
		cyl.Position.Y,
		cyl.Position.Z,
		cyl.Radius,
		cyl.Height,
	)
}
