package edgeDetection

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func ReadObjFile(filename string) ([]Vector3, []Vector2, [][3]int, [][3]int) {

	reader := open(filename)

	bytestring, err := ioutil.ReadAll(reader)
	if err != nil {
		color.Red("Cannot read input file:", err)
	}

	var tempV3 Vector3
	var verticesXYZ []Vector3

	var tempV2 Vector2
	var verticesUV []Vector2

	var temp3xyz [3]int
	var temp3uv [3]int
	var uvs [][3]int
	var xyzs [][3]int
	line := bytes.Split(bytestring, []byte("\n"))
	for _, l := range line {
		element := bytes.Split(l, []byte(" "))

		switch string(element[0]) {
		case "#":
			continue
		case "v":
			tempV3.X = parseToFloat(element[1])
			tempV3.Y = parseToFloat(element[2])
			tempV3.Z = parseToFloat(element[3])
			verticesXYZ = append(verticesXYZ, tempV3)
		case "vt":
			tempV2.X = parseToFloat(element[1])
			tempV2.Y = parseToFloat(element[2])
			verticesUV = append(verticesUV, tempV2)
		case "f":
			first := bytes.Split(element[1], []byte("/"))
			temp3xyz[0] = parseToInt(first[0]) - 1
			temp3uv[0] = parseToInt(first[1]) - 1

			second := bytes.Split(element[2], []byte("/"))
			temp3xyz[1] = parseToInt(second[0]) - 1
			temp3uv[1] = parseToInt(second[1]) - 1

			third := bytes.Split(element[3], []byte("/"))
			temp3xyz[2] = parseToInt(third[0]) - 1
			temp3uv[2] = parseToInt(third[1]) - 1

			xyzs = append(xyzs, temp3xyz)
			uvs = append(uvs, temp3uv)
		}
	}
	return verticesXYZ, verticesUV, xyzs, uvs
}

func parseToFloat(num []byte) float64 {
	flotnum, err := strconv.ParseFloat(string(num), 64)
	if err != nil {
		log.Fatal(err)
	}
	return flotnum
}

func parseToInt(num []byte) int {
	intnum, err := strconv.Atoi(string(num))
	if err != nil {
		log.Fatal("Error convert string to int: ", err)
	}
	return intnum
}

func ReadDepthFile(filename string) []float64 {

	reader := open(filename)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	var result []float64
	for scanner.Scan() {
		x, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			color.Red("Cannot read input file:", scanner.Err())
		}
		result = append(result, x)
	}
	return result
}

func open(filename string) *os.File {

	reader, err := os.Open(filename) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return reader
}

func (d *Data) writeToobjFile(filename string, k int) {

	colors, colorname := d.getcolors(k)
	outputfile, _ := os.Create(filename + "/clustered.obj")
	defer outputfile.Close()
	writer := bufio.NewWriter(outputfile)

	for i, o := range d.clusterNumber {

		bsx := fmt.Sprintf("%f", d.points[i][0])
		bsy := fmt.Sprintf("%f", d.points[i][1])
		bsz := fmt.Sprintf("%f", d.points[i][2])

		c := colors[o]
		R := fmt.Sprintf("%1.3f", c.R)
		G := fmt.Sprintf("%1.3f", c.G)
		B := fmt.Sprintf("%1.3f", c.B)
		_, err := writer.WriteString("v" + " " + bsx + " " + bsy + " " + bsz + " " + R + " " + G + " " + B + "\n")
		if err != nil {
			panic(err)
		}

	}
	// Very important to invoke after writing a large number of lines
	writer.Flush()

	for i, o := range d.mean {
		fmt.Println("Cluster: ", i, "Mean: ", o, "Color: ", colorname[i])

	}

}

func (d *Data) getcolors(k int) ([]Color, []string) {
	time := time.Now().Nanosecond()
	rand.Seed(int64(time))

	colors := make([]Color, k)
	colorname := make([]string, k)
	for i := 0; i < k; i++ {

		randnum := rand.Intn(len(Colorsrand)) // colorsrand[randnum]
		color := Colorsrand[randnum]

		colors[i], _ = IsColorName(color)
		colorname[i] = color
	}
	return colors, colorname
}
