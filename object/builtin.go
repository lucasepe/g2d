package object

import (
	"fmt"
)

// Builtin  is the builtin object type that simply holds a reference to
// a BuiltinFunction type that takes zero or more objects as arguments
// and returns an object.
type Builtin struct {
	Name string
	Fn   BuiltinFunction
	Env  *Environment
}

// Bool implements the Object Bool method
func (b *Builtin) Bool() bool {
	return true
}

// Type returns the type of the object
func (b *Builtin) Type() Type { return BUILTIN }

// Inspect returns a stringified version of the object for debugging
func (b *Builtin) Inspect() string {
	return fmt.Sprintf("<built-in function %s>", b.Name)
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Builtin) ToInterface() interface{} {
	return "<BUILTIN>"
}

func (b *Builtin) String() string { return b.Inspect() }
