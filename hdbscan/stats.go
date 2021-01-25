// Copyright 2020 Humility AI Incorporated, All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hdbscan

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
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

func (c *Clustering) distanceDistributions() {
	for i, cluster := range c.Clusters {
		// distance distribution
		ld := float64(math.MinInt64)
		var distances []float64
		for j1, p1 := range cluster.Points {
			if c.nn {
				minDistance := math.MaxFloat64
				for j2, p2 := range cluster.Points {
					if j1 != j2 {
						distance := c.distanceFunc(c.data[p1], c.data[p2])
						if distance < minDistance {
							minDistance = distance
						}
					}
				}
				distances = append(distances, minDistance)

				// largest NN-distance
				if minDistance > ld {
					ld = minDistance
				}
			} else {
				distance := c.distanceFunc(c.data[p1], cluster.Centroid)
				distances = append(distances, distance)
				if distance > ld {
					ld = distance
				}
			}
		}
		dd := &distuv.Normal{}
		dd.Fit(distances, nil)
		cluster.distanceDistribution = dd
		cluster.largestDistance = ld

		c.Clusters[i] = cluster
	}
}
