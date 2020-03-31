package hdbscan

import (
	"encoding/json"
	"fmt"
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

func TestBuildDendrogram(t *testing.T) {
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

	// graph
	mrg := clustering.mutualReachabilityGraph(EuclideanDistance)
	clustering.buildMinSpanningTree(mrg)

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram(clustering.mst.edges, []node{})

	allNodes := dendrogram.allNewNodes([]debugNode{})

	// check keys are unique and not their own children
	var uniqueKeys []int
	for i, n := range allNodes {
		for _, k := range uniqueKeys {
			if k == n.Key {
				t.Errorf("Node %+v with key %+v is repeated", i, n)
			}
		}

		for j, child := range n.Children {
			if n.Key == child {
				t.Errorf("Node %+v with key %+v is its own child", j, child)
			}
		}
	}

	fmt.Println(dendrogram)
}
