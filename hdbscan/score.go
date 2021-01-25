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
	"log"
	"math"
	"sort"

	"github.com/fatih/color"
)

func (c *Clustering) scoreClusters(optimization string) {

	if c.verbose {
		log.Println("score clusters")
	}

	switch optimization {
	case VarianceScore:
		c.varianceScores()
	case Leaf:
		c.leafScore()
	case StabilityScore:
		c.stabilityScores()
	}

	if c.verbose {
		log.Println("finished score clusters")
	}
}

func (c *Clustering) varianceScores() {
	c.setNormalizedSizes()
	c.setNormalizedVariances()
	c.Clusters.setVarianceScores()
}

func (c clusters) setVarianceScores() {
	for _, cluster := range c {
		cluster.score = cluster.size / cluster.variance
	}
}

func (c *Clustering) setNormalizedSizes() {
	// distro
	var sizes []float64
	for _, cluster := range c.Clusters {
		size := float64(len(cluster.Points))
		sizes = append(sizes, size)
		cluster.size = size
	}
}

func (c *Clustering) setNormalizedVariances() {
	// variances
	var variances []float64
	for _, cluster := range c.Clusters {
		// data
		var clusterData [][]float64
		for _, pointIndex := range cluster.Points {
			clusterData = append(clusterData, c.data[pointIndex])
		}
		// unfold reshape [][]float64 -> []float64 (reshape to list)
		// ClusterData contains the point coordinates
		variance := GeneralizedVariance(len(cluster.Points), len(clusterData[0]), unfold(clusterData))
		cluster.variance = isNum(variance)
		variances = append(variances, cluster.variance)
	}
}

func (c *Clustering) stabilityScores() {

	// Go leaves up to root
	leaves := c.Clusters.leaves()
	for _, leaf := range leaves {
		c.calculateStability(leaf)
	}
}

// https://arxiv.org/pdf/1911.02282.pdf
// https://arxiv.org/pdf/1702.08607.pdf
// https://www.arxiv-vanity.com/papers/1911.02282/
// https://arxiv.org/pdf/1705.07321.pdf
func (c *Clustering) calculateStability(leaf *cluster) {

	pointsLastFork := []int{}
	if leaf.parent == nil {
		leaf.calcClusterStability(c, pointsLastFork)
		return
	}

	parent := *leaf.parent
	for {
		currentCluster := c.Clusters.getClusterByID(parent)

		if len(currentCluster.children) == 1 {
			if currentCluster.parent == nil {
				currentCluster.calcClusterStability(c, pointsLastFork)
				break // At Root
			}
			parent = *currentCluster.parent

			// When 2 child cluster
		} else if len(currentCluster.children) == 2 {

			// Calculate score for single branch
			child1 := c.Clusters.getClusterByID(currentCluster.children[0])
			child1.calcClusterStability(c, pointsLastFork)
			child2 := c.Clusters.getClusterByID(currentCluster.children[1])
			child2.calcClusterStability(c, pointsLastFork)

			// Calculate score for fork
			currentCluster.calcClusterStability(c, pointsLastFork)

			pointsLastFork = currentCluster.Points

			if currentCluster.parent == nil {
				currentCluster.calcClusterStability(c, pointsLastFork)
				break // At root
			}

			parent = *currentCluster.parent
		} else if len(currentCluster.children) > 2 {
			color.Red("More then on child cluster!")
		}
	}

	// Check child cluster
	// if sum of score of both child cluster is bigger then the score of parent cluster
	// score Parent cluster: sum score of child cluster
	for i := 0; i < len(c.Clusters); i++ {
		if len(c.Clusters[i].children) == 2 {
			scoreParentClauster := c.Clusters[i].score
			scoreChild1 := c.Clusters.getClusterByID(c.Clusters[i].children[0]).score
			scoreChild2 := c.Clusters.getClusterByID(c.Clusters[i].children[1]).score

			scoreSumChild := scoreChild1 + scoreChild2
			if scoreParentClauster < scoreSumChild {
				c.Clusters[i].score = scoreSumChild
			}
		}
	}
}

// https://medium.com/@greenraccoon23/multi-thread-for-loops-easily-and-safely-in-go-a2e915302f8b
func (cl *cluster) calcClusterStability(c *Clustering, pointsLastFork []int) {

	// Lambda at which the cluster is born
	birth := make([]float64, len(cl.Points))
	// Lambda at which the point falls out of the cluster
	death := make([]float64, len(cl.Points))

	for i, p1Index := range cl.Points {
		lambda := make([]float64, 0)
		c.wg.Add(1)
		c.semaphore <- true
		go func(i int, p1Index int) {

			for _, p2Index := range cl.Points {
				if p1Index == p2Index {
					continue
				}
				currentLam := c.lambda[p1Index][p2Index]
				if math.IsInf(currentLam, 0) {
					continue
				}
				if currentLam == 0 {
					color.Red("calcClusterStability is Out of controll")
				}
				lambda = append(lambda, currentLam)
			}

			<-c.semaphore
			c.wg.Done()
			// FELIX compare speed of "sort.Float64s" & "your mergeSort"
			sort.Float64s(lambda)

			death[i] = lambda[len(cl.Points)-c.mcs]
			birth[i] = (lambda[0])

		}(i, p1Index)
	}
	c.wg.Wait()
	sort.Float64s(birth)
	cl.lambdaBirth = birth[0]

	var sum float64
	for i := 0; i < len(death); i++ {
		sum += (death[i]) - (cl.lambdaBirth)
	}

	// Maybe its wrong!!!!!
	c.Clusters.getClusterByID(cl.id).score = sum / float64(len(cl.Points))
}

func (c *Clustering) leafScore() {
	for _, cluster := range c.Clusters {
		cluster.size = float64(len(cluster.Points))
		cluster.score = 1
	}
}
