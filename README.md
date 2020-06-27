# HDBSCAN - Density Clustering Algorithm

HDBSCAN algorithm implementation in golang.

Written to run concurrently on CPU (uses all CPU cores by default).

A re-write of code started by the brilliant developer Edouard Belval at https://github.com/Belval/hdbscan  ... although it has been changed quite a lot from the original.

## Download

`go get -u github.com/humilityai/hdbscan`

## Use

```go
import(
    "github.com/humilityai/hdbscan"
)

func main() {
    data := [][]float64{
        []float64{1,2,3},
        []float64{3,2,1},
    }
    minimumClusterSize := len(data)
    
    clustering := hdbscan.NewClustering(data, minimumClusterSize)

    // options
    clustering = clustering.Verbose().Voronoi()

    // the final boolean argument is: minimum-spanning-tree argument
    clustering.Run(hdbscan.EuclideanDistance, hdbscan.VarianceScore, true)
}
```

### options

- `Verbose()` will log the progress of the clustering to stdout (should be called before calling other options or else some progress logs may not be printed)
- `Voronoi()` will add all points not placed in a cluster in the final clustering to their nearest cluster.
- `Subsample()` will only cluster the first `n` data points. Note: only the data points from the subsample will be in the final clustering results.
