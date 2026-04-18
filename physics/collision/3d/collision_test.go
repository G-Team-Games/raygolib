package col3d

import (
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
