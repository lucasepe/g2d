package core

import (
	"os"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Exit exit([status]) Exits the program immediately with the optional status or 0
func Exit(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("exit", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return object.NewError(err.Error())
	}

	var status int
	if len(args) == 1 {
		status = int(args[0].(*object.Integer).Value)
	}

	os.Exit(status)

	return &object.Null{}
}
