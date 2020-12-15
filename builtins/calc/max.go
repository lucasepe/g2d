package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Max returns the greatest value of a sequence of numbers.
// If the length of the sequence is zero, returns the greatest float64 value.
func Max(_ *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		return &object.Float{Value: math.MaxFloat64}
	}

	if len(args) == 1 {
		array, err := typing.ToFloatArray(args[0])
		if err != nil {
			return object.NewError(err.Error())
		}

		return &object.Float{Value: maxOf(array)}
	}

	array := make([]float64, len(args))
	for i, el := range args {
		val, err := typing.ToFloat(el)
		if err != nil {
			return object.NewError("TypeError: max() argument #%d %s", i, err.Error())
		}
		array[i] = val
	}

	return &object.Float{Value: maxOf(array)}
}

func maxOf(array []float64) float64 {
	if len(array) == 0 {
		return math.MaxFloat64
	}

	res := array[0]
	for _, el := range array {
		if el > res {
			res = el
		}
	}

	return res
}
