package tga

type UV struct {
	U float64
	V float64
}

func getTextureColor(u float64, v float64, texture *TGA) Color {
	x := int(u * float64(texture.GetWidth()-1))
	y := int(v * float64(texture.GetHeight()-1))
	return texture.GetPixel(x, y)
}
