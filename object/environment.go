package object

import (
	"unicode"

	"github.com/lucasepe/g2d/canvas"
)

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		canvas: &Screen{
			Value: canvas.NewCanvas(),
		},
	}
}

// Environment is an object that holds a mapping of names to bound objets
type Environment struct {
	canvas *Screen
	store  map[string]Object
	parent *Environment
}

// ExportedHash returns a new Hash with the names and values of every publically
// exported binding in the environment. That is every binding that starts with a
// capital letter. This is used by the module import system to wrap up the
// evaulated module into an object.
func (e *Environment) ExportedHash() *Hash {
	pairs := make(map[HashKey]HashPair)
	for k, v := range e.store {
		if unicode.IsUpper(rune(k[0])) {
			s := &String{Value: k}
			pairs[s.HashKey()] = HashPair{Key: s, Value: v}
		}
	}
	return &Hash{Pairs: pairs}
}

// Clone returns a new Environment with the parent set to the current
// environment (enclosing environment)
func (e *Environment) Clone() *Environment {
	// Create a new Environment referring to the same `cursor`
	// for this reason im not calling `NewEnvironment`
	env := &Environment{
		store:  make(map[string]Object),
		canvas: e.canvas,
	}
	env.parent = e
	return env
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Canvas returns the canvas
func (e *Environment) Canvas() *Screen { return e.canvas }
