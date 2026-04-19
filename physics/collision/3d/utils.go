package col3d

import (
	"math"

	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const epsilon float32 = 1e-7

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

func intervalGap(aMin, aMax, bMin, bMax float32) float32 {
	if aMax < bMin {
		return bMin - aMax
	}
	if bMax < aMin {
		return aMin - bMax
	}
	return 0
}

func pointRectDistanceXZ(px, pz, minX, maxX, minZ, maxZ float32) float32 {
	dx := intervalGap(px, px, minX, maxX)
	dz := intervalGap(pz, pz, minZ, maxZ)
	return combineOrthogonalGaps(dx, dz)
}

func circleRectGapXZ(cx, cz, r, minX, maxX, minZ, maxZ float32) float32 {
	pointToRect := pointRectDistanceXZ(cx, cz, minX, maxX, minZ, maxZ)
	return math32.Max(0, pointToRect-r)
}

func combineOrthogonalGaps(g1, g2 float32) float32 {
	return math32.Sqrt(g1*g1 + g2*g2)
}

func unsupportedDistance() float32 {
	return float32(math.Inf(1))
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

func vec3ToValues(vec rl.Vector3) (float32, float32, float32) {
	return vec.X, vec.Y, vec.Z
}
