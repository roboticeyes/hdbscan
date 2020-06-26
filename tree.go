package hdbscan

import (
	"math"
)

type edge struct {
	p1   int
	p2   int
	dist float64
}

type edges []edge

type tree struct {
	vertices []int
	edges    edges
}

// the edges in the minimum spanning tree will be the minimum
// mutual-reachability distances between any two connected points.
// where no two points are connected by an edge unless they are closest
// to each other relative to other points.
func (c *Clustering) addRowToMinSpanningTree(row int, data []float64) {
	if row == 0 {
		c.mst.vertices = append(c.mst.vertices, 0)
	}

	newEdge := c.mst.nearestVertice(row, data)
	if newEdge.p1 == newEdge.p2 {
		return
	}

	// add new point and new edge to mst
	c.mst.vertices = append(c.mst.vertices, newEdge.p2)
	c.mst.edges = append(c.mst.edges, newEdge)
	c.extendMinSpanningTree()
}

// nearestVertice will find the next smallest edge that is not already
// in the MST. Where "smallness" is found by finding the smallest
// mutual-reachability between two points that is not already an edge in the
// tree.
func (t *tree) nearestVertice(row int, data []float64) edge {
	minDist := math.MaxFloat64
	p1Index := 0
	p2Index := 0

	if !containsInt(t.vertices, row) {
		// check distance between point_i and all points already in MST
		for _, j := range t.vertices {
			// if distance between point_i & vertice_j is the smallest-distance
			// left in graph then this will become a new edge in the MST.
			if minDist > data[j] {
				minDist = data[j]
				p1Index = j
				p2Index = row
			}
		}
	}

	// create smallest-distance edge between mst vertice_j and point_i
	return edge{p1: p1Index, p2: p2Index, dist: data[p1Index]}
}

func (c *Clustering) extendMinSpanningTree() {
	for vertice := range c.mst.vertices {
		newEdge := edge{
			p1:   vertice,
			p2:   vertice,
			dist: c.coreDistances[vertice],
		}
		c.mst.edges = append(c.mst.edges, newEdge)
	}
}

// Len ...
func (e edges) Len() int {
	return len(e)
}

// Swap ...
func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less ...
func (e edges) Less(i, j int) bool {
	return e[j].dist > e[i].dist
}
