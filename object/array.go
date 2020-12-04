package object

import (
	"bytes"
	"strings"
)

// Array is the array literal type that holds a slice of Object(s)
type Array struct {
	Elements []Object
}

func (ar *Array) Bool() bool {
	return len(ar.Elements) > 0
}

func (ar *Array) PopLeft() Object {
	if len(ar.Elements) > 0 {
		e := ar.Elements[0]
		ar.Elements = ar.Elements[1:]
		return e
	}
	return &Null{}
}

func (ar *Array) PopRight() Object {
	if len(ar.Elements) > 0 {
		e := ar.Elements[(len(ar.Elements) - 1)]
		ar.Elements = ar.Elements[:(len(ar.Elements) - 1)]
		return e
	}
	return &Null{}
}

func (ar *Array) Prepend(obj Object) {
	ar.Elements = append([]Object{obj}, ar.Elements...)
}

func (ar *Array) Append(obj Object) {
	ar.Elements = append(ar.Elements, obj)
}

func (ar *Array) Copy() *Array {
	elements := make([]Object, len(ar.Elements))
	for i, e := range ar.Elements {
		elements[i] = e
	}
	return &Array{Elements: elements}
}

func (ar *Array) Reverse() {
	for i, j := 0, len(ar.Elements)-1; i < j; i, j = i+1, j-1 {
		ar.Elements[i], ar.Elements[j] = ar.Elements[j], ar.Elements[i]
	}
}

func (ar *Array) Len() int {
	return len(ar.Elements)
}

func (ar *Array) Swap(i, j int) {
	ar.Elements[i], ar.Elements[j] = ar.Elements[j], ar.Elements[i]
}

func (ar *Array) Less(i, j int) bool {
	if cmp, ok := ar.Elements[i].(Comparable); ok {
		return cmp.Compare(ar.Elements[j]) == -1
	}
	return false
}

func (ar *Array) Compare(other Object) int {
	if obj, ok := other.(*Array); ok {
		if len(ar.Elements) != len(obj.Elements) {
			return -1
		}
		for i, el := range ar.Elements {
			cmp, ok := el.(Comparable)
			if !ok {
				return -1
			}
			if cmp.Compare(obj.Elements[i]) != 0 {
				return cmp.Compare(obj.Elements[i])
			}
		}

		return 0
	}
	return -1
}

func (ar *Array) String() string { return ar.Inspect() }

// Type returns the type of the object
func (ar *Array) Type() Type { return ARRAY }

// Inspect returns a stringified version of the object for debugging
func (ar *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ar.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (ar *Array) ToInterface() interface{} { return ar.Elements }
