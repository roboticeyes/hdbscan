package edgeDetection

import (
	"io"

	"github.com/go-gl/mathgl/mgl64"
)

type Data struct {
	img *ImageCV

	depthFile []float64

	coordXYZ []mgl64.Vec3
	coordUV  []mgl64.Vec2

	Points [][]float64

	indexXYZ [][3]int
	indexUV  [][3]int

	clusterNumber []int
	mean          []mgl64.Vec3

	Normale    [][]float64
	Barycenter [][]float64
}

func Detection(meshReader, depthReader, imgReader io.Reader) *Data {

	verticesXYZ, verticesUV, xyzs, uvs := ReadObjFile(meshReader)
	d := &Data{
		img:       ImageControler(imgReader),
		depthFile: ReadDepthFile(depthReader),
		coordXYZ:  verticesXYZ,
		coordUV:   verticesUV,
		indexXYZ:  xyzs,
		indexUV:   uvs,
	}

	// Barycenter or raw points
	whitePoints := d.correspondendingPoints("barycenter")
	whitePoints.showImg("whitePoints")

	return d
}
