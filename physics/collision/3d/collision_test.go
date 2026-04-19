package col3d

import (
	"math"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestBoxVsBoxContact(t *testing.T) {
	t.Run("hit with penetration", func(t *testing.T) {
		a := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
		b := NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(2, 2, 2))

		hit := Collide(a, b)
		if !hit.Hit {
			t.Fatalf("expected box collision hit")
		}
		if hit.Penetration <= 0 {
			t.Fatalf("expected penetration > 0, got %f", hit.Penetration)
		}
	})

	t.Run("miss", func(t *testing.T) {
		a := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1))
		b := NewBoxColliderV(rl.NewVector3(3, 0, 0), rl.NewVector3(1, 1, 1))

		hit := Collide(a, b)
		if hit.Hit {
			t.Fatalf("expected no hit")
		}
	})
}

func TestBoxHelpers(t *testing.T) {
	box := NewBoxColliderV(rl.NewVector3(1, 2, 3), rl.NewVector3(4, 6, 8))

	center := box.Center()
	if center.X != 3 || center.Y != 5 || center.Z != 7 {
		t.Fatalf("unexpected box center: %+v", center)
	}

	min := box.Min()
	if min.X != 1 || min.Y != 2 || min.Z != 3 {
		t.Fatalf("unexpected box min: %+v", min)
	}

	max := box.Max()
	if max.X != 5 || max.Y != 8 || max.Z != 11 {
		t.Fatalf("unexpected box max: %+v", max)
	}
}

func TestRaycastBox(t *testing.T) {
	box := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
	ray := rl.NewRay(rl.NewVector3(-5, 1, 1), rl.NewVector3(1, 0, 0))

	hit := Raycast(ray, box)
	if !hit.Hit {
		t.Fatalf("expected ray hit box")
	}
	if hit.Point == (rl.Vector3{}) {
		t.Fatalf("expected hit point")
	}
}

func TestCylinderVsBoxContactAndResolve(t *testing.T) {
	cylinder := NewCylinderColliderV(rl.NewVector3(0.5, 0, 0.5), 0.5, 2)
	box := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 2, 1))

	hit := Collide(cylinder, box)
	if !hit.Hit {
		t.Fatalf("expected contact hit")
	}
	if hit.Penetration <= 0 {
		t.Fatalf("expected penetration > 0")
	}

	prev := cylinder.GetPosition()
	ResolveByMTV(cylinder.GetPosition, cylinder.SetPosition, hit)
	next := cylinder.GetPosition()
	if prev == next {
		t.Fatalf("expected position change after resolve")
	}
}

func TestPlaneAxisNormal(t *testing.T) {
	n := PlaneAxisZNeg.Normal()
	if n.X != 0 || n.Y != 0 || n.Z != -1 {
		t.Fatalf("unexpected normal: %+v", n)
	}
}

func TestCylinderCenter(t *testing.T) {
	c := NewCylinderColliderV(rl.NewVector3(2, 4, 6), 1, 8)
	center := c.Center()
	if center.X != 2 || center.Y != 8 || center.Z != 6 {
		t.Fatalf("unexpected cylinder center: %+v", center)
	}
}

func TestPlaneCenter(t *testing.T) {
	pX := NewPlaneCollider(rl.NewVector3(1, 2, 3), 6, 4, PlaneAxisXPos)
	cX := pX.Center()
	if cX.X != 1 || cX.Y != 4 || cX.Z != 6 {
		t.Fatalf("unexpected X plane center: %+v", cX)
	}

	pY := NewPlaneCollider(rl.NewVector3(1, 2, 3), 6, 4, PlaneAxisYPos)
	cY := pY.Center()
	if cY.X != 4 || cY.Y != 2 || cY.Z != 5 {
		t.Fatalf("unexpected Y plane center: %+v", cY)
	}

	pZ := NewPlaneCollider(rl.NewVector3(1, 2, 3), 6, 4, PlaneAxisZPos)
	cZ := pZ.Center()
	if cZ.X != 4 || cZ.Y != 4 || cZ.Z != 3 {
		t.Fatalf("unexpected Z plane center: %+v", cZ)
	}
}

func TestContract_CollideNormalPointsFromBToA(t *testing.T) {
	a := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
	b := NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(2, 2, 2))

	hit := a.Collide(b)
	if !hit.Hit {
		t.Fatalf("expected hit")
	}
	if hit.Penetration < 0 {
		t.Fatalf("expected non-negative penetration, got %f", hit.Penetration)
	}
	if hit.Normal.X != -1 || hit.Normal.Y != 0 || hit.Normal.Z != 0 {
		t.Fatalf("expected normal from b to a (-X), got %+v", hit.Normal)
	}
}

func TestContract_TouchingIsHitWithZeroPenetration(t *testing.T) {
	a := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1))
	b := NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(1, 1, 1))

	hit := a.Collide(b)
	if !hit.Hit {
		t.Fatalf("expected touching boxes to be a hit")
	}
	if hit.Penetration != 0 {
		t.Fatalf("expected zero penetration for touching boxes, got %f", hit.Penetration)
	}
}

func TestDistance_BoxVsBox_OverlapIsZero(t *testing.T) {
	a := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
	b := NewBoxColliderV(rl.NewVector3(1, 1, 1), rl.NewVector3(2, 2, 2))

	if dist := a.DistanceTo(b); dist != 0 {
		t.Fatalf("expected zero distance for overlapping boxes, got %f", dist)
	}
}

func TestDistance_CylinderVsCylinder_OverlapIsZero(t *testing.T) {
	a := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2)
	b := NewCylinderColliderV(rl.NewVector3(1, 0, 0), 1, 2)

	if dist := a.DistanceTo(b); dist != 0 {
		t.Fatalf("expected zero distance for overlapping cylinders, got %f", dist)
	}
}

func TestDistance_PointVsCylinder_OutsideCircleNonZero(t *testing.T) {
	cyl := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2)
	pt := NewPointV(rl.NewVector3(2, 1, 0))

	if dist := pt.DistanceTo(cyl); dist <= 0 {
		t.Fatalf("expected positive distance for outside point, got %f", dist)
	}
}

func TestDistance_SymmetryForSupportedPairs(t *testing.T) {
	const eps = 1e-5

	boxA := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
	boxB := NewBoxColliderV(rl.NewVector3(4, 1, 0), rl.NewVector3(1, 2, 1))
	cylA := NewCylinderColliderV(rl.NewVector3(1, 0, 3), 0.75, 2)
	cylB := NewCylinderColliderV(rl.NewVector3(4, 1, 3), 0.5, 1.5)
	point := NewPointV(rl.NewVector3(3, 3, 3))
	planeY := NewPlaneCollider(rl.NewVector3(0, 0, 0), 6, 6, PlaneAxisYPos)

	tests := []struct {
		name string
		a    Collider
		b    Collider
	}{
		{name: "box-box", a: boxA, b: boxB},
		{name: "box-cylinder", a: boxA, b: cylA},
		{name: "box-point", a: boxA, b: point},
		{name: "cylinder-cylinder", a: cylA, b: cylB},
		{name: "cylinder-point", a: cylA, b: point},
		{name: "cylinder-plane", a: cylA, b: planeY},
		{name: "box-plane", a: boxA, b: planeY},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d1 := tc.a.DistanceTo(tc.b)
			d2 := tc.b.DistanceTo(tc.a)
			if math.Abs(float64(d1-d2)) > eps {
				t.Fatalf("expected symmetric distance, got %f and %f", d1, d2)
			}
		})
	}
}

type dummyCollider struct{}

func (d *dummyCollider) Kind() ShapeKind { return 255 }
func (d *dummyCollider) Collide(other Collider) Contact { return Contact{} }
func (d *dummyCollider) DistanceTo(other Collider) float32 { return float32(math.Inf(1)) }
func (d *dummyCollider) BoundingBox() rl.BoundingBox { return rl.BoundingBox{} }

func TestDistance_UnsupportedPairReturnsInf(t *testing.T) {
	plane := NewPlaneCollider(rl.NewVector3(0, 0, 0), 2, 2, PlaneAxisXPos)
	dummy := &dummyCollider{}

	d := plane.DistanceTo(dummy)
	if !math.IsInf(float64(d), 1) {
		t.Errorf("expected +Inf for unsupported distance pair, got %f", d)
	}
}

func TestResolveByMTV_ReducesOverlap_BoxBoxInOneStep(t *testing.T) {
	a := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
	b := NewBoxColliderV(rl.NewVector3(1.5, 0, 0), rl.NewVector3(2, 2, 2))

	hit := Collide(a, b)
	if !hit.Hit {
		t.Fatalf("expected initial overlap")
	}
	if hit.Penetration <= 0 {
		t.Fatalf("expected penetration > 0, got %f", hit.Penetration)
	}

	ResolveByMTV(a.GetPosition, a.SetPosition, hit)
	after := Collide(a, b)

	if !after.Hit {
		t.Fatalf("expected touching state after one-step resolve")
	}
	if after.Penetration != 0 {
		t.Fatalf("expected zero penetration after resolve, got %f", after.Penetration)
	}
}

func TestCollide_ReverseOrderNormalsAreOpposite_WhenBothSupported(t *testing.T) {
	tests := []struct {
		name string
		a    Collider
		b    Collider
	}{
		{
			name: "box-box",
			a:    NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
			b:    NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(2, 2, 2)),
		},
		{
			name: "box-cylinder",
			a:    NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
			b:    NewCylinderColliderV(rl.NewVector3(1.5, 0, 1), 0.75, 2),
		},
		{
			name: "cylinder-plane",
			a:    NewCylinderColliderV(rl.NewVector3(1, 0, 1), 0.5, 2),
			b:    NewPlaneCollider(rl.NewVector3(0, 1.8, 0), 3, 3, PlaneAxisYPos),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ab := Collide(tc.a, tc.b)
			ba := Collide(tc.b, tc.a)

			if !ab.Hit || !ba.Hit {
				t.Fatalf("expected hits in both call orders, got ab=%v ba=%v", ab.Hit, ba.Hit)
			}

			if math.Abs(float64(ab.Penetration-ba.Penetration)) > 1e-5 {
				t.Fatalf("expected equal penetration, got %f vs %f", ab.Penetration, ba.Penetration)
			}

			if math.Abs(float64(ab.Normal.X+ba.Normal.X)) > 1e-5 ||
				math.Abs(float64(ab.Normal.Y+ba.Normal.Y)) > 1e-5 ||
				math.Abs(float64(ab.Normal.Z+ba.Normal.Z)) > 1e-5 {
				t.Fatalf("expected opposite normals, got ab=%+v ba=%+v", ab.Normal, ba.Normal)
			}
		})
	}
}

func TestCylinderVsHorizontalPlane_MinimalTranslationDirection(t *testing.T) {
	plane := NewPlaneCollider(rl.NewVector3(0, 0.2, 0), 4, 4, PlaneAxisYPos)
	cylinder := NewCylinderColliderV(rl.NewVector3(1, 0, 1), 0.5, 2)

	hit := Collide(cylinder, plane)
	if !hit.Hit {
		t.Fatalf("expected hit")
	}

	if hit.Normal.X != 0 || hit.Normal.Y != 1 || hit.Normal.Z != 0 {
		t.Fatalf("expected +Y normal for nearest escape, got %+v", hit.Normal)
	}
	if math.Abs(float64(hit.Penetration-0.2)) > 1e-5 {
		t.Fatalf("expected +Y penetration of 0.2, got %f", hit.Penetration)
	}
}

func TestCollide_OrderIndependence_AllSupportedUnorderedPairs(t *testing.T) {
	tests := []struct {
		name string
		a    Collider
		b    Collider
	}{
		{
			name: "box-box",
			a:    NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
			b:    NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(2, 2, 2)),
		},
		{
			name: "box-cylinder",
			a:    NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
			b:    NewCylinderColliderV(rl.NewVector3(1.4, 0, 1), 0.75, 2),
		},
		{
			name: "box-point",
			a:    NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
			b:    NewPointV(rl.NewVector3(1, 1, 1)),
		},
		{
			name: "cylinder-cylinder",
			a:    NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
			b:    NewCylinderColliderV(rl.NewVector3(1.5, 0, 0), 1, 2),
		},
		{
			name: "cylinder-point",
			a:    NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
			b:    NewPointV(rl.NewVector3(0.5, 1, 0)),
		},
		{
			name: "cylinder-plane",
			a:    NewCylinderColliderV(rl.NewVector3(1, 0, 1), 0.5, 2),
			b:    NewPlaneCollider(rl.NewVector3(0, 1.8, 0), 3, 3, PlaneAxisYPos),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ab := Collide(tc.a, tc.b)
			ba := Collide(tc.b, tc.a)

			if !ab.Hit || !ba.Hit {
				t.Fatalf("expected hits in both orders, got ab=%v ba=%v", ab.Hit, ba.Hit)
			}

			if math.Abs(float64(ab.Penetration-ba.Penetration)) > 1e-5 {
				t.Fatalf("expected equal penetration, got %f vs %f", ab.Penetration, ba.Penetration)
			}

			if math.Abs(float64(ab.Normal.X+ba.Normal.X)) > 1e-5 ||
				math.Abs(float64(ab.Normal.Y+ba.Normal.Y)) > 1e-5 ||
				math.Abs(float64(ab.Normal.Z+ba.Normal.Z)) > 1e-5 {
				t.Fatalf("expected opposite normals, got ab=%+v ba=%+v", ab.Normal, ba.Normal)
			}
		})
	}
}

func TestRaycastCylinder_SideHit(t *testing.T) {
	cyl := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 4)
	ray := rl.NewRay(rl.NewVector3(-3, 2, 0), rl.NewVector3(1, 0, 0))

	hit := Raycast(ray, cyl)
	if !hit.Hit {
		t.Fatalf("expected side hit")
	}

	if !vecApprox(hit.Point, rl.NewVector3(-1, 2, 0), 1e-4) {
		t.Fatalf("unexpected hit point: %+v", hit.Point)
	}
	if !vecApprox(hit.Normal, rl.NewVector3(-1, 0, 0), 1e-4) {
		t.Fatalf("unexpected side normal: %+v", hit.Normal)
	}
}

func TestRaycastCylinder_TopHit(t *testing.T) {
	cyl := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 4)
	ray := rl.NewRay(rl.NewVector3(0, 6, 0), rl.NewVector3(0, -1, 0))

	hit := Raycast(ray, cyl)
	if !hit.Hit {
		t.Fatalf("expected top-cap hit")
	}

	if !vecApprox(hit.Point, rl.NewVector3(0, 4, 0), 1e-4) {
		t.Fatalf("unexpected hit point: %+v", hit.Point)
	}
	if !vecApprox(hit.Normal, rl.NewVector3(0, 1, 0), 1e-4) {
		t.Fatalf("unexpected top normal: %+v", hit.Normal)
	}
}

func TestRaycastCylinder_BottomHit(t *testing.T) {
	cyl := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 4)
	ray := rl.NewRay(rl.NewVector3(0, -2, 0), rl.NewVector3(0, 1, 0))

	hit := Raycast(ray, cyl)
	if !hit.Hit {
		t.Fatalf("expected bottom-cap hit")
	}

	if !vecApprox(hit.Point, rl.NewVector3(0, 0, 0), 1e-4) {
		t.Fatalf("unexpected hit point: %+v", hit.Point)
	}
	if !vecApprox(hit.Normal, rl.NewVector3(0, -1, 0), 1e-4) {
		t.Fatalf("unexpected bottom normal: %+v", hit.Normal)
	}
}

func TestRaycastCylinder_Miss(t *testing.T) {
	cyl := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 4)
	ray := rl.NewRay(rl.NewVector3(3, 2, 0), rl.NewVector3(0, 1, 0))

	hit := Raycast(ray, cyl)
	if hit.Hit {
		t.Fatalf("expected miss, got %+v", hit)
	}
}

func TestRaycastCylinder_TallCylinderTopHit(t *testing.T) {
	cyl := NewCylinderColliderV(rl.NewVector3(0, 0, 0), 0.5, 10)
	ray := rl.NewRay(rl.NewVector3(0, 20, 0), rl.NewVector3(0, -1, 0))

	hit := Raycast(ray, cyl)
	if !hit.Hit {
		t.Fatalf("expected top hit for tall cylinder")
	}

	if !vecApprox(hit.Point, rl.NewVector3(0, 10, 0), 1e-4) {
		t.Fatalf("unexpected hit point: %+v", hit.Point)
	}
	if !vecApprox(hit.Normal, rl.NewVector3(0, 1, 0), 1e-4) {
		t.Fatalf("unexpected top normal: %+v", hit.Normal)
	}
}

func TestRaycastPlane_NormalMatchesAxis(t *testing.T) {
	tests := []struct {
		name   string
		plane  *PlaneCollider
		ray    rl.Ray
		normal rl.Vector3
	}{
		{
			name:   "x+",
			plane:  NewPlaneCollider(rl.NewVector3(0, 0, 0), 4, 4, PlaneAxisXPos),
			ray:    rl.NewRay(rl.NewVector3(-2, 2, 2), rl.NewVector3(1, 0, 0)),
			normal: rl.NewVector3(1, 0, 0),
		},
		{
			name:   "x-",
			plane:  NewPlaneCollider(rl.NewVector3(0, 0, 0), 4, 4, PlaneAxisXNeg),
			ray:    rl.NewRay(rl.NewVector3(-2, 2, 2), rl.NewVector3(1, 0, 0)),
			normal: rl.NewVector3(-1, 0, 0),
		},
		{
			name:   "y+",
			plane:  NewPlaneCollider(rl.NewVector3(0, 0, 0), 4, 4, PlaneAxisYPos),
			ray:    rl.NewRay(rl.NewVector3(2, -2, 2), rl.NewVector3(0, 1, 0)),
			normal: rl.NewVector3(0, 1, 0),
		},
		{
			name:   "y-",
			plane:  NewPlaneCollider(rl.NewVector3(0, 0, 0), 4, 4, PlaneAxisYNeg),
			ray:    rl.NewRay(rl.NewVector3(2, -2, 2), rl.NewVector3(0, 1, 0)),
			normal: rl.NewVector3(0, -1, 0),
		},
		{
			name:   "z+",
			plane:  NewPlaneCollider(rl.NewVector3(0, 0, 0), 4, 4, PlaneAxisZPos),
			ray:    rl.NewRay(rl.NewVector3(2, 2, -2), rl.NewVector3(0, 0, 1)),
			normal: rl.NewVector3(0, 0, 1),
		},
		{
			name:   "z-",
			plane:  NewPlaneCollider(rl.NewVector3(0, 0, 0), 4, 4, PlaneAxisZNeg),
			ray:    rl.NewRay(rl.NewVector3(2, 2, -2), rl.NewVector3(0, 0, 1)),
			normal: rl.NewVector3(0, 0, -1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hit := Raycast(tc.ray, tc.plane)
			if !hit.Hit {
				t.Fatalf("expected ray hit")
			}
			if !vecApprox(hit.Normal, tc.normal, 1e-5) {
				t.Fatalf("expected normal %+v, got %+v", tc.normal, hit.Normal)
			}
		})
	}
}

type pairScenario struct {
	name      string
	expectHit bool
	build     func() (Collider, Collider)
}

func TestCollisionMatrix_BoundaryAndReverseOrder(t *testing.T) {
	scenarios := []pairScenario{
		{
			name:      "box-box/face-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-box/edge-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(1, 1, 0), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-box/corner-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(1, 1, 1), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-box/fully-inside",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
					NewBoxColliderV(rl.NewVector3(0.5, 0.5, 0.5), rl.NewVector3(0.5, 0.5, 0.5))
			},
		},
		{
			name:      "box-box/separated-axis",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(2.1, 0, 0), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-box/separated-diagonal",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(2, 2, 2), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-cylinder/overlap",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
					NewCylinderColliderV(rl.NewVector3(1, 0, 1), 0.5, 2)
			},
		},
		{
			name:      "box-cylinder/separated-axis",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewCylinderColliderV(rl.NewVector3(2, 0, 0.5), 0.4, 1)
			},
		},
		{
			name:      "box-point/face-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewPointV(rl.NewVector3(1, 0.5, 0.5))
			},
		},
		{
			name:      "box-point/separated-diagonal",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewPointV(rl.NewVector3(2, 2, 2))
			},
		},
		{
			name:      "cylinder-cylinder/face-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 0.5, 1),
					NewCylinderColliderV(rl.NewVector3(1, 0, 0), 0.5, 1)
			},
		},
		{
			name:      "cylinder-cylinder/separated-diagonal",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 0.5, 1),
					NewCylinderColliderV(rl.NewVector3(2, 2, 0), 0.5, 1)
			},
		},
		{
			name:      "cylinder-point/side-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
					NewPointV(rl.NewVector3(1, 1, 0))
			},
		},
		{
			name:      "cylinder-point/separated-axis",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
					NewPointV(rl.NewVector3(0, 3, 0))
			},
		},
		{
			name:      "cylinder-plane/face-touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(2, 0, 2), 0.5, 1),
					NewPlaneCollider(rl.NewVector3(0, 1, 0), 4, 4, PlaneAxisYPos)
			},
		},
		{
			name:      "cylinder-plane/separated-axis",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(2, 2.2, 2), 0.5, 1),
					NewPlaneCollider(rl.NewVector3(0, 1, 0), 4, 4, PlaneAxisYPos)
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			a, b := scenario.build()
			ab := Collide(a, b)
			ba := Collide(b, a)

			if ab.Hit != scenario.expectHit {
				t.Fatalf("expected ab hit=%v, got %v", scenario.expectHit, ab.Hit)
			}
			if ba.Hit != scenario.expectHit {
				t.Fatalf("expected ba hit=%v, got %v", scenario.expectHit, ba.Hit)
			}

			if !scenario.expectHit {
				return
			}

			if ab.Penetration < 0 || ba.Penetration < 0 {
				t.Fatalf("penetration must be non-negative, got ab=%f ba=%f", ab.Penetration, ba.Penetration)
			}

			if math.Abs(float64(ab.Penetration-ba.Penetration)) > 1e-5 {
				t.Fatalf("expected equal penetration, got %f and %f", ab.Penetration, ba.Penetration)
			}

			if math.Abs(float64(ab.Normal.X+ba.Normal.X)) > 1e-5 ||
				math.Abs(float64(ab.Normal.Y+ba.Normal.Y)) > 1e-5 ||
				math.Abs(float64(ab.Normal.Z+ba.Normal.Z)) > 1e-5 {
				t.Fatalf("expected opposite normals, got ab=%+v ba=%+v", ab.Normal, ba.Normal)
			}
		})
	}
}

func TestDistanceInvariants_Matrix(t *testing.T) {
	const eps = 1e-5

	scenarios := []pairScenario{
		{
			name:      "box-box/overlap",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
					NewBoxColliderV(rl.NewVector3(1, 1, 1), rl.NewVector3(2, 2, 2))
			},
		},
		{
			name:      "box-box/touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(1, 0, 0), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-box/separated",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewBoxColliderV(rl.NewVector3(3, 0, 0), rl.NewVector3(1, 1, 1))
			},
		},
		{
			name:      "box-cylinder/overlap",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
					NewCylinderColliderV(rl.NewVector3(1, 0, 1), 0.75, 2)
			},
		},
		{
			name:      "box-cylinder/separated",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewCylinderColliderV(rl.NewVector3(3, 0, 0.5), 0.4, 1)
			},
		},
		{
			name:      "box-point/touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewPointV(rl.NewVector3(1, 0.2, 0.2))
			},
		},
		{
			name:      "box-point/separated",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1)),
					NewPointV(rl.NewVector3(3, 3, 3))
			},
		},
		{
			name:      "cylinder-cylinder/overlap",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
					NewCylinderColliderV(rl.NewVector3(1, 0, 0), 1, 2)
			},
		},
		{
			name:      "cylinder-cylinder/separated",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 0.5, 1),
					NewCylinderColliderV(rl.NewVector3(2, 2, 0), 0.5, 1)
			},
		},
		{
			name:      "cylinder-point/touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
					NewPointV(rl.NewVector3(1, 1, 0))
			},
		},
		{
			name:      "cylinder-point/separated",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(0, 0, 0), 1, 2),
					NewPointV(rl.NewVector3(0, 4, 0))
			},
		},
		{
			name:      "cylinder-plane/touch",
			expectHit: true,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(2, 0, 2), 0.5, 1),
					NewPlaneCollider(rl.NewVector3(0, 1, 0), 4, 4, PlaneAxisYPos)
			},
		},
		{
			name:      "cylinder-plane/separated",
			expectHit: false,
			build: func() (Collider, Collider) {
				return NewCylinderColliderV(rl.NewVector3(2, 3, 2), 0.5, 1),
					NewPlaneCollider(rl.NewVector3(0, 1, 0), 4, 4, PlaneAxisYPos)
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			a, b := scenario.build()
			d1 := a.DistanceTo(b)
			d2 := b.DistanceTo(a)

			if math.IsNaN(float64(d1)) || math.IsInf(float64(d1), 1) {
				t.Fatalf("distance a->b must be finite, got %f", d1)
			}
			if math.IsNaN(float64(d2)) || math.IsInf(float64(d2), 1) {
				t.Fatalf("distance b->a must be finite, got %f", d2)
			}

			if d1 < 0 || d2 < 0 {
				t.Fatalf("distance must be non-negative, got %f and %f", d1, d2)
			}

			if math.Abs(float64(d1-d2)) > eps {
				t.Fatalf("distance must be symmetric, got %f and %f", d1, d2)
			}

			if scenario.expectHit {
				if math.Abs(float64(d1)) > eps {
					t.Fatalf("expected zero distance for overlap/touch, got %f", d1)
				}
				return
			}

			if d1 <= 0 {
				t.Fatalf("expected positive distance for separated case, got %f", d1)
			}
		})
	}
}

func vecApprox(a, b rl.Vector3, eps float64) bool {
	return math.Abs(float64(a.X-b.X)) <= eps &&
		math.Abs(float64(a.Y-b.Y)) <= eps &&
		math.Abs(float64(a.Z-b.Z)) <= eps
}
