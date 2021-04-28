package common

import (
	"math"

	"github.com/sta-golang/go-lib-utils/log"
)

type Distance func(a, b []float64) float64

func lengthIsEq(a, b []float64) bool {
	length1 := len(a)
	length2 := len(b)

	if length1 != length2 {
		return false
	}
	return true
}

// CosSimi 余弦相似度
func CosSimi(a, b []float64) float64 {
	if !lengthIsEq(a, b) {
		log.ConsoleLogger.Warn("vector a vector b length not equal")
		return -1
	}
	rr, f1r, f2r := 0.0, 0.0, 0.0
	for i := 0; i < len(a); i++ {
		rr += a[i] * b[i]
		f1r += a[i] * a[i]
		f2r += b[i] * b[i]
	}
	return rr / (math.Sqrt(f1r) * math.Sqrt(f2r))
}

// Euclidean 欧式距离
func Euclidean(a, b []float64) float64 {
	if !lengthIsEq(a, b) {
		log.ConsoleLogger.Warn("vector a vector b length not equal")
		return -1
	}
	res := 0.0
	for i := 0; i < len(a); i++ {
		res += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(res)
}
