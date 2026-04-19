package col3d

import (
	"github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayHit struct {
	Hit    bool
	Normal rl.Vector3
	Point  rl.Vector3
}

func Raycast(ray rl.Ray, collider Collider) RayHit {
	switch c := collider.(type) {
	case *BoxCollider:
		box := c.BoundingBox()
		hit := rl.GetRayCollisionBox(ray, box)
		if !hit.Hit {
			return RayHit{}
		}
		return RayHit{Hit: true, Normal: hit.Normal, Point: hit.Point}
	case *PlaneCollider:
		p1, p2, p3, p4 := planeQuad(*c)
		hit := rl.GetRayCollisionQuad(ray, p1, p2, p3, p4)
		if !hit.Hit {
			return RayHit{}
		}
		return RayHit{Hit: true, Normal: c.Axis.Normal(), Point: hit.Point}
	case *CylinderCollider:
		return raycastCylinder(ray, *c)
	default:
		return RayHit{}
	}
}

func raycastCylinder(ray rl.Ray, cylinder CylinderCollider) RayHit {
	origin := ray.Position
	direction := ray.Direction

	cx, cz := cylinder.Position.X, cylinder.Position.Z
	baseY := cylinder.Position.Y
	topY := cylinder.Position.Y + cylinder.Height
	radius := cylinder.Radius
	radiusSq := radius * radius

	tBest := float32(-1)
	best := RayHit{}

	tryCandidate := func(t float32, normal rl.Vector3) {
		if t < 0 {
			return
		}
		if tBest >= 0 && t >= tBest {
			return
		}

		point := rl.NewVector3(
			origin.X+direction.X*t,
			origin.Y+direction.Y*t,
			origin.Z+direction.Z*t,
		)

		tBest = t
		best = RayHit{Hit: true, Point: point, Normal: normal}
	}

	// Side intersection with infinite cylinder in XZ, then clamp to finite Y range.
	ox := origin.X - cx
	oz := origin.Z - cz
	a := direction.X*direction.X + direction.Z*direction.Z
	b := 2 * (ox*direction.X + oz*direction.Z)
	c := ox*ox + oz*oz - radiusSq

	if a > epsilon {
		discriminant := b*b - 4*a*c
		if discriminant >= 0 {
			sqrtDisc := math32.Sqrt(discriminant)
			inv2A := float32(1) / (2 * a)
			t0 := (-b - sqrtDisc) * inv2A
			t1 := (-b + sqrtDisc) * inv2A

			if t1 < t0 {
				t0, t1 = t1, t0
			}

			for _, t := range []float32{t0, t1} {
				if t < 0 {
					continue
				}
				y := origin.Y + direction.Y*t
				if y < baseY-epsilon || y > topY+epsilon {
					continue
				}

				hitX := origin.X + direction.X*t
				hitZ := origin.Z + direction.Z*t
				nx := hitX - cx
				nz := hitZ - cz
				len := math32.Sqrt(nx*nx + nz*nz)
				normal := rl.NewVector3(1, 0, 0)
				if len > epsilon {
					normal = rl.NewVector3(nx/len, 0, nz/len)
				}
				tryCandidate(t, normal)
			}
		}
	}

	// Cap intersections at y=baseY and y=topY.
	if math32.Abs(direction.Y) > epsilon {
		for _, cap := range []struct {
			y      float32
			normal rl.Vector3
		}{
			{y: baseY, normal: rl.NewVector3(0, -1, 0)},
			{y: topY, normal: rl.NewVector3(0, 1, 0)},
		} {
			t := (cap.y - origin.Y) / direction.Y
			if t < 0 {
				continue
			}

			hitX := origin.X + direction.X*t
			hitZ := origin.Z + direction.Z*t
			dx := hitX - cx
			dz := hitZ - cz
			if dx*dx+dz*dz > radiusSq+epsilon {
				continue
			}

			tryCandidate(t, cap.normal)
		}
	}

	return best
}
