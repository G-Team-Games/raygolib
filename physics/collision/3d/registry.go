package col3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollisionHandler func(a, b Collider) Contact
type DistanceHandler func(a, b Collider) float32

var collisionRegistry = make(map[ShapeKind]map[ShapeKind]CollisionHandler)
var distanceRegistry = make(map[ShapeKind]map[ShapeKind]DistanceHandler)

func RegisterCollision(k1, k2 ShapeKind, h CollisionHandler) {
	if collisionRegistry[k1] == nil {
		collisionRegistry[k1] = make(map[ShapeKind]CollisionHandler)
	}
	collisionRegistry[k1][k2] = h

	if k1 != k2 {
		if collisionRegistry[k2] == nil {
			collisionRegistry[k2] = make(map[ShapeKind]CollisionHandler)
		}
		collisionRegistry[k2][k1] = func(a, b Collider) Contact {
			contact := h(b, a)
			if contact.Hit {
				contact.Normal = rl.Vector3Negate(contact.Normal)
			}
			return contact
		}
	}
}

func RegisterDistance(k1, k2 ShapeKind, h DistanceHandler) {
	if distanceRegistry[k1] == nil {
		distanceRegistry[k1] = make(map[ShapeKind]DistanceHandler)
	}
	distanceRegistry[k1][k2] = h

	if k1 != k2 {
		if distanceRegistry[k2] == nil {
			distanceRegistry[k2] = make(map[ShapeKind]DistanceHandler)
		}
		distanceRegistry[k2][k1] = func(a, b Collider) float32 {
			return h(b, a)
		}
	}
}

func init() {
	RegisterCollision(ShapeBox, ShapeBox, func(a, b Collider) Contact {
		return boxVsBoxContact(a.(*BoxCollider), b.(*BoxCollider))
	})
	RegisterCollision(ShapeBox, ShapeCylinderY, func(a, b Collider) Contact {
		contact := cylinderVsBoxContact(b.(*CylinderCollider), a.(*BoxCollider))
		if contact.Hit {
			contact.Normal = rl.Vector3Negate(contact.Normal)
		}
		return contact
	})
	RegisterCollision(ShapeBox, ShapePoint, func(a, b Collider) Contact {
		contact := boxVsPointContact(a.(*BoxCollider), b.(*PointCollider))
		if contact.Hit {
			contact.Normal = rl.Vector3Negate(contact.Normal)
		}
		return contact
	})
	RegisterCollision(ShapeBox, ShapePlaneRect, func(a, b Collider) Contact {
		return boxVsPlaneContact(a.(*BoxCollider), b.(*PlaneCollider))
	})
	RegisterCollision(ShapeCylinderY, ShapeCylinderY, func(a, b Collider) Contact {
		return cylinderVsCylinderContact(a.(*CylinderCollider), b.(*CylinderCollider))
	})
	RegisterCollision(ShapeCylinderY, ShapePoint, func(a, b Collider) Contact {
		contact := cylinderVsPointContact(a.(*CylinderCollider), b.(*PointCollider))
		if contact.Hit {
			contact.Normal = rl.Vector3Negate(contact.Normal)
		}
		return contact
	})
	RegisterCollision(ShapeCylinderY, ShapePlaneRect, func(a, b Collider) Contact {
		return cylinderVsPlaneContact(a.(*CylinderCollider), b.(*PlaneCollider))
	})
	RegisterCollision(ShapePoint, ShapePlaneRect, func(a, b Collider) Contact {
		return pointVsPlaneContact(a.(*PointCollider), b.(*PlaneCollider))
	})
	RegisterCollision(ShapePoint, ShapePoint, func(a, b Collider) Contact {
		return pointVsPointContact(a.(*PointCollider), b.(*PointCollider))
	})

	RegisterDistance(ShapeBox, ShapeBox, func(a, b Collider) float32 {
		return boxVsBoxDistance(a.(*BoxCollider), b.(*BoxCollider))
	})
	RegisterDistance(ShapeBox, ShapeCylinderY, func(a, b Collider) float32 {
		return boxVsCylinderDistance(a.(*BoxCollider), b.(*CylinderCollider))
	})
	RegisterDistance(ShapeBox, ShapePoint, func(a, b Collider) float32 {
		return boxVsPointDistance(a.(*BoxCollider), b.(*PointCollider))
	})
	RegisterDistance(ShapeBox, ShapePlaneRect, func(a, b Collider) float32 {
		return boxVsPlaneDistance(a.(*BoxCollider), b.(*PlaneCollider))
	})
	RegisterDistance(ShapeCylinderY, ShapeCylinderY, func(a, b Collider) float32 {
		return cylinderVsCylinderDistance(a.(*CylinderCollider), b.(*CylinderCollider))
	})
	RegisterDistance(ShapeCylinderY, ShapePoint, func(a, b Collider) float32 {
		return cylinderVsPointDistance(a.(*CylinderCollider), b.(*PointCollider))
	})
	RegisterDistance(ShapeCylinderY, ShapePlaneRect, func(a, b Collider) float32 {
		return cylinderVsPlaneDistance(a.(*CylinderCollider), b.(*PlaneCollider))
	})
	RegisterDistance(ShapePoint, ShapePlaneRect, func(a, b Collider) float32 {
		return pointVsPlaneDistance(a.(*PointCollider), b.(*PlaneCollider))
	})
	RegisterDistance(ShapePoint, ShapePoint, func(a, b Collider) float32 {
		return pointVsPointDistance(a.(*PointCollider), b.(*PointCollider))
	})
}
