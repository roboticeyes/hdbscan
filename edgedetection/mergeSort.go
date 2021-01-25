package edgeDetection

import "image"

// Merge Sort algorithm to sort BB by bottom right Y- value
func mergeSort(rects []image.Rectangle) []image.Rectangle {

	if len(rects) == 1 {
		return rects
	}

	mid := int(len(rects) / 2)
	var (
		left  = make([]image.Rectangle, mid)
		right = make([]image.Rectangle, len(rects)-mid)
	)
	for i := 0; i < len(rects); i++ {
		if i < mid {
			left[i] = rects[i]
		} else {
			right[i-mid] = rects[i]
		}
	}
	return merge(mergeSort(left), mergeSort(right))
}

func merge(left, right []image.Rectangle) (result []image.Rectangle) {
	result = make([]image.Rectangle, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		// ascending: < | descending >
		if left[0].Max.Y < right[0].Max.Y {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		result[i] = left[j]
		i++
	}
	for j := 0; j < len(right); j++ {
		result[i] = right[j]
		i++
	}

	return
}
