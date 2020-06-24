package hdbscan

import (
	"math"
	"sync"
)

// DistanceFunc ...
type DistanceFunc func(x1, x2 []float64) float64

// DistanceMatrix ...
type DistanceMatrix struct {
	data [][]float64
	*sync.Mutex
}

// NewDistanceMatrix ...
func NewDistanceMatrix() *DistanceMatrix {
	return &DistanceMatrix{
		data:  make([][]float64, 0),
		Mutex: &sync.Mutex{},
	}
}

// Add ...
func (d *DistanceMatrix) Add(newData []float64) {
	d.Lock()
	d.data = append(d.data, newData)
	d.Unlock()
}

// Get ...
func (d *DistanceMatrix) Get(index int) []float64 {
	if index < 0 || index >= len(d.data) {
		return []float64{}
	}

	return d.data[index]
}

// EuclideanDistance ...
var EuclideanDistance = func(v1, v2 []float64) float64 {
	acc := 0.0
	for i, v := range v1 {
		acc += math.Pow((v - v2[i]), 2)
	}
	return math.Pow(acc, 0.5)
}

// TODO: add other distance function options...
