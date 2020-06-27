package hdbscan

import (
	"log"
	"math/rand"
)

// randomSampling ...  (DO NOT USE!)
// will process a random sample of the total provided data.
// The amount of data in the sample is based on the percentage (0-1) argument provided.
// Randomly sampled data will also automatically perform a Voronoi clustering for
// all data points by using the centroids of the clusters generated from the sample.
// In other words, the final result may not be as of high quality as on a different sample
// or the entire dataset.
func (c *Clustering) randomSampling(percentage float64) *Clustering {
	if c.verbose {
		log.Println("random sampling data")
	}

	// if percentage is not valid: process total data
	if percentage < 0 || percentage > 1 {
		percentage = 1
	}

	bound := int(percentage * 100)

	var newData [][]float64
	for _, v := range c.data {
		if rand.Intn(100) < bound {
			newData = append(newData, v)
		}
	}

	c.data = newData
	c.randomSample = true
	c.voronoi = true

	if c.verbose {
		log.Println("finished random sampling data")
	}

	return c
}

// Subsampling will take the first 'n' data points and perform clustering on
// those. 'n' is a provided argument and should be between 0 and the total data size.
// Voronoi clustering will be performed after the clusters have been found for all points
// that are not in the subsample.
func (c *Clustering) Subsampling(n int) *Clustering {
	if c.verbose {
		log.Println("sub-sampling data")
	}

	if n < 0 || n > len(c.data) {
		n = len(c.data)
	}

	var newData [][]float64
	for i, v := range c.data {
		if i < n {
			newData = append(newData, v)
		}
	}

	c.data = newData
	c.subSample = true
	c.voronoi = true

	if c.verbose {
		log.Println("finished sub-sampling data")
	}

	return c
}
