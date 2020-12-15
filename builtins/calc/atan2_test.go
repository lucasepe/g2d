package calc

import (
	"fmt"
	"testing"

	"github.com/lucasepe/g2d/object"
)

func TestAtan2(t *testing.T) {
	cases := []struct {
		input []object.Object
		want  object.Object
	}{
		{
			input: []object.Object{
				&object.Integer{Value: 10},
				&object.Float{Value: 4.9790119248836735e+00},
			},
			want: &object.Float{Value: 1.1088291730037003},
		},

		{
			input: []object.Object{
				&object.Integer{Value: 10},
				&object.Float{Value: -2.7688005719200159e-01},
			},
			want: &object.Float{Value: 1.5984772603216203736068915e+00},
		},

		{
			input: []object.Object{
				&object.Integer{Value: 10},
				&object.Float{Value: -5.0106036182710749e+00},
			},
			want: &object.Float{Value: 2.0352918654092086637227327e+00},
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("atan2_%d", i), func(t *testing.T) {
			got := Atan2(nil, tt.input...)
			if got.(object.Comparable).Compare(tt.want) != 0 {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}
