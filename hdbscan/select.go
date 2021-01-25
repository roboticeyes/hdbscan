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
	"fmt"
	"log"
	"sort"

	"github.com/fatih/color"
)

func (c *Clustering) selectOptimalClustering(score string) {
	if c.verbose {
		log.Println("selecting optimal clusters")
	}

	switch score {
	case VarianceScore:
		c.setVarianceDeltas()
	case StabilityScore:
		c.setStabilityScore()
	case Leaf:
		c.selectbyLeaves()
	}

	var finalClusters clusters
	for _, cluster := range c.Clusters {
		if cluster.delta == 1 {
			finalClusters = append(finalClusters, cluster)
			color.Cyan("Selected cluster Id: %s has %s points", fmt.Sprint(cluster.id), fmt.Sprint(len(cluster.Points)))
		}
	}

	c.Clusters = finalClusters

	if c.verbose {
		log.Println("Final number of clusters: ", len(c.Clusters))
		log.Println("finished selecting optimal clusters")
	}
}

func (c *Clustering) selectbyLeaves() {

	if c.verbose {
		log.Println("Select by leaves")
	}

	leaves := c.Clusters.leaves()
	forks := c.Clusters.forks()
	if c.verbose {
		log.Println("Number of leaves: ", len(leaves))
		log.Println("Number of forks: ", len(forks))
	}

	// Should only one cluster exist
	if len(forks) == 0 {
		for _, cl := range c.Clusters {
			if cl.id == c.NumberOfClusters {
				cl.delta = 1

				if c.verbose {
					color.Magenta("only one cluster exists")
				}
				return
			}
		}
	}

	for i, leaf := range leaves {

		log.Println("Current leaf", i)
		if leaf.parent == nil {
			continue
		}
		parent := *leaf.parent
		for {
			currentCluster := c.Clusters.getClusterByID(parent)

			if len(currentCluster.children) == 1 {
				if currentCluster.parent == nil {
					currentCluster.delta = 1
					break // Vl nicht richtig
				}
				parent = *currentCluster.parent
			} else if len(currentCluster.children) == 2 {
				child1 := c.Clusters.getClusterByID(currentCluster.children[0])
				child1.delta = 1
				child2 := c.Clusters.getClusterByID(currentCluster.children[1])
				child2.delta = 1
				// Stop at first fork
				// break
				if currentCluster.parent == nil {
					currentCluster.delta = 1
					break // Vl nicht richtig
				}
				parent = *currentCluster.parent
			} else if len(currentCluster.children) > 2 {
				color.Red("More then on child cluster!")
			}
		}
	}

	if c.verbose {
		log.Println("Finished select by leaves")
	}
}

func (c *Clustering) setStabilityScore() {

	leaves := c.Clusters.leaves()

	for _, leaf := range leaves {
		if leaf.viseted == true {
			continue
		}
		leaf.viseted = true
		c.setStabilityDelta(leaf)
	}
}

func (c *Clustering) setStabilityDelta(leaf *cluster) {

	if leaf.parent == nil {
		c.Clusters.getClusterByID(leaf.id).delta = 1
		return
	}
	parent := *leaf.parent
	var hasFork bool
	for {
		currentCluster := c.Clusters.getClusterByID(parent)

		if len(currentCluster.children) == 1 {
			if currentCluster.parent == nil {
				if hasFork == false {
					currentCluster.delta = 1
				}
				break // At Root
			}
			parent = *currentCluster.parent
		} else if len(currentCluster.children) == 2 {

			// Get score for single branch
			scorechild1 := c.Clusters.getClusterByID(currentCluster.children[0]).score

			scorechild2 := c.Clusters.getClusterByID(currentCluster.children[1]).score

			if scorechild1 > scorechild2 {
				c.Clusters.getClusterByID(currentCluster.children[0]).delta = 1
			} else if scorechild1 < scorechild2 {
				c.Clusters.getClusterByID(currentCluster.children[1]).delta = 1
			}

			if currentCluster.parent == nil {
				break // At root
			}
			hasFork = true
			parent = *currentCluster.parent
		} else if len(currentCluster.children) > 2 {
			color.Red("More then on child cluster!")
		}
	}
}

func (c *Clustering) setVarianceDeltas() {
	// sort clusters by size
	sort.Sort(c.Clusters)

	for _, cluster := range c.Clusters {
		// calculate average childrens scores
		var avgScore float64
		for _, child := range cluster.children {
			avgScore += c.Clusters.getClusterByID(child).score
		}
		avgScore /= float64(len(cluster.children))

		if cluster.score <= avgScore && len(cluster.children) > 0 {
			cluster.delta = 0
		} else {
			cluster.delta = 1

			// set sub-clusters to 0
			subClusters := c.Clusters.subTree(cluster.id)
			for _, subCluster := range subClusters {
				subCluster.delta = 0
			}
		}
	}
}

func (c clusters) leaves() clusters {
	var leaves clusters
	for _, cluster := range c {
		if len(cluster.children) == 0 {
			leaves = append(leaves, cluster)
		}
	}
	return leaves
}

func (c clusters) forks() clusters {
	var forks clusters
	for _, cluster := range c {
		if len(cluster.children) == 2 {
			forks = append(forks, cluster)
		}
	}
	return forks
}

func (c clusters) allParents(clstr *cluster) clusters {
	var parents clusters

	if clstr.parent != nil {
		parentCluster := c.getClusterByID(*clstr.parent)
		allParents := c.allParents(parentCluster)
		parents = append(parents, allParents...)
	}

	return parents
}

func (c clusters) subTree(id int) clusters {
	var subTree clusters
	for _, cluster := range c {
		if cluster.parent != nil {
			if *cluster.parent == id {
				subTree = append(subTree, cluster)
				childTree := c.subTree(cluster.id)
				subTree = append(subTree, childTree...)
			}
		}
	}

	return subTree
}
