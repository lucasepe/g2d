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
	maximumSize     = 1536
	defaultWidth    = 1024
	defaultHeight   = 1024
	defaultFontSize = 12.0
)

// Canvas defines a drawing context.
type Canvas struct {
	width    int
	height   int
	font     *truetype.Font
	fontSize float64
	dc       *gg.Context
}

// NewCanvas creates a drawing context.
func NewCanvas(width, height int) *Canvas {
	res := &Canvas{
		width:    width,
		height:   height,
		fontSize: defaultFontSize,
		dc:       gg.NewContext(width, height),
	}

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

	return res
}

// Graphics returns the canvas graphic context.
func (scr *Canvas) Graphics() *gg.Context { return scr.dc }

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

// Size returns the canvas size.
func (scr *Canvas) Size() (int, int) { return scr.width, scr.height }

// Reset resets screen using the specified options.
func (scr *Canvas) Reset(width, height int) {
	scr.dc = gg.NewContext(width, height)
	scr.width = width
	scr.height = height
	scr.SetFontSize(defaultFontSize)
}

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
