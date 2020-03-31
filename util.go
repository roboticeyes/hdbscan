package hdbscan

import "math"

func max(values []float64) float64 {
	maxNum := float64(math.MinInt64)
	for _, v := range values {
		if v > maxNum {
			maxNum = v
		}
	}

	return maxNum
}

func unfold(data [][]float64) []float64 {
	var ud []float64
	for _, row := range data {
		ud = append(ud, row...)
	}
	return ud
}

func isNaN(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	}

	return value
}

func isInf(value float64) float64 {
	if math.IsInf(value, 1) {
		return math.MaxFloat64
	}

	if math.IsInf(value, -1) {
		return float64(math.MinInt64)
	}

	return value
}

func isNum(value float64) float64 {
	return isNaN(isInf(value))
}

func containsNode(list []node, n node) bool {
	for _, node := range list {
		if node.key == n.key {
			return true
		}
	}

	return false
}
