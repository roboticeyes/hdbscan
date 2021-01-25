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
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// the mutual reachability graph provides a mutual-reachability-distance matrix
// which specifies a metric of how far a point is from another point.
func (c *Clustering) mutualReachabilityGraph() edges {
	if c.verbose {
		log.Println("starting mutual reachability")
	}

	lambda := make([][]float64, len(c.data))
	// core-distances
	length := len(c.data)
	coreDistances := make([]float64, length, length)
	for i, p1 := range c.data {
		lambda[i] = make([]float64, len(c.data))
		c.wg.Add(1)
		c.semaphore <- true
		go func(i int, p1 []float64) {
			pointDistances := []float64{}
			for ii, p2 := range c.data {
				dist := c.distanceFunc(p1, p2)
				pointDistances = append(pointDistances, dist)
				// Transform in lambda
				lambda[i][ii] = (1 / dist)
			}
			sort.Float64s(pointDistances)
			coreDistances[i] = pointDistances[c.mcs-1]
			<-c.semaphore
			c.wg.Done()
		}(i, p1)
	}
	c.lambda = lambda
	c.wg.Wait()

	// mutual-reachability distances
	for i := 0; i < length; i++ {
		c.wg.Add(1)
		c.semaphore <- true
		go func(i int) {
			mutualReachabilityDistances := make([]float64, length, length)
			// the mutual reachability distance is the maximum of:
			// point_1's core-distance, point_2's core-distance, or the distance between point_1 and point_2
			// max{dcore(xp),dcore(xq),d(xp,xq)}
			for j := 0; j < length; j++ {
				mutualReachabilityDistances[j] = max([]float64{coreDistances[i], coreDistances[j], c.distanceFunc(c.data[i], c.data[j])})
			}

			if c.minTree {
				c.addRowToMinSpanningTree(i, mutualReachabilityDistances)
			} else {
				minIndex, minValue := min(mutualReachabilityDistances)
				e := edge{
					p1:   i,
					p2:   minIndex,
					dist: minValue,
				}

				// just use tree for edge storage
				c.mst.addEdge(e)
			}
			<-c.semaphore
			c.wg.Done()
		}(i)
	}
	c.wg.Wait()

	outputfile, _ := os.Create(c.directory + "debug1.txt")
	defer outputfile.Close()
	writer := bufio.NewWriter(outputfile)
	for _, p := range c.mst.edges {

		x := fmt.Sprintf("%v", p.p1)
		y := fmt.Sprintf("%v", p.p2)
		z := fmt.Sprintf("%f", p.dist)

		_, err := writer.WriteString(x + " " + y + " " + z + "\n")
		if err != nil {
			panic(err)
		}

	}
	// Very important to invoke after writing a large number of lines
	writer.Flush()

	sort.Sort(c.mst.edges)

	outputfile, _ = os.Create(c.directory + "debug2.txt")
	defer outputfile.Close()
	writer = bufio.NewWriter(outputfile)
	for _, p := range c.mst.edges {

		x := fmt.Sprintf("%v", p.p1)
		y := fmt.Sprintf("%v", p.p2)
		z := fmt.Sprintf("%f", p.dist)

		_, err := writer.WriteString(x + " " + y + " " + z + "\n")
		if err != nil {
			panic(err)
		}

	}
	// Very important to invoke after writing a large number of lines
	writer.Flush()

	if c.verbose {
		log.Println("finished mutual reachability")
	}

	return c.mst.edges
}
