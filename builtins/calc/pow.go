package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Pow returns x**y, the base-x exponential of y.
func Pow(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("pow", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: pow() argument #1 %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: pow() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Pow(x, y)}
}
