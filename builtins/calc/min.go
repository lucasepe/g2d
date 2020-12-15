package calc

import (
	"math"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Min returns the smalles value of a sequence of numbers.
// If the length of the sequence is zero, returns the smallest non zero float64 value.
func Min(_ *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		return &object.Float{Value: math.SmallestNonzeroFloat64}
	}

	if len(args) == 1 {
		array, err := typing.ToFloatArray(args[0])
		if err != nil {
			return object.NewError(err.Error())
		}

		return &object.Float{Value: minOf(array)}
	}

	array := make([]float64, len(args))
	for i, el := range args {
		val, err := typing.ToFloat(el)
		if err != nil {
			return object.NewError("TypeError: min() argument #%d %s", i, err.Error())
		}
		array[i] = val
	}

	return &object.Float{Value: minOf(array)}
}

func minOf(array []float64) float64 {
	if len(array) == 0 {
		return math.SmallestNonzeroFloat64
	}

	res := array[0]
	for _, el := range array {
		if el < res {
			res = el
		}
	}

	return res
}
