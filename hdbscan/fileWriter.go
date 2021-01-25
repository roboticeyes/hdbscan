package hdbscan

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/edgeDetection/edgeDetection"
)

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

func (c *Clustering) writeClusterToObj() {
	colors, _ := getcolors(len(c.Clusters))
	for i, cl := range c.Clusters {
		outputfile, _ := os.Create(c.directory + "cluster_" + fmt.Sprint(i) + "_.obj")
		defer outputfile.Close()
		writer := bufio.NewWriter(outputfile)

		for _, p := range cl.Points {

			x := fmt.Sprintf("%f", c.data[p][0])
			y := fmt.Sprintf("%f", c.data[p][1])
			z := fmt.Sprintf("%f", c.data[p][2])

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

			x := fmt.Sprintf("%f", c.data[p.Index][0])
			y := fmt.Sprintf("%f", c.data[p.Index][1])
			z := fmt.Sprintf("%f", c.data[p.Index][2])

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

func (c *Clustering) writeClusterToFile(position string) {
	outputfile, _ := os.Create(c.directory + position + "cluster.txt")
	defer outputfile.Close()
	writer := bufio.NewWriter(outputfile)
	for _, c := range c.Clusters {
		// fmt.Println(i)
		id := c.id
		parent := c.parent
		child := c.children
		numPoints := len(c.Points)

		if parent == nil {
			_, err := writer.WriteString("id: " + fmt.Sprint(id) + " " + "parent: " + fmt.Sprint(9999) + " " + "children: " + fmt.Sprint(child) + " " + "numP: " + fmt.Sprint(numPoints) + " " + "stability: " + fmt.Sprint(c.score) + "\n")
			// _, err = writer.WriteString("Points: " + fmt.Sprint(c.Points) + "\n")
			if err != nil {
				panic(err)
			}
			continue
		}

		_, err := writer.WriteString("id: " + fmt.Sprint(id) + " " + "parent: " + fmt.Sprint(*parent) + " " + "children: " + fmt.Sprint(child) + " " + "numP: " + fmt.Sprint(numPoints) + " " + "stability: " + fmt.Sprint(c.score) + " " + fmt.Sprint(c.delta) + "\n")
		// _, err = writer.WriteString("Points: " + fmt.Sprint(c.Points) + "\n")
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}
