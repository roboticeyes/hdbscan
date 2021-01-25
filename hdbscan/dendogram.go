package hdbscan

import (
	"log"
)

type link struct {
	id              int
	parent          *link
	children        []*link
	points          []int
	descendantCount int
}

type node struct {
	key             int
	parentKey       int
	parent          *node
	children        []*node
	descendantCount int
}

func (c *Clustering) buildDendogram(baseEdge edges) []*link {

	if c.verbose {
		log.Println("starting dendrogram")
	}

	var links []*link
	for _, e := range baseEdge {

		var p1TopLink *link
		var p2TopLink *link

		// Are the points are already in a cluster
		for _, link := range links {
			if containsInt(link.points, e.p1) && link.parent == nil {
				p1TopLink = link
			}

			if containsInt(link.points, e.p2) && link.parent == nil {
				p2TopLink = link
			}
		}

		// Both points already belong to a cluster
		// Two new branches are created in the dentogram
		// All previous points are now in both branches
		if p1TopLink != nil && p2TopLink != nil {
			var points []int
			for _, p := range p1TopLink.points {
				points = append(points, p)
			}
			for _, p := range p2TopLink.points {
				points = append(points, p)
			}

			newLink := link{
				id:       len(links),
				children: []*link{p1TopLink, p2TopLink},
				points:   points,
			}

			p1TopLink.parent = &newLink
			p2TopLink.parent = &newLink

			links = append(links, &newLink)
			// One of the two points is not in any cluster,
			// a new cluster is formed and added to the existing cluster as parent.
		} else if p1TopLink != nil && p2TopLink == nil {
			var points []int
			points = append(points, e.p2)
			for _, p := range p1TopLink.points {
				points = append(points, p)
			}

			newlink := link{
				id:       len(links),
				children: []*link{p1TopLink},
				points:   points,
			}

			p1TopLink.parent = &newlink
			links = append(links, &newlink)
			// One of the two points is not in any cluster,
			// a new cluster is formed and added to the existing cluster as parent.
		} else if p1TopLink == nil && p2TopLink != nil {
			var points []int
			points = append(points, e.p1)
			for _, p := range p2TopLink.points {
				points = append(points, p)
			}

			newlink := link{
				id:       len(links),
				children: []*link{p2TopLink},
				points:   points,
			}

			p2TopLink.parent = &newlink
			links = append(links, &newlink)

			// If both points do not belong to any cluster, a new one is formed
		} else {
			newLink := link{
				id:     len(links),
				points: []int{e.p1, e.p2},
			}

			links = append(links, &newLink)
		}
	}

	if c.verbose {
		log.Println("finished dendrogram")
	}

	return links
}
