package col3d

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestResolveMultiMTV(t *testing.T) {
	// Active box fully inside two other boxes - guaranteed to have penetration
	// Active at (0,0,0), size (1,1,1) -> occupies [0,1] on each axis
	// b1 at (-0.5,-0.5,-0.5), size (2,2,2) -> occupies [-0.5,1.5] - contains active
	// b2 at (0.5,0.5,0.5), size (2,2,2) -> occupies [0.5,2.5] - also contains active
	active := NewBoxColliderV(rl.NewVector3(0, 0, 0), rl.NewVector3(1, 1, 1))
	
	b1 := NewBoxColliderV(rl.NewVector3(-0.5, -0.5, -0.5), rl.NewVector3(2, 2, 2))
	b2 := NewBoxColliderV(rl.NewVector3(0.5, 0.5, 0.5), rl.NewVector3(2, 2, 2))
	
	others := []Collider{b1, b2}
	resolved := ResolveMultiMTV(active, others, 10)

	if !resolved {
		t.Fatal("expected ResolveMultiMTV to fully resolve with enough iterations")
	}
}
