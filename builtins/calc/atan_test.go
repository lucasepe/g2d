package calc

import (
	"fmt"
	"testing"

	"github.com/lucasepe/g2d/object"
)

func TestAtan(t *testing.T) {
	cases := []struct {
		in   object.Object
		want object.Object
	}{
		{&object.Float{Value: 4.9790119248836735e+00}, &object.Float{Value: 1.372590262129621651920085e+00}},
		{&object.Float{Value: 7.7388724745781045e+002}, &object.Float{Value: 1.5695041496118367}},
		{&object.Integer{Value: 3}, &object.Float{Value: 1.2490457723982544}},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("atan_%d", i), func(t *testing.T) {
			got := Atan(nil, tt.in)
			if got.(object.Comparable).Compare(tt.want) != 0 {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}
