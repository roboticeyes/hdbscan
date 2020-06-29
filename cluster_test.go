package hdbscan

import (
	"fmt"
	"sort"
	"testing"
)

var (
	data = [][]float64{
		// cluster-1 (0-7)
		[]float64{1, 2, 3},
		[]float64{1, 2, 4},
		[]float64{1, 2, 5},
		[]float64{1, 3, 4},
		[]float64{2, 3, 3},
		[]float64{2, 2, 4},
		[]float64{2, 2, 5},
		[]float64{2, 3, 4},
		// cluster-2 (8-15)
		[]float64{21, 15, 6},
		[]float64{22, 15, 5},
		[]float64{23, 15, 7},
		[]float64{24, 15, 8},
		[]float64{21, 15, 6},
		[]float64{22, 16, 5},
		[]float64{23, 17, 7},
		[]float64{24, 18, 8},
		// cluster-3 (16-23)
		[]float64{80, 85, 90},
		[]float64{89, 90, 91},
		[]float64{100, 100, 100}, // possible outlier
		[]float64{90, 90, 90},
		[]float64{81, 85, 90},
		[]float64{89, 91, 91},
		[]float64{100, 101, 100}, // possible outlier
		[]float64{90, 91, 90},
		// outlier
		[]float64{-2400, 2000, -30},
	}
	minimumClusterSize = 3
)

func TestMinimumSpanningTree(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	clustering.distanceFunc = EuclideanDistance
	clustering.minTree = true

	// graph
	fmt.Println(clustering.mutualReachabilityGraph())
}

func TestBuildDendrogram(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	clustering.distanceFunc = EuclideanDistance
	clustering.minTree = true

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram(clustering.mutualReachabilityGraph())

	for _, link := range dendrogram {
		t.Logf("Link %+v with points: %+v", link.id, link.points)
	}
}

func TestBuildClusters(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	clustering.distanceFunc = EuclideanDistance
	// clustering.minTree = true

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram(clustering.mutualReachabilityGraph())
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
	clustering.distanceFunc = EuclideanDistance

	// cluster-hierarchy
	dendrogram := clustering.buildDendrogram(clustering.mutualReachabilityGraph())
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

func TestClusteringNoTree(t *testing.T) {
	clustering, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	err = clustering.Run(EuclideanDistance, VarianceScore, false)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range clustering.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusteringVerbose(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Verbose()

	err = c.Run(EuclideanDistance, VarianceScore, false)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusteringSampling(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Verbose().Subsampling(16)

	err = c.Run(EuclideanDistance, VarianceScore, true)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusteringSamplingAndAssign(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Subsampling(16).OutlierDetection()

	err = c.Run(EuclideanDistance, VarianceScore, true)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	newClustering, err := c.Assign(data)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range newClustering.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusteringSamplingAndAssignAndOutlierClustering(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Subsampling(16).NearestNeighbor().OutlierClustering()

	err = c.Run(EuclideanDistance, VarianceScore, true)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	newClustering, err := c.Assign(data)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range newClustering.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusteringOutliers(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}

	c = c.OutlierDetection().NearestNeighbor()

	err = c.Run(EuclideanDistance, VarianceScore, true)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with Points %+v and outliers: %+v", cluster.id, cluster.Points, cluster.Outliers)
	}
}

func TestClusteringVoronoi(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Verbose().Voronoi()

	err = c.Run(EuclideanDistance, VarianceScore, true)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}

func TestClusteringVoronoiParts(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Verbose().Voronoi()
	c.distanceFunc = EuclideanDistance
	c.minTree = true

	edges := c.mutualReachabilityGraph()
	t.Logf("%+v\n", edges)
	dendrogram := c.buildDendrogram(edges)
	for _, link := range dendrogram {
		t.Logf("Link %+v with points: %+v", link.id, link.points)
	}

	c.buildClusters(dendrogram)
	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}

	c.scoreClusters(VarianceScore)
	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with variance %+v and score %+v and points: %+v", cluster.id, cluster.variance, cluster.score, cluster.Points)
	}

	c.selectOptimalClustering(VarianceScore)
	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with variance %+v and score %+v and points: %+v", cluster.id, cluster.variance, cluster.score, cluster.Points)
	}

	c.clusterCentroids()
	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with variance %+v and score %+v and points: %+v and Centroid %+v", cluster.id, cluster.variance, cluster.score, cluster.Points, cluster.Centroid)
	}

	c.outliersAndVoronoi()
	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with variance %+v and score %+v and points: %+v and Centroid %+v", cluster.id, cluster.variance, cluster.score, cluster.Points, cluster.Centroid)
	}
}

func TestClusteringVoronoiNoTree(t *testing.T) {
	c, err := NewClustering(data, minimumClusterSize)
	if err != nil {
		t.Errorf("clustering creation error: %+v", err)
	}
	c = c.Verbose().Voronoi()

	err = c.Run(EuclideanDistance, VarianceScore, false)
	if err != nil {
		t.Errorf("clustering run error: %+v", err)
	}

	for _, cluster := range c.Clusters {
		sort.Ints(cluster.Points)
		t.Logf("Cluster %+v with points: %+v", cluster.id, cluster.Points)
	}
}
