package lession_5

import (
	"math"
	model "tinyrender-golang/model"
	"tinyrender-golang/tga"
	"tinyrender-golang/util"
)

func DrawWithCamera(fb *tga.TGA, obj *model.Object, texture *tga.TGA) {

	zBuffer := make([]float64, fb.GetWidth()*fb.GetHeight())
	for j := 0; j < fb.GetWidth(); j++ {
		for i := 0; i < fb.GetHeight(); i++ {
			zBuffer[i+j*fb.GetWidth()] = -math.MaxFloat64
		}
	}

	eye := util.NewVec3(1, 1, 3)
	center := util.NewVec3(0, 0, 0)

	modelView := util.NewLookAtMatrix(eye, center, util.NewVec3(0, 1, 0))
	cameraVector := util.NewVec3(0, 0, 3)
	projection := util.NewIdentity(4)
	projection.Data[3][2] = -1 / cameraVector.Z
	viewPort := util.NewViewPortMatrix(fb.GetWidth()/8, fb.GetHeight()/8, fb.GetWidth()*3/4, fb.GetHeight()*3/4, 256)

	for _, face := range obj.Faces {
		var screenV [3]util.Vector3
		for i := 0; i < 3; i++ {
			v := util.NewVector3FromVertex(face.Points[i].Vertex)
			screenV[i] = viewPort.Mul(projection).Mul(modelView).Mul(util.NewFromVector3(v)).ToVector3()
		}

		tri := &util.Triangle{Points: screenV[:]}

		uvList := []tga.UV{
			{face.Points[0].Texture.U, face.Points[0].Texture.V},
			{face.Points[1].Texture.U, face.Points[1].Texture.V},
			{face.Points[2].Texture.U, face.Points[2].Texture.V},
		}
		fb.DrawTriangleWithTexture(tri, zBuffer, texture, uvList)
	}
}
