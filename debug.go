package hdbscan

import (
	"encoding/json"
)

type debugClusterVariance struct {
	Score    float64
	Variance float64
	Size     float64
	Delta    int
	Points   []int
}

type debugNode struct {
	Key      int
	Children []int
}

type debugEdge struct {
	P1, P2   int
	Distance float64
}

// NewNode will return a printable version of the node object
// used for debugging.
func (n *node) NewNode() debugNode {
	var children []int
	for _, child := range n.children {
		children = append(children, child.key)
	}

	return debugNode{
		Key:      n.key,
		Children: children,
	}
}

func (n *node) allNewNodes(output []debugNode) []debugNode {
	output = append(output, n.NewNode())

	for _, child := range n.children {
		allOutputs := child.allNewNodes([]debugNode{})
		output = append(output, allOutputs...)
	}

	return output
}

// String can be used for testing/debugging by printing
// a json version of the complete dendrogram hierarchy.
func (n *node) String() string {
	outputs := n.allNewNodes([]debugNode{})
	data, _ := json.MarshalIndent(outputs, "", "  ")
	return string(data)
}

func debugEdges(edges []edge) []debugEdge {
	var debugEdges []debugEdge
	for _, edge := range edges {
		d := debugEdge{
			P1:       edge.p1,
			P2:       edge.p2,
			Distance: edge.dist,
		}
		debugEdges = append(debugEdges, d)
	}

	return debugEdges
}

func debugPointsMap(m map[int]bool) []int {
	var points []int
	for point, ok := range m {
		if ok {
			points = append(points, point)
		}
	}

	return points
}

func debugNodeKeys(nodes []node) []int {
	var keys []int

	for _, node := range nodes {
		keys = append(keys, node.key)
	}

	return keys
}

func (c *cluster) debugClusterVariance() debugClusterVariance {
	return debugClusterVariance{
		Score:    c.score,
		Variance: c.variance,
		Size:     c.size,
		Delta:    c.delta,
		Points:   c.pointIndexes(),
	}
}

func (c *cluster) allDebugClusterVariances(output []debugClusterVariance) []debugClusterVariance {
	output = append(output, c.debugClusterVariance())

	for _, child := range c.children {
		allOutputs := child.allDebugClusterVariances([]debugClusterVariance{})
		output = append(output, allOutputs...)
	}

	return output
}

// String ...
func (c cluster) String() string {
	output := c.allDebugClusterVariances([]debugClusterVariance{})
	data, _ := json.MarshalIndent(output, "", "  ")
	return string(data)
}
