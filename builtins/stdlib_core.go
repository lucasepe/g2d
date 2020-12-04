package builtins

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Args args() Returns an array of command-line options passed to the program
func Args(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check(
		"args", args,
		typing.ExactArgs(0),
	); err != nil {
		return newError(err.Error())
	}

	elements := make([]object.Object, len(object.Arguments))
	for i, arg := range object.Arguments {
		elements[i] = &object.String{Value: arg}
	}
	return &object.Array{Elements: elements}
}

// Bool converts value to a bool
func Bool(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check(
		"bool", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.Boolean{Value: args[0].Bool()}
}

// Exit exit([status]) Exits the program immediately with the optional status or 0
func Exit(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check(
		"exit", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	var status int
	if len(args) == 1 {
		status = int(args[0].(*object.Integer).Value)
	}

	object.ExitFunction(status)

	return nil
}

// Float converts decimal value str to float. If value is invalid returns null
func Float(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check(
		"float", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
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
			return newError("could not parse string to int: %s", err)
		}
		return &object.Float{Value: n}
	default:
		return &object.Float{}
	}
}

// Int converts decimal value str to int. If value is invalid returns null
func Int(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check(
		"int", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
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
			return newError("could not parse string to int: %s", err)
		}
		return &object.Integer{Value: n}
	default:
		return &object.Integer{}
	}
}

// Len len(iterable) Returns the length of the iterable (str, array or hash).
func Len(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("len", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	if size, ok := args[0].(object.Sizeable); ok {
		return &object.Integer{Value: int64(size.Len())}
	}

	return newError("TypeError: object of type '%s' has no len()", args[0].Type())
}

// Input reads a line from standard input optionally printing prompt.
// input([prompt]) prints the prompt.
func Input(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("input", args, typing.RangeOfArgs(0, 1), typing.WithTypes(object.STRING)); err != nil {
		return newError(err.Error())
	}

	if len(args) == 1 {
		prompt := args[0].(*object.String).Value
		fmt.Fprintf(os.Stdout, prompt)
	}

	buffer := bufio.NewReader(os.Stdin)

	line, _, err := buffer.ReadLine()
	if err != nil && err != io.EOF {
		return newError(fmt.Sprintf("error reading input from stdin: %s", err))
	}
	return &object.String{Value: string(line)}
}

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
		return newError(err.Error())
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
		return newError(err.Error())
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

// Str returns the string representation of value.
func Str(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("str", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: args[0].String()}
}

// TypeOf returns a str denoting the type of value: nil, bool, int, float, str, array, hash, or fn.
func TypeOf(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("type", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: string(args[0].Type())}
}

// Push returns a new array with value pushed onto the end of array.
func Push(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("push", args, typing.ExactArgs(2), typing.WithTypes(object.ARRAY)); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.Append(args[1])

	return newArray
}

// Returns the keys of the specified hash
func hashKeys(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("hashKeys", args, typing.ExactArgs(1), typing.WithTypes(object.HASH)); err != nil {
		return newError(err.Error())
	}

	// The object we're working with
	hash := args[0].(*object.Hash)
	ents := len(hash.Pairs)

	// Create a new array for the results.
	array := make([]object.Object, ents)

	// Now copy the keys into it.
	i := 0
	for _, ent := range hash.Pairs {
		array[i] = ent.Key
		i++
	}

	// Return the array.
	return &object.Array{Elements: array}
}

// Delete a given hash-key
func hashDelete(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("hashDelete", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	// The object we're working with
	hash, ok := args[0].(*object.Hash)
	if !ok {
		return newError("TypeError: hashDelete() expected argument #1 to be `hash` got `%s`", args[0].Type())
	}

	// The key we're going to delete
	key := args[1].(object.Hashable)
	if !ok {
		return newError(
			"TypeError: hashDelete() expected argument #2 to be `string`, `int`, `float` or `bool` got `%s`",
			args[1].Type())
	}

	// Make a new hash
	newHash := make(map[object.HashKey]object.HashPair)

	// Copy the values EXCEPT the one we have.
	for k, v := range hash.Pairs {
		if k != key.HashKey() {
			newHash[k] = v
		}
	}
	return &object.Hash{Pairs: newHash}
}
