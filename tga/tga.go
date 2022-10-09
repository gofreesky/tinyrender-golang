package tga

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
)

type rawHeader struct {
	IdLength      uint8
	PaletteType   uint8
	ImageType     uint8
	PaletteFirst  uint16
	PaletteLength uint16
	PaletteBPP    uint8
	OriginX       uint16
	OriginY       uint16
	Width         uint16
	Height        uint16
	BPP           uint8
	Flags         uint8
}

type rawFooter struct {
	ExtAreaOffset uint32
	DevDirOffset  uint32
	Signature     [18]byte // tgaSignature
}

const (
	flagOriginRight   = 1 << 4
	flagOriginTop     = 1 << 5
	flagAlphaSizeMask = 0x0f
)

const (
	imageTypePaletted   = 1
	imageTypeTrueColor  = 2
	imageTypeMonoChrome = 3
	imageTypeMask       = 3
	imageTypeFlagRLE    = 1 << 3
)

const (
	tgaRawHeaderSize = 18
	tgaRawFooterSize = 26
)

const (
	extAreaAttrTypeOffset = 0x1ee
)

const (
	attrTypeAlpha              = 3
	attrTypePremultipliedAlpha = 4
)

var tgaSignature = []byte("TRUEVISION-XFILE.\x00")

func newFooter() *rawFooter {
	f := &rawFooter{}
	copy(f.Signature[:], tgaSignature)
	return f
}

func newExtArea(attrType byte) []byte {
	area := make([]byte, extAreaAttrTypeOffset+1)
	area[0], area[1] = 0xef, 0x01 // size
	area[extAreaAttrTypeOffset] = attrType
	return area
}

func CreateTga(width int, height int) *TGA {
	t := &TGA{
		r: &bytes.Reader{},
		raw: rawHeader{
			IdLength:      0,
			PaletteType:   0,
			ImageType:     2,
			PaletteFirst:  0,
			PaletteLength: 0,
			PaletteBPP:    0,
			OriginX:       0,
			OriginY:       0,
			Width:         uint16(width),
			Height:        uint16(height),
			BPP:           24,
			Flags:         0,
		},
		isPaletted:    false,
		hasAlpha:      false,
		width:         width,
		height:        height,
		pixelSize:     3,
		palette:       nil,
		paletteLength: 0,
		ColorModel:    color.NRGBAModel,
		tmp:           [4]byte{},
		pixels:        make([]byte, 4*width*height),
		decode:        nil,
	}

	return t

}

func (tga *TGA) SaveToFile(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	rect := image.Rect(0, 0, tga.width, tga.height)
	var img image.Image
	if tga.ColorModel == color.NRGBAModel {
		im := image.NewNRGBA(rect)
		im.Pix = tga.pixels
		img = im
	} else {
		im := image.NewRGBA(rect)
		im.Pix = tga.pixels
		img = im
	}

	return Encode(f, img)
}

// applyExtensions reads extensions section (if it exists) and parses attribute type.
func (tga *TGA) applyExtensions() (err error) {
	var rawFooter rawFooter

	if _, err = tga.r.Seek(int64(-tgaRawFooterSize), 2); err != nil {
		return
	} else if err = binary.Read(tga.r, binary.LittleEndian, &rawFooter); err != nil {
		return
	} else if bytes.Equal(rawFooter.Signature[:], tgaSignature[:]) && rawFooter.ExtAreaOffset != 0 {
		offset := int64(rawFooter.ExtAreaOffset + extAreaAttrTypeOffset)

		var n int64
		var t byte

		if n, err = tga.r.Seek(offset, 0); err != nil || n != offset {
			return
		} else if t, err = tga.r.ReadByte(); err != nil {
			return
		} else if t == attrTypeAlpha {
			// alpha
			tga.hasAlpha = true
		} else if t == attrTypePremultipliedAlpha {
			// premultiplied alpha
			tga.hasAlpha = true
			tga.ColorModel = color.RGBAModel
		} else {
			// attribute is not an alpha channel value, ignore it
			tga.hasAlpha = false
		}
	}

	return
}

// flip flips pixels of image based on its origin.
func (tga *TGA) flip(out []byte) {
	flipH := tga.raw.Flags&flagOriginRight != 0
	flipV := tga.raw.Flags&flagOriginTop == 0
	rowSize := tga.width * 4

	if flipH {
		for y := 0; y < tga.height; y++ {
			for x, offset := 0, y*rowSize; x < tga.width/2; x++ {
				a := out[offset+x*4:]
				b := out[offset+(tga.width-x-1)*4:]

				a[0], a[1], a[2], a[3], b[0], b[1], b[2], b[3] = b[0], b[1], b[2], b[3], a[0], a[1], a[2], a[3]
			}
		}
	}

	if flipV {
		for y := 0; y < tga.height/2; y++ {
			for x := 0; x < tga.width; x++ {
				a := out[y*rowSize+x*4:]
				b := out[(tga.height-y-1)*rowSize+x*4:]

				a[0], a[1], a[2], a[3], b[0], b[1], b[2], b[3] = b[0], b[1], b[2], b[3], a[0], a[1], a[2], a[3]
			}
		}
	}
}

// getHeader reads and validates TGA header.
func (tga *TGA) getHeader() (err error) {
	if err = binary.Read(tga.r, binary.LittleEndian, &tga.raw); err != nil {
		return
	}

	if tga.raw.ImageType&imageTypeFlagRLE != 0 {
		tga.decode = decodeRLE
	} else {
		tga.decode = decodeRaw
	}

	tga.raw.ImageType &= imageTypeMask
	alphaSize := tga.raw.Flags & flagAlphaSizeMask

	if alphaSize != 0 && alphaSize != 1 && alphaSize != 8 {
		err = ErrAlphaSize
		return
	}

	tga.hasAlpha = ((alphaSize != 0 || tga.raw.BPP == 32) ||
		(tga.raw.ImageType == imageTypeMonoChrome && tga.raw.BPP == 16) ||
		(tga.raw.ImageType == imageTypePaletted && tga.raw.PaletteBPP == 32))

	tga.width = int(tga.raw.Width)
	tga.height = int(tga.raw.Height)
	tga.pixelSize = int(tga.raw.BPP) >> 3

	// default is NOT premultiplied alpha model
	tga.ColorModel = color.NRGBAModel

	if err = tga.applyExtensions(); err != nil {
		return
	}

	var formatIsInvalid bool

	switch tga.raw.ImageType {
	case imageTypePaletted:
		formatIsInvalid = (tga.raw.PaletteType != 1 ||
			tga.raw.BPP != 8 ||
			tga.raw.PaletteFirst >= tga.raw.PaletteLength ||
			(tga.raw.PaletteBPP != 15 && tga.raw.PaletteBPP != 16 && tga.raw.PaletteBPP != 24 && tga.raw.PaletteBPP != 32))
		tga.isPaletted = true

	case imageTypeTrueColor:
		formatIsInvalid = (tga.raw.BPP != 32 &&
			tga.raw.BPP != 16 &&
			(tga.raw.BPP != 24 || tga.hasAlpha))

	case imageTypeMonoChrome:
		formatIsInvalid = ((tga.hasAlpha && tga.raw.BPP != 16) ||
			(!tga.hasAlpha && tga.raw.BPP != 8))

	default:
		err = fmt.Errorf("TGA: unknown image type %d", tga.raw.ImageType)
	}

	if err == nil && formatIsInvalid {
		err = ErrFormat
	}

	return
}

func (tga *TGA) getPixel(dst []byte) (err error) {
	var R, G, B, A uint8 = 0xff, 0xff, 0xff, 0xff
	src := tga.tmp

	if _, err = io.ReadFull(tga.r, src[0:tga.pixelSize]); err != nil {
		return
	}

	switch tga.pixelSize {
	case 4:
		if tga.hasAlpha {
			A = src[3]
		}
		fallthrough

	case 3:
		B, G, R = src[0], src[1], src[2]

	case 2:
		if tga.raw.ImageType == imageTypeMonoChrome {
			B, G, R = src[0], src[0], src[0]

			if tga.hasAlpha {
				A = src[1]
			}
		} else {
			word := uint16(src[0]) | (uint16(src[1]) << 8)
			B, G, R = wordToBGR(word)

			if tga.hasAlpha && (word&(1<<15)) == 0 {
				A = 0
			}
		}

	case 1:
		if tga.isPaletted {
			index := int(src[0])

			if int(index) >= tga.paletteLength {
				return ErrPaletteIndex
			}

			var m int

			if tga.raw.PaletteBPP == 24 {
				m = index * 3
				B, G, R = tga.palette[m+0], tga.palette[m+1], tga.palette[m+2]
			} else if tga.raw.PaletteBPP == 32 {
				m = index * 4
				B, G, R = tga.palette[m+0], tga.palette[m+1], tga.palette[m+2]

				if tga.hasAlpha {
					A = tga.palette[m+3]
				}
			} else if tga.raw.PaletteBPP == 16 {
				m = index * 2
				word := uint16(tga.palette[m+0]) | (uint16(tga.palette[m+1]) << 8)
				B, G, R = wordToBGR(word)
			}
		} else {
			B, G, R = src[0], src[0], src[0]
		}
	}

	dst[0], dst[1], dst[2], dst[3] = R, G, B, A

	return nil
}

// wordToBGR converts 15-bit color to BGR
func wordToBGR(word uint16) (B, G, R uint8) {
	B = uint8((word >> 0) & 31)
	B = uint8((B << 3) + (B >> 2))
	G = uint8((word >> 5) & 31)
	G = uint8((G << 3) + (G >> 2))
	R = uint8((word >> 10) & 31)
	R = uint8((R << 3) + (R >> 2))
	return
}
