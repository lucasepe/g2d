package calc

import (
	"fmt"
	"testing"

	"github.com/lucasepe/g2d/object"
)

func TestLerp(t *testing.T) {
	cases := []struct {
		input []object.Object
		want  object.Object
	}{
		{
			input: []object.Object{
				&object.Integer{Value: 20},
				&object.Float{Value: 80},
				&object.Float{Value: 0.2},
			},
			want: &object.Float{Value: 32},
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("lerp_%d", i), func(t *testing.T) {
			got := Lerp(nil, tt.input...)
			if got.(object.Comparable).Compare(tt.want) != 0 {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}
