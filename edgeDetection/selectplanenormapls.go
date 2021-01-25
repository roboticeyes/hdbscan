package edgeDetection

import (
	"fmt"

	"gocv.io/x/gocv"
)

func (d *Data) correspondendingPlaneNormals() *ImageCV {

	minHeight := d.getminheight()
	fmt.Println("min height: ", minHeight)
	upperZThreshold := 0.5 //20cm
	lowerZThreshold := 0.02
	whitePoints := d.findNormalsAndBarycenter(upperZThreshold, lowerZThreshold, minHeight)

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

func (d *Data) calcBarycenterFacenormal(indexXYZ [3]int, upperZThreshold, lowerZThreshold, minHeight float64) bool {

	vertex1 := d.coordXYZ[indexXYZ[0]]
	vertex2 := d.coordXYZ[indexXYZ[1]]
	vertex3 := d.coordXYZ[indexXYZ[2]]

	vec21 := vertex2.Sub3(vertex1)
	vec31 := vertex3.Sub3(vertex1)

	barycenter := vertex1.Add3(vertex2).Add3(vertex3).MultiplyByScalar3(1. / 3.)
	if (barycenter.Z <= minHeight+lowerZThreshold) || (barycenter.Z >= minHeight+upperZThreshold) {
		return false
	}

	cross := vec21.Cross3(vec31)
	norm := cross.Normalize3()

	d.Normale = append(d.Normale, []float64{norm.X, norm.Y, norm.Z})
	d.Barycenter = append(d.Barycenter, []float64{barycenter.X, barycenter.Y, barycenter.Z})

	return true
}
