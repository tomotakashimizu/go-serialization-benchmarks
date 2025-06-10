package utils

import "slices"

// CalculateAverage calculates the average of a slice of int64 values
func CalculateAverage(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	var sum int64
	for _, v := range values {
		sum += v
	}
	return sum / int64(len(values))
}

// CalculateMedian calculates the median of a slice of int64 values
func CalculateMedian(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int64, len(values))
	copy(sorted, values)
	slices.Sort(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}
