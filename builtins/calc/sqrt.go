package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Sqrt returns the square root of x.
func Sqrt(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("sqrt", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: sqrt() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Sqrt(val)}
}
