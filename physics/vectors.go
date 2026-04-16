package physics

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func GetRotationX(vector rl.Vector2) float32 {
	return rl.Vector2Angle(rl.NewVector2(1, 0), vector)*rl.Rad2deg + 180
}
