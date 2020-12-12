package geometry

import (
	"fmt"
	"math"
)

// Line is a 2D line segment, between points A and B.
type Line struct {
	A, B Vec
}

// L creates and returns a new Line.
func L(from, to Vec) Line {
	return Line{
		A: from,
		B: to,
	}
}

// Center will return the point at center of the line; that is, the point equidistant from either end.
func (l Line) Center() Vec {
	return l.A.Add(l.A.To(l.B).Scaled(0.5))
}

// Closest will return the point on the line which is closest to the Vec provided.
func (l Line) Closest(v Vec) Vec {
	// between is a helper function which determines whether x is greater than min(a, b) and less than max(a, b)
	between := func(a, b, x float64) bool {
		min := math.Min(a, b)
		max := math.Max(a, b)
		return min < x && x < max
	}

	// Closest point will be on a line which perpendicular to this line.
	// If and only if the infinite perpendicular line intersects the segment.
	m, b := l.Formula()

	// Account for horizontal lines
	if m == 0 {
		x := v.X
		y := l.A.Y

		// check if the X coordinate of v is on the line
		if between(l.A.X, l.B.X, v.X) {
			return V(x, y)
		}

		// Otherwise get the closest endpoint
		if l.A.To(v).Len() < l.B.To(v).Len() {
			return l.A
		}
		return l.B
	}

	// Account for vertical lines
	if math.IsInf(math.Abs(m), 1) {
		x := l.A.X
		y := v.Y

		// check if the Y coordinate of v is on the line
		if between(l.A.Y, l.B.Y, v.Y) {
			return V(x, y)
		}

		// Otherwise get the closest endpoint
		if l.A.To(v).Len() < l.B.To(v).Len() {
			return l.A
		}
		return l.B
	}

	perpendicularM := -1 / m
	perpendicularB := v.Y - (perpendicularM * v.X)

	// Coordinates of intersect (of infinite lines)
	x := (perpendicularB - b) / (m - perpendicularM)
	y := m*x + b

	// Check if the point lies between the x and y bounds of the segment
	if !between(l.A.X, l.B.X, x) && !between(l.A.Y, l.B.Y, y) {
		// Not within bounding box
		toStart := v.To(l.A)
		toEnd := v.To(l.B)

		if toStart.Len() < toEnd.Len() {
			return l.A
		}
		return l.B
	}

	return V(x, y)
}

// Contains returns whether the provided Vec lies on the line.
func (l Line) Contains(v Vec) bool {
	return l.Closest(v).Eq(v)
}

// Formula will return the values that represent the line in the formula: y = mx + b
// This function will return math.Inf+, math.Inf- for a vertical line.
func (l Line) Formula() (m, b float64) {
	// Account for horizontal lines
	if l.B.Y == l.A.Y {
		return 0, l.A.Y
	}

	m = (l.B.Y - l.A.Y) / (l.B.X - l.A.X)
	b = l.A.Y - (m * l.A.X)

	return m, b
}

// Intersect will return the point of intersection for the two line segments.  If the line segments do not intersect,
// this function will return the zero-vector and false.
func (l Line) Intersect(k Line) (Vec, bool) {
	// Check if the lines are parallel
	lDir := l.A.To(l.B)
	kDir := k.A.To(k.B)
	if lDir.X == kDir.X && lDir.Y == kDir.Y {
		return ZV, false
	}

	// The lines intersect - but potentially not within the line segments.
	// Get the intersection point for the lines if they were infinitely long, check if the point exists on both of the
	// segments
	lm, lb := l.Formula()
	km, kb := k.Formula()

	// Account for vertical lines
	if math.IsInf(math.Abs(lm), 1) && math.IsInf(math.Abs(km), 1) {
		// Both vertical, therefore parallel
		return ZV, false
	}

	var x, y float64

	if math.IsInf(math.Abs(lm), 1) || math.IsInf(math.Abs(km), 1) {
		// One line is vertical
		intersectM := lm
		intersectB := lb
		verticalLine := k

		if math.IsInf(math.Abs(lm), 1) {
			intersectM = km
			intersectB = kb
			verticalLine = l
		}

		y = intersectM*verticalLine.A.X + intersectB
		x = verticalLine.A.X
	} else {
		// Coordinates of intersect
		x = (kb - lb) / (lm - km)
		y = lm*x + lb
	}

	if l.Contains(V(x, y)) && k.Contains(V(x, y)) {
		// The intersect point is on both line segments, they intersect.
		return V(x, y), true
	}

	return ZV, false
}

// IntersectCircle will return the shortest Vec such that moving the Line by that Vec will cause the Line and Circle
// to no longer intesect.  If they do not intersect at all, this function will return a zero-vector.
func (l Line) IntersectCircle(c Circle) Vec {
	// Get the point on the line closest to the center of the circle.
	closest := l.Closest(c.Center)
	cirToClosest := c.Center.To(closest)

	if cirToClosest.Len() >= c.Radius {
		return ZV
	}

	return cirToClosest.Scaled(cirToClosest.Len() - c.Radius)
}

// Len returns the length of the line segment.
func (l Line) Len() float64 {
	return l.A.To(l.B).Len()
}

// Moved will return a line moved by the delta Vec provided.
func (l Line) Moved(delta Vec) Line {
	return Line{
		A: l.A.Add(delta),
		B: l.B.Add(delta),
	}
}

// Rotated will rotate the line around the provided Vec.
func (l Line) Rotated(around Vec, angle float64) Line {
	// Move the line so we can use `Vec.Rotated`
	lineShifted := l.Moved(around.Scaled(-1))

	lineRotated := Line{
		A: lineShifted.A.Rotated(angle),
		B: lineShifted.B.Rotated(angle),
	}

	return lineRotated.Moved(around)
}

// Scaled will return the line scaled around the center point.
func (l Line) Scaled(scale float64) Line {
	return l.ScaledXY(l.Center(), scale)
}

// ScaledXY will return the line scaled around the Vec provided.
func (l Line) ScaledXY(around Vec, scale float64) Line {
	toA := around.To(l.A).Scaled(scale)
	toB := around.To(l.B).Scaled(scale)

	return Line{
		A: around.Add(toA),
		B: around.Add(toB),
	}
}

func (l Line) String() string {
	return fmt.Sprintf("Line(%v, %v)", l.A, l.B)
}
