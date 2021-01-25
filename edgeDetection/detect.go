package edgeDetection

import "io"

// github.com/go-gl/mathgl/mgl64
type Data struct {
	img *ImageCV

	depthFile []float64

	coordXYZ []Vector3
	coordUV  []Vector2

	Points [][]float64

	indexXYZ [][3]int
	indexUV  [][3]int

	clusterNumber []int
	mean          []Vector3

	Normale    [][]float64
	Barycenter [][]float64
}

type Vector3 struct {
	X, Y, Z float64
}

type Vector2 struct {
	X, Y float64
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

	whitePoints := d.correspondendingPlaneNormals()
	whitePoints.showImg("whitePoints")

	// whitePoints := d.findCorrespondendingPoints()
	// whitePoints.showImg("whitePoints")

	return d
}
