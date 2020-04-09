package hdbscan

func (c *Clustering) selectOptimalClustering(score string) {
	switch score {
	case VarianceScore:
		c.setVarianceDeltas()
	default:
		// setStabilityDelta(hierarchy)
	}

	var finalClusters clusters
	for _, cluster := range c.Clusters {
		if cluster.delta == 1 {
			finalClusters = append(finalClusters, cluster)
		}
	}

	c.Clusters = finalClusters
}

func (c *Clustering) setVarianceDeltas() {
	// var one bool
	for _, cluster := range c.Clusters {
		// calculate average childrens scores
		var avgScore float64
		for _, child := range cluster.children {
			// calculate childrens scores
			avgScore += c.Clusters.getClusterByID(child).score
		}
		avgScore /= float64(len(cluster.children))

		if cluster.score <= avgScore && len(cluster.children) > 0 {
			cluster.delta = 0
		} else {
			cluster.delta = 1

			// check if any subclusters are already delta-1
			var subDeltaOne bool
			subClusters := c.Clusters.subTree(cluster.id)
			for _, subCluster := range subClusters {
				if subCluster.delta == 1 {
					subDeltaOne = true
				}
			}

			if subDeltaOne {
				cluster.delta = 0

				// set parents to 0
				parents := c.Clusters.allParents(cluster)
				for _, parent := range parents {
					parent.delta = 0
				}
			}

			if cluster.parent != nil {
				// if parent already calculated to be better
				if c.Clusters.getClusterByID(*cluster.parent).delta == 1 {
					cluster.delta = 0
				}
			}
		}
	}
}

func (c clusters) root() *cluster {
	for _, cluster := range c {
		if cluster.parent == nil {
			return cluster
		}
	}

	return nil
}

func (c clusters) leaves() clusters {
	var leaves clusters
	for _, cluster := range c {
		if len(cluster.children) == 0 {
			leaves = append(leaves, cluster)
		}
	}
	return leaves
}

func (c clusters) allParents(clstr *cluster) clusters {
	var parents clusters

	if clstr.parent != nil {
		parentCluster := c.getClusterByID(*clstr.parent)
		allParents := c.allParents(parentCluster)
		parents = append(parents, allParents...)
	}

	return parents
}

func (c clusters) subTree(id int) clusters {
	var subTree clusters
	for _, cluster := range c {
		if cluster.parent != nil {
			if *cluster.parent == id {
				subTree = append(subTree, cluster)
				childTree := c.subTree(cluster.id)
				subTree = append(subTree, childTree...)
			}
		}
	}

	return subTree
}
