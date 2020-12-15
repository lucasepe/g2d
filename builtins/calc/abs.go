package calc

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Abs returns the absolute value (magnitude) of a number.
// The absolute value of a number is always positive.
func Abs(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("abs", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	if args[0].Type() == object.INTEGER {
		value := args[0].(*object.Integer).Value
		if value < 0 {
			value = value * -1
		}
		return &object.Integer{Value: value}
	}

	if args[0].Type() == object.FLOAT {
		value := args[0].(*object.Float).Value
		if value < 0 {
			value = value * -1
		}
		return &object.Float{Value: value}
	}

	return object.NewError("TypeError: abs() argument #1 expected to be `int` or `float` got `%s`", args[0].Type())
}
