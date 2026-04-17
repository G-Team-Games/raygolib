package col3d

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Collision query result between two shapes or ray and shape
type Contact struct {
	Hit         bool
	Normal      rl.Vector3
	Distance    float32
	Penetration float32
}


func boxVsBoxContact(a, b BoxCollider) Contact {
	// Penetration on each axis
	px := overlap1D(a.Position.X, a.Position.X+a.Size.X, b.Position.X, b.Position.X+b.Size.X)
	py := overlap1D(a.Position.Y, a.Position.Y+a.Size.Y, b.Position.Y, b.Position.Y+b.Size.Y)
	pz := overlap1D(a.Position.Z, a.Position.Z+a.Size.Z, b.Position.Z, b.Position.Z+b.Size.Z)

	// If no penetration return distance
	if px <= 0 || py <= 0 || pz <= 0 {
		return Contact{
			Distance: math32.Max(px, math32.Max(py, pz)),
		}
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

func cylinderVsBoxContact(cylinder CylinderCollider, box BoxCollider) Contact {
	const epsilon float32 = 1e-7

	closestX := math32.Min(box.Position.X+box.Size.X, math32.Max(box.Position.X, cylinder.Position.X))
	closestZ := math32.Min(box.Position.Z+box.Size.Z, math32.Max(box.Position.Z, cylinder.Position.Z))

	dx := cylinder.Position.X - closestX
	dz := cylinder.Position.Z - closestZ

	distXZ := math32.Sqrt(dx*dx + dz*dz)
	penetrationXZ := cylinder.Radius - distXZ

	penetrationY := overlap1D(
		cylinder.Position.Y,
		cylinder.Position.Y+cylinder.Height,
		box.Position.Y,
		box.Position.Y+box.Size.Y,
	)

	if penetrationXZ <= epsilon || penetrationY <= epsilon {
		return Contact{
			Hit:      false,
			Distance: math32.Max(penetrationXZ, penetrationY),
		}
	}

	var normalXZ rl.Vector3
	if distXZ > epsilon {
		normalXZ = rl.NewVector3(dx/distXZ, 0, dz/distXZ)
	} else {

		centerBoxX := box.Position.X + box.Size.X*0.5
		centerBoxZ := box.Position.Z + box.Size.Z*0.5

		fallback := rl.NewVector2(
			cylinder.Position.X-centerBoxX,
			cylinder.Position.Z-centerBoxZ,
		)

		n := safeNormalize2(fallback)
		if rl.Vector2LengthSqr(n) <= epsilon {
			n = rl.NewVector2(1, 0)
		}

		normalXZ = rl.NewVector3(n.X, 0, n.Y)
	}

	var normal rl.Vector3
	var penetration float32

	if penetrationXZ < penetrationY {
		// Side collision
		normal = normalXZ
		penetration = penetrationXZ
	} else {
		// Top/bottom collision
		centerC := cylinder.Position.Y + cylinder.Height*0.5
		centerB := box.Position.Y + box.Size.Y*0.5
		normal = rl.NewVector3(0, sign(centerC-centerB), 0)
		penetration = penetrationY
	}

	return Contact{
		Hit:         true,
		Normal:      normal,
		Penetration: penetration,
	}
}

// cylinderVsCylinderContact computes contact between two upright cylinders.
func cylinderVsCylinderContact(a, b CylinderCollider) Contact {
	difference := rl.Vector2Subtract(rl.NewVector2(a.Position.X, a.Position.Z), rl.NewVector2(b.Position.X, b.Position.Z))
	distanceXZ := rl.Vector2Length(difference)
	penetrationXZ := (a.Radius + b.Radius) - distanceXZ

	distanceY1 := a.Position.Y - (b.Position.Y + b.Height)
	distanceY2 := b.Position.Y - (a.Position.Y + a.Height)

	if penetrationXZ < 0 || distanceY1 > 0 || distanceY2 > 0 {
		return Contact{}
	}

	if penetrationXZ > -distanceY1 && penetrationXZ > -distanceY2 {
		normalXZ := safeNormalize2(difference)
		if rl.Vector2Length(normalXZ) == 0 {
			normalXZ = rl.NewVector2(1, 0)
		}
		normal := rl.NewVector3(normalXZ.X, 0, normalXZ.Y)
		return Contact{Hit: true, Normal: normal, Penetration: penetrationXZ}
	}

	if -distanceY1 < -distanceY2 {
		return Contact{
			Hit:         true,
			Normal:      rl.NewVector3(0, 1, 0),
			Penetration: -distanceY1,
		}
	}

	return Contact{
		Hit:         true,
		Normal:      rl.NewVector3(0, -1, 0),
		Penetration: -distanceY2,
	}
}
