package lession_1

import (
	model "tinyrender-golang/model"
	"tinyrender-golang/tga"
	"tinyrender-golang/util"
)

func drawMesh(fb *tga.TGA, obj *model.Object) {
	for _, face := range obj.Faces {
		for i := 0; i < 3; i++ {
			v1 := face.Points[i]
			v2 := face.Points[(i+1)%3]
			p1 := util.Point{
				X: int((v1.Vertex.X + 1.0) * float64(fb.GetWidth()) / 2.0),
				Y: int((v1.Vertex.Y + 1.0) * float64(fb.GetHeight()) / 2.0),
			}
			p2 := util.Point{
				X: int((v2.Vertex.X + 1.0) * float64(fb.GetWidth()) / 2.0),
				Y: int((v2.Vertex.Y + 1.0) * float64(fb.GetHeight()) / 2.0),
			}
			fb.DrawLine(p1, p2, tga.NewColor(255, 255, 255, 0))
		}
	}

}
