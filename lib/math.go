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

func NewLinearScale(domain, base []float32) func(value float32) float32 {
	if len(domain) != len(base) || len(domain) < 2 {
		panic("domain and base must be of equal length and have at least two points")
	}

	return func(value float32) float32 {
		// Handle values outside the domain range (extrapolation)
		if value <= domain[0] {
			return base[0]
		}
		if value >= domain[len(domain)-1] {
			return base[len(base)-1]
		}

		// Find the correct interval
		for i := 0; i < len(domain)-1; i++ {
			dStart, dEnd := domain[i], domain[i+1]
			if value >= dStart && value <= dEnd {
				// Corresponding base interval
				bStart, bEnd := base[i], base[i+1]
				percent := (value - dStart) / (dEnd - dStart)
				return bStart + percent*(bEnd-bStart)
			}
		}

		// Should never reach here if domain is sorted and value is within range
		panic("value out of interpolation range")
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
