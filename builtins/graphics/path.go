package graphics

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// BeginPath starts a new path.
func BeginPath(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.NewSubPath()
	return &object.Null{}
}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
func ClosePath(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.ClosePath()
	return &object.Null{}
}

// RouteTo adds a line segment to the current path starting at the current point.
// If there is no current point, it is equivalent to MoveTo(x, y)
func RouteTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("routeTo", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	d, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: routeTo() argument #1 `distance` %s", err.Error())
	}

	a, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: routeTo() argument #2 `angle` %s", err.Error())
	}

	// displacements in x and y directions
	dx, dy := d*math.Cos(a), d*math.Sin(a)

	dc := env.Canvas().Value
	if pt, ok := dc.GetCurrentPoint(); !ok {
		dc.MoveTo(dx, dy)
	} else {
		dc.LineTo(pt.X+dx, pt.Y+dy)
	}

	return &object.Null{}
}

// MoveTo starts a new subpath within the current path starting at the specified point.
func MoveTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("moveTo", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: moveTo() argument #1 %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: moveTo() argument #2 %s", err.Error())
	}

	env.Canvas().Value.MoveTo(x, y)
	return &object.Null{}
}

// LineTo adds a line segment to the current path starting at the current point.
// If there is no current point, it is equivalent to MoveTo(x, y)
func LineTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("lineTo", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: lineTo() argument #1 %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: lineTo() argument #2 %s", err.Error())
	}

	env.Canvas().Value.LineTo(x, y)
	return &object.Null{}
}

// ArcTo adds a circular arc to the current sub-path, using the given
// control points and radius.
// The arc is automatically connected to the path's latest point
// with a straight line, if necessary for the specified parameters.
func ArcTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("arcTo", args, typing.ExactArgs(5)); err != nil {
		return object.NewError(err.Error())
	}

	x1, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: arcTo() argument #1 `x1` %s", err.Error())
	}

	y1, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: arcTo() argument #2 `y1` %s", err.Error())
	}

	x2, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: arcTo() argument #3 `x2` %s", err.Error())
	}

	y2, err := typing.ToFloat(args[3])
	if err != nil {
		return object.NewError("TypeError: arcTo() argument #4 `y2` %s", err.Error())
	}

	r, err := typing.ToFloat(args[4])
	if err != nil {
		return object.NewError("TypeError: arcTo() argument #5 `r` %s", err.Error())
	}

	env.Canvas().Value.ArcTo(x1, y1, x2, y2, r)
	return &object.Null{}
}

// QuadraticCurveTo adds a quadratic Bézier curve to the current sub-path.
// It requires two points: the first one is a control point and the second one is the end point.
// The starting point is the latest point in the current path, which can be
// changed using `moveTo()`` before creating the quadratic Bézier curve.
func QuadraticCurveTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("quadraticCurveTo", args, typing.ExactArgs(4)); err != nil {
		return object.NewError(err.Error())
	}

	x1, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: quadraticCurveTo() argument #1 %s", err.Error())
	}

	y1, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: quadraticCurveTo() argument #2 %s", err.Error())
	}

	x2, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: quadraticCurveTo() argument #3 %s", err.Error())
	}

	y2, err := typing.ToFloat(args[3])
	if err != nil {
		return object.NewError("TypeError: quadraticCurveTo() argument #4 %s", err.Error())
	}

	env.Canvas().Value.QuadraticTo(x1, y1, x2, y2)
	return &object.Null{}
}
