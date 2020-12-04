package object

// Package object implements the object system (or value system) of Monkey
// used to both represent values as the evaluator encounters and constructs
// them as well as how the user interacts with values.

import (
	"fmt"
)

const (
	// INTEGER is the Integer object type
	INTEGER = "int"

	// FLOAT is the Float object type
	FLOAT = "float"

	// STRING is the String object type
	STRING = "str"

	// BOOLEAN is the Boolean object type
	BOOLEAN = "bool"

	// NULL is the Null object type
	NULL = "null"

	// RETURN is the Return object type
	RETURN = "return"

	// ERROR is the Error object type
	ERROR = "error"

	// FUNCTION is the Function object type
	FUNCTION = "fn"

	// BUILTIN is the Builtin object type
	BUILTIN = "builtin"

	// ARRAY is the Array object type
	ARRAY = "array"

	// HASH is the Hash object type
	HASH = "hash"

	// MODULE is the Module object type
	MODULE = "module"

	// CANVAS is the Canvas object type
	CANVAS = "canvas"
)

// Comparable is the interface for comparing two Object and their underlying
// values. It is the responsibility of the caller (left) to check for types.
// Returns `true` iif the types and values are identical, `false` otherwise.
type Comparable interface {
	Compare(other Object) int
}

// Sizeable is the interface for returning the size of an Object.
// Object(s) that have a valid size must implement  this interface and the
// Len() method.
type Sizeable interface {
	Len() int
}

// Immutable is the interface for all immutable objects which must implement
// the Clone() method used by binding names to values.
type Immutable interface {
	Clone() Object
}

// Hashable is the interface for all hashable objects which must implement
// the HashKey() method which reutrns a HashKey result.
type Hashable interface {
	HashKey() HashKey
}

// BuiltinFunction represents the builtin function type
type BuiltinFunction func(env *Environment, args ...Object) Object

// Type represents the type of an object
type Type string

// Object represents a value and implementations are expected to implement
// `Type()` and `Inspect()` functions
type Object interface {
	fmt.Stringer
	Type() Type
	Bool() bool
	Inspect() string
	ToInterface() interface{}
}
