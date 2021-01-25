package edgeDetection

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

type ImageCV struct {
	mat gocv.Mat

	rects  []image.Rectangle
	height int
	width  int
	center []int //[0]width [1]hight
}

func ImageControler(filename string) *ImageCV {

	// Read Original image
	org := readImg(filename)
	// Blur image (sigmaX, sigmaY, kernel size)
	blured := org.gauSSianBlur(0, 0, 7)
	// Auto canny
	edge := blured.canny(0.5)
	// Dilatation
	dilate := edge.dilatation(7)
	// Find contours for dilate
	dilate.findContours()
	// Show BB
	// dilate.loopContours()

	// NonMaximaSuppression for BB
	// dilate.nonMaxSuppression(10000)

	// Calculate image information's
	dilate.getcenter()
	return dilate
}

func (i *ImageCV) findContours() {
	rects := make([]image.Rectangle, 0)
	contours := gocv.FindContours(i.mat, 0, 2)

	for _, c := range contours {
		area := gocv.ContourArea(c)
		// Filter bb with less then 5000 image Points
		if area > 5000 {
			rect := gocv.BoundingRect(c)
			rects = append(rects, rect)
		}
	}
	i.rects = rects
}

func (i *ImageCV) loopContours() {

	input := i.mat.Clone()
	for _, r := range i.rects {
		gocv.Rectangle(&input, r, color.RGBA{255, 255, 255, 255}, 3)
	}
	window := gocv.NewWindow("Contours")
	window.IMShow(input)
	window.ResizeWindow(1024, 1024)
	window.WaitKey(0)

}

func (i *ImageCV) canny(sigma float64) *ImageCV {

	e := gocv.NewMat()
	img := i.mat
	mean := i.mat.Mean()
	lower := math.Max(0, (1.0-sigma)*mean.Val1)
	upper := math.Min(255, (1.0+sigma)*mean.Val1)
	gocv.Canny(img, &e, float32(lower), float32(upper))

	return &ImageCV{mat: e}
}

func (i *ImageCV) gauSSianBlur(sigmaX float64, sigmaY float64, ksize int) *ImageCV {

	src := i.mat
	dst := gocv.NewMat()
	kernel := image.Point{ksize, ksize}
	gocv.GaussianBlur(src, &dst, kernel, sigmaX, sigmaY, 1)

	return &ImageCV{mat: dst}
}

func initKernel(dim int) gocv.Mat {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(dim, dim))
	return kernel
}

func (i *ImageCV) dilatation(kernelsize int) *ImageCV {

	kernel := initKernel(kernelsize)
	dilate := gocv.NewMat()
	gocv.Dilate(i.mat, &dilate, kernel)
	return &ImageCV{mat: dilate}
}

func (icv *ImageCV) showImg(windownName string) {

	window := gocv.NewWindow(windownName)
	window.IMShow(icv.mat)
	window.ResizeWindow(1024, 1024)
	window.WaitKey(0)
}

func readImg(filename string) *ImageCV {
	return &ImageCV{mat: gocv.IMRead(filename, 0)}
}

func (i *ImageCV) resizeImage(fx float64, fy float64) *ImageCV {

	dst := gocv.NewMat()
	gocv.Resize(i.mat, &dst, image.Point{}, fx, fy, 0)

	return &ImageCV{mat: dst}
}

func (i *ImageCV) crop(left, top, right, bottom int) *ImageCV {
	fmt.Println("Cropping", left, top, right, bottom)
	croppedMat := i.mat.Region(image.Rect(left, top, right, bottom))
	return &ImageCV{mat: croppedMat}
}

func (i *ImageCV) getcenter() {

	dim := i.mat.Size()
	i.height = dim[0]
	i.width = dim[1]
	i.center = []int{i.width / 2, i.height / 2}
}
