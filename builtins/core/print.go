package core

import (
	"fmt"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Print output a string to stdout
func Print(env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	return &object.Null{}
}

// Printf ...
func Printf(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("printf", args, typing.MinimumArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	// Convert to the formatted version, via our `sprintf`
	// function.
	out := Sprintf(env, args...)

	// If that returned a string then we can print it
	if out.Type() == object.STRING {
		fmt.Print(out.(*object.String).Value)

	}

	return &object.Null{}
}

// Sprintf is the implementation of our `sprintf` function.
func Sprintf(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("sprintf", args, typing.MinimumArgs(1)); err != nil {
		return object.NewError(err.Error())
	}

	// We expect 1+ arguments
	if len(args) < 1 {
		return &object.String{Value: args[0].String()}
	}

	// Type-check
	if args[0].Type() != object.STRING {
		return &object.Null{}
	}

	// Get the format-string.
	fs := args[0].(*object.String).Value

	// Convert the arguments to something go's sprintf
	// code will understand.

	argLen := len(args)
	fmtArgs := make([]interface{}, argLen-1)

	// Here we convert and assign.
	for i, v := range args[1:] {
		fmtArgs[i] = v.ToInterface()
	}

	// Call the helper
	out := fmt.Sprintf(fs, fmtArgs...)

	// And now return the value.
	return &object.String{Value: out}
}
