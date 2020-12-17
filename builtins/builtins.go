package builtins

import (
	"fmt"
	"sort"

	"github.com/lucasepe/g2d/builtins/calc"
	"github.com/lucasepe/g2d/builtins/core"
	"github.com/lucasepe/g2d/builtins/graphics"
	"github.com/lucasepe/g2d/object"
)

// Builtins ...
var Builtins = map[string]*object.Builtin{
	// Core
	"args":    &object.Builtin{Name: "args", Fn: core.Args},
	"exit":    &object.Builtin{Name: "exit", Fn: core.Exit},
	"input":   &object.Builtin{Name: "input", Fn: core.Input},
	"print":   &object.Builtin{Name: "print", Fn: core.Print},
	"printf":  &object.Builtin{Name: "printf", Fn: core.Printf},
	"sprintf": &object.Builtin{Name: "sprintf", Fn: core.Sprintf},
	"bool":    &object.Builtin{Name: "bool", Fn: core.Bool},
	"float":   &object.Builtin{Name: "float", Fn: core.Float},
	"int":     &object.Builtin{Name: "int", Fn: core.Int},
	"str":     &object.Builtin{Name: "str", Fn: core.Str},
	"len":     &object.Builtin{Name: "len", Fn: core.Len},
	"append":  &object.Builtin{Name: "append", Fn: core.Append},
	"keys":    &object.Builtin{Name: "keys", Fn: core.Keys},
	"delete":  &object.Builtin{Name: "delete", Fn: core.Delete},
	"type":    &object.Builtin{Name: "type", Fn: core.TypeOf},

	// Calculation
	"abs":     &object.Builtin{Name: "abs", Fn: calc.Abs},
	"atan":    &object.Builtin{Name: "atan", Fn: calc.Atan},
	"atan2":   &object.Builtin{Name: "atan2", Fn: calc.Atan2},
	"cos":     &object.Builtin{Name: "cos", Fn: calc.Cos},
	"degrees": &object.Builtin{Name: "degrees", Fn: calc.Degrees},
	"hypot":   &object.Builtin{Name: "hypot", Fn: calc.Hypot},
	"lerp":    &object.Builtin{Name: "lerp", Fn: calc.Lerp},
	"map":     &object.Builtin{Name: "map", Fn: calc.Map},
	"max":     &object.Builtin{Name: "max", Fn: calc.Max},
	"min":     &object.Builtin{Name: "min", Fn: calc.Min},
	"pow":     &object.Builtin{Name: "pow", Fn: calc.Pow},
	"radians": &object.Builtin{Name: "radians", Fn: calc.Radians},
	"randf":   &object.Builtin{Name: "randf", Fn: calc.RandomFloat},
	"randi":   &object.Builtin{Name: "randi", Fn: calc.RandomInt},
	"sin":     &object.Builtin{Name: "sin", Fn: calc.Sin},
	"sqrt":    &object.Builtin{Name: "sqrt", Fn: calc.Sqrt},

	// Graphic Context
	"size":     &object.Builtin{Name: "size", Fn: graphics.Size},
	"clear":    &object.Builtin{Name: "clear", Fn: graphics.Clear},
	"dashes":   &object.Builtin{Name: "dashes", Fn: graphics.Dashes},
	"pencolor": &object.Builtin{Name: "pencolor", Fn: graphics.PenColor},
	"pensize":  &object.Builtin{Name: "pensize", Fn: graphics.PenSize},
	"xpos":     &object.Builtin{Name: "xpos", Fn: graphics.GetCurrentX},
	"ypos":     &object.Builtin{Name: "ypos", Fn: graphics.GetCurrentY},
	"snapshot": &object.Builtin{Name: "snapshot", Fn: graphics.Snapshot},
	"width":    &object.Builtin{Name: "width", Fn: graphics.Width},
	"height":   &object.Builtin{Name: "height", Fn: graphics.Height},
	"push":     &object.Builtin{Name: "push", Fn: graphics.Push},
	"pop":      &object.Builtin{Name: "pop", Fn: graphics.Pop},
	"stroke":   &object.Builtin{Name: "stroke", Fn: graphics.Stroke},
	"fill":     &object.Builtin{Name: "fill", Fn: graphics.Fill},
	"viewport": &object.Builtin{Name: "viewport", Fn: graphics.Viewport},

	// Path
	"beginPath":        &object.Builtin{Name: "beginPath", Fn: graphics.BeginPath},
	"closePath":        &object.Builtin{Name: "closePath", Fn: graphics.ClosePath},
	"quadraticCurveTo": &object.Builtin{Name: "quadraticCurveTo", Fn: graphics.QuadraticCurveTo},
	"arcTo":            &object.Builtin{Name: "arcTo", Fn: graphics.ArcTo},
	"lineTo":           &object.Builtin{Name: "lineTo", Fn: graphics.LineTo},
	"moveTo":           &object.Builtin{Name: "moveTo", Fn: graphics.MoveTo},
	"routeTo":          &object.Builtin{Name: "routeTo", Fn: graphics.RouteTo},

	// Transform
	"rotate":    &object.Builtin{Name: "rotate", Fn: graphics.RotateAbout},
	"scale":     &object.Builtin{Name: "scale", Fn: graphics.ScaleAbout},
	"translate": &object.Builtin{Name: "translate", Fn: graphics.Translate},
	"identity":  &object.Builtin{Name: "identity", Fn: graphics.Identity},
	"transform": &object.Builtin{Name: "transform", Fn: graphics.Transform},

	// 2D Primitives
	"arc":      &object.Builtin{Name: "arc", Fn: graphics.Arc},
	"circle":   &object.Builtin{Name: "circle", Fn: graphics.Circle},
	"ellipse":  &object.Builtin{Name: "ellipse", Fn: graphics.Ellipse},
	"line":     &object.Builtin{Name: "line", Fn: graphics.Line},
	"point":    &object.Builtin{Name: "point", Fn: graphics.Point},
	"quad":     &object.Builtin{Name: "quad", Fn: graphics.Quad},
	"rect":     &object.Builtin{Name: "rect", Fn: graphics.Rect},
	"triangle": &object.Builtin{Name: "triangle", Fn: graphics.Triangle},

	// Text
	"text":       &object.Builtin{Name: "text", Fn: graphics.Text},
	"textWidth":  &object.Builtin{Name: "textWidth", Fn: graphics.TextWidth},
	"fontHeight": &object.Builtin{Name: "fontHeight", Fn: graphics.FontHeight},

	// Images
	"loadPNG": &object.Builtin{Name: "loadPNG", Fn: graphics.LoadPNG},
	"image":   &object.Builtin{Name: "image", Fn: graphics.ImageAnchored},
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
