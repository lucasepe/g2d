package builtins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// ScreenSize returns or sets the screen size.
// screensize() - returns the screen current size.
// screensize(size) - sets the screen size.
func ScreenSize(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("screensize", args, typing.RangeOfArgs(0, 2)); err != nil {
		return newError(err.Error())
	}

	if len(args) == 0 {
		w, h := env.Canvas().Value.Size()
		return &object.Array{
			Elements: []object.Object{
				&object.Integer{Value: int64(w)},
				&object.Integer{Value: int64(h)},
			},
		}
	}

	w, err := typing.ToInt(args[0])
	if err != nil {
		return newError("TypeError: screensize() argument #1 `width` %s", err.Error())
	}

	h := w
	if len(args) > 1 {
		if h, err = typing.ToInt(args[1]); err != nil {
			return newError("TypeError: screensize() argument #2 `height` %s", err.Error())
		}
	}

	env.Canvas().Value.Reset(w, h)
	return &object.Null{}
}

// WorldCoords sets up user-defined coordinate system.
// This performs a screen reset, all drawings are cleared.
func WorldCoords(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("worldcoords", args, typing.ExactArgs(4)); err != nil {
		return newError(err.Error())
	}

	xMin, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: worldcoords() argument #1 `xMin` %s", err.Error())
	}

	xMax, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: worldcoords() argument #2 `xMax` %s", err.Error())
	}

	yMin, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: worldcoords() argument #3 `yMin` %s", err.Error())
	}

	yMax, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: worldcoords() argument #4 `yMax` %s", err.Error())
	}

	canvas := env.Canvas().Value
	if err := canvas.SetWorldCoordinates(xMin, xMax, yMin, yMax); err != nil {
		return newError(err.Error())
	}

	return &object.Null{}
}

// PenColor returns or sets the cursor pen color.
// pencolor(hexcolor) - sets the pen color to `hexcolor`.
// pencolor(r, g, b) - sets the pen color to `r,g,b` values - should be between 0 and 1, inclusive. Alpha will be set to 1 (fully opaque).
// pencolor(r, g, b, a) -  sets the pen color to `r,g,b,a` values - should be between 0 and 1, inclusive.
func PenColor(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return newError("pencolor() expects one or four arguments")
	}

	if len(args) == 1 {
		if args[0].Type() == object.STRING {
			color := args[0].(*object.String).Value
			env.Canvas().Value.Graphics().SetHexColor(color)
			return &object.Null{}
		}

		return newError("TypeError: pencolor() argument #1 expected to be `string` got `%s`", args[0].Type())
	}

	if err := typing.Check("pencolor", args, typing.RangeOfArgs(3, 4)); err != nil {
		return newError(err.Error())
	}

	r, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: pencolor() argument #1 `r` %s", err.Error())
	}

	g, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: pencolor() argument #2 `g` %s", err.Error())
	}

	b, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: pencolor() argument #3 `b` %s", err.Error())
	}

	switch len(args) {
	case 3:
		env.Canvas().Value.Graphics().SetRGB(r, g, b)
	case 4:
		a, err := typing.ToFloat(args[3])
		if err != nil {
			return newError("TypeError: pencolor() argument #4 `a` %s", err.Error())
		}
		env.Canvas().Value.Graphics().SetRGBA(r, g, b, a)
	}

	return &object.Null{}
}

// PenSize returns or sets the pen line thickness.
// pensize(width) - sets the pen line thickness to `width`.
func PenSize(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("pensize", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	lw, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: pensize() argument #1 %s", err.Error())
	}

	env.Canvas().Value.Graphics().SetLineWidth(lw)
	return &object.Null{}
}

// Stroke strokes the current path with the current color and line width
// the path is cleared after this operation.
// If preserve is true the path will be preserved.
func Stroke(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("stroke", args, typing.RangeOfArgs(0, 1)); err != nil {
		return newError(err.Error())
	}

	preserve := false
	if (len(args) > 0) && (args[0].Type() == object.BOOLEAN) {
		preserve = args[0].(*object.Boolean).Value
	}

	if preserve {
		env.Canvas().Value.Graphics().StrokePreserve()
	} else {
		env.Canvas().Value.Graphics().Stroke()
	}

	return &object.Null{}
}

// Fill fills the current path with the current color.
// Open subpaths are implicity closed. The path is cleared after this operation.
// If preserve is true the path is preserved after this operation.
func Fill(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("fill", args, typing.RangeOfArgs(0, 1)); err != nil {
		return newError(err.Error())
	}

	preserve := false
	if (len(args) > 0) && (args[0].Type() == object.BOOLEAN) {
		preserve = args[0].(*object.Boolean).Value
	}

	if preserve {
		env.Canvas().Value.Graphics().FillPreserve()
	} else {
		env.Canvas().Value.Graphics().Fill()
	}

	return &object.Null{}
}

// DrawPoint ...
func DrawPoint(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("point", args, typing.ExactArgs(3)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: point() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: point() argument #2 `y` %s", err.Error())
	}

	r, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: point() argument #3 `r` %s", err.Error())
	}

	dc := env.Canvas().Value.Graphics()
	dc.DrawPoint(x, y, r)

	return &object.Null{}
}

// DrawCircle draws a circle centered at [x, y] coordinates and with the radius `r`.
func DrawCircle(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("circle", args, typing.ExactArgs(3)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: circle() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: circle() argument #2 `y` %s", err.Error())
	}

	r, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: circle() argument #3 `r` %s", err.Error())
	}

	dc := env.Canvas().Value.Graphics()
	dc.DrawCircle(x, y, r)

	return &object.Null{}
}

// DrawEllipse draws an ellipse centered at [x, y] coordinates and with the radii `rx` and `ry`.
func DrawEllipse(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("ellipse", args, typing.RangeOfArgs(3, 4)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: ellipse() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: ellipse() argument #2 `y` %s", err.Error())
	}

	rx, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: ellipse() argument #3 `rx` %s", err.Error())
	}

	if len(args) == 3 {
		dc := env.Canvas().Value.Graphics()
		dc.DrawCircle(x, y, rx)
		return &object.Null{}
	}

	ry, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: ellipse() argument #4 `ry` %s", err.Error())
	}

	dc := env.Canvas().Value.Graphics()
	dc.DrawEllipse(x, y, rx, ry)
	return &object.Null{}
}

// DrawRoundedRectangle draws a (w x h) rectangle with upper left corner located at (x, y).
// rectangle(x, y, w, h, [r]) if radius `r` is specified, the rectangle will have rounded corners.
func DrawRoundedRectangle(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("rectangle", args, typing.RangeOfArgs(4, 5)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: rectangle() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: rectangle() argument #2 `y` %s", err.Error())
	}

	w, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: rectangle() argument #3 `w` %s", err.Error())
	}

	h, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: rectangle() argument #4 `h` %s", err.Error())
	}

	if len(args) == 4 {
		dc := env.Canvas().Value.Graphics()
		dc.DrawRectangle(x, y, w, h)
		return &object.Null{}
	}

	r, err := typing.ToFloat(args[4])
	if err != nil {
		return newError("TypeError: rectangle() argument #5 `r` %s", err.Error())
	}

	dc := env.Canvas().Value.Graphics()
	dc.DrawRoundedRectangle(x, y, w, h, r)
	return &object.Null{}
}

// DrawRegularPolygon draws a regular polygon of `n` sides, centered at (x,y) with radius `r` and `angle` rotation
func DrawRegularPolygon(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("polygon", args, typing.ExactArgs(5)); err != nil {
		return newError(err.Error())
	}

	n, err := typing.ToInt(args[0])
	if err != nil {
		return newError("TypeError: polygon() argument #1 `n` %s", err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: polygon() argument #2 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: polygon() argument #3 `y` %s", err.Error())
	}

	r, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: polygon() argument #4 `r` %s", err.Error())
	}

	rad, err := typing.ToFloat(args[4])
	if err != nil {
		return newError("TypeError: polygon() argument #5 `deg` %s", err.Error())
	}

	dc := env.Canvas().Value.Graphics()
	dc.DrawRegularPolygon(n, x, y, r, rad)

	return &object.Null{}
}

// MoveTo starts a new subpath within the current path starting at the specified point.
func MoveTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("moveTo", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: moveTo() argument #1 %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: moveTo() argument #2 %s", err.Error())
	}

	env.Canvas().Value.Graphics().MoveTo(x, y)
	return &object.Null{}
}

// LineTo adds a line segment to the current path starting at the current point.
// If there is no current point, it is equivalent to MoveTo(x, y)
func LineTo(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("lineTo", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: lineTo() argument #1 %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: lineTo() argument #2 %s", err.Error())
	}

	env.Canvas().Value.Graphics().LineTo(x, y)
	return &object.Null{}
}

// DrawLine draws a line from point (x1, y1) to point (x2, y2)
func DrawLine(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("line", args, typing.ExactArgs(4)); err != nil {
		return newError(err.Error())
	}

	x1, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: line() argument #1 `x1` %s", err.Error())
	}

	y1, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: line() argument #2 `y1` %s", err.Error())
	}

	x2, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: line() argument #3 `x2` %s", err.Error())
	}

	y2, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: line() argument #4 `y2` %s", err.Error())
	}

	env.Canvas().Value.Graphics().DrawLine(x1, y1, x2, y2)
	return &object.Null{}
}

// DrawArc draws a circular arc centered at `x, y` with a radius of `r`.
// The path starts at `angle1`, ends at `angle2`, and travels in the direction given by anticlockwise.
func DrawArc(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("arc", args, typing.ExactArgs(5)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: arc() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: arc() argument #2 `y` %s", err.Error())
	}

	r, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: arc() argument #3 `r` %s", err.Error())
	}

	rad1, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: arc() argument #4 `degrees1` %s", err.Error())
	}

	rad2, err := typing.ToFloat(args[4])
	if err != nil {
		return newError("TypeError: arc() argument #4 `degrees2` %s", err.Error())
	}

	env.Canvas().Value.Graphics().DrawArc(x, y, r, rad1, rad2)
	return &object.Null{}
}

// DrawEllipticalArc draws an elliptical arc centered at `x, y` with a radius of `rx` in x direction and `ry` for y direction.
// The path starts at `angle1`, ends at `angle2`, and travels in the direction given by anticlockwise.
func DrawEllipticalArc(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("ellipticalArc", args, typing.ExactArgs(6)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #2 `y` %s", err.Error())
	}

	rx, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #3 `rx` %s", err.Error())
	}

	ry, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #4 `ry` %s", err.Error())
	}

	rad1, err := typing.ToFloat(args[4])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #5 `degrees1` %s", err.Error())
	}

	rad2, err := typing.ToFloat(args[5])
	if err != nil {
		return newError("TypeError: ellipticalArc() argument #6 `degrees2` %s", err.Error())
	}

	env.Canvas().Value.Graphics().DrawEllipticalArc(x, y, rx, ry, rad1, rad2)
	return &object.Null{}
}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
func ClosePath(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Graphics().ClosePath()
	return &object.Null{}
}

// ClearPath clears the current path. There is no current point after this
// operation.
func ClearPath(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Graphics().ClearPath()
	return &object.Null{}
}

// DrawStringAnchored draws the specified text at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
func DrawStringAnchored(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("text", args, typing.RangeOfArgs(3, 5)); err != nil {
		return newError(err.Error())
	}

	txt, err := typing.ToString(args[0])
	if err != nil {
		return newError("TypeError: text() argument #1 `string` %s", err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: text() argument #2 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: text() argument #3 `y` %s", err.Error())
	}

	ax, ay := 0.5, 0.5
	if len(args) == 5 {
		if ax, err = typing.ToFloat(args[3]); err != nil {
			return newError("TypeError: text() argument #4 `ax` %s", err.Error())
		}

		if ay, err = typing.ToFloat(args[4]); err != nil {
			return newError("TypeError: text() argument #5 `ay` %s", err.Error())
		}
	}

	dc := env.Canvas().Value.Graphics()
	dc.Push()
	dc.DrawStringAnchored(txt, x, y, ax, ay)
	dc.Pop()

	return &object.Null{}
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func MeasureString(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("measureText", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	txt, err := typing.ToString(args[0])
	if err != nil {
		return newError("TypeError: measureText() argument #1 `str` %s", err.Error())
	}

	w, h := env.Canvas().Value.Graphics().MeasureString(txt)
	return &object.Array{
		Elements: []object.Object{
			&object.Float{Value: w},
			&object.Float{Value: h},
		},
	}
}

// FontSize returns or sets the font size.
// fontsize() - returns the font current size.
// fontsize(size) - sets the font current size.
func FontSize(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("fontsize", args, typing.RangeOfArgs(0, 1)); err != nil {
		return newError(err.Error())
	}

	if len(args) == 0 {
		size := env.Canvas().Value.FontSize()
		return &object.Float{Value: size}
	}

	fs, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: fontsize() argument #1 %s", err.Error())
	}

	env.Canvas().Value.SetFontSize(fs)
	return &object.Null{}
}

// Snapshot creates a png image with the current drawings.
// snapshot() - saves the png image in the source code folder.
func Snapshot(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("input", args, typing.RangeOfArgs(0, 1), typing.WithTypes(object.STRING)); err != nil {
		return newError(err.Error())
	}

	if len(args) == 1 {
		filename := args[0].(*object.String).Value
		dc := env.Canvas().Value.Graphics()
		if err := dc.SavePNG(filename); err != nil {
			return newError(err.Error())
		}

		return &object.Null{}
	}

	object.SaveCounter = object.SaveCounter + 1

	folder := filepath.Join(object.WorkDir, object.SourceFile)
	if err := mkdir(folder); err != nil {
		return newError(err.Error())
	}

	filename := filepath.Join(folder,
		fmt.Sprintf("%s_%03d.png", object.SourceFile, object.SaveCounter))

	dc := env.Canvas().Value.Graphics()
	if err := dc.SavePNG(filename); err != nil {
		return newError(err.Error())
	}

	return &object.Null{}
}

// SaveState saves the current state of the canvas by pushin it onto a stack. These can be nested.
func SaveState(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Graphics().Push()
	return &object.Null{}
}

// RestoreState restores the last saved canvas state from the stack.
func RestoreState(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Graphics().Pop()
	return &object.Null{}
}

// Clear fills the entire image with the current color.
func Clear(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Graphics().Clear()
	return &object.Null{}
}

// RotateAbout updates the current matrix with a anticlockwise rotation.
// rotate(angle) - rotation occurs about the origin.
// rotate(angle, x, y) - rotation occurs about the specified point.
// Angle is specified in radians.
func RotateAbout(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return newError("rotate() expects one or three arguments")
	}

	rad, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: rotate() argument #1 `angle` %s", err.Error())
	}

	if len(args) == 1 {
		env.Canvas().Value.Graphics().Rotate(rad)
		return &object.Null{}
	}

	if err := typing.Check("rotate", args, typing.ExactArgs(3)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: rotate() argument #2 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: rotate() argument #3 `y` %s", err.Error())
	}

	env.Canvas().Value.Graphics().RotateAbout(rad, x, y)
	return &object.Null{}
}

// ScaleAbout updates the current matrix with a scaling factor.
// scale(sx, sy) - scaling occurs about the origin.
// scale(sx, sy, x, y) - scaling occurs about the specified point.
func ScaleAbout(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 2 {
		return newError("scale() expects two or four arguments")
	}

	sx, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: scale() argument #1 `sx` %s", err.Error())
	}

	sy, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: scale() argument #2 `sy` %s", err.Error())
	}

	if len(args) == 2 {
		env.Canvas().Value.Graphics().Scale(sx, sy)
		return &object.Null{}
	}

	if err := typing.Check("scale", args, typing.ExactArgs(4)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: scale() argument #3 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: scale() argument #4 `y` %s", err.Error())
	}

	env.Canvas().Value.Graphics().ScaleAbout(sx, sy, x, y)
	return &object.Null{}
}

// Translate updates the current matrix with a translation.
func Translate(env *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("translate", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: translate() argument #1 `x` %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: translate() argument #2 `y` %s", err.Error())
	}

	env.Canvas().Value.Graphics().Translate(x, y)
	return &object.Null{}
}

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func Identity(env *object.Environment, args ...object.Object) object.Object {
	env.Canvas().Value.Graphics().Identity()
	return &object.Null{}
}

func mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return nil
}
