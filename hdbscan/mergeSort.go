package hdbscan

// "Merge Sort" sorts an [][]float64 slice by given column
func mergeSort(points [][]float64, sortBy int) [][]float64 {

	if len(points) == 1 {
		return points
	}

	mid := int(len(points) / 2)
	var (
		left  = make([][]float64, mid)
		right = make([][]float64, len(points)-mid)
	)
	for i := 0; i < len(points); i++ {
		if i < mid {
			left[i] = points[i]
		} else {
			right[i-mid] = points[i]
		}
	}
	return merge(mergeSort(left, sortBy), mergeSort(right, sortBy), sortBy)
}

func merge(left, right [][]float64, sortBy int) (result [][]float64) {
	result = make([][]float64, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		// ascending: < | descending >
		if left[0][sortBy] < right[0][sortBy] {
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
