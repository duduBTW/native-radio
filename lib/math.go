package lib

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Max(value, max float32) float32 {
	if value < max {
		return max
	}

	return value
}

func Min(value, min float32) float32 {
	if value > min {
		return min
	}

	return value
}

func MinInt32(value, min int32) int32 {
	if value > min {
		return min
	}

	return value
}
func MinInt(value, min int) int {
	if value > min {
		return min
	}

	return value
}

func MaxInt32(value, max int32) int32 {
	if value < max {
		return max
	}

	return value
}

func MaxInt(value, max int) int {
	if value < max {
		return max
	}

	return value
}

func NewLinearScale(domain, base [2]float32) func(value float32) float32 {
	dStart := domain[0]
	dEnd := domain[1]

	bStart := base[0]
	bEnd := base[1]

	return func(value float32) float32 {
		percent := (value - dStart) / (dEnd - dStart)
		return bStart + percent*(bEnd-bStart)
	}
}

func Clamp(value, min, max float32) float32 {
	if value >= max {
		return max
	}

	if value <= min {
		return min
	}

	return value
}

func CheckCollisionPointCircle(centerX, centerY int32, radius float32, mousePoint rl.Vector2) bool {
	dx := centerX - int32(mousePoint.X)
	dy := centerY - int32(mousePoint.Y)
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= int32(radius*radius)
}

func RandomRange(min, max int) int {
	return rand.IntN(max-min) + min
}
