package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Hypot returns Sqrt(p*p + q*q), taking care to avoid
// unnecessary overflow and underflow.
func Hypot(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("hypot", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	p, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: sqrt() argument #1 %s", err.Error())
	}

	q, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: sqrt() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Hypot(p, q)}
}
