package core

import (
	"math"
	"strconv"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Bool converts value to a bool
func Bool(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("bool", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	return &object.Boolean{Value: args[0].Bool()}
}

// Float converts decimal value str to float. If value is invalid returns null
func Float(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check(
		"float", args,
		typing.ExactArgs(1),
	); err != nil {
		return object.NewError(err.Error())
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		if arg.Value {
			return &object.Float{Value: 1}
		}
		return &object.Float{Value: 0}
	case *object.Float:
		return arg
	case *object.Integer:
		return &object.Float{Value: float64(arg.Value)}
	case *object.String:
		n, err := strconv.ParseFloat(arg.Value, 64)
		if err != nil {
			return object.NewError("could not parse string to int: %s", err)
		}
		return &object.Float{Value: n}
	default:
		return &object.Float{}
	}
}

// Int converts decimal value str to int. If value is invalid returns null
func Int(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("int", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		if arg.Value {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case *object.Integer:
		return arg
	case *object.Float:
		return &object.Integer{Value: int64(math.Round(arg.Value))}
	case *object.String:
		n, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			return object.NewError("could not parse string to int: %s", err)
		}
		return &object.Integer{Value: n}
	default:
		return &object.Integer{}
	}
}

// Str returns the string representation of value.
func Str(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("str", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	return &object.String{Value: args[0].Inspect()}
}

// TypeOf returns a str denoting the type of value: nil, bool, int, float, str, array, hash, or fn.
func TypeOf(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("type", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	return &object.String{Value: string(args[0].Type())}
}
