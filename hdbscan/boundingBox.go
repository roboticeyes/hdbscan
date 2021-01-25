package hdbscan

import (
	"github.com/edgeDetection/edgeDetection"
)

func (c *Clustering) clusterToImg() {
	// Read files
	// verticesXYZ, verticesUV, xyzs, uvs := edgeDetection.ReadObjFile(c.directory + "/mesh.obj")
	// // Read Image
	// img := gocv.IMRead(c.directory+"/texture.png", 1)
	// imgDim := img.Size()
	// height := imgDim[0]
	// width := imgDim[1]

	// whitePointsNew := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	// // Iterate clusters
	// for ci, cl := range c.Clusters {
	// 	fmt.Println(ci)
	// 	whitePointsNew = gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	// 	var imgPoints []image.Point
	// 	// Iterate cluster points
	// 	for _, p := range cl.Points {
	// 		clusterPoint := c.storePointInVector(p)
	// 		// Iterate coord index
	// 		for i, xyz_ := range xyzs {

	// 			for _, xyz := range xyz_ {
	// 				point1 := verticesXYZ[xyz]

	// 				dist := point1.Sub3(clusterPoint).Length3()

	// 				if dist < 0.0001 {
	// 					correspondingImgPixel := meanUVCoords(uvs[i], verticesUV)

	// 					u := correspondingImgPixel.X
	// 					v := correspondingImgPixel.Y

	// 					row := int((1-v)*float64(height) - 1) // Rows
	// 					col := int(u * float64(width))        // Columns
	// 					whitePointsNew.SetUCharAt(row, col, uint8(255))
	// 					imgPoints = append(imgPoints, image.Point{X: col, Y: row})
	// 				}
	// 			}
	// 		}
	// 	}
	// 	rect := gocv.BoundingRect(imgPoints)

	// 	gocv.Rectangle(&img, rect, color.RGBA{255, 255, 255, 255}, 3)
	// }
	// window := gocv.NewWindow("whitePointsNew")
	// window.IMShow(img)
	// window.ResizeWindow(1024, 1024)
	// window.WaitKey(0)

	// gocv.IMWrite(c.directory+"/bb.png", img)
}

func meanUVCoords(uvs [3]int, verticesUV []edgeDetection.Vector2) edgeDetection.Vector2 {
	vUV1 := verticesUV[uvs[0]]
	vUV2 := verticesUV[uvs[1]]
	vUV3 := verticesUV[uvs[2]]
	mean := vUV1.Add2(vUV2).Add2(vUV3).MultiplyByScalar2(1. / 3)
	return mean
}

func (c *Clustering) storePointInVector(p int) edgeDetection.Vector3 {
	var clusterPoint edgeDetection.Vector3
	clusterPoint.X = c.data[p][0]
	clusterPoint.Y = c.data[p][1]
	clusterPoint.Z = c.data[p][2]
	return clusterPoint

}
