package utils

func CalculateGrowth(current, previous int) float64 {
	if previous == 0 {
		return 0
	}

	return float64(current-previous) / float64(previous) * 100
}
