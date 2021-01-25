# HDBSCAN - Density Clustering Algorithm

HDBSCAN algorithm implementation in golang.

Written to run concurrently on CPU (uses all CPU cores by default).

This repository uses the great hdbscan algorithm from Humility AI (https://github.com/humilityai/hdbscan.git) and has been extended with some features.
Further description follows!

### options for hdbscan clustering

- `Verbose()` will log the progress of the clustering to stdout.
- `Voronoi()` will add all points not placed in a cluster in the final clustering to their nearest cluster. All unassigned data points outliers will be added to their nearest cluster.
- `OutlierDetection()` will mark all unassigned data points as outliers of their nearest cluster and provide a `NormalizedDistance` value for each outlier that can be interpreted as the probability that the data point is an outlier of that cluster.
- `NearestNeighbor()` specifies if an unassigned points "nearness" to a cluster should be based on it's nearest assigned neighboring data point in that cluster (default "nearness" is based on distance to centroid of cluster).
- `Subsample(n int)` specifies to only use the first `n` data points in the clustering process. This speeds up the clustering. The remaining data points can be added to clusters using the `Assign(data [][]float64)` method after a successful clustering.
- `OutlierClustering()` will create a new cluster for the outliers of an existing cluster if the number of outliers is equal to or greater than the specified minimum-cluster-size.

<!-- TODO: random sampling option -->
