package calc

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Lerp calculates a number between two numbers at a specific increment.
// The amt parameter is the amount to interpolate between the two values
// where 0.0 equal to the first point, 0.1 is very near the first
// point, 0.5 is half-way in between, etc.
func Lerp(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("lerp", args, typing.ExactArgs(3)); err != nil {
		return object.NewError(err.Error())
	}

	start, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: lerp() argument #1 %s", err.Error())
	}

	stop, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: lerp() argument #2 %s", err.Error())
	}

	amt, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: lerp() argument #3 %s", err.Error())
	}

	delta := stop - start
	return &object.Float{Value: start + delta*amt}
}
