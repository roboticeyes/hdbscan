package hdbscan

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

// GeneralizedVariance will return the determinant of the covariance matrix
// of the supplied data.
// The supplied data is a list of 'rows' observations of length 'columns'.
func GeneralizedVariance(rows, columns int, data []float64) float64 {
	covMatrix := &mat.SymDense{}
	matrix := mat.NewDense(rows, columns, data)
	stat.CovarianceMatrix(covMatrix, matrix, nil)
	det, _ := mat.LogDet(covMatrix)
	return math.Abs(det)
}
