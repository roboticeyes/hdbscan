package main

import (
	"log"
	"os"

	"github.com/edgeDetection/edgeDetection"
	"github.com/edgeDetection/hdbscan"
)

func main() {

	if len(os.Args) > 1 {

		argument := os.Args[1]
		log.Println(argument)
		points := edgeDetection.Detection(argument)
		log.Println("Number of points to cluster: ", len(points))

		// hdbscan
		minimumClusterSize := 100
		minimumSpanningTree := true
		clustering, err := hdbscan.NewClustering(points, minimumClusterSize, argument)
		if err != nil {
			panic(err)
		}
		// Set options for clustering
		clustering = clustering.Verbose().OutlierDetection().NearestNeighbor()
		clustering.Run(hdbscan.EuclideanDistance, hdbscan.StabilityScore, minimumSpanningTree)

	} else {
		panic("No file founded!")
	}
}
