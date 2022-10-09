package util

import "math"

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func NewVec3(x, y, z float64) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v1 Vector3) Normalize() Vector3 {
	dr := v1.DotProduct(v1)
	drSqrt := math.Sqrt(dr)
	if drSqrt == 0 {
		drSqrt = 0.000001
	}
	return Vector3{
		X: v1.X / drSqrt,
		Y: v1.Y / drSqrt,
		Z: v1.Z / drSqrt,
	}
}

func (v1 Vector3) CrossProduct(v2 Vector3) Vector3 {
	x := v1.Y*v2.Z - v1.Z*v2.Y
	y := v1.Z*v2.X - v1.X*v2.Z
	z := v1.X*v2.Y - v1.Y*v2.X

	return Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v1 Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

func (v1 Vector3) DotProduct(v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

type Point struct {
	X int
	Y int
}

type Triangle struct {
	Points []Vector3
}

// InTriangle 利用了重心坐标
func InTriangle(tri *Triangle, p Point) bool {

	b := Barycentric(tri, p)

	return !(b.X < 0 || b.Y < 0 || b.Z < 0)
}

func Barycentric(tri *Triangle, p Point) Vector3 {
	// ACx
	v1x := tri.Points[2].X - tri.Points[0].X
	// ABx
	v1y := tri.Points[1].X - tri.Points[0].X
	// PAx
	v1z := tri.Points[0].X - float64(p.X)
	v1 := Vector3{
		X: float64(v1x),
		Y: float64(v1y),
		Z: float64(v1z),
	}

	// ACy
	v2x := tri.Points[2].Y - tri.Points[0].Y
	// ABy
	v2y := tri.Points[1].Y - tri.Points[0].Y
	// PAx
	v2z := tri.Points[0].Y - float64(p.Y)
	v2 := Vector3{
		X: float64(v2x),
		Y: float64(v2y),
		Z: float64(v2z),
	}

	v3 := v1.CrossProduct(v2)
	if math.Abs(v3.Z) < 0.000001 {
		return Vector3{
			X: -1,
			Y: 1,
			Z: 1,
		}
	}

	a := float64(1) - float64(v3.X+v3.Y)/float64(v3.Z)
	b := float64(v3.Y) / float64(v3.Z)
	c := float64(v3.X) / float64(v3.Z)

	return Vector3{
		X: a,
		Y: b,
		Z: c,
	}
}
