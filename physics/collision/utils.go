package collision

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Normalizes vector and returns zero for tiny lengths
func safeNormalize2(v rl.Vector2) rl.Vector2 {
	if rl.Vector2LengthSqr(v) <= 1e-8 {
		return rl.NewVector2(0, 0)
	}
	return rl.Vector2Normalize(v)
}

func overlap1D(aMin, aMax, bMin, bMax float32) float32 {
	return math32.Min(aMax, bMax) - math32.Max(aMin, bMin)
}

func sign(x float32) float32 {
	if x >= 0 {
		return 1
	}

	return -1
}

// planeQuad builds quad corners for finite axis-aligned plane.
func planeQuad(plane PlaneCollider) (rl.Vector3, rl.Vector3, rl.Vector3, rl.Vector3) {
	switch plane.Axis {
	case PlaneAxisXPos, PlaneAxisXNeg:
		return rl.NewVector3(plane.Position.X, plane.Position.Y, plane.Position.Z),
			rl.NewVector3(plane.Position.X, plane.Position.Y+plane.Height, plane.Position.Z),
			rl.NewVector3(plane.Position.X, plane.Position.Y+plane.Height, plane.Position.Z+plane.Width),
			rl.NewVector3(plane.Position.X, plane.Position.Y, plane.Position.Z+plane.Width)
	case PlaneAxisYPos, PlaneAxisYNeg:
		return rl.NewVector3(plane.Position.X, plane.Position.Y, plane.Position.Z),
			rl.NewVector3(plane.Position.X+plane.Width, plane.Position.Y, plane.Position.Z),
			rl.NewVector3(plane.Position.X+plane.Width, plane.Position.Y, plane.Position.Z+plane.Height),
			rl.NewVector3(plane.Position.X, plane.Position.Y, plane.Position.Z+plane.Height)
	default:
		return rl.NewVector3(plane.Position.X, plane.Position.Y, plane.Position.Z),
			rl.NewVector3(plane.Position.X+plane.Width, plane.Position.Y, plane.Position.Z),
			rl.NewVector3(plane.Position.X+plane.Width, plane.Position.Y+plane.Height, plane.Position.Z),
			rl.NewVector3(plane.Position.X, plane.Position.Y+plane.Height, plane.Position.Z)
	}
}
