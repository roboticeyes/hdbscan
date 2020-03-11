package hdbscan

import (
	"fmt"
	"math"
)

type edge struct {
	p1   int
	p2   int
	dist float64
}

type tree struct {
	vertices map[int]bool
	edges    []edge
}

// the edges in the minimum spanning tree will be the minimum
// mutual-reachability distances between any two connected points.
// where no two points are connected by an edge unless they are closest
// to each other relative to other points.
func (c *Clustering) buildMinSpanningTree(graph *graph) {
	mrg := graph.data

	c.mst.vertices[0] = true
	for len(c.mst.edges) < len(mrg) {
		newEdge := c.mst.nearestVertice(mrg)
		if newEdge.p1 == newEdge.p2 {
			break
		}

		fmt.Println("edge: ", newEdge)

		// add new point and new edge to mst
		c.mst.vertices[newEdge.p2] = true
		c.mst.edges = append(c.mst.edges, newEdge)
	}
}

// nearestVertice will find the next smallest edge that is not already
// in the MST. Where "smallness" is found by finding the smallest
// mutual-reachability between two points that is not already an edge in the
// tree.
func (t *tree) nearestVertice(mrg [][]float64) edge {
	minDist := math.MaxFloat64
	p1Index := 0
	p2Index := 0

	for i := 0; i < len(mrg); i++ {
		// if point_i is NOT already a vertex in the mst
		if _, ok := t.vertices[i]; !ok {
			// check distance between point_i and all points already in MST
			for j := range t.vertices {
				// if distance between point_i & vertice_j is smallest-distance
				// left in graph then this will become a new edge in the MST.
				if minDist > mrg[i][j] {
					minDist = mrg[i][j]
					p1Index = j
					p2Index = i
				}
			}
		}
	}

	// create smallest-distance edge between mst vertice_j and point_i
	return edge{p1: p1Index, p2: p2Index, dist: mrg[p1Index][p2Index]}
}
