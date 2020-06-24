package hdbscan

import (
	"sort"
	"sync"
)

type graph struct {
	data [][]float64
	*sync.Mutex
}

func newGraph() *graph {
	return &graph{
		data:  make([][]float64, 0),
		Mutex: &sync.Mutex{},
	}
}

func (g *graph) add(newData []float64) {
	g.Lock()
	g.data = append(g.data, newData)
	g.Unlock()
}

// the mutual reachability graph provides a mutual-reachability-distance matrix
// which specifies a metric of how far a point is from another point.
func (c *Clustering) mutualReachabilityGraph(distanceFunc DistanceFunc) {
	length := len(c.data)
	mutualReachabililtyGraph := newGraph()

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
	for i := 0; i < length; i++ {
		// c.sempahore <- true
		// c.wg.Add(1)
		// go func(i int) {
		mutualReachabilityDistances := []float64{}

		// the mutual reachability distance is the maximum of:
		// point_1's core-distance, point_2's core-distance, or the distance between point_1 and point_2
		for j := 0; j < length; j++ {
			mutualReachabilityDistances = append(mutualReachabilityDistances, max([]float64{coreDistances[i], coreDistances[j], distanceMatrix.Get(i)[j]}))
		}

		mutualReachabililtyGraph.add(mutualReachabilityDistances)

		// 	<-c.sempahore
		// 	c.wg.Done()
		// }(i)
	}
	// c.wg.Wait()

	c.distanceMatrix = mutualReachabililtyGraph
}
