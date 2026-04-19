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
	cylinder := NewCylinderCollider(rl.NewVector3(0.5, 0, 0.5), 0.5, 2)
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
	c := NewCylinderCollider(rl.NewVector3(2, 4, 6), 1, 8)
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
	a := NewCylinderCollider(rl.NewVector3(0, 0, 0), 1, 2)
	b := NewCylinderCollider(rl.NewVector3(1, 0, 0), 1, 2)

	if dist := a.DistanceTo(b); dist != 0 {
		t.Fatalf("expected zero distance for overlapping cylinders, got %f", dist)
	}
}

func TestDistance_PointVsCylinder_OutsideCircleNonZero(t *testing.T) {
	cyl := NewCylinderCollider(rl.NewVector3(0, 0, 0), 1, 2)
	pt := NewPointV(rl.NewVector3(2, 1, 0))

	if dist := pt.DistanceTo(cyl); dist <= 0 {
		t.Fatalf("expected positive distance for outside point, got %f", dist)
	}
}

func TestDistance_SymmetryForSupportedPairs(t *testing.T) {
	const eps = 1e-5

	boxA := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2))
	boxB := NewBoxColliderV(rl.NewVector3(4, 1, 0), rl.NewVector3(1, 2, 1))
	cylA := NewCylinderCollider(rl.NewVector3(1, 0, 3), 0.75, 2)
	cylB := NewCylinderCollider(rl.NewVector3(4, 1, 3), 0.5, 1.5)
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

func TestDistance_UnsupportedPairReturnsInf(t *testing.T) {
	plane := NewPlaneCollider(rl.NewVector3(0, 0, 0), 2, 2, PlaneAxisXPos)
	point := NewPointV(rl.NewVector3(0, 0, 0))

	d := plane.DistanceTo(point)
	if !math.IsInf(float64(d), 1) {
		t.Fatalf("expected +Inf for unsupported distance pair, got %f", d)
	}

	d = point.DistanceTo(plane)
	if !math.IsInf(float64(d), 1) {
		t.Fatalf("expected +Inf for unsupported distance pair, got %f", d)
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
			b:    NewCylinderCollider(rl.NewVector3(1.5, 0, 1), 0.75, 2),
		},
		{
			name: "cylinder-plane",
			a:    NewCylinderCollider(rl.NewVector3(1, 0, 1), 0.5, 2),
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
	cylinder := NewCylinderCollider(rl.NewVector3(1, 0, 1), 0.5, 2)

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
			b:    NewCylinderCollider(rl.NewVector3(1.4, 0, 1), 0.75, 2),
		},
		{
			name: "box-point",
			a:    NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(2, 2, 2)),
			b:    NewPointV(rl.NewVector3(1, 1, 1)),
		},
		{
			name: "cylinder-cylinder",
			a:    NewCylinderCollider(rl.NewVector3(0, 0, 0), 1, 2),
			b:    NewCylinderCollider(rl.NewVector3(1.5, 0, 0), 1, 2),
		},
		{
			name: "cylinder-point",
			a:    NewCylinderCollider(rl.NewVector3(0, 0, 0), 1, 2),
			b:    NewPointV(rl.NewVector3(0.5, 1, 0)),
		},
		{
			name: "cylinder-plane",
			a:    NewCylinderCollider(rl.NewVector3(1, 0, 1), 0.5, 2),
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
