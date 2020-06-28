package hdbscan

import (
	"math"
)

// DistanceFunc ...
type DistanceFunc func(x1, x2 []float64) float64

// EuclideanDistance ...
var EuclideanDistance = func(v1, v2 []float64) float64 {
	acc := 0.0
	for i, v := range v1 {
		acc += math.Pow((v - v2[i]), 2)
	}
	return math.Pow(acc, 0.5)
}

// TODO: add other distance function options...
