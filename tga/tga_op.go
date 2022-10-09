package tga

import "tinyrender-golang/util"

type Color struct {
	R byte
	G byte
	B byte
	A byte
}

func NewColor(r, g, b, a byte) Color {
	return Color{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (tga *TGA) GetWidth() int {
	return tga.width
}

func (tga *TGA) GetHeight() int {
	return tga.height
}

func (tga *TGA) SetPixel(x int, y int, c Color) {
	pixel := tga.pixels[(y*tga.width+x)*4:]
	pixel[0] = c.R
	pixel[1] = c.G
	pixel[2] = c.B
	pixel[3] = c.A
}

func (tga *TGA) GetPixel(x int, y int) Color {
	pixel := tga.pixels[(y*tga.width+x)*4:]
	return Color{
		R: pixel[0],
		G: pixel[1],
		B: pixel[2],
		A: pixel[3],
	}
}

func (tga *TGA) FlipVertical() {
	for i := 0; i < tga.height/2; i++ {
		for j := 0; j < tga.width; j++ {
			p1 := tga.GetPixel(j, i)
			p2 := tga.GetPixel(j, tga.height-i-1)
			tga.SetPixel(j, i, p2)
			tga.SetPixel(j, tga.height-i-1, p1)
		}
	}
}

func (tga *TGA) DrawLine(p1 util.Point, p2 util.Point, c Color) {
	if abs(p1.X-p2.X) > abs(p1.Y-p2.Y) {
		tga.drawLineByX(p1, p2, c)
	} else {
		tga.drawLineByY(p1, p2, c)
	}
}

func (tga *TGA) drawLineByX(p1 util.Point, p2 util.Point, c Color) {
	x1 := p1.X
	x2 := p2.X
	y1 := p1.Y
	y2 := p2.Y
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	ratio := float64(y2-y1) / float64(x2-x1)

	for i := x1; i < x2; i++ {
		y := y1 + int(float64(i-x1)*ratio)
		if y < 0 {
			y = 0
		}
		if y >= tga.height {
			y = tga.height - 1
		}
		tga.SetPixel(i, y, c)
	}
}

func (tga *TGA) drawLineByY(p1 util.Point, p2 util.Point, c Color) {
	x1 := p1.X
	x2 := p2.X
	y1 := p1.Y
	y2 := p2.Y
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	ratio := float64(x2-x1) / float64(y2-y1)

	for i := y1; i < y2; i++ {
		x := x1 + int(float64(i-y1)*ratio)
		if x < 0 {
			x = 0
		}
		if x >= tga.width {
			x = tga.width - 1
		}
		tga.SetPixel(x, i, c)
	}
}

func (tga *TGA) DrawTriangle(tri *util.Triangle, c Color) {
	minX := tga.width - 1
	minY := tga.height - 1
	maxX := 0
	maxY := 0

	for i := 0; i < len(tri.Points); i++ {
		p := tri.Points[i]
		px := int(p.X)
		py := int(p.Y)
		maxX = Max(maxX, px)
		maxY = Max(maxY, py)
		minX = Min(minX, px)
		minY = Min(minY, py)
	}

	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			p := util.Point{
				X: i,
				Y: j,
			}
			if util.InTriangle(tri, p) {
				tga.SetPixel(i, j, c)
			}
		}
	}
}

func (tga *TGA) DrawTriangleWithZBuffer(tri *util.Triangle, c Color, zBuffer []float64) {
	minX := tga.width - 1
	minY := tga.height - 1
	maxX := 0
	maxY := 0

	for i := 0; i < len(tri.Points); i++ {
		p := tri.Points[i]
		px := int(p.X)
		py := int(p.Y)
		maxX = Max(maxX, px)
		maxY = Max(maxY, py)
		minX = Min(minX, px)
		minY = Min(minY, py)
	}

	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			p := util.Point{
				X: i,
				Y: j,
			}
			bary := util.Barycentric(tri, p)
			if bary.X < 0 || bary.Y < 0 || bary.Z < 0 {
				continue
			}
			z := tri.Points[0].Z*bary.X + tri.Points[1].Z*bary.Y + tri.Points[2].Z*bary.Z
			if zBuffer[i+j*tga.width] < z {
				zBuffer[i+j*tga.width] = z
				tga.SetPixel(i, j, c)
			}
		}
	}
}

func (tga *TGA) DrawTriangleWithTexture(tri *util.Triangle, zBuffer []float64, texture *TGA, uvList []UV) {

	minX := tga.width - 1
	minY := tga.height - 1
	maxX := 0
	maxY := 0

	for i := 0; i < len(tri.Points); i++ {
		p := tri.Points[i]
		px := int(p.X)
		py := int(p.Y)
		maxX = Max(maxX, px)
		maxY = Max(maxY, py)
		minX = Min(minX, px)
		minY = Min(minY, py)
	}

	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			p := util.Point{
				X: i,
				Y: j,
			}
			bary := util.Barycentric(tri, p)
			if bary.X < 0 || bary.Y < 0 || bary.Z < 0 {
				continue
			}
			z := tri.Points[0].Z*bary.X + tri.Points[1].Z*bary.Y + tri.Points[2].Z*bary.Z
			u := uvList[0].U*bary.X + uvList[1].U*bary.Y + uvList[2].U*bary.Z
			v := uvList[0].V*bary.X + uvList[1].V*bary.Y + uvList[2].V*bary.Z
			c := getTextureColor(u, v, texture)

			if zBuffer[i+j*tga.width] < z {
				zBuffer[i+j*tga.width] = z
				tga.SetPixel(i, j, c)
			}
		}
	}
}
