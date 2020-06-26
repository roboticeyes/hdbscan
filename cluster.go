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

type cluster struct {
	id       int
	parent   *int
	children []int
	Points   []int
	score    float64
	delta    int
	size     float64
	variance float64
	lMin     float64
}

type clusters []*cluster

// Clustering ...
type Clustering struct {
	data [][]float64
	// distanceMatrix *graph
	mcs int // minimum cluster size

	// minimum spanning tree
	minTree       bool
	mst           *tree
	coreDistances []float64

	// optimal-clustering
	score    string
	Clusters clusters

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
			vertices: []int{},
		},
		sempahore: make(chan bool, runtime.NumCPU()),
		wg:        &sync.WaitGroup{},
	}, nil
}

// Run ...
func (c *Clustering) Run(distanceFunc DistanceFunc, score string, mst bool) error {
	c.minTree = mst

	dendrogram := c.buildDendrogram(c.mutualReachabilityGraph(distanceFunc))
	c.buildClusters(dendrogram)
	c.scoreClusters(score)
	c.selectOptimalClustering(score)

	return nil
}

// the clusters hierarchy will not contain clusters that are smaller than the minimum cluster size
// every leaf-cluster is unique subset of points.
func (c *Clustering) buildClusters(links []*link) {
	var clusters clusters

	for i, link := range links {
		if len(link.points) >= c.mcs {
			var children []int
			for _, childLink := range link.children {
				children = append(children, childLink.id)
			}

			newCluster := &cluster{
				id:       i,
				Points:   link.points,
				children: children,
			}

			if link.parent == nil {
				newCluster.parent = nil
			} else {
				id := link.parent.id
				newCluster.parent = &id
			}

			clusters = append(clusters, newCluster)
		}
	}

	for _, cluster := range clusters {
		var newChildren []int
		for _, child := range cluster.children {
			if clusters.getClusterByID(child) != nil {
				newChildren = append(newChildren, child)
			}
		}
		cluster.children = newChildren
	}

	c.Clusters = clusters
}

func (c clusters) getClusterByID(id int) *cluster {
	for _, cluster := range c {
		if cluster.id == id {
			return cluster
		}
	}

	return nil
}
