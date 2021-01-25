package edgeDetection

import (
	"fmt"
	"log"
	"sort"

	"gocv.io/x/gocv"
)

func (d *Data) findCorrespondendingPoints() *ImageCV {

	minHeight := d.getminheight()
	log.Println("min height: ", minHeight)
	upperZThreshold := 0.5
	lowerZThreshold := 0.02
	whitePoints := d.findPoints(upperZThreshold, lowerZThreshold, minHeight)

	return &ImageCV{mat: whitePoints}
}

func (d *Data) findPoints(upperZThreshold, lowerZThreshold, minHeight float64) gocv.Mat {

	whitePoints := gocv.NewMatWithSize(d.img.mat.Rows(), d.img.mat.Cols(), gocv.MatTypeCV8U)
	wp := 0

	for i := 0; i < len(d.indexUV); i += 3 {
		//

		iXYZ := d.indexXYZ[i]
		for ii := 0; ii < len(iXYZ); ii++ {

			pixel, row, col := d.checkCurrentPixel(d.indexUV[i][ii])
			vertex := d.coordXYZ[iXYZ[ii]]

			if pixel && (vertex.Z >= minHeight+lowerZThreshold) && (vertex.Z <= minHeight+upperZThreshold) {
				wp++
				whitePoints.SetUCharAt(row, col, uint8(255))
				d.points = append(d.points, []float64{vertex.X, vertex.Y, vertex.Z})
			}
		}
	}

	fmt.Println("Number of white points: ", wp)
	return whitePoints
}

func (d *Data) checkCurrentPixel(iuv int) (bool, int, int) {

	u := d.coordUV[iuv].X
	v := d.coordUV[iuv].Y

	row := int((1-v)*float64(d.img.height) - 1) // Rows
	col := int(u * float64(d.img.width))        // Columns

	pixValue := d.img.mat.GetUCharAt(row, col)
	if pixValue != uint8(0) {
		for _, r := range d.img.rects {
			if (row >= r.Min.Y) && (row <= r.Max.Y) && (col >= r.Min.X) && (col <= r.Max.X) {
				return true, row, col
			}
		}
	} else {
		return false, 0, 0
	}
	return false, 0, 0
}

func (d *Data) getminheight() float64 {

	height := make([]float64, len(d.coordXYZ))
	for i := 0; i < len(height); i++ {
		height[i] = d.coordXYZ[i].Z
	}
	return average(height, 1000)
}

func average(height []float64, count int) float64 {
	sort.Float64s(height)
	sum := 0.0
	for i := 0; i < count; i++ {
		sum += height[i]
	}
	sum = sum / float64(count)
	return sum
}
