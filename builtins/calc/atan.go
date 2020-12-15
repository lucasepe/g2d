package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Atan returns the arctangent, in radians, of x.
func Atan(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("atan", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: atan() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Atan(val)}
}
