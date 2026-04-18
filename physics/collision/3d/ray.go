package col3d

import (
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
		return RayHit{Hit: true, Normal: hit.Normal, Point: hit.Point}
	case *CylinderCollider:
		return raycastCylinder(ray, *c)
	default:
		return RayHit{}
	}
}

func raycastCylinder(ray rl.Ray, cylinder CylinderCollider) RayHit {
	hit := rl.GetRayCollisionSphere(ray, cylinder.Position, cylinder.Radius)
	if hit.Hit {
		if hit.Point.Y >= cylinder.Position.Y && hit.Point.Y <= cylinder.Position.Y+cylinder.Height {
			return RayHit{Hit: true, Normal: hit.Normal, Point: hit.Point}
		}
	}

	bottomCenter := cylinder.Position
	topCenter := rl.NewVector3(cylinder.Position.X, cylinder.Position.Y+cylinder.Height, cylinder.Position.Z)

	hitBottom := rl.GetRayCollisionQuad(
		ray,
		rl.NewVector3(bottomCenter.X-cylinder.Radius, bottomCenter.Y, bottomCenter.Z-cylinder.Radius),
		rl.NewVector3(bottomCenter.X+cylinder.Radius, bottomCenter.Y, bottomCenter.Z-cylinder.Radius),
		rl.NewVector3(bottomCenter.X+cylinder.Radius, bottomCenter.Y, bottomCenter.Z+cylinder.Radius),
		rl.NewVector3(bottomCenter.X-cylinder.Radius, bottomCenter.Y, bottomCenter.Z+cylinder.Radius),
	)
	if hitBottom.Hit {
		delta := rl.Vector2Subtract(rl.NewVector2(hitBottom.Point.X, hitBottom.Point.Z), rl.NewVector2(bottomCenter.X, bottomCenter.Z))
		if rl.Vector2Length(delta) <= cylinder.Radius {
			return RayHit{Hit: true, Normal: hitBottom.Normal, Point: hitBottom.Point}
		}
	}

	hitTop := rl.GetRayCollisionQuad(
		ray,
		rl.NewVector3(topCenter.X-cylinder.Radius, topCenter.Y, topCenter.Z-cylinder.Radius),
		rl.NewVector3(topCenter.X+cylinder.Radius, topCenter.Y, topCenter.Z-cylinder.Radius),
		rl.NewVector3(topCenter.X+cylinder.Radius, topCenter.Y, topCenter.Z+cylinder.Radius),
		rl.NewVector3(topCenter.X-cylinder.Radius, topCenter.Y, topCenter.Z+cylinder.Radius),
	)
	if hitTop.Hit {
		delta := rl.Vector2Subtract(rl.NewVector2(hitTop.Point.X, hitTop.Point.Z), rl.NewVector2(topCenter.X, topCenter.Z))
		if rl.Vector2Length(delta) <= cylinder.Radius {
			return RayHit{Hit: true, Normal: hitTop.Normal, Point: hitTop.Point}
		}
	}

	return RayHit{}
}
