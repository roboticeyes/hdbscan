package hdbscan

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

// a node in the minimum spanning tree
// store the index of a point_i as "key"
// and the index of it's parent vertice point_j
// as "parentKey", it also has a pointer to it's parent
// node as well as a list of child nodes.
type node struct {
	key             int
	parentKey       int
	parent          *node
	distToParent    float64
	children        []*node
	descendantCount int
}

func (c *Clustering) buildHierarchy(remainingEdges []edge, leaves []node) node {
	if len(remainingEdges) == 0 {
		// return root node
		return leaves[0]
	}

	firstPoints := make(map[int]bool)
	for _, e := range remainingEdges {
		// add starting vertices to list
		if _, ok := firstPoints[e.p1]; !ok && e.p1 != e.p2 {
			firstPoints[e.p1] = true
		}
	}

	stillRemainingEdges := []edge{}
	potentialLeaves := []node{}
	for _, e := range c.mst.edges {
		// if potential leaf (p2 is not a starting vertex)
		if _, ok := firstPoints[e.p2]; !ok {
			potentialLeaves = append(potentialLeaves,
				node{
					key:             e.p2,
					parentKey:       e.p1,
					parent:          nil,
					distToParent:    e.dist,
					children:        []*node{},
					descendantCount: 0,
				},
			)
		} else {
			// if not a potential leaf
			stillRemainingEdges = append(stillRemainingEdges, e)
		}
	}

	for i, nl := range potentialLeaves {
		for j, ol := range leaves {
			// if potential leaf is actually a parent of another node
			// update the potential leaf and the other node.
			if ol.parentKey == nl.key {
				potentialLeaves[i].children = append(potentialLeaves[i].children, &leaves[j])
				potentialLeaves[i].descendantCount = potentialLeaves[i].descendantCount + ol.descendantCount + 1
				leaves[j].parent = &potentialLeaves[i]
			}
		}
	}

	// add all disconnected nodes to list of potential leaves
	for _, ol := range leaves {
		if ol.parent == nil {
			potentialLeaves = append(potentialLeaves, ol)
		}
	}

	return c.buildHierarchy(stillRemainingEdges, potentialLeaves)
}

// the clusters hierarchy will not contain clusters that are smaller than the minimum cluster size
// every leaf-cluster is unique subset of points.
func (c *Clustering) clustersHierarchy(root *node, parentCluster *cluster) cluster {
	// set starting cluster
	if parentCluster == nil {
		parentCluster = &cluster{
			parent:           nil,
			points:           []int{},
			pointsDistParent: []float64{},
			children:         []*cluster{},
		}
	}

	// traverse minimum spanning tree (from top node)
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

			// find and build all sub-clusters of new cluster
			c.clustersHierarchy(childNode, subCluster)
		} else {
			// if sub-tree is not large enough to be a cluster
			// add current point to current cluster
			parentCluster.points = append(parentCluster.points, childNode.key)
			parentCluster.pointsDistParent = append(parentCluster.pointsDistParent, childNode.distToParent)
		}
	}

	return *parentCluster
}
