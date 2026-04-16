package collision

import rl "github.com/gen2brain/raylib-go/raylib"

type PlaneAxis int8

const (
	// PlaneAxisXPos points toward positive X.
	PlaneAxisXPos PlaneAxis = iota
	// PlaneAxisXNeg points toward negative X.
	PlaneAxisXNeg
	// PlaneAxisYPos points toward positive Y.
	PlaneAxisYPos
	// PlaneAxisYNeg points toward negative Y.
	PlaneAxisYNeg
	// PlaneAxisZPos points toward positive Z.
	PlaneAxisZPos
	// PlaneAxisZNeg points toward negative Z.
	PlaneAxisZNeg
)

// Normal returns unit normal vector for axis direction.
func (a PlaneAxis) Normal() rl.Vector3 {
	switch a {
	case PlaneAxisXPos:
		return rl.NewVector3(1, 0, 0)
	case PlaneAxisXNeg:
		return rl.NewVector3(-1, 0, 0)
	case PlaneAxisYPos:
		return rl.NewVector3(0, 1, 0)
	case PlaneAxisYNeg:
		return rl.NewVector3(0, -1, 0)
	case PlaneAxisZPos:
		return rl.NewVector3(0, 0, 1)
	case PlaneAxisZNeg:
		return rl.NewVector3(0, 0, -1)
	default:
		return rl.NewVector3(0, 0, 0)
	}
}
