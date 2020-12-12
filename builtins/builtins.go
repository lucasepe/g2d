package builtins

import (
	"fmt"
	"sort"

	"github.com/lucasepe/g2d/object"
)

// Builtins ...
var Builtins = map[string]*object.Builtin{
	// Core
	"args":    &object.Builtin{Name: "args", Fn: Args},
	"exit":    &object.Builtin{Name: "exit", Fn: Exit},
	"input":   &object.Builtin{Name: "input", Fn: Input},
	"print":   &object.Builtin{Name: "print", Fn: Print},
	"printf":  &object.Builtin{Name: "printf", Fn: Printf},
	"sprintf": &object.Builtin{Name: "sprintf", Fn: Sprintf},
	"bool":    &object.Builtin{Name: "bool", Fn: Bool},
	"float":   &object.Builtin{Name: "float", Fn: Float},
	"int":     &object.Builtin{Name: "int", Fn: Int},
	"str":     &object.Builtin{Name: "str", Fn: Str},
	"len":     &object.Builtin{Name: "len", Fn: Len},
	"push":    &object.Builtin{Name: "push", Fn: Push},
	"keys":    &object.Builtin{Name: "keys", Fn: hashKeys},
	"delete":  &object.Builtin{Name: "delete", Fn: hashDelete},
	"type":    &object.Builtin{Name: "type", Fn: TypeOf},

	// Math
	"abs":     &object.Builtin{Name: "abs", Fn: Abs},
	"atan":    &object.Builtin{Name: "atan", Fn: Atan},
	"atan2":   &object.Builtin{Name: "atan2", Fn: Atan2},
	"cos":     &object.Builtin{Name: "cos", Fn: Cos},
	"hypot":   &object.Builtin{Name: "hypot", Fn: Hypot},
	"pow":     &object.Builtin{Name: "pow", Fn: Pow},
	"sin":     &object.Builtin{Name: "sin", Fn: Sin},
	"sqrt":    &object.Builtin{Name: "sqrt", Fn: Sqrt},
	"randf":   &object.Builtin{Name: "randf", Fn: RandomFloat},
	"randi":   &object.Builtin{Name: "randi", Fn: RandomInt},
	"radians": &object.Builtin{Name: "radians", Fn: Radians},
	"degrees": &object.Builtin{Name: "degrees", Fn: Degrees},

	// Canvas
	"screensize":  &object.Builtin{Name: "screensize", Fn: ScreenSize},
	"clear":       &object.Builtin{Name: "clear", Fn: Clear},
	"worldcoords": &object.Builtin{Name: "worldcoords", Fn: WorldCoords},
	"pencolor":    &object.Builtin{Name: "pencolor", Fn: PenColor},
	"pensize":     &object.Builtin{Name: "pensize", Fn: PenSize},
	//"fontsize":    &object.Builtin{Name: "fontsize", Fn: FontSize},
	"stroke": &object.Builtin{Name: "stroke", Fn: Stroke},
	"fill":   &object.Builtin{Name: "fill", Fn: Fill},

	"polygon": &object.Builtin{Name: "polygon", Fn: DrawRegularPolygon},
	"moveTo":  &object.Builtin{Name: "moveTo", Fn: MoveTo},
	"lineTo":  &object.Builtin{Name: "lineTo", Fn: LineTo},

	"ellArc":      &object.Builtin{Name: "ellArc", Fn: DrawEllipticalArc},
	"text":        &object.Builtin{Name: "text", Fn: DrawStringAnchored},
	"drawImage":   &object.Builtin{Name: "drawImage", Fn: DrawImageAnchored},
	"measureText": &object.Builtin{Name: "measureText", Fn: MeasureString},

	"beginPath":    &object.Builtin{Name: "beginPath", Fn: BeginPath},
	"closePath":    &object.Builtin{Name: "closePath", Fn: ClosePath},
	"saveState":    &object.Builtin{Name: "saveState", Fn: SaveState},
	"restoreState": &object.Builtin{Name: "restoreState", Fn: RestoreState},
	"rotate":       &object.Builtin{Name: "rotate", Fn: RotateAbout},
	"scale":        &object.Builtin{Name: "scale", Fn: ScaleAbout},
	"translate":    &object.Builtin{Name: "translate", Fn: Translate},
	"identity":     &object.Builtin{Name: "identity", Fn: Identity},

	// 2D Primitives
	"arc":      &object.Builtin{Name: "arc", Fn: DrawArc},
	"circle":   &object.Builtin{Name: "circle", Fn: DrawCircle},
	"ellipse":  &object.Builtin{Name: "ellipse", Fn: DrawEllipse},
	"line":     &object.Builtin{Name: "line", Fn: DrawLine},
	"point":    &object.Builtin{Name: "point", Fn: DrawPoint},
	"quad":     &object.Builtin{Name: "quad", Fn: DrawQuadrilateral},
	"rect":     &object.Builtin{Name: "rect", Fn: DrawRoundedRectangle},
	"triangle": &object.Builtin{Name: "triangle", Fn: DrawTriangle},

	"snapshot": &object.Builtin{Name: "snapshot", Fn: Snapshot},
	"loadPNG":  &object.Builtin{Name: "loadPNG", Fn: LoadPNG},
	"width":    &object.Builtin{Name: "width", Fn: ImageWidth},
	"height":   &object.Builtin{Name: "height", Fn: ImageHeight},
}

// BuiltinsIndex ...
var BuiltinsIndex []*object.Builtin

func init() {
	var keys []string
	for k := range Builtins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		BuiltinsIndex = append(BuiltinsIndex, Builtins[k])
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
