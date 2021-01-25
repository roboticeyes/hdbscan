package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

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

		detections := edgeDetection.Detection(meshReader, depthReader, imgReader)
		log.Println("Number of points to cluster: ", len(detections.Normale))

		// hdbscan
		minimumClusterSize := 20
		minimumSpanningTree := true

		clustering, err := hdbscan.NewClustering(detections.Normale, minimumClusterSize, argument)
		if err != nil {
			panic(err)
		}

		// Set options for clustering
		clustering = clustering.Verbose().OutlierDetection().NearestNeighbor()
		clustering.Run(hdbscan.AngleVector, hdbscan.VarianceScore, minimumSpanningTree)

		writeClusterToObj(clustering, detections, argument)
	} else {
		panic("No file founded!")
	}
}

func writeClusterToObj(c *hdbscan.Clustering, d *edgeDetection.Data, argument string) {
	colors, _ := getcolors(len(c.Clusters))
	for i, cl := range c.Clusters {
		outputfile, _ := os.Create(argument + "cluster_" + fmt.Sprint(i) + "_.obj")
		defer outputfile.Close()
		writer := bufio.NewWriter(outputfile)

		for _, p := range cl.Points {

			x := fmt.Sprintf("%f", d.Barycenter[p][0])
			y := fmt.Sprintf("%f", d.Barycenter[p][1])
			z := fmt.Sprintf("%f", d.Barycenter[p][2])

			c := colors[i]
			R := fmt.Sprintf("%1.3f", c.R)
			G := fmt.Sprintf("%1.3f", c.G)
			B := fmt.Sprintf("%1.3f", c.B)
			_, err := writer.WriteString("v" + " " + x + " " + y + " " + z + " " + R + " " + G + " " + B + "\n")
			if err != nil {
				panic(err)
			}

		}
		for _, p := range cl.Outliers {

			x := fmt.Sprintf("%f", d.Barycenter[p.Index][0])
			y := fmt.Sprintf("%f", d.Barycenter[p.Index][1])
			z := fmt.Sprintf("%f", d.Barycenter[p.Index][2])

			// c := colors[i]
			R := fmt.Sprintf("%1.3f", 0.502)
			G := fmt.Sprintf("%1.3f", 0.502)
			B := fmt.Sprintf("%1.3f", 0.502)
			_, err := writer.WriteString("v" + " " + x + " " + y + " " + z + " " + R + " " + G + " " + B + "\n")
			if err != nil {
				panic(err)
			}

		}
		// Very important to invoke after writing a large number of lines
		writer.Flush()
	}
}

func getcolors(k int) ([]edgeDetection.Color, []string) {
	time := time.Now().Nanosecond()
	rand.Seed(int64(time))

	colors := make([]edgeDetection.Color, k)
	colorname := make([]string, k)
	for i := 0; i < k; i++ {

		randnum := rand.Intn(len(edgeDetection.Colorsrand)) // colorsrand[randnum]
		color := edgeDetection.Colorsrand[randnum]

		colors[i], _ = edgeDetection.IsColorName(color)
		colorname[i] = color
	}
	return colors, colorname
}
