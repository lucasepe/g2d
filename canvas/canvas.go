package canvas

import (
	"fmt"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	minimumSize     = 256
	maximumSize     = 1536
	sizeIncrement   = 64
	defaultSize     = 1024
	defaultFontSize = 12.0
)

// Canvas defines a drawing context.
type Canvas struct {
	size       int
	font       *truetype.Font
	fontSize   float64
	xMin, xMax float64
	yMin, yMax float64
	dc         *gg.Context
}

// Option defines a screen parameter.
type Option func(*Canvas)

// Size defines the canvas size option.
func Size(size int64) Option {
	return func(scr *Canvas) {
		if size < minimumSize {
			size = minimumSize
		}

		if size > maximumSize {
			size = maximumSize
		}

		div := math.Round(float64(size) / float64(sizeIncrement))
		mul, _ := math.Modf(div)
		scr.size = int(mul) * sizeIncrement
	}
}

// NewCanvas creates a drawing context.
func NewCanvas(opts ...Option) *Canvas {
	res := &Canvas{
		size:     defaultSize,
		fontSize: defaultFontSize,
	}

	for _, op := range opts {
		op(res)
	}

	res.dc = gg.NewContext(res.size, res.size)

	res.dc.SetHexColor("#ffffff")
	res.dc.Clear()

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %s\n", err.Error())
	} else {
		res.font = font
		face := truetype.NewFace(font, &truetype.Options{Size: res.fontSize})
		res.dc.SetFontFace(face)
	}

	delta := float64(res.size)
	res.xMin = -0.5 * delta
	res.xMax = 0.5 * delta
	res.yMin = -0.5 * delta
	res.yMax = 0.5 * delta

	mapWorldCoordinates(res.dc, res.xMin, res.xMax, res.yMin, res.yMax)

	return res
}

// Graphics returns the canvas graphic context.
func (scr *Canvas) Graphics() *gg.Context { return scr.dc }

// Reset resets screen using the specified options.
func (scr *Canvas) Reset(opts ...Option) {
	for _, op := range opts {
		op(scr)
	}

	scr.dc = gg.NewContext(scr.size, scr.size)

	scr.dc.SetHexColor("#ffffff")
	scr.dc.Clear()

	scr.SetFontSize(12)

	scr.xMin = 0
	scr.xMax = float64(scr.size)
	scr.yMin = 0
	scr.yMax = float64(scr.size)

	delta := float64(scr.size)
	scr.xMin = -0.5 * delta
	scr.xMax = 0.5 * delta
	scr.yMin = -0.5 * delta
	scr.yMax = 0.5 * delta

	mapWorldCoordinates(scr.dc, scr.xMin, scr.xMax, scr.yMin, scr.yMax)
}

// SetFontSize sets the font's size.
func (scr *Canvas) SetFontSize(size float64) {
	if scr.font != nil {
		scr.fontSize = size
		face := truetype.NewFace(scr.font, &truetype.Options{Size: size})
		scr.dc.SetFontFace(face)
	}
}

// FontSize returns font's size.
func (scr *Canvas) FontSize() float64 { return scr.fontSize }

// Size returns the screen size.
func (scr *Canvas) Size() int { return scr.size }

// Xmin returns the screen xMin coords.
func (scr *Canvas) Xmin() float64 { return scr.xMin }

// Xmax returns the screen xMax coords.
func (scr *Canvas) Xmax() float64 { return scr.xMax }

// Ymin returns the screen yMin coords.
func (scr *Canvas) Ymin() float64 { return scr.yMin }

// Ymax returns the screen yMax coords.
func (scr *Canvas) Ymax() float64 { return scr.yMax }

// SetWorldCoordinates sets up user-defined coordinate system.
//
// xMin: x-coordinate of lower left corner of canvas.
// xMax: x-coordinate of upper right corner of canvas.
// yMin: y-coordinate of lower left corner of canvas.
// yMax: y-coordinate of upper right corner of canvas.
func (scr *Canvas) SetWorldCoordinates(xMin, xMax float64, yMin, yMax float64) error {
	if xMax <= xMin {
		return fmt.Errorf("xMax must be greater then xMin")
	}

	if yMax <= yMin {
		return fmt.Errorf("yMax must be greater then yMin")
	}

	scr.xMin = xMin
	scr.xMax = xMax
	scr.yMin = yMin
	scr.yMax = yMax
	mapWorldCoordinates(scr.dc, xMin, xMax, yMin, yMax)

	return nil
}

// SavePNG saves the picture as PNG to the specified filename.
func (scr *Canvas) SavePNG(path string) error {
	return scr.dc.SavePNG(path)
}

func mapWorldCoordinates(dc *gg.Context, xMin, xMax, yMin, yMax float64) {
	w, h := float64(dc.Width()), float64(dc.Height())

	displayAspect := math.Abs(h / w)
	windowAspect := math.Abs((yMax - yMin) / (xMax - xMin))

	if displayAspect > windowAspect {
		// Expand the viewport vertically.
		excess := (yMax - yMin) * (displayAspect/windowAspect - 1)
		yMax = yMax + excess/2
		yMin = yMin - excess/2
	} else if displayAspect < windowAspect {
		// Expand the viewport vertically.
		excess := (xMax - xMin) * (windowAspect/displayAspect - 1)
		xMax = xMax + excess/2
		xMin = xMin - excess/2
	}

	sx, sy := w/(xMax-xMin), h/(yMin-yMax)
	tx, ty := -xMin, -yMax

	dc.Identity()
	dc.Scale(sx, sy)
	dc.Translate(tx, ty)
}
