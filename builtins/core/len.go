package core

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Len len(iterable) Returns the length of the iterable (str, array or hash).
func Len(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("len", args, typing.ExactArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	if size, ok := args[0].(object.Sizeable); ok {
		return &object.Integer{Value: int64(size.Len())}
	}

	return object.NewError("TypeError: object of type '%s' has no len()", args[0].Type())
}
