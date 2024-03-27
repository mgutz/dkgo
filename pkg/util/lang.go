package util

import "strconv"

func ParseIntOr(x string, or int) int {
	if x == "" {
		return or
	}
	i, err := strconv.Atoi(x)
	if err != nil {
		return or
	}
	return i
}
