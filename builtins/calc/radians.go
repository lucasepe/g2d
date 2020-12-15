package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Radians converts degrees to radians.
// radians(angle) - angle in degrees is the value that you want to convert
func Radians(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("radians", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	degrees, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: radians() argument #1 %s", err.Error())
	}

	res := degrees * math.Pi / 180.0
	return &object.Float{Value: res}
}
