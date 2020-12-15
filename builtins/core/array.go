package core

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Append returns a new array with value pushed onto the end of array.
func Append(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("append", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return object.NewError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.Append(args[1])

	return newArray
}
