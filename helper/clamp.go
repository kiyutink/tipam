package helper

import "golang.org/x/exp/constraints"

func Clamp[T constraints.Ordered](val T, min T, max T) T {
	if val < min {
		return min
	}
	if val > max {
		return max
	}

	return val
}
