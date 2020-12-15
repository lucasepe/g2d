package calc

import (
	"fmt"
	"testing"

	"github.com/lucasepe/g2d/object"
)

func TestAbs(t *testing.T) {
	cases := []struct {
		in   object.Object
		want object.Object
	}{
		{&object.Float{Value: -4.933}, &object.Float{Value: 4.933}},
		{&object.Integer{Value: -2}, &object.Integer{Value: 2}},
		{&object.Float{Value: -44}, &object.Float{Value: 44}},
		{&object.Integer{Value: 8}, &object.Integer{Value: 8}},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("abs_%d", i), func(t *testing.T) {
			got := Abs(nil, tt.in)
			if got.(object.Comparable).Compare(tt.want) != 0 {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}
