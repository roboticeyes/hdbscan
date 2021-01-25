// Copyright 2020 Humility AI Incorporated, All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hdbscan

import (
	"log"
	"runtime"
	"sync"

	"gonum.org/v1/gonum/stat/distuv"
)

var (
	// VarianceScore will select an optimal clustering
	// that minimizes the generalized variance across each cluster.
	VarianceScore = "variance_score"
	// StabilityScore will select an optimal clustering that
	// maximized the stability across all clusters.
	StabilityScore = "stability_score"

	Leaf = "leaf"
)

type cluster struct {
	id int

	parent               *int
	children             []int
	score                float64
	delta                int
	size                 float64
	variance             float64
	lMin                 float64
	lambdaBirth          float64
	distanceDistribution *distuv.Normal
	largestDistance      float64
	viseted              bool
	//public
	Centroid []float64
	Points   []int
	Outliers Outliers
}

type clusters []*cluster

// Outlier struct is used to provide information
// about an outlier data point.
type Outlier struct {
	Index              int
	NormalizedDistance float64
}

// Outliers is an array of outlier points
// for a given cluster.
type Outliers []Outlier

// Clustering struct which holds
// all final results.
type Clustering struct {
	data      [][]float64
	directory string

	// settings
	mcs          int
	minTree      bool
	verbose      bool
	randomSample bool
	subSample    bool
	voronoi      bool
	nn           bool // NearestNeighbor
	od           bool // Outlier detection
	oc           bool // Outlier Clustering
	sampleBound  int
	distanceFunc DistanceFunc

	// minimum spanning tree
	mst *tree

	// optimal-clustering
	score            string
	Clusters         clusters
	ClustersReverse  clusters
	NumberOfClusters int
	lambda           [][]float64

	// Multithreading
	semaphore chan bool
	wg        *sync.WaitGroup
}

// NewClustering creates (a pointer to) a new clustering struct.
// This function does not automatically start the clustering
// process. The `Run` method needs to be called to do that.
// Make sure to apply all options *before* calling `Run`.
func NewClustering(data [][]float64, minimumClusterSize int, directory string) (*Clustering, error) {
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
		data:      data,
		directory: directory,
		mcs:       minimumClusterSize,
		mst:       newTree(),
		semaphore: make(chan bool, runtime.NumCPU()),
		wg:        &sync.WaitGroup{},
	}, nil
}

// Run will run the clustering.
func (c *Clustering) Run(distanceFunc DistanceFunc, score string, mst bool) error {
	c.distanceFunc = distanceFunc
	c.minTree = mst
	if c.verbose && !c.minTree {
		log.Println("not using minimum spanning tree")
	}

	// Will not be used right now
	c.sample()
	// Calculate "Mutual Reachability Graph" and build minimum spaning tree
	edges := c.mutualReachabilityGraph()
	// Filter edges by MAD - Use only if point distance is equidistant
	// edges = c.filterEdges(edges)
	// Plot minimum spanning tree
	c.plotminimumSpanningTree(edges)
	// Build dendogram
	dendogram := c.buildDendogram(edges)
	// Build Clusters
	c.buildClusters(dendogram)

	// Calculate stability
	c.scoreClusters(score)
	// Write Clusters to file before selecting the clusters
	// c.writeClusterToFile("before")
	// Select Clusters
	c.selectOptimalClustering(score)
	// Write Clusters to file after selecting the clusters
	// c.writeClusterToFile("after")
	// Calculate centroids for every cluster
	c.clusterCentroids()
	// Outlier detection
	c.outliersAndVoronoi()
	// If oc (outlier clustering) is true
	// all outliers from a cluster become a cluster of their own
	c.outlierClustering()

	return nil
}

// the clusters hierarchy will not contain clusters that are smaller than the minimum cluster size
// every leaf-cluster is unique subset of points.
func (c *Clustering) buildClusters(links []*link) {
	if c.verbose {
		log.Println("building clusters")
	}

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
				Outliers: make(Outliers, 0),
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
		c.NumberOfClusters = cluster.id
	}

	c.Clusters = clusters
	if c.verbose {
		log.Println("finished building clusters, Number of clusters: ", len(c.Clusters))
	}
}

func (c *Clustering) clusterCentroids() {
	if c.verbose {
		log.Println("calculating cluster centroids")
	}

	for i, cluster := range c.Clusters {
		avg := make([]float64, len(c.data[0]), len(c.data[0]))
		for _, index := range cluster.Points {
			vec := c.data[index]
			if len(vec) == len(avg) {
				for j, v := range vec {
					avg[j] += v
				}
			}
		}

		for k, v := range avg {
			v /= float64(len(cluster.Points))
			avg[k] = v
		}

		cluster.Centroid = avg
		c.Clusters[i] = cluster
	}

	if c.verbose {
		log.Println("finished calculating cluster centroids")
	}
}
