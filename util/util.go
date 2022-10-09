package util

import (
	obj "tinyrender-golang/model"
)

func NewVector3FromVertex(v *obj.Vertex) Vector3 {
	return Vector3{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}

func NewTriangleFromFace(face obj.Face, width float64, height float64) *Triangle {
	tri := &Triangle{Points: []Vector3{
		{
			X: ((face.Points[0].Vertex.X + 1.0) * width / 2.0) + 0.5,
			Y: ((face.Points[0].Vertex.Y + 1.0) * height / 2.0) + 0.5,
			Z: face.Points[0].Vertex.Z,
		},
		{
			X: ((face.Points[1].Vertex.X + 1.0) * width / 2.0) + 0.5,
			Y: ((face.Points[1].Vertex.Y + 1.0) * height / 2.0) + 0.5,
			Z: face.Points[1].Vertex.Z,
		},
		{
			X: ((face.Points[2].Vertex.X + 1.0) * width / 2.0) + 0.5,
			Y: ((face.Points[2].Vertex.Y + 1.0) * height / 2.0) + 0.5,
			Z: face.Points[2].Vertex.Z,
		},
	}}
	return tri
}

func NewViewPortMatrix(x int, y int, w int, h int, depth float64) *Matrix {
	m := NewIdentity(4)
	m.Data[0][3] = float64(x + w/2)
	m.Data[1][3] = float64(y + h/2)
	m.Data[2][3] = depth / 2

	m.Data[0][0] = float64(w / 2)
	m.Data[1][1] = float64(h / 2)
	m.Data[2][2] = depth / 2
	return m
}

func NewLookAtMatrix(eye Vector3, center Vector3, up Vector3) *Matrix {
	z := eye.Sub(center).Normalize()
	x := up.CrossProduct(z).Normalize()
	y := z.CrossProduct(x).Normalize()
	r := NewIdentity(4)

	r.Data[0][0] = x.X
	r.Data[1][0] = y.X
	r.Data[2][0] = z.X

	r.Data[0][1] = x.Y
	r.Data[1][1] = y.Y
	r.Data[2][1] = z.Y

	r.Data[0][2] = x.Z
	r.Data[1][2] = y.Z
	r.Data[2][2] = z.Z

	r.Data[0][3] = -center.X
	r.Data[1][3] = -center.Y
	r.Data[2][3] = -center.Z

	return r
}
