# HDBSCAN - Density Clustering Algorithm

HDBSCAN algorithm implementation in golang.

## Download

`go get -u github.com/humilityai/hdbscan`

## Use

A re-write of code started by the brilliant developer Edouard Belval at https://github.com/Belval/hdbscan

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

    // the final boolean argument is: minimum-spanning-tree argument
    clustering.Run(hdbscan.EuclideanDistance, hdbscan.VarianceScore, true)
}
```