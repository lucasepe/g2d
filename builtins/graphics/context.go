package graphics

import (
	"fmt"
	"os"
	"path/filepath"

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

	env.Canvas().Reset(w, h)
	return &object.Null{}
}

// Clear fills the entire image with the current color. Clear all drawings.
func Clear(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("clear", args, typing.ExactArgs(0)); err != nil {
		return object.NewError(err.Error())
	}

	env.Canvas().Value.Clear()
	return &object.Null{}
}

// PenColor returns or sets the cursor pen color.
// pencolor(hexcolor) - sets the pen color to `hexcolor`.
// pencolor(r, g, b) - sets the pen color to `r,g,b` values - should be between 0 and 1, inclusive. Alpha will be set to 1 (fully opaque).
// pencolor(r, g, b, a) -  sets the pen color to `r,g,b,a` values - should be between 0 and 1, inclusive.
func PenColor(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError("pencolor() expects one or four arguments")
	}

	if len(args) == 1 {
		if args[0].Type() == object.STRING {
			color := args[0].(*object.String).Value
			env.Canvas().Value.SetHexColor(color)
			return &object.Null{}
		}

		return object.NewError("TypeError: pencolor() argument #1 expected to be `string` got `%s`", args[0].Type())
	}

	if err := typing.Check("pencolor", args, typing.RangeOfArgs(3, 4)); err != nil {
		return object.NewError(err.Error())
	}

	r, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: pencolor() argument #1 `r` %s", err.Error())
	}

	g, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: pencolor() argument #2 `g` %s", err.Error())
	}

	b, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: pencolor() argument #3 `b` %s", err.Error())
	}

	switch len(args) {
	case 3:
		env.Canvas().Value.SetRGB(r, g, b)
	case 4:
		a, err := typing.ToFloat(args[3])
		if err != nil {
			return object.NewError("TypeError: pencolor() argument #4 `a` %s", err.Error())
		}
		env.Canvas().Value.SetRGBA(r, g, b, a)
	}

	return &object.Null{}
}

// PenSize returns or sets the pen line thickness.
// pensize() - returns the current pen line thickness.
// pensize(width) - sets the pen line thickness to `width`.
func PenSize(env *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		lw := env.Canvas().Value.LineWidth()
		return &object.Float{Value: lw}
	}

	lw, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: pensize() argument #1 %s", err.Error())
	}

	env.Canvas().Value.SetLineWidth(lw)
	return &object.Null{}
}

// Dashes sets the current dash pattern to use. Call with zero arguments to
// disable dashes. The values specify the lengths of each dash, with
// alternating on and off lengths.
func Dashes(env *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		env.Canvas().Value.SetDash()
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

		env.Canvas().Value.SetDash(dashes...)
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

	env.Canvas().Value.SetDash(dashes...)
	return &object.Null{}
}

// GetCurrentX returns the current X position if there is a current point.
func GetCurrentX(env *object.Environment, args ...object.Object) object.Object {
	dc := env.Canvas().Value
	if pt, ok := dc.GetCurrentPoint(); ok {
		return &object.Float{Value: pt.X}
	}

	return object.NewError("Error: xpos() - there is no current point after `stroke` or `fill`")
}

// GetCurrentY returns the current Y position if there is a current point.
func GetCurrentY(env *object.Environment, args ...object.Object) object.Object {
	dc := env.Canvas().Value
	if pt, ok := dc.GetCurrentPoint(); ok {
		return &object.Float{Value: pt.Y}
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
		w := env.Canvas().Value.Width()
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
		h := env.Canvas().Value.Height()
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
	env.Canvas().Value.Push()
	return &object.Null{}
}

// Pop restores the last saved context state from the stack.
func Pop(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Pop()
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

	if len(args) == 1 {
		filename := args[0].(*object.String).Value
		if err := env.Canvas().Value.SavePNG(filename); err != nil {
			return object.NewError(err.Error())
		}

		return &object.Null{}
	}

	object.SaveCounter = object.SaveCounter + 1

	folder := filepath.Join(object.WorkDir, object.SourceFile)
	if err := mkdir(folder); err != nil {
		return object.NewError(err.Error())
	}

	filename := filepath.Join(folder,
		fmt.Sprintf("%s_%03d.png", object.SourceFile, object.SaveCounter))

	if err := env.Canvas().Value.SavePNG(filename); err != nil {
		return object.NewError(err.Error())
	}

	return &object.Null{}
}

func mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return nil
}

// Stroke strokes the current path with the current color and line width
// the path is cleared after this operation.
// If preserve is true the path will be preserved.
func Stroke(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("stroke", args, typing.RangeOfArgs(0, 1)); err != nil {
		return object.NewError(err.Error())
	}

	preserve := false
	if (len(args) > 0) && (args[0].Type() == object.BOOLEAN) {
		preserve = args[0].(*object.Boolean).Value
	}

	if preserve {
		env.Canvas().Value.StrokePreserve()
	} else {
		env.Canvas().Value.Stroke()
	}

	return &object.Null{}
}

// Fill fills the current path with the current color.
// Open subpaths are implicity closed. The path is cleared after this operation.
// If preserve is true the path is preserved after this operation.
func Fill(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("fill", args, typing.RangeOfArgs(0, 1)); err != nil {
		return object.NewError(err.Error())
	}

	preserve := false
	if (len(args) > 0) && (args[0].Type() == object.BOOLEAN) {
		preserve = args[0].(*object.Boolean).Value
	}

	if preserve {
		env.Canvas().Value.FillPreserve()
	} else {
		env.Canvas().Value.Fill()
	}

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

	env.Canvas().Value.SetWorldCoordinates(xMin, xMax, yMin, yMax, xOffset, yOffset)
	return &object.Null{}
}
