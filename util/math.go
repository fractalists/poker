package util

func Min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

func Max(a, b float32) float32 {
	if a >= b {
		return a
	}

	return b
}
