package hdbscan

func (c *Clustering) selectOptimalClustering(hierarchy *cluster, score string) {
	switch score {
	case VarianceScore:
		setVarianceDelta(hierarchy)
	default:
		// setStabilityDelta(hierarchy)
	}

	var finalClusters []*cluster
	finalClusters = getDelta(hierarchy, finalClusters)
	for i, finalCluster := range finalClusters {
		finalClusters[i].FinalPoints = finalCluster.pointIndexes()
	}
	c.OptimalClustering = finalClusters
}

func getDelta(c *cluster, list []*cluster) []*cluster {
	if c.delta == 1 {
		list = append(list, c)
	}

	for _, childCluster := range c.children {
		subTreeList := getDelta(childCluster, []*cluster{})
		list = append(list, subTreeList...)
	}

	return list
}

func setVarianceDelta(c *cluster) {
	for _, childCluster := range c.children {
		if len(childCluster.children) > 0 {
			setVarianceDelta(childCluster)
		}

		calculateVarianceDelta(childCluster)
	}

	calculateVarianceDelta(c)
}

func calculateVarianceDelta(c *cluster) {
	// if any children are 0 delta, set delta to 0
	var zero bool
	for _, childCluster := range c.children {
		if childCluster.delta == 0 {
			zero = true
		}
	}
	if zero {
		c.delta = 0
		return
	}

	// set delta
	var childrenAverageScore float64
	if len(c.children) > 0 {
		for _, childCluster := range c.children {
			childrenAverageScore += childCluster.score
		}
		childrenAverageScore /= float64(len(c.children))
	}

	if c.score < childrenAverageScore {
		c.delta = 0
	} else {
		c.delta = 1

		for _, childCluster := range c.children {
			setSubTreeDelta(childCluster, 0)
		}
	}
}

func setSubTreeDelta(c *cluster, delta int) {
	c.delta = delta

	for _, childCluster := range c.children {
		setSubTreeDelta(childCluster, delta)
	}
}
