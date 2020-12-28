package graphics

import (
	"image"
	"path/filepath"
	"reflect"

	"github.com/lucasepe/g2d/gg/img"
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Size sets the image size.
// `size(d)` creates a squared image where width and height are equals to `d`.
// `size(w, h) creates an image where width is equals to `w` and height to `h`.
func Size(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("size", args, typing.RangeOfArgs(1, 2)); err != nil {
		return object.NewError(err.Error())
	}

	w, err := typing.ToInt(args[0])
	if err != nil {
		return object.NewError("TypeError: size() argument #1 %s", err.Error())
	}

	h := w
	if len(args) > 1 {
		if h, err = typing.ToInt(args[1]); err != nil {
			return object.NewError("TypeError: size() argument #2 %s", err.Error())
		}
	}

	ctx := img.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, w, h)))
	env.SetGraphicContext(ctx)
	return &object.Null{}
}

// Clear fills the entire image with the current color. Clear all drawings.
func Clear(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("clear", args, typing.ExactArgs(0)); err != nil {
		return object.NewError(err.Error())
	}

	env.GraphicContext().Clear()
	return &object.Null{}
}

// StrokeColor returns or sets the stroke color.
// strokeColor(hexcolor) - sets the stroke color to `hexcolor`.
// strokeColor(r, g, b) - sets the stroke color to `r,g,b` values.
// strokeColor(r, g, b, a) -  sets the stroke color to `r,g,b,a` values.
func StrokeColor(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("strokeColor() expects one or four arguments")
	}

	if len(args) == 1 {
		if args[0].Type() == object.STRING {
			color := args[0].(*object.String).Value
			env.GraphicContext().SetStrokeColor(parseHexColor(color))
			return &object.Null{}
		}

		return object.NewError("TypeError: strokeColor() argument #1 expected to be `string` got `%s`", args[0].Type())
	}

	if err := typing.Check("strokeColor", args, typing.RangeOfArgs(3, 4)); err != nil {
		return object.NewError(err.Error())
	}

	r, err := typing.ToInt(args[0])
	if err != nil {
		return object.NewError("TypeError: strokeColor() argument #1 `r` %s", err.Error())
	}

	g, err := typing.ToInt(args[1])
	if err != nil {
		return object.NewError("TypeError: strokeColor() argument #2 `g` %s", err.Error())
	}

	b, err := typing.ToInt(args[2])
	if err != nil {
		return object.NewError("TypeError: strokeColor() argument #3 `b` %s", err.Error())
	}

	a := 255
	if len(args) == 4 {
		a, _ = typing.ToInt(args[3])
	}

	env.GraphicContext().SetStrokeColor(r, g, b, a)
	return &object.Null{}
}

// FillColor returns or sets the fill color.
// fillColor(hexcolor) - sets the fill color to `hexcolor`.
// fillColor(r, g, b) - sets the fill color to `r,g,b` values.
// fillColor(r, g, b, a) -  sets the fill color to `r,g,b,a` values.
func FillColor(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("fillColor() expects one or four arguments")
	}

	if len(args) == 1 {
		if args[0].Type() == object.STRING {
			color := args[0].(*object.String).Value
			env.GraphicContext().SetFillColor(parseHexColor(color))
			return &object.Null{}
		}

		return object.NewError("TypeError: fillColor() argument #1 expected to be `string` got `%s`", args[0].Type())
	}

	if err := typing.Check("fillColor", args, typing.RangeOfArgs(3, 4)); err != nil {
		return object.NewError(err.Error())
	}

	r, err := typing.ToInt(args[0])
	if err != nil {
		return object.NewError("TypeError: fillColor() argument #1 `r` %s", err.Error())
	}

	g, err := typing.ToInt(args[1])
	if err != nil {
		return object.NewError("TypeError: fillColor() argument #2 `g` %s", err.Error())
	}

	b, err := typing.ToInt(args[2])
	if err != nil {
		return object.NewError("TypeError: fillColor() argument #3 `b` %s", err.Error())
	}

	a := 255
	if len(args) == 4 {
		a, _ = typing.ToInt(args[3])
	}

	env.GraphicContext().SetFillColor(r, g, b, a)
	return &object.Null{}
}

// StrokeWeight returns or sets the stroke line thickness.
// strokeWeight() - returns the current stroke line thickness.
// strokeWeight(width) - sets the stroke thickness to `width`.
func StrokeWeight(env *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		lw := env.GraphicContext().StrokeWeight()
		return &object.Float{Value: lw}
	}

	lw, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: strokeWeight() argument #1 %s", err.Error())
	}

	env.GraphicContext().SetStrokeWeight(lw)
	return &object.Null{}
}

// Dashes sets the current dash pattern to use. Call with zero arguments to
// disable dashes. The values specify the lengths of each dash, with
// alternating on and off lengths.
func Dashes(env *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		env.GraphicContext().SetLineDash()
		return &object.Null{}
	}

	if args[0].Type() == object.ARRAY {
		info := args[0].(*object.Array)

		dashes := []float64{}
		for _, item := range info.Elements {
			el, err := typing.ToFloat(item)
			if err != nil {
				return object.NewError("TypeError: dashes() %s", err.Error())
			}
			dashes = append(dashes, el)
		}

		env.GraphicContext().SetLineDash(dashes...)
		return &object.Null{}
	}

	dashes := []float64{}
	for i, el := range args {
		val, err := typing.ToFloat(el)
		if err != nil {
			return object.NewError("TypeError: dashes() argument #%d %s", (i + 1), err.Error())
		}
		dashes = append(dashes, val)
	}

	env.GraphicContext().SetLineDash(dashes...)
	return &object.Null{}
}

// GetCurrentX returns the current X position if there is a current point.
func GetCurrentX(env *object.Environment, args ...object.Object) object.Object {
	if x, _, ok := env.GraphicContext().CurrentPoint(); ok {
		return &object.Float{Value: x}
	}

	return object.NewError("Error: xpos() - there is no current point after `stroke` or `fill`")
}

// GetCurrentY returns the current Y position if there is a current point.
func GetCurrentY(env *object.Environment, args ...object.Object) object.Object {
	if _, y, ok := env.GraphicContext().CurrentPoint(); ok {
		return &object.Float{Value: y}
	}

	return object.NewError("Error: ypos() - there is no current point after `stroke` or `fill`")
}

// Width returns the image width
// `width([img])` - if the img is not specified the
func Width(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("width", args, typing.RangeOfArgs(0, 1)); err != nil {
		return object.NewError(err.Error())
	}

	if len(args) == 0 {
		w := env.GraphicContext().Width()
		return &object.Integer{Value: int64(w)}
	}

	if args[0].Type() == object.IMAGE {
		val := args[0].(*object.Image).Value
		return &object.Integer{Value: int64(val.Bounds().Size().X)}
	}

	return &object.Null{}
}

// Height returns the image height
func Height(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("height", args, typing.RangeOfArgs(0, 1)); err != nil {
		return object.NewError(err.Error())
	}

	if len(args) == 0 {
		h := env.GraphicContext().Height()
		return &object.Integer{Value: int64(h)}
	}

	if args[0].Type() == object.IMAGE {
		val := args[0].(*object.Image).Value
		return &object.Integer{Value: int64(val.Bounds().Size().Y)}
	}

	return &object.Null{}
}

// Push saves the current state of the context by pushin it onto a stack. These can be nested.
func Push(env *object.Environment, args ...object.Object) object.Object {
	env.GraphicContext().Push()
	return &object.Null{}
}

// Pop restores the last saved context state from the stack.
func Pop(env *object.Environment, args ...object.Object) object.Object {
	env.GraphicContext().Pop()
	return &object.Null{}
}

// Snapshot encodes as PNG the image with the current drawings and saves the image on the filesystem.
// If file name is omitted it will be autogenerated adn the PNG saved in the .g2d file folder.
func Snapshot(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("snapshot", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return object.NewError(err.Error())
	}

	ctx, ok := env.GraphicContext().(*img.Context)
	if !ok {
		return object.NewError(
			"expected graphic context of type img.Context, got: %v (not impl yet)",
			reflect.TypeOf(env.GraphicContext()))
	}

	filename := env.SnapshotFilename()
	if len(args) == 1 {
		filename = args[0].(*object.String).Value
	}

	if folder := env.SnapshotFolder(); folder != "" {
		/*
			if err := utils.Mkdir(folder); err != nil {
				return object.NewError(err.Error())
			}*/

		filename = filepath.Join(folder, filename)
	}

	if err := savePNG(filename, ctx.Image()); err != nil {
		return object.NewError(err.Error())
	}

	return &object.Null{}
}

// Stroke strokes the current path with the current color and line width
// the path is cleared after this operation.
func Stroke(env *object.Environment, args ...object.Object) object.Object {
	env.GraphicContext().Stroke()
	return &object.Null{}
}

// Fill fills the current path with the current color.
// Open subpaths are implicity closed. The path is cleared after this operation.
func Fill(env *object.Environment, args ...object.Object) object.Object {
	env.GraphicContext().Fill()
	return &object.Null{}
}

// FillAndStroke first fills the current path and than strokes it
func FillAndStroke(env *object.Environment, args ...object.Object) object.Object {
	env.GraphicContext().FillAndStroke()
	return &object.Null{}
}

// Viewport sets up user-defined coordinate system.
// This performs a screen reset, all drawings are cleared.
func Viewport(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("viewport", args, typing.RangeOfArgs(4, 6)); err != nil {
		return object.NewError(err.Error())
	}

	xMin, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: viewport() argument #1 `xMin` %s", err.Error())
	}

	xMax, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: viewport() argument #2 `xMax` %s", err.Error())
	}

	if xMax <= xMin {
		return object.NewError("RangeError: viewport() xMax must be greater then xMin")
	}

	yMin, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: viewport() argument #3 `yMin` %s", err.Error())
	}

	yMax, err := typing.ToFloat(args[3])
	if err != nil {
		return object.NewError("TypeError: viewport() argument #4 `yMax` %s", err.Error())
	}

	if yMax <= yMin {
		return object.NewError("RangeError: viewport() yMax must be greater then yMin")
	}

	xOffset := 0.0
	if len(args) == 5 {
		xOffset, _ = typing.ToFloat(args[4])
	}

	yOffset := 0.0
	if len(args) == 6 {
		yOffset, _ = typing.ToFloat(args[5])
	}

	setWorldCoordinates(env.GraphicContext(), xMin, xMax, yMin, yMax, xOffset, yOffset)
	return &object.Null{}
}
