package hdbscan

import (
	"encoding/json"
	"testing"
)

func TestClustering(t *testing.T) {
	data := [][]float64{
		// cluster-1
		[]float64{1, 2, 3},
		[]float64{1, 2, 4},
		[]float64{1, 2, 5},
		[]float64{1, 3, 4},
		// cluster-2
		[]float64{4, 5, 6},
		[]float64{4, 5, 5},
		// cluster-3
		[]float64{80, 85, 90},
		[]float64{89, 90, 91},
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
