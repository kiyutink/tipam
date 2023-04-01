package main

import "golang.org/x/exp/constraints"

func clamp[T constraints.Ordered](val T, min T, max T) T {
	if val < min {
		return min
	}
	if val > max {
		return max
	}

	return val
}
