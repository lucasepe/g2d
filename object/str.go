package object

import (
	"unicode/utf8"
)

// String is the string type used to represent string literals and holds
// an internal string value
type String struct {
	Value string
}

// Type returns the type of the object
func (s *String) Type() Type { return STRING }

// Inspect returns a stringified version of the object for debugging
func (s *String) Inspect() string { return s.Value }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (s *String) ToInterface() interface{} { return s.Value }

// Bool implements the Object Bool method
func (s *String) Bool() bool {
	return s.Value != ""
}

// Len complies with Sizeable interface
func (s *String) Len() int {
	return utf8.RuneCountInString(s.Value)
}

// Compare complies with Comparable interface
func (s *String) Compare(other Object) int {
	if obj, ok := other.(*String); ok {
		switch {
		case s.Value < obj.Value:
			return -1
		case s.Value > obj.Value:
			return 1
		default:
			return 0
		}
	}
	return 1
}

// Clone creates a new copy (complies with Immutable interface)
func (s *String) Clone() Object {
	return &String{Value: s.Value}
}

func (s *String) String() string { return s.Value }
