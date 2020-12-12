package canvas

import (
	"math"

	"github.com/fogleman/gg"
)

func RoundedRect(dc *gg.Context, x, y, w, h, tl, tr, br, bl float64) {
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
	dc.NewSubPath()
	dc.MoveTo(x+tl, y)
	ArcTo(dc, x+w, y, x+w, y+h, tr)
	ArcTo(dc, x+w, y+h, x, y+h, br)
	ArcTo(dc, x, y+h, x, y, bl)
	ArcTo(dc, x, y, x+w, y, tl)
	dc.ClosePath()
}

// https://github.com/WebKit/webkit/blob/main/Source/WebCore/platform/graphics/cairo/PathCairo.cpp#L204
func ArcTo(dc *gg.Context, x1, y1, x2, y2, radius float64) {
	// Get current point
	p0, _ := dc.GetCurrentPoint()

	// Draw only a straight line to p1 if any of the points are equal or the radius is zero
	// or the points are collinear (triangle that the points form has area of zero value).
	if (x1 == p0.X && y1 == p0.Y) || (x1 == x2 && y1 == y2) || (radius == 0) {
		dc.LineTo(x1, y1)
		return
	}

	p1p0 := gg.Point{X: p0.X - x1, Y: p0.Y - y1}
	p1p2 := gg.Point{X: x2 - x1, Y: y2 - y1}

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
		ep := gg.Point{X: (p0.X + factorMax*p1p0.X), Y: (p0.Y + factorMax*p1p0.Y)}
		dc.LineTo(ep.X, ep.Y)
		return
	}

	tangent := radius / math.Tan(math.Acos(cosPhi)/2)
	factorP1P0 := tangent / p1p0Length
	tP1P0 := gg.Point{
		X: x1 + factorP1P0*p1p0.X,
		Y: y1 + factorP1P0*p1p0.Y,
	}

	orthP1P0 := gg.Point{X: p1p0.Y, Y: -p1p0.X}
	orthP1P0Length := math.Hypot(orthP1P0.X, orthP1P0.Y)
	factorRa := radius / orthP1P0Length

	// angle between orth_p1p0 and p1p2 to get the right vector orthographic to p1p0
	cosAlpha := (orthP1P0.X*p1p2.X + orthP1P0.Y*p1p2.Y) / (orthP1P0Length * p1p2Length)
	if cosAlpha < 0 {
		orthP1P0 = gg.Point{X: -orthP1P0.X, Y: -orthP1P0.Y}
	}

	p := gg.Point{
		X: tP1P0.X + factorRa*orthP1P0.X,
		Y: tP1P0.Y + factorRa*orthP1P0.Y,
	}

	// calculate angles for addArc
	orthP1P0 = gg.Point{X: -orthP1P0.X, Y: -orthP1P0.Y}
	sa := math.Acos(orthP1P0.X / orthP1P0Length)
	if orthP1P0.Y < 0 {
		sa = 2*math.Pi - sa
	}

	factorP1P2 := tangent / p1p2Length
	tP1P2 := gg.Point{
		X: x1 + factorP1P2*p1p2.X,
		Y: y1 + factorP1P2*p1p2.Y,
	}

	orthP1P2 := gg.Point{
		X: tP1P2.X - p.X,
		Y: tP1P2.Y - p.Y,
	}
	orthP1P2Length := math.Hypot(orthP1P2.X, orthP1P2.Y)

	ea := math.Acos(orthP1P2.X / orthP1P2Length)
	if orthP1P2.Y <= 0 {
		ea = 2*math.Pi - ea
	}

	/* anticlockwise logic
	anticlockwise := false
	if (sa > ea) && ((sa - ea) < math.Pi) {
		anticlockwise = true
	}
	if (sa < ea) && ((ea - sa) > math.Pi) {
		anticlockwise = true
	}
	*/
	dc.LineTo(tP1P0.X, tP1P0.Y)

	//fmt.Println("anticlockwise =", anticlockwise, " r = ", radius, " sa = ", gg.Degrees(sa), " ea = ", gg.Degrees(ea))
	/*
		if anticlockwise && (math.Pi*2 != radius) {
			// we don't have anticlockwise draw arc
			return
		}*/

	dc.DrawArc(p.X, p.Y, radius, sa, ea)
}
