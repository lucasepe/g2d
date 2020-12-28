package img

import (
	"math"

	"golang.org/x/image/math/fixed"
)

// Point holds a coordinates pair
type Point struct {
	X, Y float64
}

// Fixed returns a fixed-point coordinate pair
func (a Point) Fixed() fixed.Point26_6 {
	return fixp(a.X, a.Y)
}

// Distance returns the distance from another point
func (a Point) Distance(b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

// Interpolate returns a new interpolated point
func (a Point) Interpolate(b Point, t float64) Point {
	x := a.X + (b.X-a.X)*t
	y := a.Y + (b.Y-a.Y)*t
	return Point{x, y}
}
