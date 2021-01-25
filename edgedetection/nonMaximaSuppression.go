package edgeDetection

import "image"

func (i *ImageCV) nonMaxSuppression(overlapThreshold float64) {
	if len(i.rects) == 0 {
		return
	}

	// Sort BB by bottom- right y- coordinate
	i.rects = mergeSort(i.rects)

	bb := make([][]float64, len(i.rects))
	area := make([]float64, len(i.rects))

	for i, b := range i.rects {
		bb[i] = []float64{float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y)}
		area[i] = (float64(b.Max.X) - float64(b.Min.X+1)) / (float64(b.Max.Y) - float64(b.Min.Y+1))
	}

	suppress := make([]int, 0)
	for len(bb) > 0 {
		last := len(bb) - 1
		var xx1, yy1, xx2, yy2, w, h float64
		for j := 0; j < last; j++ {

			// find the largest (x, y) coordinates for the start of
			// the bounding box and the smallest (x, y) coordinates
			// for the end of the bounding box
			xx1 = maxValue(bb[last][0], bb[j][0])
			yy1 = maxValue(bb[last][1], bb[j][1])
			xx2 = minValue(bb[last][2], bb[j][2])
			yy2 = minValue(bb[last][3], bb[j][3])

			// compute the width and height of the bounding box
			w = maxValue(0, xx2-xx1+1)
			h = maxValue(0, yy2-yy1+1)

			overlap := (w * h) / area[j]

			if overlap > overlapThreshold {
				contain := containsInt(suppress, j)
				if !contain {
					suppress = append(suppress, j)
				}
			}
		}
		bb = bb[:last]
	}

	rectNew := make([]image.Rectangle, 0)
	for k := 0; k < len(i.rects); k++ {
		if !containsInt(suppress, k) {
			rectNew = append(rectNew, i.rects[k])
		}
	}
	i.rects = rectNew
}

func containsInt(list []int, i int) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}

	return false
}

func maxValue(x float64, y float64) float64 {

	if x > y {
		return x
	} else {
		return y
	}
}

func minValue(x float64, y float64) float64 {

	if x > y {
		return y
	} else {
		return x
	}
}
