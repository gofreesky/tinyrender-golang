package lession_3

import (
	"math"
	model "tinyrender-golang/model"
	"tinyrender-golang/tga"
	"tinyrender-golang/util"
)

func drawZBufferTriangle(fb *tga.TGA, obj *model.Object) {

	lightDir := util.Vector3{
		X: 0,
		Y: 0,
		Z: -1,
	}
	zBuffer := make([]float64, fb.GetWidth()*fb.GetHeight())
	for j := 0; j < fb.GetWidth(); j++ {
		for i := 0; i < fb.GetHeight(); i++ {
			zBuffer[i+j*fb.GetWidth()] = -math.MaxFloat64
		}
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
		fb.DrawTriangleWithZBuffer(tri, tga.NewColor(gray, gray, gray, 0), zBuffer)
	}
}

func DrawWithTexture(fb *tga.TGA, obj *model.Object, texture *tga.TGA) {

	zBuffer := make([]float64, fb.GetWidth()*fb.GetHeight())
	for j := 0; j < fb.GetWidth(); j++ {
		for i := 0; i < fb.GetHeight(); i++ {
			zBuffer[i+j*fb.GetWidth()] = -math.MaxFloat64
		}
	}

	for _, face := range obj.Faces {
		tri := util.NewTriangleFromFace(face, float64(fb.GetWidth()-1), float64(fb.GetHeight()-1))
		uvList := []tga.UV{
			{face.Points[0].Texture.U, face.Points[0].Texture.V},
			{face.Points[1].Texture.U, face.Points[1].Texture.V},
			{face.Points[2].Texture.U, face.Points[2].Texture.V},
		}
		fb.DrawTriangleWithTexture(tri, zBuffer, texture, uvList)
	}
}
