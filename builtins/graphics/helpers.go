package graphics

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/lucasepe/g2d/gg"
)

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func drawArc(dc gg.GraphicContext, x, y, r, angle1, angle2 float64) {
	dc.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

func drawEllipse(dc gg.GraphicContext, x, y, rx, ry float64) {
	dc.BeginPath()
	dc.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
	dc.ClosePath()
}

func drawRectangle(dc gg.GraphicContext, x, y, w, h float64) {
	dc.BeginPath()
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

func drawRoundedRectangle(dc gg.GraphicContext, x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.BeginPath()
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.DrawEllipticalArc(x2, y1, r, r, radians(270), radians(360))
	dc.LineTo(x3, y2)
	dc.DrawEllipticalArc(x2, y2, r, r, radians(0), radians(90))
	dc.LineTo(x1, y3)
	dc.DrawEllipticalArc(x1, y2, r, r, radians(90), radians(180))
	dc.LineTo(x0, y1)
	dc.DrawEllipticalArc(x1, y1, r, r, radians(180), radians(270))
	dc.ClosePath()
}

func drawRoundedRectangleExtended(dc gg.GraphicContext, x, y, w, h, tl, tr, br, bl float64) {
	// corner rounding must always be positive
	absW, absH := math.Abs(w), math.Abs(h)
	hw, hh := absW/2, absH/2

	// Clip radii
	if absW < 2*tl {
		tl = hw
	}
	if absH < 2*tl {
		tl = hh
	}
	if absW < 2*tr {
		tr = hw
	}
	if absH < 2*tr {
		tr = hh
	}
	if absW < 2*br {
		br = hw
	}
	if absH < 2*br {
		br = hh
	}
	if absW < 2*bl {
		bl = hw
	}
	if absH < 2*bl {
		bl = hh
	}

	// Draw shape
	dc.BeginPath()
	dc.MoveTo(x+tl, y)
	dc.ArcTo(x+w, y, x+w, y+h, tr)
	dc.ArcTo(x+w, y+h, x, y+h, br)
	dc.ArcTo(x, y+h, x, y, bl)
	dc.ArcTo(x, y, x+w, y, tl)
	dc.ClosePath()
}

func drawCircle(dc gg.GraphicContext, x, y, r float64) {
	dc.BeginPath()
	dc.DrawEllipticalArc(x, y, r, r, 0, 2*math.Pi)
	dc.ClosePath()
}

func drawTriangle(dc gg.GraphicContext, x1, y1, x2, y2, x3, y3 float64) {
	dc.BeginPath()
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	dc.ClosePath()
}

func drawRegularPolygon(dc gg.GraphicContext, n int, x, y, r, rotation float64) {
	angle := 2 * math.Pi / float64(n)
	rotation -= math.Pi / 2
	if n%2 == 0 {
		rotation += angle / 2
	}
	dc.BeginPath()
	for i := 0; i < n; i++ {
		a := rotation + angle*float64(i)
		dc.LineTo(x+r*math.Cos(a), y+r*math.Sin(a))
	}
	dc.ClosePath()
}

// drawQuadrilateral draws a quadrilateral, a four sided polygon.
// It is similar to a rectangle, but the angles between its edges
// are not constrained to ninety degrees.
func drawQuadrilateral(dc gg.GraphicContext, x1, y1, x2, y2, x3, y3, x4, y4 float64) {
	dc.BeginPath()
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	dc.LineTo(x4, y4)
	dc.ClosePath()
}

func drawStar(dc gg.GraphicContext, cx, cy float64, spikes int, outerRadius, innerRadius float64) {
	rot := math.Pi / 2.0 * 3
	step := math.Pi / float64(spikes)

	dc.BeginPath()
	dc.MoveTo(cx, cy-outerRadius)
	for i := 0; i < spikes; i++ {
		x := cx + math.Cos(rot)*outerRadius
		y := cy + math.Sin(rot)*outerRadius
		dc.LineTo(x, y)
		rot = rot + step

		x = cx + math.Cos(rot)*innerRadius
		y = cy + math.Sin(rot)*innerRadius
		dc.LineTo(x, y)
		rot = rot + step
	}
	dc.LineTo(cx, cy-outerRadius)
	dc.ClosePath()
}

// rotateAbout updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the specified point. Angle is specified in radians.
func rotateAbout(dc gg.GraphicContext, angle, x, y float64) {
	dc.Translate(x, y)
	dc.Rotate(angle)
	dc.Translate(-x, -y)
}

// scaleAbout updates the current matrix with a scaling factor.
// Scaling occurs about the specified point.
func scaleAbout(dc gg.GraphicContext, sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Scale(sx, sy)
	dc.Translate(-x, -y)
}

// setWorldCoordinates sets up user-defined coordinate system.
//
// xMin: x-coordinate of lower left corner of canvas.
// xMax: x-coordinate of upper right corner of canvas.
// yMin: y-coordinate of lower left corner of canvas.
// yMax: y-coordinate of upper right corner of canvas.
func setWorldCoordinates(dc gg.GraphicContext, xMin, xMax, yMin, yMax float64, xOffset, yOffset float64) {
	w := float64(dc.Width()) - 2*xOffset
	h := float64(dc.Height()) - 2*yOffset

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

	dc.Translate(xOffset, yOffset)
	dc.Scale(sx, sy)
	dc.Translate(tx, ty)
}

func parseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}
	if len(x) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	}
	if len(x) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b, &a)
	}
	return
}

func savePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}
