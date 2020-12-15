package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Atan2 returns the arc tangent of y/x, using
// the signs of the two to determine the quadrant
// of the return value.
func Atan2(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("atan2", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	y, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: atan2() argument #1 %s", err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: atan2() argument #2 %s", err.Error())
	}

	return &object.Float{Value: math.Atan2(y, x)}
}
