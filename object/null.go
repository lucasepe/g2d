package object

// Null is the null type and used to represent the absence of a value
type Null struct{}

// Bool implements the Object Bool method
func (n *Null) Bool() bool {
	return false
}

// Compare complies with Comparable interface
func (n *Null) Compare(other Object) int {
	if _, ok := other.(*Null); ok {
		return 0
	}
	return 1
}

// Type returns the type of the object
func (n *Null) Type() Type { return NULL }

// Inspect returns a stringified version of the object for debugging
func (n *Null) Inspect() string { return "null" }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (n *Null) ToInterface() interface{} { return "<NULL>" }

func (n *Null) String() string { return n.Inspect() }
