package hdbscan

import (
	"runtime"
	"sync"
)

var (
	// VarianceScore will select an optimal clustering that minimizes the generalized variance across each cluster.
	VarianceScore = "variance_score"
	// StabilityScore will select an optimal clustering that maximized the stability across all clusters.
	StabilityScore = "stability_score"
)

// Clustering ...
type Clustering struct {
	data [][]float64
	mcs  int // minimum cluster size

	// minimum spanning tree
	mst *tree

	// optimal-clustering
	score             string
	OptimalClustering []*cluster

	sempahore chan bool
	wg        *sync.WaitGroup
}

// NewClustering ...
func NewClustering(data [][]float64, minimumClusterSize int) (*Clustering, error) {
	if minimumClusterSize < 1 {
		return &Clustering{}, ErrMCS
	}

	if len(data) < minimumClusterSize {
		return &Clustering{}, ErrDataLen
	}

	dataLength := len(data[0])

	for _, row := range data {
		if len(row) != dataLength {
			return &Clustering{}, ErrRowLength
		}
	}

	return &Clustering{
		data: data,
		mcs:  minimumClusterSize,
		mst: &tree{
			vertices: make(map[int]bool),
		},
		sempahore: make(chan bool, runtime.NumCPU()),
		wg:        &sync.WaitGroup{},
	}, nil
}

// Run ...
func (c *Clustering) Run(distanceFunc DistanceFunc, score string) error {
	// graph
	mrg := c.mutualReachabilityGraph(distanceFunc)
	c.buildMinSpanningTree(mrg)

	// cluster-hierarchy
	hierarchy := c.buildHierarchy(c.mst.edges, []node{})
	clusterHierarchy := c.clustersHierarchy(&hierarchy, nil)

	// optimal-clustering
	err := c.scoreClusters(&clusterHierarchy, score)
	if err != nil {
		return err
	}

	c.selectOptimalClustering(&clusterHierarchy, score)

	return nil
}
