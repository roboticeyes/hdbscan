package hdbscan

import (
	"errors"
	"log"
	"math"
	"math/rand"
)

func (c *Clustering) sample() {
	if c.randomSample {
		if c.verbose {
			log.Println("randomly sampling data")
		}

		var newData [][]float64
		for _, v := range c.data {
			if rand.Intn(100) < c.sampleBound {
				newData = append(newData, v)
			}
		}

		c.data = newData
	} else if c.subSample {
		if c.verbose {
			log.Println("sub-sampling data")
		}

		c.data = c.data[:c.sampleBound]
	}

	if c.verbose {
		log.Println("finished sampling data")
	}
}

// Assign will assign a list of data points to an existing cluster.
// If the original clustering had OutlierDetection option enabled
// then it will perform outlier detection based on existing outliers.
// The results are returned as a new clustering object with only the
// indexes from the supplied data. All clusters returned have the same ID
// as they had in the original clustering.
// This method can be useful if a sampling was used for the initial clustering
// and the data points outside of the sample need to be assigned to a cluster
// as well.
func (c *Clustering) Assign(data [][]float64) (*Clustering, error) {
	if c.verbose {
		log.Println("assigning data")
	}

	newClustering, err := NewClustering(data, c.mcs)
	if err != nil {
		return newClustering, err
	}

	if len(c.Clusters) == 0 {
		return newClustering, errors.New("no clusters")
	}

	if !c.od {
		c.distanceDistributions()
	}

	// create new clusters
	for _, clust := range c.Clusters {
		newCluster := &cluster{
			id:       clust.id,
			Points:   make([]int, 0),
			Outliers: make(Outliers, 0),
		}

		newClustering.Clusters = append(newClustering.Clusters, newCluster)
	}

	// assign data
	for i, v := range data {
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

		if c.od {
			if len(c.Clusters[nearestClusterIndex].Outliers) == 0 {
				if minDistance > c.Clusters[nearestClusterIndex].largestDistance {
					prob := c.Clusters[nearestClusterIndex].distanceDistribution.CDF(minDistance)
					newOutlier := Outlier{
						Index:              i,
						NormalizedDistance: prob,
					}
					newClustering.Clusters[nearestClusterIndex].Outliers = append(newClustering.Clusters[nearestClusterIndex].Outliers, newOutlier)
				} else {
					newClustering.Clusters[nearestClusterIndex].Points = append(newClustering.Clusters[nearestClusterIndex].Points, i)
				}
			}

			minOutlier := c.Clusters[nearestClusterIndex].Outliers.MinProb()
			prob := c.Clusters[nearestClusterIndex].distanceDistribution.CDF(minDistance)
			if prob > minOutlier.NormalizedDistance {
				newOutlier := Outlier{
					Index:              i,
					NormalizedDistance: prob,
				}
				newClustering.Clusters[nearestClusterIndex].Outliers = append(newClustering.Clusters[nearestClusterIndex].Outliers, newOutlier)
			}

			if c.voronoi {
				newClustering.Clusters[nearestClusterIndex].Points = append(newClustering.Clusters[nearestClusterIndex].Points, i)
			}
		} else {
			newClustering.Clusters[nearestClusterIndex].Points = append(newClustering.Clusters[nearestClusterIndex].Points, i)
		}
	}

	// outlier-clustering
	newClustering.outlierClustering()

	if c.verbose {
		log.Println("finished assigning data")
	}

	return newClustering, nil
}
