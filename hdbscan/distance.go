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

	"github.com/golang/geo/r3"
)

// DistanceFunc ...
type DistanceFunc func(x1, x2 []float64) float64

// EuclideanDistance ...
var EuclideanDistance = func(v1, v2 []float64) float64 {
	acc := 0.0
	for i, v := range v1 {
		acc += math.Pow((v - v2[i]), 2)
	}
	return math.Pow(acc, 0.5)
}

// EuclideanDistance ...
var AngleVector = func(v1, v2 []float64) float64 {
	vec1 := r3.Vector{X: v1[0], Y: v1[1], Z: v1[2]}
	vec2 := r3.Vector{X: v2[0], Y: v2[1], Z: v2[2]}
	return float64(vec1.Angle(vec2).Radians())
	// dotproduct := v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
	// v1Length := EuclideanDistance(v1, v1)
	// v2Length := EuclideanDistance(v2, v2)
	// theta := dotproduct / (v1Length * v2Length)
	// return math.Acos(theta)
	// return math.Acos(clamp(theta, -1, 1))
}
