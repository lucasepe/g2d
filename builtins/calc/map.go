package calc

import (
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

// Map re-maps a number from one range to another.
// establishes a proportion between two ranges of values
func Map(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("map", args, typing.ExactArgs(5)); err != nil {
		return object.NewError(err.Error())
	}

	value, err := typing.ToFloat(args[0])
	if err != nil {
		return object.NewError("TypeError: map() argument #1 `value` %s", err.Error())
	}

	istart, err := typing.ToFloat(args[1])
	if err != nil {
		return object.NewError("TypeError: map() argument #2 `istart` %s", err.Error())
	}

	istop, err := typing.ToFloat(args[2])
	if err != nil {
		return object.NewError("TypeError: map() argument #3 `istop` %s", err.Error())
	}

	ostart, err := typing.ToFloat(args[3])
	if err != nil {
		return object.NewError("TypeError: map() argument #4 `ostart` %s", err.Error())
	}

	ostop, err := typing.ToFloat(args[4])
	if err != nil {
		return object.NewError("TypeError: map() argument #5 `ostop` %s", err.Error())
	}

	res := ostart + (ostop-ostart)*((value-istart)/(istop-istart))
	return &object.Float{Value: res}
}
