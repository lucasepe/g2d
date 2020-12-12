package geometry

import (
	"fmt"
	"math"
)

// Vec is a 2D vector type with X and Y coordinates.
//
// Create vectors with the V constructor:
//
//   u := geometry.V(1, 2)
//   v := geometry.V(8, -3)
//
// Use various methods to manipulate them:
//
//   w := u.Add(v)
//   fmt.Println(w)        // Vec(9, -1)
//   fmt.Println(u.Sub(v)) // Vec(-7, 5)
//   u = geometry.V(2, 3)
//   v = geometry.V(8, 1)
//   if u.X < 0 {
//	     fmt.Println("this won't happen")
//   }
//   x := u.Unit().Dot(v.Unit())
type Vec struct {
	X, Y float64
}

// ZV is a zero vector.
var ZV = Vec{0, 0}

// V returns a new 2D vector with the given coordinates.
func V(x, y float64) Vec {
	return Vec{x, y}
}

// Eq will compare two vectors and return whether they are equal accounting for rounding errors.  At worst, the result
// is correct to 7 significant digits.
func (u Vec) Eq(v Vec) bool {
	return nearlyEqual(u.X, v.X) && nearlyEqual(u.Y, v.Y)
}

// Unit returns a vector of length 1 facing the given angle.
func Unit(angle float64) Vec {
	return Vec{1, 0}.Rotated(angle)
}

// String returns the string representation of the vector u.
//
//   u := geometry.V(4.5, -1.3)
//   u.String()     // returns "Vec(4.5, -1.3)"
//   fmt.Println(u) // Vec(4.5, -1.3)
func (u Vec) String() string {
	return fmt.Sprintf("Vec(%v, %v)", u.X, u.Y)
}

// XY returns the components of the vector in two return values.
func (u Vec) XY() (x, y float64) {
	return u.X, u.Y
}

// Add returns the sum of vectors u and v.
func (u Vec) Add(v Vec) Vec {
	return Vec{
		u.X + v.X,
		u.Y + v.Y,
	}
}

// Sub returns the difference betweeen vectors u and v.
func (u Vec) Sub(v Vec) Vec {
	return Vec{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// Floor converts x and y to their integer equivalents.
func (u Vec) Floor() Vec {
	return Vec{
		math.Floor(u.X),
		math.Floor(u.Y),
	}
}

// To returns the vector from u to v. Equivalent to v.Sub(u).
func (u Vec) To(v Vec) Vec {
	return Vec{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// Scaled returns the vector u multiplied by c.
func (u Vec) Scaled(c float64) Vec {
	return Vec{u.X * c, u.Y * c}
}

// ScaledXY returns the vector u multiplied by the vector v component-wise.
func (u Vec) ScaledXY(v Vec) Vec {
	return Vec{u.X * v.X, u.Y * v.Y}
}

// Len returns the length of the vector u.
func (u Vec) Len() float64 {
	return math.Hypot(u.X, u.Y)
}

// Angle returns the angle between the vector u and the x-axis. The result is in range [-Pi, Pi].
func (u Vec) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

// Unit returns a vector of length 1 facing the direction of u (has the same angle).
func (u Vec) Unit() Vec {
	if u.X == 0 && u.Y == 0 {
		return Vec{1, 0}
	}
	return u.Scaled(1 / u.Len())
}

// Rotated returns the vector u rotated by the given angle in radians.
func (u Vec) Rotated(angle float64) Vec {
	sin, cos := math.Sincos(angle)
	return Vec{
		u.X*cos - u.Y*sin,
		u.X*sin + u.Y*cos,
	}
}

// Normal returns a vector normal to u. Equivalent to u.Rotated(math.Pi / 2), but faster.
func (u Vec) Normal() Vec {
	return Vec{-u.Y, u.X}
}

// Dot returns the dot product of vectors u and v.
func (u Vec) Dot(v Vec) float64 {
	return u.X*v.X + u.Y*v.Y
}

// Cross return the cross product of vectors u and v.
func (u Vec) Cross(v Vec) float64 {
	return u.X*v.Y - v.X*u.Y
}

// Project returns a projection (or component) of vector u in the direction of vector v.
//
// Behaviour is undefined if v is a zero vector.
func (u Vec) Project(v Vec) Vec {
	len := u.Dot(v) / v.Len()
	return v.Unit().Scaled(len)
}

// Map applies the function f to both x and y components of the vector u and returns the modified
// vector.
//
//   u := geometry.V(10.5, -1.5)
//   v := u.Map(math.Floor)   // v is Vec(10, -2), both components of u floored
func (u Vec) Map(f func(float64) float64) Vec {
	return Vec{
		f(u.X),
		f(u.Y),
	}
}

// Lerp returns a linear interpolation between vectors a and b.
//
// This function basically returns a point along the line between a and b and t chooses which one.
// If t is 0, then a will be returned, if t is 1, b will be returned. Anything between 0 and 1 will
// return the appropriate point between a and b and so on.
func Lerp(a, b Vec, t float64) Vec {
	return a.Scaled(1 - t).Add(b.Scaled(t))
}

// nearlyEqual compares two float64s and returns whether they are equal, accounting for rounding errors.At worst, the
// result is correct to 7 significant digits.
func nearlyEqual(a, b float64) bool {
	epsilon := 0.000001

	if a == b {
		return true
	}

	diff := math.Abs(a - b)

	if a == 0.0 || b == 0.0 || diff < math.SmallestNonzeroFloat64 {
		return diff < (epsilon * math.SmallestNonzeroFloat64)
	}

	absA := math.Abs(a)
	absB := math.Abs(b)

	return diff/math.Min(absA+absB, math.MaxFloat64) < epsilon
}
