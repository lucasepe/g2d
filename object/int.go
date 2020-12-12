package object

import (
	"fmt"
)

// Integer is the integer type used to represent integer literals and holds
// an internal int64 value
type Integer struct {
	Value int64
}

func (i *Integer) Bool() bool {
	return i.Value != 0
}

func (i *Integer) Compare(other Object) int {
	if obj, ok := other.(*Integer); ok {
		switch {
		case i.Value < obj.Value:
			return -1
		case i.Value > obj.Value:
			return 1
		default:
			return 0
		}
	}

	if obj, ok := other.(*Float); ok {
		return compareFloats(float64(i.Value), obj.Value)
	}

	return -1
}

func (i *Integer) String() string {
	return i.Inspect()
}

// Clone creates a new copy
func (i *Integer) Clone() Object {
	return &Integer{Value: i.Value}
}

// Type returns the type of the object
func (i *Integer) Type() Type { return INTEGER }

// Inspect returns a stringified version of the object for debugging
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (i *Integer) ToInterface() interface{} { return i.Value }
