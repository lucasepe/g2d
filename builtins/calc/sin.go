package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Sin returns the sine of the radian argument x.
func Sin(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("sin", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: sin() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Sin(val)}
}
