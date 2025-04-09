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
