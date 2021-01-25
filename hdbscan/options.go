// Copyright 2020 Humility AI Incorporated, All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hdbscan

// OutlierDetection will track all unassigned
// points as outliers of their nearest cluster.
// It provides a `NormalizedDistance` value for
// each outlier which can be interpreted as the
// probability of the point being an outlier
// (relative to all other outliers).
func (c *Clustering) OutlierDetection() *Clustering {
	c.od = true
	return c
}

// Verbose will set verbosity to true for clustering process
// and the internals of a clustering run will be logged to stdout.
func (c *Clustering) Verbose() *Clustering {
	c.verbose = true
	return c
}

// Voronoi will set voronoi-clustering to true, and
// after density clustering is performed,
// all points not assigned to a cluster will be placed
// into their nearest cluster (by centroid distance).
func (c *Clustering) Voronoi() *Clustering {
	c.voronoi = true
	return c
}

// NearestNeighbor specifies if nearest-neighbor
// distances should be used for outlier detection
// and for voronoi clustering instead of centroid-based
// distances.
// NearestNeighbor will find the closest assigned data
// point to an unassigned data point and consider the
// unassigned data point to be of that same cluster (as an outlier and/or a point).
func (c *Clustering) NearestNeighbor() *Clustering {
	c.nn = true
	return c
}

// OutlierClustering is an option to group the outliers of a cluster
// into a new cluster if there are at least a minimum-cluster-size
// number of them.
// This option will automatically perform outlier detection on the clustering
// as well.
func (c *Clustering) OutlierClustering() *Clustering {
	c.od = true
	c.oc = true
	return c
}
