package hdbscan

import (
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (c *Clustering) plotminimumSpanningTree(baseEdge edges) {

	p, err := plot.New()
	if err != nil {
		log.Println("Plot error: ", err)
	}

	xmin, ymin, xmax, ymax := c.findAxisLim()

	p.Title.Text = "Minimum Spanning Tree"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.X.Min = xmin
	p.X.Max = ymin
	p.Y.Min = xmax
	p.Y.Max = ymax

	p.Add(plotter.NewGrid())

	for _, e := range baseEdge {
		pts := c.getCoords(e)
		s, err := plotter.NewScatter(pts)
		s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

		l, err := plotter.NewLine(pts)
		l.LineStyle.Color = color.RGBA{R: 128, G: 128, B: 128, A: 255}
		l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

		if err != nil {
			log.Println("Plot error: ", err)
		}

		p.Add(s, l)
	}

	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, c.directory+"minimumSpannngTree.png"); err != nil {
		panic(err)
	}

	if c.verbose {
		log.Println("Minimum spanning tree plot is saved as:", c.directory+"minimumSpannngTree.png")
	}

}

func (c *Clustering) getCoords(e edge) plotter.XYs {

	pts := make(plotter.XYs, 2)
	p1 := e.p1
	p2 := e.p2
	for i, p := range c.data {
		if i == p1 {
			pts[0].X = p[0]
			pts[0].Y = p[1]
		}

		if i == p2 {
			pts[1].X = p[0]
			pts[1].Y = p[1]
		}
	}
	return pts
}

func (c *Clustering) findAxisLim() (float64, float64, float64, float64) {
	xSorted := mergeSort(c.data, 0)
	ySorted := mergeSort(c.data, 1)

	return xSorted[0][0], ySorted[0][1], xSorted[len(xSorted)-1][0], ySorted[len(xSorted)-1][1]
}
