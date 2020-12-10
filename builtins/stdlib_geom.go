package builtins

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// PartitionLine partition a line into equals parts.
// partitionLine(x1, y1, x2, y2, k) divides the line (x1,y1) - (x2,y2) in k equals parts
func PartitionLine(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("partitionLine", args, typing.ExactArgs(5)); err != nil {
		return newError(err.Error())
	}

	x1, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: partitionLine() argument #1 `x1` %s", err.Error())
	}

	y1, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: partitionLine() argument #2 `y1` %s", err.Error())
	}

	x2, err := typing.ToFloat(args[2])
	if err != nil {
		return newError("TypeError: partitionLine() argument #3 `x1` %s", err.Error())
	}

	y2, err := typing.ToFloat(args[3])
	if err != nil {
		return newError("TypeError: partitionLine() argument #4 `y1` %s", err.Error())
	}

	parts, err := typing.ToInt(args[4])
	if err != nil {
		return newError("TypeError: partitionLine() argument #5 `parts` %s", err.Error())
	}
	if parts <= 0 {
		return newError("TypeError: partitionLine() argument #5 `parts` must be >= 0")
	}

	res := &object.Array{
		Elements: []object.Object{},
	}

	// Cx = Ax * (1-t) + Bx * t
	// Cy = Ay * (1-t) + By * t
	//
	// When t=0, you get A.
	// When t=1, you get B.
	// When t=.25, you a point 25% of the way from A to B
	// etc.
	// So, to divide the line into k equal parts, make a loop and find C, for t=0/k, t=1/k, t=2/k, ... , t=k/k
	for i := 0; i <= parts; i++ {
		t := float64(i) / float64(parts)

		x := x1*(1-t) + x2*t
		y := y1*(1-t) + y2*t

		pt := &object.Array{
			Elements: []object.Object{
				&object.Float{Value: x},
				&object.Float{Value: y},
			},
		}

		res.Elements = append(res.Elements, pt)
	}

	return res
}
