package eval

import (
	"image"
	"image/color"

	"github.com/golang/freetype/truetype"

	"github.com/lucasepe/g2d/gg"
)

// MockGraphicContext implements a graphic context for tests only
type MockGraphicContext struct{}

// Image returns the image that has been drawn by this context.
func (dc *MockGraphicContext) Image() image.Image { return nil }

// Width returns the width of the image in pixels.
func (dc *MockGraphicContext) Width() float64 { return 0 }

// Height returns the height of the image in pixels.
func (dc *MockGraphicContext) Height() float64 { return 0 }

// BeginPath starts a new subpath within the current path. There is no current
// point after this operation.
func (dc *MockGraphicContext) BeginPath() {}

// MoveTo starts a new subpath within the current path starting at the
// specified point.
func (dc *MockGraphicContext) MoveTo(x, y float64) {}

// LineTo adds a line segment to the current path starting at the current
// point. If there is no current point, it is equivalent to MoveTo(x, y)
func (dc *MockGraphicContext) LineTo(x, y float64) {}

// QuadraticTo adds a quadratic bezier curve to the current path starting at
// the current point. If there is no current point, it first performs
// MoveTo(x1, y1)
func (dc *MockGraphicContext) QuadraticTo(x1, y1, x2, y2 float64) {}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
func (dc *MockGraphicContext) ClosePath() {}

// CurrentPoint returns the current point and if there is a current point.
// The point will have been transformed by the context's transformation matrix.
func (dc *MockGraphicContext) CurrentPoint() (float64, float64, bool) { return 0, 0, false }

// SetStrokeColor sets the current stroke color. r, g, b, a
// values should be between 0 and 255, inclusive.
func (dc *MockGraphicContext) SetStrokeColor(r, g, b, a int) {}

// SetFillColor sets the current stroke color. r, g, b, a values should be between 0 and
// 255, inclusive.
func (dc *MockGraphicContext) SetFillColor(r, g, b, a int) {}

// SetFillRule sets the current fill rule
func (dc *MockGraphicContext) SetFillRule(fillRule gg.FillRule) {}

// ArcTo adds a circular arc to the current sub-path, using
// the given control points and radius.
func (dc *MockGraphicContext) ArcTo(x1, y1, x2, y2, radius float64) {}

// SetFillStyle sets current fill style
func (dc *MockGraphicContext) SetFillStyle(pattern gg.Pattern) {}

// SetStrokeStyle sets current stroke style
func (dc *MockGraphicContext) SetStrokeStyle(pattern gg.Pattern) {}

// SetStrokeWeight sets the lineWidth.
func (dc *MockGraphicContext) SetStrokeWeight(lineWidth float64) {}

// StrokeWeight returns the current line width
func (dc *MockGraphicContext) StrokeWeight() float64 { return 0 }

// SetLineCap sets the current line cap
func (dc *MockGraphicContext) SetLineCap(lineCap gg.LineCap) {}

// SetLineJoin sets the current line join
func (dc *MockGraphicContext) SetLineJoin(lineJoin gg.LineJoin) {}

// SetLineDash sets the current dash
func (dc *MockGraphicContext) SetLineDash(dashes ...float64) {}

// SetLineDashOffset sets the initial offset into the dash pattern
func (dc *MockGraphicContext) SetLineDashOffset(offset float64) {}

// SetFontSize sets the current font size
func (dc *MockGraphicContext) SetFontSize(points float64) {}

// FontSize returns font's size
func (dc *MockGraphicContext) FontSize() float64 { return 0 }

// SetFont sets the current font
func (dc *MockGraphicContext) SetFont(font *truetype.Font) {}

// Push saves the current state of the context for later retrieval. These
// can be nested.
func (dc *MockGraphicContext) Push() {}

// Pop restores the last saved context state from the stack.
func (dc *MockGraphicContext) Pop() {}

// Clear fills the entire image with the current fill color.
func (dc *MockGraphicContext) Clear() {}

// Stroke strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is cleared after this
// operation.
func (dc *MockGraphicContext) Stroke() {}

// FillPreserve fills the current path with the current color. Open subpaths
// are implicity closed. The path is preserved after this operation.
func (dc *MockGraphicContext) FillPreserve() {}

// Fill fills the current path with the current color. Open subpaths
// are implicity closed. The path is cleared after this operation.
func (dc *MockGraphicContext) Fill() {}

// FillAndStroke first fills the paths and than strokes them
func (dc *MockGraphicContext) FillAndStroke() {}

// ClearPath clears the current path. There is no current point after this
// operation.
func (dc *MockGraphicContext) ClearPath() {}

// ClipPreserve updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is preserved after this operation.
func (dc *MockGraphicContext) ClipPreserve() {}

// Clip updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is cleared after this operation.
func (dc *MockGraphicContext) Clip() {}

// ResetClip clears the clipping region.
func (dc *MockGraphicContext) ResetClip() {}

// Transformation Matrix Operations

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func (dc *MockGraphicContext) Identity() {}

// Translate updates the current matrix with a translation.
func (dc *MockGraphicContext) Translate(x, y float64) {}

// Scale updates the current matrix with a scaling factor.
// Scaling occurs about the origin.
func (dc *MockGraphicContext) Scale(x, y float64) {}

// Rotate updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the origin. Angle is specified in radians.
func (dc *MockGraphicContext) Rotate(angle float64) {}

// TransformPoint multiplies the specified point by the current matrix,
// returning a transformed position.
func (dc *MockGraphicContext) TransformPoint(x, y float64) (tx, ty float64) { return 0, 0 }

// Convenient Drawing Functions

// SetPixelColor sets the color of the specified pixel using the current color.
func (dc *MockGraphicContext) SetPixelColor(c color.Color, x, y int) {}

// DrawPoint draws a point
func (dc *MockGraphicContext) DrawPoint(x, y float64) {}

// DrawLine draws a line
func (dc *MockGraphicContext) DrawLine(x1, y1, x2, y2 float64) {}

// DrawEllipticalArc draws an elliptical arc
func (dc *MockGraphicContext) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {}

// DrawImageAnchored draws the specified image at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image. Use ax=0.5, ay=0.5 to center the image at the specified point.
func (dc *MockGraphicContext) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {}

// DrawStringAnchored draws the specified text at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
func (dc *MockGraphicContext) DrawStringAnchored(s string, x, y, ax, ay float64) {}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func (dc *MockGraphicContext) MeasureString(s string) (w, h float64) { return 0, 0 }
