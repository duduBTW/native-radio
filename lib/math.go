package lib

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
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}
