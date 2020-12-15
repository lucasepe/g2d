package builtins

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// DrawRegularPolygon draws a regular polygon of `n` sides, centered at (x,y) with radius `r` and `angle` rotation
func DrawRegularPolygon(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("polygon", args, typing.ExactArgs(5)); err != nil {
		return newError(err.Error())
	}

	n, err := typing.ToInt(args[0])
	if err != nil {
		return newError("TypeError: polygon() argument #1 `n` %s", err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: polygon() argument #2 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: polygon() argument #3 `y` %s", err.Error())
	}

	r, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: polygon() argument #4 `r` %s", err.Error())
	}

	rad, err := typing.ToFloat(args[4])
	if err != nil {
		return newError("TypeError: polygon() argument #5 `deg` %s", err.Error())
	}

	env.Canvas().Value.DrawRegularPolygon(n, x, y, r, rad)
	return &object.Null{}
}

// DrawEllipticalArc draws an elliptical arc centered at `x, y` with a radius of `rx` in x direction and `ry` for y direction.
// The path starts at `angle1`, ends at `angle2`, and travels in the direction given by anticlockwise.
func DrawEllipticalArc(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("ellipticalArc", args, typing.ExactArgs(6)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #2 `y` %s", err.Error())
	}

	rx, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #3 `rx` %s", err.Error())
	}

	ry, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #4 `ry` %s", err.Error())
	}

	rad1, err := typing.ToFloat(args[4])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #5 `degrees1` %s", err.Error())
	}

	rad2, err := typing.ToFloat(args[5])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #6 `degrees2` %s", err.Error())
	}

	env.Canvas().Value.DrawEllipticalArc(x, y, rx, ry, rad1, rad2)
	return &object.Null{}
}

// FontSize returns or sets the font size.
// fontsize() - returns the font current size.
// fontsize(size) - sets the font current size.
/* TO FIX
func FontSize(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("fontsize", args, typing.RangeOfArgs(0, 1)); err != nil {
		return newError(err.Error())
	}

	if len(args) == 0 {
		size := env.Canvas().Value.FontHeight()
		return &object.Float{Value: size}
	}

	fs, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: fontsize() argument #1 %s", err.Error())
	}

	env.Canvas().Value.SetFontSize(fs)
	return &object.Null{}
}
*/
