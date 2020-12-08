package builtins

import (
	"math"
	"math/rand"
	"time"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Abs returns the absolute value of x.
func Abs(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("abs", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	if args[0].Type() == object.INTEGER {
		value := args[0].(*object.Integer).Value
		if value < 0 {
			value = value * -1
		}
		return &object.Integer{Value: value}
	}

	if args[0].Type() == object.FLOAT {
		value := args[0].(*object.Float).Value
		if value < 0 {
			value = value * -1
		}
		return &object.Float{Value: value}
	}

	return newError("TypeError: abs() argument #1 expected to be `int` or `float` got `%s`", args[0].Type())
}

// Sqrt returns the square root of x.
func Sqrt(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("sqrt", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: sqrt() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Sqrt(val)}
}

// Hypot returns Sqrt(p*p + q*q), taking care to avoid
// unnecessary overflow and underflow.
func Hypot(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("hypot", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	p, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: sqrt() argument #1 %s", err.Error())
	}

	q, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: sqrt() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Hypot(p, q)}
}

// Pow returns x**y, the base-x exponential of y.
func Pow(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("pow", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	x, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: pow() argument #1 %s", err.Error())
	}

	y, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: pow() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Pow(x, y)}
}

// Atan returns the arctangent, in radians, of x.
func Atan(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("atan", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: atan() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Atan(val)}
}

// Atan2 returns the arc tangent of y/x, using
// the signs of the two to determine the quadrant
// of the return value.
func Atan2(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("atan2", args, typing.ExactArgs(2)); err != nil {
		return newError(err.Error())
	}

	y, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: atan2() argument #1 %s", err.Error())
	}

	x, err := typing.ToFloat(args[1])
	if err != nil {
		return newError("TypeError: atan2() argument #2 %s", err.Error())
	}

	return &object.Float{Value: math.Atan2(y, x)}
}

// Sin returns the sine of the radian argument x.
func Sin(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("sin", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: sin() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Sin(val)}
}

// Cos returns the cosine of the radian argument x.
func Cos(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("cos", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	val, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: cos() argument #1 %s", err.Error())
	}

	return &object.Float{Value: math.Cos(val)}
}

// RandomFloat returns a random float.
// randf() returns a random float between 0.0 and 1.0
// randf(max) returns a random float between 0.0 and max
// randf(min, max) returns a random float between min and max
func RandomFloat(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("randf", args, typing.RangeOfArgs(0, 2)); err != nil {
		return newError(err.Error())
	}

	if len(args) == 1 {
		max, err := typing.ToFloat(args[0])
		if err != nil {
			return newError("TypeError: randf() argument #1 %s", err.Error())
		}
		if max <= 0 {
			return newError("ValueError: randf() argument #1 must be > 0")
		}
		return &object.Float{Value: rand.Float64() * max}
	}

	if len(args) == 2 {
		min, err := typing.ToFloat(args[0])
		if err != nil {
			return newError("TypeError: randf() argument #1 `min` %s", err.Error())
		}

		max, err := typing.ToFloat(args[1])
		if err != nil {
			return newError("TypeError: randf() argument #2 `max` %s", err.Error())
		}

		if max < min {
			return newError("ValueError: randf() argument #1 `min` must be > argument #2 `max`")
		}
		return &object.Float{
			Value: min + rand.Float64()*(max-min),
		}
	}

	return &object.Float{Value: rand.Float64()}
}

// RandomInt returns a ramdom int
// randi() returns a random integer
// randi(max) returns a random integer between 0 and max
// randi(min, max) returns a random integer between min and max
func RandomInt(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("randi", args,
		typing.RangeOfArgs(0, 2),
		typing.WithTypes(object.INTEGER, object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	if len(args) == 1 {
		max := args[0].(*object.Integer).Value
		if max <= 0 {
			return newError("ValueError: randi() argument #1 must be > 0")
		}
		return &object.Integer{Value: rand.Int63n(max + 1)}
	}

	if len(args) == 2 {
		min := args[0].(*object.Integer).Value
		max := args[1].(*object.Integer).Value
		if min > max {
			return newError("ValueError: randi() argument #1 `min` must be > argument #2 `max`")
		}
		return &object.Integer{Value: rand.Int63n(max-min+1) + min}
	}

	return &object.Integer{Value: rand.Int63()}
}

// Radians converts degrees to radians.
// radians(angle) - angle in degrees is the value that you want to convert
func Radians(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("radians", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	degrees, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: radians() argument #1 %s", err.Error())
	}

	res := degrees * math.Pi / 180.0
	return &object.Float{Value: res}
}

// Degrees converts radians into degrees.
// degrees(angle) - angle in radians is the value that you want to convert
func Degrees(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("degrees", args, typing.ExactArgs(1)); err != nil {
		return newError(err.Error())
	}

	radians, err := typing.ToFloat(args[0])
	if err != nil {
		return newError("TypeError: degrees() argument #1 %s", err.Error())
	}

	res := radians * 180.0 / math.Pi
	return &object.Float{Value: res}
}
