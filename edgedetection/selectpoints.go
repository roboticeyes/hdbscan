package edgedetection

import (
	"fmt"
	"sort"

	"gocv.io/x/gocv"
)

func (d *Data) correspondendingPoints(modus string) *ImageCV {

	minHeight := d.getminheight()
	fmt.Println("min height: ", minHeight)
	upperZThreshold := 0.5 //20cm
	lowerZThreshold := 0.02

	var whitePoints gocv.Mat
	switch string(modus) {
	case "barycenter":
		whitePoints = d.findNormalsAndBarycenter(upperZThreshold, lowerZThreshold, minHeight)
	case "rawpoints":
		whitePoints = d.findRawPoints(upperZThreshold, lowerZThreshold, minHeight)
	}

	return &ImageCV{mat: whitePoints}
}

func (d *Data) findNormalsAndBarycenter(upperZThreshold, lowerZThreshold, minHeight float64) gocv.Mat {

	whitePoints := gocv.NewMatWithSize(d.img.mat.Rows(), d.img.mat.Cols(), gocv.MatTypeCV8U)
	wp := 0
	for i := 0; i < len(d.indexUV); i += 2 {

		pixel1, row1, col1 := d.checkCurrentPixel(d.indexUV[i][0])
		pixel2, row2, col2 := d.checkCurrentPixel(d.indexUV[i][1])
		pixel3, row3, col3 := d.checkCurrentPixel(d.indexUV[i][2])

		if pixel1 || pixel2 || pixel3 {
			check := d.calcBarycenterFacenormal(d.indexXYZ[i], upperZThreshold, lowerZThreshold, minHeight)
			if check {

				whitePoints.SetUCharAt(row1, col1, uint8(255))
				whitePoints.SetUCharAt(row2, col2, uint8(255))
				whitePoints.SetUCharAt(row3, col3, uint8(255))
				wp++
			}
		}
	}
	fmt.Println("Number of white points: ", wp)
	return whitePoints
}

func (d *Data) checkCurrentPixel(iuv int) (bool, int, int) {

	u := d.coordUV[iuv].X()
	v := d.coordUV[iuv].Y()

	row := int((1-v)*float64(d.img.height) - 1)
	col := int(u * float64(d.img.width))

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

func (d *Data) findRawPoints(upperZThreshold, lowerZThreshold, minHeight float64) gocv.Mat {

	whitePoints := gocv.NewMatWithSize(d.img.mat.Rows(), d.img.mat.Cols(), gocv.MatTypeCV8U)
	wp := 0

	for i := 0; i < len(d.indexUV); i += 3 {
		//

		iXYZ := d.indexXYZ[i]
		for ii := 0; ii < len(iXYZ); ii++ {

			pixel, row, col := d.checkCurrentPixel(d.indexUV[i][ii])
			vertex := d.coordXYZ[iXYZ[ii]]

			if pixel && (vertex.Z() >= minHeight+lowerZThreshold) && (vertex.Z() <= minHeight+upperZThreshold) {
				wp++
				whitePoints.SetUCharAt(row, col, uint8(255))
				d.Points = append(d.Points, []float64{vertex.X(), vertex.Y(), vertex.Z()})
			}
		}
	}

	fmt.Println("Number of white points: ", wp)
	return whitePoints
}

func (d *Data) calcBarycenterFacenormal(indexXYZ [3]int, upperZThreshold, lowerZThreshold, minHeight float64) bool {

	vertex1 := d.coordXYZ[indexXYZ[0]]
	vertex2 := d.coordXYZ[indexXYZ[1]]
	vertex3 := d.coordXYZ[indexXYZ[2]]

	vec21 := vertex2.Sub(vertex1)
	vec31 := vertex3.Sub(vertex1)

	barycenter := vertex1.Add(vertex2).Add(vertex3).Mul(1. / 3.)
	if (barycenter.Z() <= minHeight+lowerZThreshold) || (barycenter.Z() >= minHeight+upperZThreshold) {
		return false
	}

	cross := vec21.Cross(vec31)
	norm := cross.Normalize()

	d.Normale = append(d.Normale, []float64{norm.X(), norm.Y(), norm.Z()})
	d.Barycenter = append(d.Barycenter, []float64{barycenter.X(), barycenter.Y(), barycenter.Z()})

	return true
}

func (d *Data) getminheight() float64 {

	height := make([]float64, len(d.coordXYZ))
	for i := 0; i < len(height); i++ {
		height[i] = d.coordXYZ[i].Z()
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
