package hdbscan

import (
	"log"
	"math"
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
)

type cluster struct {
	id       int
	parent   *int
	children []int
	Centroid []float64
	Points   []int
	Outliers Outliers
	score    float64
	delta    int
	size     float64
	variance float64
	lMin     float64
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
	data [][]float64

	// settings
	mcs          int
	minTree      bool
	verbose      bool
	randomSample bool
	subSample    bool
	voronoi      bool
	nn           bool
	od           bool

	// minimum spanning tree
	mst *tree

	// optimal-clustering
	score    string
	Clusters clusters

	semaphore chan bool
	wg        *sync.WaitGroup
}

// NewClustering creates (a pointer to) a new clustering struct.
// This function does not automatically start the clustering
// process. The `Run` method needs to be called to do that.
// Make sure to apply all options *before* calling `Run`.
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
		data:      data,
		mcs:       minimumClusterSize,
		mst:       newTree(),
		semaphore: make(chan bool, runtime.NumCPU()),
		wg:        &sync.WaitGroup{},
	}, nil
}

// Run will run the clustering.
func (c *Clustering) Run(distanceFunc DistanceFunc, score string, mst bool) error {
	c.minTree = mst
	if c.verbose && !c.minTree {
		log.Println("not using minimum spanning tree")
	}

	edges := c.mutualReachabilityGraph(distanceFunc)
	dendrogram := c.buildDendrogram(edges)
	c.buildClusters(dendrogram)
	c.scoreClusters(score)
	c.selectOptimalClustering(score)
	c.clusterCentroids()
	c.outliersAndVoronoi(distanceFunc)

	return nil
}

// OutlierDetection will track all unassigned
// points as outliers of their nearest cluster.
// It provides a `NormalizedDistance` value for
// each outlier which can be interpreted as the
// probability of the point being an outlier
// (relative to all other outliers).
func (c *Clustering) OutlierDetection() *Clustering {
	c.od = true
	return c
}

// Verbose will set verbosity to true for clustering process
// and the internals of a clustering run will be logged to stdout.
func (c *Clustering) Verbose() *Clustering {
	c.verbose = true
	return c
}

// Voronoi will set voronoi-clustering to true, and
// after density clustering is performed,
// all points not assigned to a cluster will be placed
// into their nearest cluster (by centroid distance).
func (c *Clustering) Voronoi() *Clustering {
	c.voronoi = true
	return c
}

// NearestNeighbor specifies if nearest-neighbor
// distances should be used for outlier detection
// and for voronoi clustering instead of centroid-based
// distances.
// NearestNeighbor will find the closest assigned data
// point to an unassigned data point and consider the
// unassigned data point to be of that same cluster (as an outlier and/or a point).
func (c *Clustering) NearestNeighbor() *Clustering {
	c.nn = true
	return c
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
	}

	c.Clusters = clusters

	if c.verbose {
		log.Println("finished building clusters")
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

func (c *Clustering) outliersAndVoronoi(distanceFunc DistanceFunc) {
	if !c.od && !c.voronoi {
		return
	}

	if c.verbose {
		if c.od {
			log.Println("finding outliers")
		}

		if c.voronoi {
			log.Println("starting voronoi clustering")
		}
	}

	for i, v := range c.data {
		var exists bool
		for _, cluster := range c.Clusters {
			for _, point := range cluster.Points {
				if point == i {
					exists = true
					break
				}
			}

			if exists {
				break
			}
		}

		if !exists {
			// calculate nearest cluster
			minDistance := math.MaxFloat64
			var nearestClusterIndex int
			for i, cluster := range c.Clusters {
				if c.nn {
					for _, p := range cluster.Points {
						distance := distanceFunc(c.data[p], v)
						if distance < minDistance {
							minDistance = distance
							nearestClusterIndex = i
						}
					}
				} else {
					distance := distanceFunc(cluster.Centroid, v)
					if distance < minDistance {
						minDistance = distance
						nearestClusterIndex = i
					}
				}
			}

			// voronoi cluster
			if c.voronoi && len(c.Clusters) > 0 {
				c.Clusters[nearestClusterIndex].Points = append(c.Clusters[nearestClusterIndex].Points, i)
			}

			// outlier detection
			if c.od && len(c.Clusters) > 0 {
				c.Clusters[nearestClusterIndex].Outliers = append(c.Clusters[nearestClusterIndex].Outliers, Outlier{
					Index:              i,
					NormalizedDistance: minDistance,
				})
			}
		}
	}

	// normalize outlier distances
	if c.od {
		for i, cluster := range c.Clusters {
			if len(cluster.Outliers) == 0 {
				continue
			}

			// distance distribution
			var distances []float64
			for j1, p1 := range cluster.Points {
				for j2, p2 := range cluster.Points {
					if j1 != j2 {
						distance := distanceFunc(c.data[p1], c.data[p2])
						distances = append(distances, distance)
					}
				}
			}

			dd := distuv.Normal{}
			dd.Fit(distances, nil)

			for j, outlier := range cluster.Outliers {
				outlier.NormalizedDistance = isNum(dd.CDF(outlier.NormalizedDistance))
				cluster.Outliers[j] = outlier
			}

			c.Clusters[i] = cluster
		}
	}

	if c.verbose {
		if c.od {
			log.Println("finished finding outliers")
		}

		if c.voronoi {
			log.Println("finished voronoi clustering")
		}
	}
}

func (c clusters) getClusterByID(id int) *cluster {
	for _, cluster := range c {
		if cluster.id == id {
			return cluster
		}
	}

	return nil
}

// Len ...
func (c clusters) Len() int {
	return len(c)
}

// Swap ...
func (c clusters) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less ...
func (c clusters) Less(i, j int) bool {
	return len(c[i].Points) < len(c[j].Points)
}
