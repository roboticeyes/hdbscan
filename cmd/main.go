package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/edgeDetection/edgeDetection"
	"github.com/edgeDetection/hdbscan"
	"github.com/fatih/color"
)

const (
	Filenames = "/home/felix/workspace_go/hdbscan/filenames.json"
)

type Files struct {
	Mesh  string `json:"objfile"`
	Depth string `json:"depthfile"`
	Img   string `json:"texture"`
}

func main() {

	if len(os.Args) > 1 {

		argument := os.Args[1]
		fmt.Println(argument)
		jsonReader, err := os.Open(Filenames)
		if err != nil {
			color.Red("Cannot read json file:", err)
		}
		defer jsonReader.Close()

		jsonByte, err := ioutil.ReadAll(jsonReader)
		if err != nil {
			color.Red("Cannot load json file:", err)
		}

		var files Files
		err = json.Unmarshal([]byte(jsonByte), &files)
		if err != nil {
			color.Red("Cannot unmarshal json file:", err)
		}

		meshReader, _ := os.Open(argument + files.Mesh)
		depthReader, _ := os.Open(argument + files.Depth)
		imgReader, _ := os.Open(argument + files.Img)

		defer meshReader.Close()
		defer depthReader.Close()
		defer imgReader.Close()

		points := edgeDetection.Detection(meshReader, depthReader, imgReader)
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
