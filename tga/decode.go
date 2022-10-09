package tga

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
)

type TGA struct {
	r             *bytes.Reader
	raw           rawHeader
	isPaletted    bool
	hasAlpha      bool
	width         int
	height        int
	pixelSize     int
	palette       []byte
	paletteLength int
	ColorModel    color.Model
	tmp           [4]byte
	pixels        []byte
	decode        func(tga *TGA, out []byte) (err error)
}

var (
	ErrAlphaSize    = errors.New("TGA: invalid alpha size")
	ErrFormat       = errors.New("TGA: invalid format")
	ErrPaletteIndex = errors.New("TGA: palette index out of range")
)

// Decode decodes a TARGA image.
func Decode(r io.Reader) (outImage image.Image, err error) {
	var tga TGA
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	tga.r = bytes.NewReader(data.Bytes())

	if err = tga.getHeader(); err != nil {
		return
	}

	// skip header
	if _, err = tga.r.Seek(int64(tgaRawHeaderSize+tga.raw.IdLength), 0); err != nil {
		return
	}

	if tga.isPaletted {
		// read palette
		entrySize := int((tga.raw.PaletteBPP + 1) >> 3)
		tga.paletteLength = int(tga.raw.PaletteLength - tga.raw.PaletteFirst)
		tga.palette = make([]byte, entrySize*tga.paletteLength)

		// skip to colormap
		if _, err = tga.r.Seek(int64(entrySize)*int64(tga.raw.PaletteFirst), 1); err != nil {
			return
		}

		if _, err = io.ReadFull(tga.r, tga.palette); err != nil {
			return
		}
	}

	rect := image.Rect(0, 0, tga.width, tga.height)
	var pixels []byte

	// choose a right color model
	if tga.ColorModel == color.NRGBAModel {
		im := image.NewNRGBA(rect)
		outImage = im
		pixels = im.Pix
	} else {
		im := image.NewRGBA(rect)
		outImage = im
		pixels = im.Pix
	}

	if err = tga.decode(&tga, pixels); err == nil {
		tga.flip(pixels)
	}

	return
}

func DecodeToTga(r io.Reader) (tgaV *TGA, err error) {

	tgaV = &TGA{}
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	tgaV.r = bytes.NewReader(data.Bytes())

	if err = tgaV.getHeader(); err != nil {
		return
	}

	// skip header
	if _, err = tgaV.r.Seek(int64(tgaRawHeaderSize+tgaV.raw.IdLength), 0); err != nil {
		return
	}

	if tgaV.isPaletted {
		// read palette
		entrySize := int((tgaV.raw.PaletteBPP + 1) >> 3)
		tgaV.paletteLength = int(tgaV.raw.PaletteLength - tgaV.raw.PaletteFirst)
		tgaV.palette = make([]byte, entrySize*tgaV.paletteLength)

		// skip to colormap
		if _, err = tgaV.r.Seek(int64(entrySize)*int64(tgaV.raw.PaletteFirst), 1); err != nil {
			return
		}

		if _, err = io.ReadFull(tgaV.r, tgaV.palette); err != nil {
			return
		}
	}

	tgaV.pixels = make([]byte, 4*tgaV.width*tgaV.height)

	if err = tgaV.decode(tgaV, tgaV.pixels); err == nil {
		tgaV.flip(tgaV.pixels)
	}

	return
}

// DecodeConfig decodes a header of TARGA image and returns its configuration.
func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var tga TGA
	var data bytes.Buffer

	if _, err = data.ReadFrom(r); err != nil {
		return
	}

	tga.r = bytes.NewReader(data.Bytes())

	if err = tga.getHeader(); err == nil {
		cfg = image.Config{
			ColorModel: tga.ColorModel,
			Width:      tga.width,
			Height:     tga.height,
		}
	}

	return
}

func init() {
	image.RegisterFormat("TGA", "", Decode, DecodeConfig)
}

// decodeRaw decodes a raw (uncompressed) data.
func decodeRaw(tga *TGA, out []byte) (err error) {
	for i := 0; i < len(out) && err == nil; i += 4 {
		err = tga.getPixel(out[i:])
	}

	return
}

// decodeRLE decodes run-length encoded data.
func decodeRLE(tga *TGA, out []byte) (err error) {
	size := tga.width * tga.height * 4

	for i := 0; i < size && err == nil; {
		var b byte

		if b, err = tga.r.ReadByte(); err != nil {
			break
		}

		count := uint(b)

		if count&(1<<7) != 0 {
			// encoded packet
			count &= ^uint(1 << 7)

			if err = tga.getPixel(tga.tmp[:]); err == nil {
				for count++; count > 0 && i < size; count-- {
					copy(out[i:], tga.tmp[:])
					i += 4
				}
			}
		} else {
			// raw packet
			for count++; count > 0 && i < size && err == nil; count-- {
				err = tga.getPixel(out[i:])
				i += 4
			}
		}
	}

	return
}
