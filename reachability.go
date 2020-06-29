package hdbscan

import (
	"log"
	"sort"
)

// the mutual reachability graph provides a mutual-reachability-distance matrix
// which specifies a metric of how far a point is from another point.
func (c *Clustering) mutualReachabilityGraph() edges {
	if c.verbose {
		log.Println("starting mutual reachability")
	}

	// core-distances
	length := len(c.data)
	coreDistances := make([]float64, length, length)
	for i, p1 := range c.data {
		c.wg.Add(1)
		c.semaphore <- true
		go func(i int, p1 []float64) {
			pointDistances := []float64{}
			for _, p2 := range c.data {
				pointDistances = append(pointDistances, c.distanceFunc(p1, p2))
			}
			sort.Float64s(pointDistances)
			coreDistances[i] = pointDistances[c.mcs-1]
			<-c.semaphore
			c.wg.Done()
		}(i, p1)
	}
	c.wg.Wait()

	// mutual-reachability distances
	for i := 0; i < length; i++ {
		c.wg.Add(1)
		c.semaphore <- true
		go func(i int) {
			mutualReachabilityDistances := make([]float64, length, length)
			// the mutual reachability distance is the maximum of:
			// point_1's core-distance, point_2's core-distance, or the distance between point_1 and point_2
			for j := 0; j < length; j++ {
				mutualReachabilityDistances[j] = max([]float64{coreDistances[i], coreDistances[j], c.distanceFunc(c.data[i], c.data[j])})
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

				// just use tree for edge storage
				c.mst.addEdge(e)
			}
			<-c.semaphore
			c.wg.Done()
		}(i)
	}
	c.wg.Wait()

	sort.Sort(c.mst.edges)

	if c.verbose {
		log.Println("finished mutual reachability")
	}

	return c.mst.edges
}
