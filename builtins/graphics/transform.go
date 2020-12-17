package graphics

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// RotateAbout updates the current matrix with a anticlockwise rotation.
// rotate(angle) - rotation occurs about the origin.
// rotate(angle, x, y) - rotation occurs about the specified point.
// Angle is specified in radians.
func RotateAbout(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("rotate() expects one or three arguments")
	}

	rad, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: rotate() argument #1 `angle` %s", err.Error())
	}

	if len(args) == 1 {
		env.Canvas().Value.Rotate(rad)
		return &object.Null{}
	}

	if err := typing.Check("rotate", args, typing.ExactArgs(3)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: rotate() argument #2 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: rotate() argument #3 `y` %s", err.Error())
	}

	env.Canvas().Value.RotateAbout(rad, x, y)
	return &object.Null{}
}

// ScaleAbout updates the current matrix with a scaling factor.
// scale(sx, sy) - scaling occurs about the origin.
// scale(sx, sy, x, y) - scaling occurs about the specified point.
func ScaleAbout(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError("scale() expects two or four arguments")
	}

	sx, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: scale() argument #1 `sx` %s", err.Error())
	}

	sy, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: scale() argument #2 `sy` %s", err.Error())
	}

	if len(args) == 2 {
		env.Canvas().Value.Scale(sx, sy)
		return &object.Null{}
	}

	if err := typing.Check("scale", args, typing.ExactArgs(4)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: scale() argument #3 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[3])
	if err != nil {
		return object.NewError("TypeError: scale() argument #4 `y` %s", err.Error())
	}

	env.Canvas().Value.ScaleAbout(sx, sy, x, y)
	return &object.Null{}
}

// Translate updates the current matrix with a translation.
func Translate(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("translate", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: translate() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: translate() argument #2 `y` %s", err.Error())
	}

	env.Canvas().Value.Translate(x, y)
	return &object.Null{}
}

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func Identity(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Identity()
	return &object.Null{}
}

// Transform multiplies the specified point by the current matrix,
// returning a transformed position.
func Transform(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("transform", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: transform() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: transform() argument #2 `y` %s", err.Error())
	}

	tx, ty := env.Canvas().Value.TransformPoint(x, y)
	return &object.Array{
		Elements: []object.Object{
			&object.Float{Value: tx},
			&object.Float{Value: ty},
		},
	}
}
