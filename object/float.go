package object

import (
	"math/big"
	"strconv"
)

// Float is the float type used to represent float literals and holds
// an internal float64 value
type Float struct {
	Value float64
}

func (f *Float) Bool() bool {
	return f.Value != 0
}

func (f *Float) Compare(other Object) int {
	if obj, ok := other.(*Float); ok {
		return compareFloats(f.Value, obj.Value)
	}
	if obj, ok := other.(*Integer); ok {
		return compareFloats(f.Value, float64(obj.Value))
	}
	return -1
}

func (f *Float) String() string {
	return f.Inspect()
}

// Clone creates a new copy
func (f *Float) Clone() Object {
	return &Float{Value: f.Value}
}

// Type returns the type of the object
func (f *Float) Type() Type { return FLOAT }

// Inspect returns a stringified version of the object for debugging
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (f *Float) ToInterface() interface{} { return f.Value }

func compareFloats(a, b float64) int {
	return big.NewFloat(a).Cmp(big.NewFloat(b))
}
