package hdbscan

import (
	"encoding/json"
	"testing"
)

func TestClustering(t *testing.T) {
	data := [][]float64{
		[]float64{1, 2, 3},
		[]float64{3, 2, 1},
		[]float64{4, 5, 6},
		[]float64{6, 5, 4},
		[]float64{7, 8, 9},
		[]float64{9, 8, 7},
		[]float64{0, 1, 0},
		[]float64{1, 0, 1},
		[]float64{100, 100, 100},
		[]float64{90, 90, 90},
	}
	minimumClusterSize := 2

	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	err = clustering.Run(EuclideanDistance, VarianceScore)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	jsonData, _ := json.MarshalIndent(clustering.OptimalClustering, "", "  ")

	t.Log(string(jsonData))
}
