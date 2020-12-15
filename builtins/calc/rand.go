package calc

import (
	"math/rand"
	"time"

	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/typing"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFloat returns a random float.
// randf() returns a random float between 0.0 and 1.0
// randf(max) returns a random float between 0.0 and max
// randf(min, max) returns a random float between min and max
func RandomFloat(_ *object.Environment, args ...object.Object) object.Object {
	if err := typing.Check("randf", args, typing.RangeOfArgs(0, 2)); err != nil {
		return object.NewError(err.Error())
	}

	if len(args) == 1 {
		max, err := typing.ToFloat(args[0])
		if err != nil {
			return object.NewError("TypeError: randf() argument #1 %s", err.Error())
		}
		if max <= 0 {
			return object.NewError("ValueError: randf() argument #1 must be > 0")
		}
		return &object.Float{Value: rand.Float64() * max}
	}

	if len(args) == 2 {
		min, err := typing.ToFloat(args[0])
		if err != nil {
			return object.NewError("TypeError: randf() argument #1 `min` %s", err.Error())
		}

		max, err := typing.ToFloat(args[1])
		if err != nil {
			return object.NewError("TypeError: randf() argument #2 `max` %s", err.Error())
		}

		if max < min {
			return object.NewError("ValueError: randf() argument #1 `min` must be > argument #2 `max`")
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
		return object.NewError(err.Error())
	}

	if len(args) == 1 {
		max := args[0].(*object.Integer).Value
		if max <= 0 {
			return object.NewError("ValueError: randi() argument #1 must be > 0")
		}
		return &object.Integer{Value: rand.Int63n(max + 1)}
	}

	if len(args) == 2 {
		min := args[0].(*object.Integer).Value
		max := args[1].(*object.Integer).Value
		if min > max {
			return object.NewError("ValueError: randi() argument #1 `min` must be > argument #2 `max`")
		}
		return &object.Integer{Value: rand.Int63n(max-min+1) + min}
	}

	return &object.Integer{Value: rand.Int63()}
}
