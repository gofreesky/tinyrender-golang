package lession_2

import (
	"math/rand"
	model "tinyrender-golang/model"
	"tinyrender-golang/tga"
	"tinyrender-golang/util"
)

func drawRandomColorTriangle(fb *tga.TGA, obj *model.Object) {
	for _, face := range obj.Faces {
		tri := util.NewTriangleFromFace(face, float64(fb.GetWidth()-1), float64(fb.GetHeight()-1))
		fb.DrawTriangle(tri, tga.NewColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), 0))
	}
}

func drawModelTriangle(fb *tga.TGA, obj *model.Object) {
	lightDir := util.Vector3{
		X: 0,
		Y: 0,
		Z: -1,
	}
	for _, face := range obj.Faces {
		tri := util.NewTriangleFromFace(face, float64(fb.GetWidth()-1), float64(fb.GetHeight()-1))
		v0 := util.NewVector3FromVertex(face.Points[0].Vertex)
		v1 := util.NewVector3FromVertex(face.Points[1].Vertex)
		v2 := util.NewVector3FromVertex(face.Points[2].Vertex)

		n := v2.Sub(v0).CrossProduct(v1.Sub(v0))
		nNor := n.Normalize()
		var intensity = nNor.DotProduct(lightDir)
		if intensity < 0 {
			continue
		}
		gray := byte(intensity * 255)
		fb.DrawTriangle(tri, tga.NewColor(gray, gray, gray, 0))
	}
}
