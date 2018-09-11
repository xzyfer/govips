package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import "io"

type imageRef struct {
	i *C.VipsImage
}

type Image struct {
	ref *imageRef
}

func LoadImage2(r io.Reader) (*Image, error) {
	return nil, nil
}

func (i *Image) Free() {
}

// Format returns the initial format of the vips image when loaded
func (i *Image) Format() ImageType {
	return ImageTypeUnknown
}

// Width returns the width of this image
func (i *Image) Width() int {
	return int(i.ref.i.Xsize)
}

// Height returns the height of this iamge
func (i *Image) Height() int {
	return int(i.ref.i.Ysize)
}

// Bands returns the number of bands for this image
func (i *Image) Bands() int {
	return int(i.ref.i.Bands)
}

// ResolutionX returns the X resolution
func (i *Image) ResolutionX() float64 {
	return float64(i.ref.i.Xres)
}

// ResolutionY returns the Y resolution
func (i *Image) ResolutionY() float64 {
	return float64(i.ref.i.Yres)
}

// OffsetX returns the X offset
func (i *Image) OffsetX() int {
	return int(i.ref.i.Xoffset)
}

// OffsetY returns the Y offset
func (i *Image) OffsetY() int {
	return int(i.ref.i.Yoffset)
}

// BandFormat returns the current band format
func (i *Image) BandFormat() BandFormat {
	return BandFormat(int(i.ref.i.BandFmt))
}

// Coding returns the image coding
func (i *Image) Coding() Coding {
	return Coding(int(i.ref.i.Coding))
}

// Interpretation returns the current interpretation
func (i *Image) Interpretation() Interpretation {
	return Interpretation(int(i.ref.i.Type))
}
