package gg

import (
	"image"
	"image/color"

	"github.com/golang/freetype/truetype"
)

type LineCap int

const (
	LineCapRound LineCap = iota
	LineCapButt
	LineCapSquare
)

type LineJoin int

const (
	LineJoinRound LineJoin = iota
	LineJoinBevel
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// GraphicContext describes the interface for the various backends (images, pdf, opengl, ...)
type GraphicContext interface {
	// Width returns the width of the context
	Width() float64
	// Height returns the height of the context
	Height() float64

	// BeginPath creates a new path
	BeginPath()
	// MoveTo creates a new subpath that start at the specified point
	MoveTo(x, y float64)
	// LineTo adds a line to the current subpath
	LineTo(x, y float64)
	// QuadraticTo adds a quadratic Bézier curve to the current subpath
	QuadraticTo(x1, y1, x2, y2 float64)
	// ArcTo adds a circular arc to the current sub-path, using
	// the given control points and radius.
	ArcTo(x1, y1, x2, y2, radius float64)
	// Close creates a line from the current point to the last MoveTo
	// point (if not the same) and mark the path as closed so the
	// first and last lines join nicely.
	ClosePath()
	// CurrentPoint returns the current point of the current sub path
	CurrentPoint() (float64, float64, bool)

	// SetStrokeColor sets the current stroke color
	SetStrokeColor(r, g, b, a int)
	// SetFillColor sets the current fill color
	SetFillColor(r, g, b, a int)
	// SetFillRule sets the current fill rule
	SetFillRule(f FillRule)

	// SetFillStyle sets current fill style
	SetFillStyle(pattern Pattern)

	// SetStrokeStyle sets current stroke style
	SetStrokeStyle(pattern Pattern)

	// SetStrokeWeight sets the current line width
	SetStrokeWeight(lineWidth float64)
	// StrokeWeight returns the current line width
	StrokeWeight() float64

	// SetLineCap sets the current line cap
	SetLineCap(cap LineCap)
	// SetLineJoin sets the current line join
	SetLineJoin(join LineJoin)
	// SetLineDash sets the current dash
	SetLineDash(dashes ...float64)
	// SetLineDashOffset sets the initial offset into
	// the dash pattern to use when stroking dashed paths.
	SetLineDashOffset(dashOffset float64)

	// SetFontSize sets the current font size
	SetFontSize(fontSize float64)
	// FontSize returns the current font size
	FontSize() float64
	// SetFont sets the current font
	SetFont(f *truetype.Font)

	// Identity resets the current transformation matrix to the identity matrix.
	// This results in no translating, scaling, rotating, or shearing.
	Identity()
	// Rotate applies a rotation to the current transformation matrix. angle is in radian.
	Rotate(angle float64)
	// Translate applies a translation to the current transformation matrix.
	Translate(tx, ty float64)
	// Scale applies a scale to the current transformation matrix.
	Scale(sx, sy float64)
	// TransformPoint multiplies the specified point by the current matrix,
	// returning a transformed position.
	TransformPoint(x, y float64) (float64, float64)

	// Push saves the current state of the context for later retrieval. These
	// can be nested.
	Push()
	// Pop restores the last saved context state from the stack
	Pop()

	// Clear fills the current canvas with a default transparent color
	Clear()

	// Stroke strokes the paths with the color specified by SetStrokeColor
	Stroke()
	// Fill fills the paths with the color specified by SetFillColor
	Fill()
	// FillAndStroke first fills the current path and than strokes it
	FillAndStroke()

	// SetPixelColor sets the color of the specified pixel using the current color.
	SetPixelColor(c color.Color, x, y int)

	// DrawPoint draws a point
	DrawPoint(x, y float64)

	// DrawLine draws a line
	DrawLine(x1, y1, x2, y2 float64)

	// DrawEllipticalArc draws an elliptical arc
	DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64)

	// DrawImageAnchored draws the specified image at the specified anchor point
	DrawImageAnchored(im image.Image, x, y int, ax, ay float64)

	// DrawStringAnchored draws the specified text at the specified anchor point.
	// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
	// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
	DrawStringAnchored(s string, x, y, ax, ay float64)

	// MeasureString returns the rendered width and height of the specified text
	// given the current font face.
	MeasureString(s string) (w, h float64)

	// Clip updates the clipping region by intersecting the current
	// clipping region with the current path as it would be filled by dc.Fill().
	// The path is cleared after this operation.
	Clip()

	// ResetClip clears the clipping region.
	ResetClip()
}
