package graphics

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Text draws the specified text at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
func Text(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("text", args, typing.RangeOfArgs(3, 5)); err != nil {
		return object.NewError(err.Error())
	}

	txt, err := typing.ToString(args[0])
	if err != nil {
		return object.NewError("TypeError: text() argument #1 %s", err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: text() argument #2 %s", err.Error())
	}

	y, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: text() argument #3 %s", err.Error())
	}

	ax, ay := 0.5, 0.5
	if len(args) == 5 {
		if ax, err = typing.ToFloat(args[3]); err != nil {
			return object.NewError("TypeError: text() argument #4 %s", err.Error())
		}

		if ay, err = typing.ToFloat(args[4]); err != nil {
			return object.NewError("TypeError: text() argument #5 %s", err.Error())
		}
	}

	env.Canvas().Value.DrawStringAnchored(txt, x, y, ax, ay)
	return &object.Null{}
}

// TextWidth returns the rendered width of the specified text
// given the current font face.
func TextWidth(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("textWidth", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return object.NewError(err.Error())
	}

	txt, err := typing.ToString(args[0])
	if err != nil {
		return object.NewError("TypeError: textWidth() argument #1 `str` %s", err.Error())
	}

	w, _ := env.Canvas().Value.MeasureString(txt)
	return &object.Float{Value: w}
	/*
		return &object.Array{
			Elements: []object.Object{
				&object.Float{Value: w},
				&object.Float{Value: h},
			},
		}
	*/
}
