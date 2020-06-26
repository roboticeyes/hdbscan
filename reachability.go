package hdbscan

import (
	"sort"
	"sync"
)

type graph struct {
	data []float64
	*sync.Mutex
}

func newGraph() *graph {
	return &graph{
		data:  make([]float64, 0),
		Mutex: &sync.Mutex{},
	}
}

func (g *graph) add(newData []float64) {
	g.Lock()
	g.data = append(g.data, newData...)
	g.Unlock()
}

// the mutual reachability graph provides a mutual-reachability-distance matrix
// which specifies a metric of how far a point is from another point.
func (c *Clustering) mutualReachabilityGraph(distanceFunc DistanceFunc) edges {
	var bases edges

	// distance matrix
	distanceMatrix := NewDistanceMatrix()
	for _, p1 := range c.data {
		pointDistances := []float64{}

		for _, p2 := range c.data {
			pointDistances = append(pointDistances, distanceFunc(p1, p2))
		}

		distanceMatrix.Add(pointDistances)
	}

	// core distances = distance from point to its mcs-1 nearest-neighbor
	length := len(c.data)
	coreDistances := []float64{}
	for i := 0; i < length; i++ {
		pointDistances := []float64{}
		pointDistances = append(pointDistances, distanceMatrix.Get(i)...)
		sort.Float64s(pointDistances)
		coreDistances = append(coreDistances, pointDistances[c.mcs-1])
	}

	// set core distances for later use
	c.coreDistances = coreDistances

	// mutualReachabililtyGraph
	mutualReachabilityDistances := make([]float64, length, length)
	for i := 0; i < length; i++ {
		// the mutual reachability distance is the maximum of:
		// point_1's core-distance, point_2's core-distance, or the distance between point_1 and point_2
		for j := 0; j < length; j++ {
			mutualReachabilityDistances[j] = max([]float64{coreDistances[i], coreDistances[j], distanceMatrix.Get(i)[j]})
		}

		if c.minTree {
			c.addRowToMinSpanningTree(i, mutualReachabilityDistances)
		} else {
			minIndex, minValue := min(mutualReachabilityDistances)
			e := edge{
				p1:   i,
				p2:   minIndex,
				dist: minValue,
			}
			bases = append(bases, e)
		}
	}

	if c.minTree {
		bases = c.mst.edges
	}

	sort.Sort(bases)

	return bases
}
