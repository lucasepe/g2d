package typing

import (
	"fmt"
	"math"

	"github.com/lucasepe/g2d/object"
)

type CheckFunc func(name string, args []object.Object) error

func Check(name string, args []object.Object, checks ...CheckFunc) error {
	for _, check := range checks {
		if err := check(name, args); err != nil {
			return err
		}
	}
	return nil
}

func ExactArgs(n int) CheckFunc {
	return func(name string, args []object.Object) error {
		if len(args) != n {
			return fmt.Errorf(
				"TypeError: %s() takes exactly %d argument (%d given)",
				name, n, len(args),
			)
		}
		return nil
	}
}

func MinimumArgs(n int) CheckFunc {
	return func(name string, args []object.Object) error {
		if len(args) < n {
			return fmt.Errorf(
				"TypeError: %s() takes a minimum %d arguments (%d given)",
				name, n, len(args),
			)
		}
		return nil
	}
}

func RangeOfArgs(n, m int) CheckFunc {
	return func(name string, args []object.Object) error {
		if len(args) < n || len(args) > m {
			return fmt.Errorf(
				"TypeError: %s() takes at least %d arguments at most %d (%d given)",
				name, n, m, len(args),
			)
		}
		return nil
	}
}

func WithTypes(types ...object.Type) CheckFunc {
	return func(name string, args []object.Object) error {
		for i, t := range types {
			if i < len(args) && args[i].Type() != t {
				return fmt.Errorf(
					"TypeError: %s() expected argument #%d to be `%s` got `%s`",
					name, (i + 1), t, args[i].Type(),
				)
			}
		}
		return nil
	}
}

func ToFloat(obj object.Object) (float64, error) {
	if obj.Type() == object.INTEGER {
		val := obj.(*object.Integer).Value
		return float64(val), nil
	}

	if obj.Type() == object.FLOAT {
		val := obj.(*object.Float).Value
		return val, nil
	}

	return math.NaN(), fmt.Errorf("expected to be `int` or `float` got `%s`", obj.Type())
}

func ToInt(obj object.Object) (int, error) {
	if obj.Type() == object.INTEGER {
		val := obj.(*object.Integer).Value
		return int(val), nil
	}

	if obj.Type() == object.FLOAT {
		val := obj.(*object.Float).Value
		return int(math.Round(val)), nil
	}

	return -1, fmt.Errorf("expected to be `int` got `%s`", obj.Type())
}

func ToString(obj object.Object) (string, error) {
	if obj.Type() == object.STRING {
		val := obj.(*object.String).Value
		return val, nil
	}

	return "", fmt.Errorf("expected to be `string` got `%s`", obj.Type())
}

func ToBool(obj object.Object) (bool, error) {
	if obj.Type() == object.BOOLEAN {
		val := obj.(*object.Boolean).Value
		return val, nil
	}

	return false, fmt.Errorf("expected to be `bool` got `%s`", obj.Type())
}
