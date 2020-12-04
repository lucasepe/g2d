package object

import (
	"fmt"
)

// Module is the module type used to represent a collection of variabels.
type Module struct {
	Name  string
	Attrs Object
}

func (m *Module) Bool() bool {
	return true
}

func (m *Module) Compare(other Object) int {
	return 1
}

func (m *Module) String() string {
	return m.Inspect()
}

// Type returns the type of the object
func (m *Module) Type() Type { return MODULE }

// Inspect returns a stringified version of the object for debugging
func (m *Module) Inspect() string { return fmt.Sprintf("<module '%s'>", m.Name) }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (m *Module) ToInterface() interface{} {
	return "<MODULE>"
}
