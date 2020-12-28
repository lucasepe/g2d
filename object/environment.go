package object

import (
	"fmt"
	"math"
	"strings"

	"github.com/lucasepe/g2d/gg"
)

var (
	snapshotCounter = 0
)

const (
	keySnapshotFolder = "__SNAPSHOT_FOLDER__"
	keySnapshotPrefix = "__SNAPSHOT_PREFIX__"
)

// EnvironmentOption defines a functional option for the environment creation
type EnvironmentOption func(*Environment)

// WithOutputDir sets the snapshot output directory
func WithOutputDir(dir string) EnvironmentOption {
	return func(env *Environment) {
		env.store[keySnapshotFolder] = &String{Value: dir}
	}
}

// WithSnapshotPrefix sets the snapshot filename prefix
func WithSnapshotPrefix(prefix string) EnvironmentOption {
	return func(env *Environment) {
		env.store[keySnapshotPrefix] = &String{Value: prefix}
	}
}

// Environment is an object that holds a mapping of names to bound objets
type Environment struct {
	gContext gg.GraphicContext
	store    map[string]Object
	parent   *Environment
}

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment(ctx gg.GraphicContext, opts ...EnvironmentOption) *Environment {
	res := &Environment{
		store:    make(map[string]Object),
		gContext: ctx,
	}

	res.store["PI"] = &Float{Value: math.Pi}
	res.store["HALF_PI"] = &Float{Value: math.Pi / 2}
	res.store["QUARTER_PI"] = &Float{Value: math.Pi / 4}
	res.store["TWO_PI"] = &Float{Value: 2.0 * math.Pi}
	res.store["WIDTH"] = &Float{Value: ctx.Width()}
	res.store["HEIGHT"] = &Float{Value: ctx.Height()}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *Environment as the argument
		opt(res)
	}

	return res
}

// Clone returns a new Environment with the parent set to the current
// environment (enclosing environment)
func (e *Environment) Clone() *Environment {
	// Create a new Environment referring to the same `canvas`
	env := &Environment{
		store:    make(map[string]Object),
		gContext: e.gContext,
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
func (e *Environment) Set(name string, val Object) (Object, bool) {
	if isReserved(name) {
		return nil, false
	}

	e.store[name] = val
	return val, true
}

// GraphicContext returns the graphics context
func (e *Environment) GraphicContext() gg.GraphicContext { return e.gContext }

// SetGraphicContext sets the graphics context
func (e *Environment) SetGraphicContext(ctx gg.GraphicContext) {
	e.gContext = ctx
	if ctx != nil {
		e.store["WIDTH"] = &Float{Value: ctx.Width()}
		e.store["HEIGHT"] = &Float{Value: ctx.Height()}
	}
}

// SnapshotFilename returns the snapshot filename
func (e *Environment) SnapshotFilename() string {
	pattern := "frame_%04d.png"
	if obj, ok := e.Get(keySnapshotPrefix); ok {
		prefix := obj.(*String).Value
		pattern = fmt.Sprintf("%s_%%04d.png", prefix)
	}

	snapshotCounter = snapshotCounter + 1
	return fmt.Sprintf(pattern, snapshotCounter)
}

// SnapshotFolder returns the snapshot output folder
func (e *Environment) SnapshotFolder() string {
	obj, ok := e.Get(keySnapshotFolder)
	if !ok {
		return ""
	}

	return obj.(*String).Value
}

func isReserved(key string) bool {
	if strings.HasPrefix(key, "__") && strings.HasSuffix(key, "__") {
		return true
	}

	switch key {
	case
		"PI",
		"HALF_PI",
		"QUARTER_PI",
		"TWO_PI",
		"WIDTH",
		"HEIGHT":
		return true
	}
	return false
}
