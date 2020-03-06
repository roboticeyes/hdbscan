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
	if len(v1) != len(v2) {
		return 0
	}

	var total float64
	for i, v := range v1 {
		diff := v - v2[i]
		diffSquared := diff * diff
		total += diffSquared
	}

	return math.Sqrt(total)
}

// TODO: add other distance function options...
