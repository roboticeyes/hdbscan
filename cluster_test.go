package hdbscan

import (
	"fmt"
	"sort"
	"testing"
)

var (
	data = [][]float64{
		// cluster-1
		[]float64{1, 2, 3},
		[]float64{1, 2, 4},
		[]float64{1, 2, 5},
		[]float64{1, 3, 4},
		[]float64{2, 3, 3},
		[]float64{2, 2, 4},
		[]float64{2, 2, 5},
		[]float64{2, 3, 4},
		// cluster-2
		[]float64{21, 15, 6},
		[]float64{22, 15, 5},
		[]float64{23, 15, 7},
		[]float64{24, 15, 8},
		[]float64{21, 15, 6},
		[]float64{22, 16, 5},
		[]float64{23, 17, 7},
		[]float64{24, 18, 8},
		// cluster-3
		[]float64{80, 85, 90},
		[]float64{89, 90, 91},
		[]float64{100, 100, 100},
		[]float64{90, 90, 90},
		[]float64{81, 85, 90},
		[]float64{89, 91, 91},
		[]float64{100, 101, 100},
		[]float64{90, 91, 90},
	}
	minimumClusterSize = 3
)

func TestMinimumSpanningTree(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	clustering.minTree = true

	// graph
	clustering.mutualReachabilityGraph(EuclideanDistance)
	clustering.buildMinSpanningTree()

	fmt.Println(clustering.mst.edges)
}

func TestBuildDendrogram(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	// graph
	clustering.mutualReachabilityGraph(EuclideanDistance)
	clustering.buildMinSpanningTree()

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram()

	for _, link := range dendrogram {
		t.Logf("Link %+v with points: %+v", link.id, link.points)
	}
}

func TestBuildClusters(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	// graph
	clustering.mutualReachabilityGraph(EuclideanDistance)
	clustering.buildMinSpanningTree()

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram()
	clustering.buildClusters(dendrogram)

	for _, cluster := range clustering.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusterScoring(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	// graph
	clustering.mutualReachabilityGraph(EuclideanDistance)
	clustering.buildMinSpanningTree()

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram()
	clustering.buildClusters(dendrogram)
	clustering.scoreClusters(VarianceScore)

	for _, cluster := range clustering.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with variance %+v and score %+v and points: %+v", cluster.id, cluster.variance, cluster.score, cluster.Points)
	}
}

func TestClustering(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	err = clustering.Run(EuclideanDistance, VarianceScore, true)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range clustering.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}
