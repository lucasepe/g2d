package object

import (
	"image"
)

type Image struct {
	Value image.Image
}

func (im *Image) String() string { return im.Inspect() }

func (im *Image) Bool() bool { return im.Value != nil }

// Type returns the type of the object
func (im *Image) Type() Type { return IMAGE }

// Inspect returns a stringified version of the object for debugging
func (im *Image) Inspect() string { return "<IMAGE>" }

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (im *Image) ToInterface() interface{} { return "<IMAGE>" }

// Clone creates a new copy
func (im *Image) Clone() Object {
	return &Image{Value: im.Value}
}
