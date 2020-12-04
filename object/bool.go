package object

import (
	"fmt"
)

// Boolean is the boolean type and used to represent boolean literals and
// holds an interval bool value
type Boolean struct {
	Value bool
}

func (b *Boolean) Bool() bool {
	return b.Value
}

func (b *Boolean) Int() int {
	if b.Value {
		return 1
	}
	return 0
}

func (b *Boolean) Compare(other Object) int {
	if obj, ok := other.(*Boolean); ok {
		return b.Int() - obj.Int()
	}
	return 1
}

func (b *Boolean) String() string {
	return b.Inspect()
}

// Clone creates a new copy
func (b *Boolean) Clone() Object {
	return &Boolean{Value: b.Value}
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (b *Boolean) ToInterface() interface{} { return b.Value }
