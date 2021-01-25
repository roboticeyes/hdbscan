package hdbscan

import (
	"log"
	"math"
)

func (c *Clustering) outliersAndVoronoi() {
	if !c.od && !c.voronoi {
		return
	}

	if len(c.Clusters) == 0 {
		return
	}

	if c.verbose {
		if c.od {
			log.Println("finding outliers")
		}

		if c.voronoi {
			log.Println("starting voronoi clustering")
		}
	}

	for i, v := range c.data {
		var exists bool
		for _, cluster := range c.Clusters {
			for _, point := range cluster.Points {
				if point == i {
					exists = true
					break
				}
			}

			if exists {
				break
			}
		}

		if !exists {
			// calculate nearest cluster
			minDistance := math.MaxFloat64
			var nearestClusterIndex int
			for i, cluster := range c.Clusters {
				if c.nn {
					for _, p := range cluster.Points {
						distance := c.distanceFunc(c.data[p], v)
						if distance < minDistance {
							minDistance = distance
							nearestClusterIndex = i
						}
					}
				} else {
					distance := c.distanceFunc(cluster.Centroid, v)
					if distance < minDistance {
						minDistance = distance
						nearestClusterIndex = i
					}
				}
			}

			// voronoi cluster
			if c.voronoi {
				c.Clusters[nearestClusterIndex].Points = append(c.Clusters[nearestClusterIndex].Points, i)
			}

			// outlier detection
			if c.od {
				c.Clusters[nearestClusterIndex].Outliers = append(c.Clusters[nearestClusterIndex].Outliers, Outlier{
					Index:              i,
					NormalizedDistance: minDistance,
				})
			}
		}
	}

	// normalize outlier distances
	if c.od {
		c.distanceDistributions()

		for _, cluster := range c.Clusters {
			for j, outlier := range cluster.Outliers {
				outlier.NormalizedDistance = isNum(cluster.distanceDistribution.CDF(outlier.NormalizedDistance))
				cluster.Outliers[j] = outlier
			}
		}
	}

	if c.verbose {
		if c.od {
			log.Println("finished finding outliers")
		}

		if c.voronoi {
			log.Println("finished voronoi clustering")
		}
	}
}

func (c *Clustering) outlierClustering() {
	if !c.oc {
		return
	}

	maxID := c.Clusters.maxID()
	var newClusters clusters
	for i, clust := range c.Clusters {
		if len(clust.Outliers) >= c.mcs {
			newCluster := &cluster{
				id:     i + maxID + 1,
				Points: make([]int, 0),
			}

			for _, o := range clust.Outliers {
				newCluster.Points = append(newCluster.Points, o.Index)
			}

			newClusters = append(newClusters, newCluster)
		}
	}

	c.Clusters = append(c.Clusters, newClusters...)
}
