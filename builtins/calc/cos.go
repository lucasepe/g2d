package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Cos returns the cosine of the radian argument x.
func Cos(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("cos", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: cos() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Cos(val)}
}
