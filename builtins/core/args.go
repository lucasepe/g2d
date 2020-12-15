package core

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Args args() Returns an array of command-line options passed to the program
func Args(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("args", args, typing.ExactArgs(0)); err != nil {
		return object.NewError(err.Error())
	}

	elements := make([]object.Object, len(object.Arguments))
	for i, arg := range object.Arguments {

		elements[i] = &object.String{Value: arg}
	}
	return &object.Array{Elements: elements}
}
