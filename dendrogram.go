package hdbscan

import (
	"fmt"
)

type cluster struct {
	// hierarchy
	parent   *cluster
	children []*cluster

	// data
	points           []int
	pointsDistParent []float64
	FinalPoints      []int

	// score
	score float64
	delta int

	// params
	size     float64
	variance float64
	lMin     float64
}

// a node describes the position
// of a point in the dendrogram of a full
// point-hierarchy.
type node struct {
	key             int
	parentKey       int
	parent          *node
	distToParent    float64
	children        []*node
	descendantCount int
}

func (c *Clustering) buildDendrogram(edgesToProcess []edge, nodes []node) *node {
	if len(edgesToProcess) == 0 {
		// return root node
		for i, nd := range nodes {
			if nd.parent == nil {
				return &nodes[i]
			}
		}
	}

	// find all unique starting points (of remaining edges)
	startingPoints := make(map[int]bool)
	for _, e := range edgesToProcess {
		if _, ok := startingPoints[e.p1]; !ok && e.p1 != e.p2 {
			startingPoints[e.p1] = true
		}
	}

	// find all child nodes
	remainingEdges := []edge{}
	childNodes := []node{}
	for _, e := range c.mst.edges {
		// child-only node ("leaf")
		if _, ok := startingPoints[e.p2]; !ok {
			n := node{
				key:          e.p2,
				parentKey:    e.p1,
				parent:       nil,
				distToParent: e.dist,
				children:     []*node{},
			}

			childNodes = append(childNodes, n)
		} else {
			// edge not processed
			remainingEdges = append(remainingEdges, e)
		}
	}

	// find all nodes that childNodes are parents of
	for i, nl := range childNodes {
		for j, ol := range nodes {
			if ol.parentKey == nl.key && nodes[j].parent == nil {
				childNodes[i].children = append(childNodes[i].children, &nodes[j])
				childNodes[i].descendantCount = childNodes[i].descendantCount + ol.descendantCount + 1
				nodes[j].parent = &childNodes[i]
			}
		}
	}

	// pass root into next nodes
	for _, ol := range nodes {
		if ol.parent == nil {
			childNodes = append(childNodes, ol)
		}
	}

	return c.buildDendrogram(remainingEdges, childNodes)
}

// the clusters hierarchy will not contain clusters that are smaller than the minimum cluster size
// every leaf-cluster is unique subset of points.
func (c *Clustering) buildClusters(root *node, parentCluster *cluster) cluster {
	// set starting cluster
	if parentCluster == nil {
		parentCluster = &cluster{
			parent:           nil,
			points:           []int{},
			pointsDistParent: []float64{},
			children:         []*cluster{},
		}
	}

	// traverse dendrogram tree (from top node)
	for _, childNode := range root.children {
		// if sub-tree is large enough to be a cluster
		// create new (sub-)cluster ...
		if childNode.descendantCount >= c.mcs {
			subCluster := &cluster{
				parent:           parentCluster,
				points:           []int{},
				pointsDistParent: []float64{},
				children:         []*cluster{},
			}
			subCluster.points = append(subCluster.points, childNode.key)
			subCluster.pointsDistParent = append(subCluster.pointsDistParent, childNode.distToParent)
			parentCluster.children = append(parentCluster.children, subCluster)

			c.buildClusters(childNode, subCluster)
		} else {
			// if sub-tree is not large enough to be a cluster
			// add current point to parent cluster
			parentCluster.points = append(parentCluster.points, childNode.key)
			parentCluster.pointsDistParent = append(parentCluster.pointsDistParent, childNode.distToParent)
		}
	}

	fmt.Println("cluster points: ", parentCluster.points)

	return *parentCluster
}
