package helper

import (
	"math"
)

func CalculateTotalPages(total int, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}

func FloatToPtr(f float64) *float64 {
	if f <= 0 {
		return nil
	}
	return &f
}

func IntToPtr(i int) *int {
	if i <= 0 {
		return nil
	}
	return &i
}
