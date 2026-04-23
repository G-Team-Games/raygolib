package col3d

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestPlaneDrawBox_AppliesThicknessOnFlatAxis(t *testing.T) {
	tests := []struct {
		name      string
		axis      PlaneAxis
		wantSizeX float32
		wantSizeY float32
		wantSizeZ float32
	}{
		{
			name:      "x axis plane has thickness on x",
			axis:      PlaneAxisXPos,
			wantSizeX: defaultPlaneThickness,
			wantSizeY: 4,
			wantSizeZ: 2,
		},
		{
			name:      "y axis plane has thickness on y",
			axis:      PlaneAxisYPos,
			wantSizeX: 2,
			wantSizeY: defaultPlaneThickness,
			wantSizeZ: 4,
		},
		{
			name:      "z axis plane has thickness on z",
			axis:      PlaneAxisZPos,
			wantSizeX: 2,
			wantSizeY: 4,
			wantSizeZ: defaultPlaneThickness,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plane := NewPlaneCollider(rl.NewVector3(1, 2, 3), 2, 4, tt.axis)
			_, size := planeDrawBox(plane)

			if size.X != tt.wantSizeX {
				t.Fatalf("unexpected size.X: got %v, want %v", size.X, tt.wantSizeX)
			}
			if size.Y != tt.wantSizeY {
				t.Fatalf("unexpected size.Y: got %v, want %v", size.Y, tt.wantSizeY)
			}
			if size.Z != tt.wantSizeZ {
				t.Fatalf("unexpected size.Z: got %v, want %v", size.Z, tt.wantSizeZ)
			}
		})
	}
}

func TestPlaneDrawBox_CenterUsesExpandedSize(t *testing.T) {
	plane := NewPlaneCollider(rl.NewVector3(10, 20, 30), 2, 6, PlaneAxisXPos)
	center, size := planeDrawBox(plane)

	wantCenter := rl.NewVector3(10+size.X*0.5, 20+size.Y*0.5, 30+size.Z*0.5)
	if center != wantCenter {
		t.Fatalf("unexpected center: got %v, want %v", center, wantCenter)
	}
}
