package col3d

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Contact struct {
	// Hit reports whether colliders overlap or touch.
	Hit bool
	// Normal points from collider b toward collider a for a.Collide(b).
	// It is the direction used by ResolveByMTV to move a out of b.
	Normal rl.Vector3
	// Penetration is overlap depth along Normal and is >= 0 when Hit is true.
	// For touching cases, penetration is 0.
	Penetration float32
}

func boxVsBoxContact(a, b *BoxCollider) Contact {
	// Penetration on each axis
	px := overlap1D(a.Position.X, a.Position.X+a.Size.X, b.Position.X, b.Position.X+b.Size.X)
	py := overlap1D(a.Position.Y, a.Position.Y+a.Size.Y, b.Position.Y, b.Position.Y+b.Size.Y)
	pz := overlap1D(a.Position.Z, a.Position.Z+a.Size.Z, b.Position.Z, b.Position.Z+b.Size.Z)

	// If no penetration, no collision
	if px < 0 || py < 0 || pz < 0 {
		return Contact{}
	}

	centerA := rl.Vector3Add(a.Position, rl.Vector3Scale(a.Size, 0.5))
	centerB := rl.Vector3Add(b.Position, rl.Vector3Scale(b.Size, 0.5))

	// Min penetration axis for normal
	penetration := px
	normal := rl.NewVector3(sign(centerA.X-centerB.X), 0, 0)

	if py < penetration {
		penetration = py
		normal = rl.NewVector3(0, sign(centerA.Y-centerB.Y), 0)
	}
	if pz < penetration {
		penetration = pz
		normal = rl.NewVector3(0, 0, sign(centerA.Z-centerB.Z))
	}

	return Contact{
		Hit:         true,
		Normal:      normal,
		Penetration: penetration,
	}
}

func cylinderVsBoxContact(cyl *CylinderCollider, box *BoxCollider) Contact {
	boxMin := box.Position
	boxMax := box.Max()
	pushUp := boxMax.Y - cyl.Position.Y
	pushDown := (cyl.Position.Y + cyl.Height) - boxMin.Y

	penY := pushUp
	normalY := rl.NewVector3(0, 1, 0)
	if pushDown < penY {
		penY = pushDown
		normalY = rl.NewVector3(0, -1, 0)
	}

	closestX := math32.Max(boxMin.X, math32.Min(cyl.Position.X, boxMax.X))
	closestZ := math32.Max(boxMin.Z, math32.Min(cyl.Position.Z, boxMax.Z))

	dx := cyl.Position.X - closestX
	dz := cyl.Position.Z - closestZ
	distXZ := math32.Sqrt(dx*dx + dz*dz)

	penSide := cyl.Radius - distXZ
	normalSide := rl.Vector3{}

	if distXZ > epsilon {
		normalSide = rl.NewVector3(dx/distXZ, 0, dz/distXZ)
	} else {
		// Center is inside/on box XZ projection, choose nearest face.
		left := cyl.Position.X - boxMin.X
		right := boxMax.X - cyl.Position.X
		back := cyl.Position.Z - boxMin.Z
		front := boxMax.Z - cyl.Position.Z

		penSide = cyl.Radius + left
		normalSide = rl.NewVector3(-1, 0, 0)

		if p := cyl.Radius + right; p < penSide {
			penSide = p
			normalSide = rl.NewVector3(1, 0, 0)
		}
		if p := cyl.Radius + back; p < penSide {
			penSide = p
			normalSide = rl.NewVector3(0, 0, -1)
		}
		if p := cyl.Radius + front; p < penSide {
			penSide = p
			normalSide = rl.NewVector3(0, 0, 1)
		}
	}

	if penSide < 0 || penY < 0 {
		return Contact{}
	}

	if penY < penSide {
		return Contact{Hit: true, Normal: normalY, Penetration: penY}
	}

	return Contact{Hit: true, Normal: normalSide, Penetration: penSide}
}

// cylinderVsCylinderContact computes contact between two upright cylinders.
func cylinderVsCylinderContact(a, b *CylinderCollider) Contact {
	dx := a.Position.X - b.Position.X
	dz := a.Position.Z - b.Position.Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)

	penSide := (a.Radius + b.Radius) - distXZ
	normalSide := rl.NewVector3(1, 0, 0)
	if distXZ > epsilon {
		normalSide = rl.NewVector3(dx/distXZ, 0, dz/distXZ)
	}

	pushUp := b.Position.Y + b.Height - a.Position.Y
	pushDown := a.Position.Y + a.Height - b.Position.Y

	penY := pushUp
	normalY := rl.NewVector3(0, 1, 0)
	if pushDown < penY {
		penY = pushDown
		normalY = rl.NewVector3(0, -1, 0)
	}

	if penSide < 0 || penY < 0 {
		return Contact{}
	}

	if penY < penSide {
		return Contact{Hit: true, Normal: normalY, Penetration: penY}
	}

	return Contact{Hit: true, Normal: normalSide, Penetration: penSide}
}

// boxVsPlaneContact computes contact between box and finite plane.
func boxVsPlaneContact(box *BoxCollider, plane *PlaneCollider) Contact {
	bMin := box.Min()
	bMax := box.Max()
	pBox := plane.BoundingBox()

	px := overlap1D(bMin.X, bMax.X, pBox.Min.X, pBox.Max.X)
	py := overlap1D(bMin.Y, bMax.Y, pBox.Min.Y, pBox.Max.Y)
	pz := overlap1D(bMin.Z, bMax.Z, pBox.Min.Z, pBox.Max.Z)

	if px < 0 || py < 0 || pz < 0 {
		return Contact{}
	}

	var normal rl.Vector3
	var pen float32

	switch plane.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		pen = px
		normal = rl.NewVector3(1, 0, 0)
		if box.Center().X < pBox.Min.X {
			normal.X = -1
		}
	case PlaneAxisYPos, PlaneAxisYNeg:
		pen = py
		normal = rl.NewVector3(0, 1, 0)
		if box.Center().Y < pBox.Min.Y {
			normal.Y = -1
		}
	case PlaneAxisZPos, PlaneAxisZNeg:
		pen = pz
		normal = rl.NewVector3(0, 0, 1)
		if box.Center().Z < pBox.Min.Z {
			normal.Z = -1
		}
	}

	return Contact{Hit: true, Normal: normal, Penetration: pen}
}

func boxVsPointContact(box *BoxCollider, pt *PointCollider) Contact {
	bMinX, bMinY, bMinZ := vec3ToValues(box.Position)
	bMaxX, bMaxY, bMaxZ := vec3ToValues(box.Max())
	ptX, ptY, ptZ := vec3ToValues(pt.Position)

	inside := (ptX >= bMinX && ptX <= bMaxX) && (ptY >= bMinY && ptY <= bMaxY) && (ptZ >= bMinZ && ptZ <= bMaxZ)
	if !inside {
		return Contact{}
	}

	// penetration and coresponding normals
	deltasAndNormals := []struct {
		delta  float32
		normal rl.Vector3
	}{
		{ptX - bMinX, rl.NewVector3(-1, 0, 0)}, // left
		{bMaxX - ptX, rl.NewVector3(1, 0, 0)},  // right
		{ptY - bMinY, rl.NewVector3(0, -1, 0)}, // bottom
		{bMaxY - ptY, rl.NewVector3(0, 1, 0)},  // top
		{ptZ - bMinZ, rl.NewVector3(0, 0, -1)}, // back
		{bMaxZ - ptZ, rl.NewVector3(0, 0, 1)},  // front
	}

	// Finding min penetration
	minPen := deltasAndNormals[0]
	for i := 1; i < len(deltasAndNormals); i++ {
		if deltasAndNormals[i].delta < minPen.delta {
			minPen = deltasAndNormals[i]
		}
	}

	return Contact{
		Hit:         true,
		Normal:      minPen.normal,
		Penetration: minPen.delta,
	}
}

func cylinderVsPointContact(cylinder *CylinderCollider, pt *PointCollider) Contact {
	dx := pt.Position.X - cylinder.Position.X
	dz := pt.Position.Z - cylinder.Position.Z
	distXZ := math32.Sqrt(dx*dx + dz*dz)
	penetrationXZ := cylinder.Radius - distXZ

	dBottom := pt.Position.Y - cylinder.Position.Y
	dTop := (cylinder.Position.Y + cylinder.Height) - pt.Position.Y
	penetrationY := math32.Min(dBottom, dTop)

	if penetrationXZ < 0 || penetrationY < 0 {
		return Contact{}
	}

	sideNormal := rl.NewVector3(1, 0, 0)
	if distXZ > epsilon {
		sideNormal = rl.NewVector3(dx/distXZ, 0, dz/distXZ)
	}

	if penetrationXZ < penetrationY {
		return Contact{Hit: true, Normal: sideNormal, Penetration: penetrationXZ}
	}

	if dBottom < dTop {
		return Contact{Hit: true, Normal: rl.NewVector3(0, -1, 0), Penetration: dBottom}
	}
	return Contact{Hit: true, Normal: rl.NewVector3(0, 1, 0), Penetration: dTop}
}

// cylinderVsPlaneContact computes contact between cylinder and finite plane.
func cylinderVsPlaneContact(cylinder *CylinderCollider, plane *PlaneCollider) Contact {
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
		planeY := plane.Position.Y
		cylMinY := cylinder.Position.Y
		cylMaxY := cylinder.Position.Y + cylinder.Height

		if penetrationXZ < 0 || planeY < cylMinY || planeY > cylMaxY {
			return Contact{}
		}

		moveUp := planeY - cylMinY
		moveDown := cylMaxY - planeY

		// Deterministic fallback for exact mid-plane ties: prefer +Y.
		if moveUp <= moveDown {
			return Contact{Hit: true, Normal: rl.NewVector3(0, 1, 0), Penetration: moveUp}
		}

		return Contact{Hit: true, Normal: rl.NewVector3(0, -1, 0), Penetration: moveDown}
	}

	return Contact{}
}

func pointVsPlaneContact(pt *PointCollider, plane *PlaneCollider) Contact {
	pBox := plane.BoundingBox()
	
	px := overlap1D(pt.Position.X, pt.Position.X, pBox.Min.X, pBox.Max.X)
	py := overlap1D(pt.Position.Y, pt.Position.Y, pBox.Min.Y, pBox.Max.Y)
	pz := overlap1D(pt.Position.Z, pt.Position.Z, pBox.Min.Z, pBox.Max.Z)

	if px < 0 || py < 0 || pz < 0 {
		return Contact{}
	}

	var normal rl.Vector3
	var pen float32

	switch plane.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		pen = px
		normal = rl.NewVector3(1, 0, 0)
		if pt.Position.X < pBox.Min.X {
			normal.X = -1
		}
	case PlaneAxisYPos, PlaneAxisYNeg:
		pen = py
		normal = rl.NewVector3(0, 1, 0)
		if pt.Position.Y < pBox.Min.Y {
			normal.Y = -1
		}
	case PlaneAxisZPos, PlaneAxisZNeg:
		pen = pz
		normal = rl.NewVector3(0, 0, 1)
		if pt.Position.Z < pBox.Min.Z {
			normal.Z = -1
		}
	}

	return Contact{Hit: true, Normal: normal, Penetration: pen}
}

func pointVsPointContact(a, b *PointCollider) Contact {
	if a.Position.X == b.Position.X && a.Position.Y == b.Position.Y && a.Position.Z == b.Position.Z {
		return Contact{Hit: true, Normal: rl.NewVector3(0, 1, 0), Penetration: 0}
	}
	return Contact{}
}

