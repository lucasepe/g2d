package graphics

import (
	"image/png"
	"os"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// LoadPNG loads a PNG image file
func LoadPNG(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("loadPNG", args,
		typing.ExactArgs(1), typing.WithTypes(object.STRING),
	); err != nil {
		return object.NewError(err.Error())
	}

	name, err := typing.ToString(args[0])
	if err != nil {
		return object.NewError("TypeError: loadPNG() argument #1 %s", err.Error())
	}

	fd, err := os.Open(name)
	if err != nil {
		return object.NewError("TypeError: loadPNG() - %s", err.Error())
	}
	defer fd.Close()

	im, err := png.Decode(fd)
	if err != nil {
		return object.NewError("DecodeError: loadPNG() - %s", err.Error())
	}

	return &object.Image{Value: im}
}

// ImageAnchored draws the specified image at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image. Use ax=0.5, ay=0.5 to center the image at the specified point.
func ImageAnchored(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("image", args, typing.RangeOfArgs(3, 5)); err != nil {
		return object.NewError(err.Error())
	}

	im, err := typing.ToImage(args[0])
	if err != nil {
		return object.NewError("TypeError: image() argument #1 %s", err.Error())
	}

	x, err := typing.ToInt(args[1])
	if err != nil {
		return object.NewError("TypeError: image() argument #2 %s", err.Error())
	}

	y, err := typing.ToInt(args[2])
	if err != nil {
		return object.NewError("TypeError: image() argument #3 %s", err.Error())
	}

	ax, ay := 0.5, 0.5
	if len(args) == 5 {
		ax, err = typing.ToFloat(args[3])
		if err != nil {
			return object.NewError("TypeError: image() argument #4 %s", err.Error())
		}

		ay, err = typing.ToFloat(args[4])
		if err != nil {
			return object.NewError("TypeError: image() argument #5 %s", err.Error())
		}
	}

	env.Canvas().Value.DrawImageAnchored(im, x, y, ax, ay)
	return &object.Null{}
}
