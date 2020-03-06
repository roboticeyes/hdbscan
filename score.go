package hdbscan

import (
	"gonum.org/v1/gonum/stat/distuv"
)

func (c *Clustering) scoreClusters(clusterHierarchy *cluster, optimization string) error {
	var err error
	switch optimization {
	case VarianceScore:
		err = c.varianceScores(clusterHierarchy)
	default:
		err = c.stabilityScores(clusterHierarchy)
	}

	return err
}

func (c *Clustering) stabilityScores(hierarchy *cluster) error {
	// TODO: implement
	return nil
}

// optimizing for largest cluster sizes with smallest variances
func (c *Clustering) varianceScores(hierarchy *cluster) error {
	// normalized cluster sizes
	setSizes(hierarchy)
	sizes := getSizes(hierarchy, []float64{})
	sizeDistro := distuv.Normal{}
	sizeDistro.Fit(sizes, nil)
	normalizeSizes(hierarchy, &sizeDistro)

	// normalized cluster general-variances
	setVariances(hierarchy, c.data)
	variances := getVariances(hierarchy, []float64{})
	varianceDistro := distuv.Normal{}
	varianceDistro.Fit(variances, nil)
	normalizeVariances(hierarchy, &varianceDistro)

	setVarianceScore(hierarchy)

	return nil
}

func setVarianceScore(c *cluster) {
	if c.variance == 0 {
		c.variance = 1
	}
	c.score = c.size / c.variance

	for _, childCluster := range c.children {
		setVarianceScore(childCluster)
	}
}

func setSizes(c *cluster) {
	// leaf
	if len(c.points) > 0 {
		c.size = float64(len(c.points))
		return
	}

	// parent
	var size float64
	for _, childCluster := range c.children {
		setSizes(childCluster)
		size += childCluster.size
	}

	c.size = size
}

func normalizeSizes(c *cluster, nd *distuv.Normal) {
	c.size = nd.CDF(c.size)

	for _, childCluster := range c.children {
		normalizeSizes(childCluster, nd)
	}
}

func getSizes(c *cluster, sizes []float64) []float64 {
	sizes = append(sizes, c.size)

	for _, childCluster := range c.children {
		sizes = append(sizes, getSizes(childCluster, sizes)...)
	}

	return sizes
}

func setVariances(hierarchy *cluster, data [][]float64) {
	// children
	for _, childCluster := range hierarchy.children {
		setVariances(childCluster, data)
	}

	// this
	hierarchy.calculateVariance(data)
}

func normalizeVariances(c *cluster, nd *distuv.Normal) {
	c.score = nd.CDF(c.score)

	for _, childCluster := range c.children {
		normalizeSizes(childCluster, nd)
	}
}

func getVariances(c *cluster, variances []float64) []float64 {
	variances = append(variances, c.score)

	for _, childCluster := range c.children {
		variances = append(variances, getVariances(childCluster, variances)...)
	}

	return variances
}

func (c *cluster) calculateVariance(data [][]float64) {
	pointIndices := c.pointIndexes()

	var clusterData [][]float64
	for _, pointIndex := range pointIndices {
		clusterData = append(clusterData, data[pointIndex])
	}

	if len(clusterData) > 0 {
		c.score = GeneralizedVariance(len(clusterData), len(clusterData[0]), unfold(clusterData))
	}
}

func (c *cluster) pointIndexes() []int {
	if len(c.points) > 0 {
		return c.points
	}

	var points []int
	for _, childCluster := range c.children {
		childPoints := childCluster.pointIndexes()
		points = append(points, childPoints...)
	}

	return points
}

func (c *cluster) calculateStability(mrg *graph) float64 {
	if len(c.points) > 0 {
		// var sum float64
		// for _, pIndex := range c.points {

		// }
		// calculate sum of points (1 / e_min) - (1 / e_max)
		// e_min = points mrg
		return c.score
	}

	var stability float64
	for _, childCluster := range c.children {
		childStability := childCluster.calculateStability(mrg)
		stability += childStability
	}

	return stability
}

func potentialStability(c *cluster) float64 {
	// if leaf node: return stability
	// else: return max(stability, sum-of-children-stabilities)
	return 0
}

// func listClusters(c *cluster, list []*cluster) []*cluster {
// 	list = append(list, c)

// 	for _, childCluster := range c.children {
// 		list = append(list, listClusters(childCluster, list)...)
// 	}

// 	return list
// }
