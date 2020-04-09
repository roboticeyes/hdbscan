package hdbscan

func (c *Clustering) scoreClusters(optimization string) {
	switch optimization {
	case VarianceScore:
		c.varianceScores()
	default:
		c.stabilityScores()
	}
}

func (c *Clustering) varianceScores() {
	c.setNormalizedSizes()
	c.setNormalizedVariances()
	c.Clusters.setVarianceScores()
}

func (c clusters) setVarianceScores() {
	for _, cluster := range c {
		cluster.score = cluster.size / cluster.variance
	}
}

func (c *Clustering) setNormalizedSizes() {
	// distro
	var sizes []float64
	for _, cluster := range c.Clusters {
		size := float64(len(cluster.Points))
		sizes = append(sizes, size)
		cluster.size = size
	}
}

func (c *Clustering) setNormalizedVariances() {
	// variances
	var variances []float64
	for _, cluster := range c.Clusters {
		// data
		var clusterData [][]float64
		for _, pointIndex := range cluster.Points {
			clusterData = append(clusterData, c.data[pointIndex])
		}

		variance := GeneralizedVariance(len(cluster.Points), len(clusterData[0]), unfold(clusterData))
		cluster.variance = isNum(variance)
		variances = append(variances, cluster.variance)
	}
}

func (c *Clustering) stabilityScores() {
	// TODO: implement
}

func (c *cluster) calculateStability() float64 {
	if len(c.Points) > 0 {
		// var sum float64
		// for _, pIndex := range c.points {

		// }
		// calculate sum of points (1 / e_min) - (1 / e_max)
		// e_min = points mrg
		return c.score
	}

	var stability float64
	// for _, childCluster := range c.children {
	// 	childStability := childCluster.calculateStability(mrg)
	// 	stability += childStability
	// }

	return stability
}

func potentialStability(c *cluster) float64 {
	// if leaf node: return stability
	// else: return max(stability, sum-of-children-stabilities)
	return 0
}
