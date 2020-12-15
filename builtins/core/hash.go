package core

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Keys returns the keys of the specified hash
func Keys(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("keys", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.HASH),
	); err != nil {
		return object.NewError(err.Error())
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
func Delete(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("delete", args, typing.ExactArgs(2)); err != nil {
		return object.NewError(err.Error())
	}

	// The object we're working with
	hash, ok := args[0].(*object.Hash)
	if !ok {
		return object.NewError("TypeError: hashDelete() expected argument #1 to be `hash` got `%s`", args[0].Type())
	}

	// The key we're going to delete
	key := args[1].(object.Hashable)
	if !ok {
		return object.NewError(
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
