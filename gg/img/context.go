package img

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/f64"

	"github.com/lucasepe/g2d/gg"
)

// Context implements the graphic context for drawing on images
type Context struct {
	rasterizer  *raster.Rasterizer
	im          *image.RGBA
	mask        *image.Alpha
	fillColor   color.Color
	fillPattern gg.Pattern
	fillPath    raster.Path

	strokeColor   color.Color
	strokePattern gg.Pattern
	strokePath    raster.Path

	start      Point
	current    Point
	hasCurrent bool
	dashes     []float64
	dashOffset float64
	lineWidth  float64
	lineCap    gg.LineCap
	lineJoin   gg.LineJoin
	fillRule   gg.FillRule
	font       *truetype.Font
	fontSize   float64
	matrix     gg.Matrix
	stack      []*Context
}

// NewContextForRGBA prepares a context for rendering onto the specified image.
// No copy is made.
func NewContextForRGBA(im *image.RGBA) *Context {
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y

	res := &Context{
		rasterizer:    raster.NewRasterizer(w, h),
		im:            im,
		fillColor:     color.Transparent,
		fillPattern:   gg.NewSolidPattern(color.White),
		strokeColor:   color.Black,
		strokePattern: gg.NewSolidPattern(color.Black),
		lineWidth:     1,
		fillRule:      gg.FillRuleWinding,
		fontSize:      14,
		matrix:        gg.Identity(),
	}

	return res
}

// Image returns the image that has been drawn by this context.
func (dc *Context) Image() image.Image { return dc.im }

// Width returns the width of the image in pixels.
func (dc *Context) Width() float64 {
	res := dc.im.Bounds().Size().X
	return float64(res)
}

// Height returns the height of the image in pixels.
func (dc *Context) Height() float64 {
	res := dc.im.Bounds().Size().Y
	return float64(res)
}

// BeginPath starts a new subpath within the current path. There is no current
// point after this operation.
func (dc *Context) BeginPath() {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	dc.hasCurrent = false
}

// MoveTo starts a new subpath within the current path starting at the
// specified point.
func (dc *Context) MoveTo(x, y float64) {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	x, y = dc.TransformPoint(x, y)
	p := Point{x, y}
	dc.strokePath.Start(p.Fixed())
	dc.fillPath.Start(p.Fixed())
	dc.start = p
	dc.current = p
	dc.hasCurrent = true
}

// LineTo adds a line segment to the current path starting at the current
// point. If there is no current point, it is equivalent to MoveTo(x, y)
func (dc *Context) LineTo(x, y float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x, y)
	} else {
		x, y = dc.TransformPoint(x, y)
		p := Point{x, y}
		dc.strokePath.Add1(p.Fixed())
		dc.fillPath.Add1(p.Fixed())
		dc.current = p
	}
}

// QuadraticTo adds a quadratic bezier curve to the current path starting at
// the current point. If there is no current point, it first performs
// MoveTo(x1, y1)
func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	p1 := Point{x1, y1}
	p2 := Point{x2, y2}
	dc.strokePath.Add2(p1.Fixed(), p2.Fixed())
	dc.fillPath.Add2(p1.Fixed(), p2.Fixed())
	dc.current = p2
}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
func (dc *Context) ClosePath() {
	if dc.hasCurrent {
		dc.strokePath.Add1(dc.start.Fixed())
		dc.fillPath.Add1(dc.start.Fixed())
		dc.current = dc.start
	}
}

// CurrentPoint returns the current point and if there is a current point.
// The point will have been transformed by the context's transformation matrix.
func (dc *Context) CurrentPoint() (float64, float64, bool) {
	if dc.hasCurrent {
		return dc.current.X, dc.current.Y, true
	}
	return 0, 0, false
}

// SetStrokeColor sets the current stroke color. r, g, b, a
// values should be between 0 and 255, inclusive.
func (dc *Context) SetStrokeColor(r, g, b, a int) {
	dc.strokeColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	dc.strokePattern = gg.NewSolidPattern(dc.strokeColor)
}

// SetFillColor sets the current stroke color. r, g, b, a values should be between 0 and
// 255, inclusive.
func (dc *Context) SetFillColor(r, g, b, a int) {
	dc.fillColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	dc.fillPattern = gg.NewSolidPattern(dc.fillColor)
}

// SetFillRule sets the current fill rule
func (dc *Context) SetFillRule(fillRule gg.FillRule) {
	dc.fillRule = fillRule
}

// ArcTo adds a circular arc to the current sub-path, using
// the given control points and radius.
// The arc is automatically connected to the path's latest
// point with a straight line, if necessary for the specified parameters.
//
// This method is commonly used for making rounded corners.
// https://github.com/WebKit/webkit/blob/main/Source/WebCore/platform/graphics/cairo/PathCairo.cpp#L204
func (dc *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	// Get current point
	x0, y0, _ := dc.CurrentPoint()

	// Draw only a straight line to p1 if any of the points are equal or the radius is zero
	// or the points are collinear (triangle that the points form has area of zero value).
	if (x1 == x0 && y1 == y0) || (x1 == x2 && y1 == y2) || (radius == 0) {
		dc.LineTo(x1, y1)
		return
	}

	p1p0 := Point{X: x0 - x1, Y: y0 - y1}
	p1p2 := Point{X: x2 - x1, Y: y2 - y1}

	p1p0Length := math.Hypot(p1p0.X, p1p0.Y)
	p1p2Length := math.Hypot(p1p2.X, p1p2.Y)

	cosPhi := (p1p0.X*p1p2.X + p1p0.Y*p1p2.Y) / (p1p0Length * p1p2Length)
	// all points on a line logic
	if cosPhi == -1 {
		dc.LineTo(x1, y1)
		return
	}

	if cosPhi == 1 {
		// add infinite far away point
		maxLength := 65535.0
		factorMax := maxLength / p1p0Length
		ep := Point{
			X: x0 + factorMax*p1p0.X,
			Y: y0 + factorMax*p1p0.Y,
		}
		dc.LineTo(ep.X, ep.Y)
		return
	}

	tangent := radius / math.Tan(math.Acos(cosPhi)/2)
	factorP1P0 := tangent / p1p0Length
	tP1P0 := Point{
		X: x1 + factorP1P0*p1p0.X,
		Y: y1 + factorP1P0*p1p0.Y,
	}

	orthP1P0 := Point{X: p1p0.Y, Y: -p1p0.X}
	orthP1P0Length := math.Hypot(orthP1P0.X, orthP1P0.Y)
	factorRa := radius / orthP1P0Length

	// angle between orth_p1p0 and p1p2 to get the right vector orthographic to p1p0
	cosAlpha := (orthP1P0.X*p1p2.X + orthP1P0.Y*p1p2.Y) / (orthP1P0Length * p1p2Length)
	if cosAlpha < 0 {
		orthP1P0 = Point{X: -orthP1P0.X, Y: -orthP1P0.Y}
	}

	p := Point{
		X: tP1P0.X + factorRa*orthP1P0.X,
		Y: tP1P0.Y + factorRa*orthP1P0.Y,
	}

	// calculate angles for addArc
	orthP1P0 = Point{X: -orthP1P0.X, Y: -orthP1P0.Y}
	sa := math.Acos(orthP1P0.X / orthP1P0Length)
	if orthP1P0.Y < 0 {
		sa = 2*math.Pi - sa
	}

	factorP1P2 := tangent / p1p2Length
	tP1P2 := Point{
		X: x1 + factorP1P2*p1p2.X,
		Y: y1 + factorP1P2*p1p2.Y,
	}

	orthP1P2 := Point{
		X: tP1P2.X - p.X,
		Y: tP1P2.Y - p.Y,
	}
	orthP1P2Length := math.Hypot(orthP1P2.X, orthP1P2.Y)

	ea := math.Acos(orthP1P2.X / orthP1P2Length)
	if orthP1P2.Y <= 0 {
		ea = 2*math.Pi - ea
	}

	dc.LineTo(tP1P0.X, tP1P0.Y)
	dc.DrawEllipticalArc(p.X, p.Y, radius, radius, sa, ea)
}

// SetFillStyle sets current fill style
func (dc *Context) SetFillStyle(pattern gg.Pattern) {
	// if pattern is SolidPattern, also change dc.color(for dc.Clear, dc.drawString)
	if fillStyle, ok := pattern.(*gg.SolidPattern); ok {
		dc.fillColor = fillStyle.Color
	}
	dc.fillPattern = pattern
}

// SetStrokeStyle sets current stroke style
func (dc *Context) SetStrokeStyle(pattern gg.Pattern) {
	dc.strokePattern = pattern
}

// SetStrokeWeight sets the lineWidth.
func (dc *Context) SetStrokeWeight(lineWidth float64) { dc.lineWidth = lineWidth }

// StrokeWeight returns the current line width
func (dc *Context) StrokeWeight() float64 { return dc.lineWidth }

// SetLineCap sets the current line cap
func (dc *Context) SetLineCap(lineCap gg.LineCap) { dc.lineCap = lineCap }

// SetLineJoin sets the current line join
func (dc *Context) SetLineJoin(lineJoin gg.LineJoin) { dc.lineJoin = lineJoin }

// SetLineDash sets the current dash
func (dc *Context) SetLineDash(dashes ...float64) {
	dc.dashes = dashes
}

// SetLineDashOffset sets the initial offset into the dash pattern
func (dc *Context) SetLineDashOffset(offset float64) {
	dc.dashOffset = offset
}

// SetFontSize sets the current font size
func (dc *Context) SetFontSize(points float64) {
	dc.fontSize = points * 72 / 96
}

// FontSize returns font's size
func (dc *Context) FontSize() float64 {
	return dc.fontSize
}

// SetFont sets the current font
func (dc *Context) SetFont(font *truetype.Font) {
	dc.font = font
}

// Push saves the current state of the context for later retrieval. These
// can be nested.
func (dc *Context) Push() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Pop restores the last saved context state from the stack.
func (dc *Context) Pop() {
	if dc.stack == nil {
		return
	}
	before := *dc
	s := dc.stack
	x, s := s[len(s)-1], s[:len(s)-1]
	*dc = *x
	dc.mask = before.mask
	dc.strokePath = before.strokePath
	dc.fillPath = before.fillPath
	dc.start = before.start
	dc.current = before.current
	dc.hasCurrent = before.hasCurrent
}

// Clear fills the entire image with the current fill color.
func (dc *Context) Clear() {
	src := image.NewUniform(dc.fillColor)
	draw.Draw(dc.im, dc.im.Bounds(), src, image.ZP, draw.Src)
}

// Stroke strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is cleared after this
// operation.
func (dc *Context) Stroke() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.strokePattern.(*gg.SolidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.Color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.strokePattern)
	}
	dc.stroke(painter)

	dc.ClearPath()
}

// FillPreserve fills the current path with the current color. Open subpaths
// are implicity closed. The path is preserved after this operation.
func (dc *Context) FillPreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.fillPattern.(*gg.SolidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.Color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.fillPattern)
	}
	dc.fill(painter)
}

// Fill fills the current path with the current color. Open subpaths
// are implicity closed. The path is cleared after this operation.
func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.ClearPath()
}

// FillAndStroke first fills the paths and than strokes them
func (dc *Context) FillAndStroke() {
	dc.FillPreserve()
	dc.Stroke()
}

// ClearPath clears the current path. There is no current point after this
// operation.
func (dc *Context) ClearPath() {
	dc.strokePath.Clear()
	dc.fillPath.Clear()
	dc.hasCurrent = false
}

// ClipPreserve updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is preserved after this operation.
func (dc *Context) ClipPreserve() {
	width, height := dc.im.Bounds().Size().X, dc.im.Bounds().Size().Y

	clip := image.NewAlpha(image.Rect(0, 0, width, height))
	painter := raster.NewAlphaOverPainter(clip)
	dc.fill(painter)
	if dc.mask == nil {
		dc.mask = clip
	} else {
		mask := image.NewAlpha(image.Rect(0, 0, width, height))
		draw.DrawMask(mask, mask.Bounds(), clip, image.ZP, dc.mask, image.ZP, draw.Over)
		dc.mask = mask
	}
}

// Clip updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is cleared after this operation.
func (dc *Context) Clip() {
	dc.ClipPreserve()
	dc.ClearPath()
}

// ResetClip clears the clipping region.
func (dc *Context) ResetClip() {
	dc.mask = nil
}

// Transformation Matrix Operations

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func (dc *Context) Identity() {
	dc.matrix = gg.Identity()
}

// Translate updates the current matrix with a translation.
func (dc *Context) Translate(x, y float64) {
	dc.matrix = dc.matrix.Translate(x, y)
}

// Scale updates the current matrix with a scaling factor.
// Scaling occurs about the origin.
func (dc *Context) Scale(x, y float64) {
	dc.matrix = dc.matrix.Scale(x, y)
}

// Rotate updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the origin. Angle is specified in radians.
func (dc *Context) Rotate(angle float64) {
	dc.matrix = dc.matrix.Rotate(angle)
}

// TransformPoint multiplies the specified point by the current matrix,
// returning a transformed position.
func (dc *Context) TransformPoint(x, y float64) (tx, ty float64) {
	return dc.matrix.TransformPoint(x, y)
}

// Convenient Drawing Functions

// SetPixelColor sets the color of the specified pixel using the current color.
func (dc *Context) SetPixelColor(c color.Color, x, y int) {
	dc.im.Set(x, y, c)
}

// DrawPoint draws a point
func (dc *Context) DrawPoint(x, y float64) {
	r := math.Max(1, dc.lineWidth)

	dc.Push()
	tx, ty := dc.TransformPoint(x, y)
	dc.Identity()

	dc.BeginPath()
	dc.DrawEllipticalArc(tx, ty, r, r, 0, 2*math.Pi)
	dc.ClosePath()

	dc.Pop()
}

// DrawLine draws a line
func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

// DrawEllipticalArc draws an elliptical arc
func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := 0; i < n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := angle1 + (angle2-angle1)*p1
		a2 := angle1 + (angle2-angle1)*p2
		x0 := x + rx*math.Cos(a1)
		y0 := y + ry*math.Sin(a1)
		x1 := x + rx*math.Cos((a1+a2)/2)
		y1 := y + ry*math.Sin((a1+a2)/2)
		x2 := x + rx*math.Cos(a2)
		y2 := y + ry*math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			if dc.hasCurrent {
				dc.LineTo(x0, y0)
			} else {
				dc.MoveTo(x0, y0)
			}
		}
		dc.QuadraticTo(cx, cy, x2, y2)
	}
}

// DrawImageAnchored draws the specified image at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image. Use ax=0.5, ay=0.5 to center the image at the specified point.
func (dc *Context) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	s := im.Bounds().Size()
	x -= int(ax * float64(s.X))
	y -= int(ay * float64(s.Y))
	transformer := draw.BiLinear
	fx, fy := float64(x), float64(y)
	m := dc.matrix.Translate(fx, fy)
	s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
	if dc.mask == nil {
		transformer.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, nil)
	} else {
		transformer.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, &draw.Options{
			DstMask:  dc.mask,
			DstMaskP: image.ZP,
		})
	}
}

// DrawStringAnchored draws the specified text at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64) {
	width, height := dc.im.Bounds().Size().X, dc.im.Bounds().Size().Y

	w, h := dc.MeasureString(s)
	x -= ax * w
	y += ay * h
	if dc.mask == nil {
		dc.drawString(dc.im, s, x, y)
	} else {
		im := image.NewRGBA(image.Rect(0, 0, width, height))
		dc.drawString(im, s, x, y)
		draw.DrawMask(dc.im, dc.im.Bounds(), im, image.ZP, dc.mask, image.ZP, draw.Over)
	}
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func (dc *Context) MeasureString(s string) (w, h float64) {

	ff, err := dc.currentFontFace()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning img.Context MeasureString error: %s", err.Error())
		return 0, 0
	}

	d := &font.Drawer{Face: ff}
	a := d.MeasureString(s)
	return float64(a >> 6), dc.fontSize
}

// Path Drawing

func (dc *Context) capper() raster.Capper {
	switch dc.lineCap {
	case gg.LineCapButt:
		return raster.ButtCapper
	case gg.LineCapRound:
		return raster.RoundCapper
	case gg.LineCapSquare:
		return raster.SquareCapper
	}
	return nil
}

func (dc *Context) joiner() raster.Joiner {
	switch dc.lineJoin {
	case gg.LineJoinBevel:
		return raster.BevelJoiner
	case gg.LineJoinRound:
		return raster.RoundJoiner
	}
	return nil
}

func (dc *Context) stroke(painter raster.Painter) {
	path := dc.strokePath
	if len(dc.dashes) > 0 {
		path = dashed(path, dc.dashes, dc.dashOffset)
	} else {
		// TODO: this is a temporary workaround to remove tiny segments
		// that result in rendering issues
		path = rasterPath(flattenPath(path))
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = true
	r.Clear()
	r.AddStroke(path, fix(dc.lineWidth), dc.capper(), dc.joiner())
	r.Rasterize(painter)
}

func (dc *Context) fill(painter raster.Painter) {
	path := dc.fillPath
	if dc.hasCurrent {
		path = make(raster.Path, len(dc.fillPath))
		copy(path, dc.fillPath)
		path.Add1(dc.start.Fixed())
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = dc.fillRule == gg.FillRuleWinding
	r.Clear()
	r.AddPath(path)
	r.Rasterize(painter)
}

// Text drawings

func (dc *Context) currentFontFace() (font.Face, error) {
	if dc.font == nil {
		val, err := truetype.Parse(gomono.TTF)
		if err != nil {
			return nil, err
		}
		dc.font = val
	}

	res := truetype.NewFace(dc.font, &truetype.Options{Size: dc.fontSize})
	return res, nil
}

func (dc *Context) drawString(im *image.RGBA, s string, x, y float64) {
	ff, err := dc.currentFontFace()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning img.Context drawString error: %s", err.Error())
		return
	}

	d := &font.Drawer{
		Dst:  im,
		Src:  image.NewUniform(dc.strokeColor),
		Face: ff,
		Dot:  fixp(x, y),
	}
	// based on Drawer.DrawString() in golang.org/x/image/font/font.go
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			// TODO: is falling back on the U+FFFD glyph the responsibility of
			// the Drawer or the Face?
			// TODO: set prevC = '\ufffd'?
			continue
		}
		sr := dr.Sub(dr.Min)
		transformer := draw.BiLinear
		fx, fy := float64(dr.Min.X), float64(dr.Min.Y)
		m := dc.matrix.Translate(fx, fy)
		s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		transformer.Transform(d.Dst, s2d, d.Src, sr, draw.Over, &draw.Options{
			SrcMask:  mask,
			SrcMaskP: maskp,
		})
		d.Dot.X += advance
		prevC = c
	}
}
