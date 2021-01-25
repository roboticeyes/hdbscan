package hdbscan

import (
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/fatih/color"
)

func (c *Clustering) filterEdges(baseEdges edges) edges {
	if c.verbose {
		log.Println("filter edges")
	}
	newEdges := make([]edge, 0)
	sort.Sort(baseEdges)
	mad := madEdges(baseEdges)
	// Sigma is MAD * 1.4826
	sig := mad * 1.4826
	for _, e := range baseEdges {
		// Use only edges with a length between:
		if (e.dist > mad-sig) && (e.dist < mad+sig) {
			newEdges = append(newEdges, e)
		}
	}
	color.Blue("Min: %s Max: %s MAD: %s Sigma: %s", fmt.Sprintf("%1.3f", baseEdges[0].dist), fmt.Sprintf("%1.3f", baseEdges[len(baseEdges)-1].dist), fmt.Sprintf("%1.3f", mad), fmt.Sprintf("%1.3f", sig))
	color.Blue("Lower limit: %s Upper limit: %s", fmt.Sprint(mad-sig), fmt.Sprint(mad+sig))
	color.Blue("Number of edges: %s Number of deleted edges: %s Remaining edges: %s", fmt.Sprint(len(baseEdges)), fmt.Sprint(len(baseEdges)-len(newEdges)), fmt.Sprint(len(newEdges)))

	c.mst.edges = newEdges

	if c.verbose {
		log.Println("finished filter eges")
	}

	return newEdges
}

func madEdges(baseEdges edges) float64 {
	var mad float64
	numEdges := len(baseEdges)
	med := median(baseEdges, numEdges)
	var sum float64
	for _, e := range baseEdges {
		sum += math.Abs(e.dist - med)
	}
	mad = sum / float64(numEdges)
	return mad
}

func median(baseEdges edges, numEdges int) float64 {

	var med float64
	half := numEdges / 2
	// Even
	if numEdges%2 == 0 {
		med = 0.5 * (baseEdges[half].dist + baseEdges[half+1].dist)
		// Odd
	} else {
		med = baseEdges[half+1].dist
	}

	return med
}
