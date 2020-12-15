package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Degrees converts radians into degrees.
// degrees(angle) - angle in radians is the value that you want to convert
func Degrees(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("degrees", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	radians, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: degrees() argument #1 %s", err.Error())
	}

	res := radians * 180.0 / math.Pi
	return &object.Float{Value: res}
}
