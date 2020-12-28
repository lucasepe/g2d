package gg

import (
	"image"
	"image/color"
)

type RepeatOp int

const (
	RepeatBoth RepeatOp = iota
	RepeatX
	RepeatY
	RepeatNone
)

type Pattern interface {
	ColorAt(x, y int) color.Color
}

// Solid Pattern
type SolidPattern struct {
	Color color.Color
}

func (p *SolidPattern) ColorAt(x, y int) color.Color {
	return p.Color
}

func NewSolidPattern(color color.Color) Pattern {
	return &SolidPattern{Color: color}
}

// Surface Pattern
type SurfacePattern struct {
	im image.Image
	op RepeatOp
}

func (p *SurfacePattern) ColorAt(x, y int) color.Color {
	b := p.im.Bounds()
	switch p.op {
	case RepeatX:
		if y >= b.Dy() {
			return color.Transparent
		}
	case RepeatY:
		if x >= b.Dx() {
			return color.Transparent
		}
	case RepeatNone:
		if x >= b.Dx() || y >= b.Dy() {
			return color.Transparent
		}
	}
	x = x%b.Dx() + b.Min.X
	y = y%b.Dy() + b.Min.Y
	return p.im.At(x, y)
}

func NewSurfacePattern(im image.Image, op RepeatOp) Pattern {
	return &SurfacePattern{im: im, op: op}
}
