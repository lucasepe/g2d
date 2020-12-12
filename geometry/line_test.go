package geometry

import (
	"math"
	"reflect"
	"testing"
)

func TestLine_Center(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   Vec
	}{
		{
			name:   "Positive slope",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			want:   V(5, 5),
		},
		{
			name:   "Negative slope",
			fields: fields{A: V(10, 10), B: V(0, 0)},
			want:   V(5, 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.Center(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.Center() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Closest(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		v Vec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Vec
	}{
		{
			name:   "Point on line",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{v: V(5, 5)},
			want:   V(5, 5),
		},
		{
			name:   "Point on next to line",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{v: V(0, 10)},
			want:   V(5, 5),
		},
		{
			name:   "Point on next to vertical line",
			fields: fields{A: V(5, 0), B: V(5, 10)},
			args:   args{v: V(6, 5)},
			want:   V(5, 5),
		},
		{
			name:   "Point on next to horizontal line",
			fields: fields{A: V(0, 5), B: V(10, 5)},
			args:   args{v: V(5, 6)},
			want:   V(5, 5),
		},
		{
			name:   "Point far from line",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{v: V(80, -70)},
			want:   V(5, 5),
		},
		{
			name:   "Point on inline with line",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{v: V(20, 20)},
			want:   V(10, 10),
		},
		{
			name:   "Vertical line",
			fields: fields{A: V(0, -10), B: V(0, 10)},
			args:   args{v: V(-1, 0)},
			want:   V(0, 0),
		},
		{
			name:   "Horizontal line",
			fields: fields{A: V(-10, 0), B: V(10, 0)},
			args:   args{v: V(0, -1)},
			want:   V(0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.Closest(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.Closest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Contains(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		v Vec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Point on line",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{v: V(5, 5)},
			want:   true,
		},
		{
			name:   "Point on negative sloped line",
			fields: fields{A: V(0, 10), B: V(10, 0)},
			args:   args{v: V(5, 5)},
			want:   true,
		},
		{
			name:   "Point not on line",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{v: V(0, 10)},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.Contains(tt.args.v); got != tt.want {
				t.Errorf("Line.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Formula(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	tests := []struct {
		name   string
		fields fields
		wantM  float64
		wantB  float64
	}{
		{
			name:   "Getting formula - 45 degs",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			wantM:  1,
			wantB:  0,
		},
		{
			name:   "Getting formula - 90 degs",
			fields: fields{A: V(0, 0), B: V(0, 10)},
			wantM:  math.Inf(1),
			wantB:  math.NaN(),
		},
		{
			name:   "Getting formula - 0 degs",
			fields: fields{A: V(0, 0), B: V(10, 0)},
			wantM:  0,
			wantB:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			gotM, gotB := l.Formula()
			if gotM != tt.wantM {
				t.Errorf("Line.Formula() gotM = %v, want %v", gotM, tt.wantM)
			}
			if gotB != tt.wantB {
				if math.IsNaN(tt.wantB) && !math.IsNaN(gotB) {
					t.Errorf("Line.Formula() gotB = %v, want %v", gotB, tt.wantB)
				}
			}
		})
	}
}

func TestLine_Intersect(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		k Line
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Vec
		want1  bool
	}{
		{
			name:   "Lines intersect",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{k: L(V(0, 10), V(10, 0))},
			want:   V(5, 5),
			want1:  true,
		},
		{
			name:   "Lines intersect 2",
			fields: fields{A: V(5, 1), B: V(1, 1)},
			args:   args{k: L(V(2, 0), V(2, 3))},
			want:   V(2, 1),
			want1:  true,
		},
		{
			name:   "Line intersect with vertical",
			fields: fields{A: V(5, 0), B: V(5, 10)},
			args:   args{k: L(V(0, 0), V(10, 10))},
			want:   V(5, 5),
			want1:  true,
		},
		{
			name:   "Line intersect with horizontal",
			fields: fields{A: V(0, 5), B: V(10, 5)},
			args:   args{k: L(V(0, 0), V(10, 10))},
			want:   V(5, 5),
			want1:  true,
		},
		{
			name:   "Lines don't intersect",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{k: L(V(0, 10), V(1, 20))},
			want:   ZV,
			want1:  false,
		},
		{
			name:   "Lines don't intersect 2",
			fields: fields{A: V(1, 1), B: V(1, 5)},
			args:   args{k: L(V(-5, 0), V(-2, 2))},
			want:   ZV,
			want1:  false,
		},
		{
			name:   "Lines don't intersect 3",
			fields: fields{A: V(2, 0), B: V(2, 3)},
			args:   args{k: L(V(1, 5), V(5, 5))},
			want:   ZV,
			want1:  false,
		},
		{
			name:   "Lines parallel",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{k: L(V(0, 1), V(10, 11))},
			want:   ZV,
			want1:  false,
		}, {
			name:   "Lines intersect",
			fields: fields{A: V(600, 600), B: V(925, 150)},
			args:   args{k: L(V(740, 255), V(925, 255))},
			want:   V(849.1666666666666, 255),
			want1:  true,
		},
		{
			name:   "Lines intersect",
			fields: fields{A: V(600, 600), B: V(925, 150)},
			args:   args{k: L(V(740, 255), V(925, 255.0001))},
			want:   V(849.1666240490657, 255.000059008986),
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			got, got1 := l.Intersect(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.Intersect() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Line.Intersect() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLine_IntersectCircle(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		c Circle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Vec
	}{
		{
			name:   "Cirle intersects",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{c: C(V(6, 4), 2)},
			want:   V(0.5857864376269049, -0.5857864376269049),
		},
		{
			name:   "Cirle doesn't intersects",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{c: C(V(0, 5), 1)},
			want:   ZV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.IntersectCircle(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.IntersectCircle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Len(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "End right-up of start",
			fields: fields{A: V(0, 0), B: V(3, 4)},
			want:   5,
		},
		{
			name:   "End left-up of start",
			fields: fields{A: V(0, 0), B: V(-3, 4)},
			want:   5,
		},
		{
			name:   "End right-down of start",
			fields: fields{A: V(0, 0), B: V(3, -4)},
			want:   5,
		},
		{
			name:   "End left-down of start",
			fields: fields{A: V(0, 0), B: V(-3, -4)},
			want:   5,
		},
		{
			name:   "End same as start",
			fields: fields{A: V(0, 0), B: V(0, 0)},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.Len(); got != tt.want {
				t.Errorf("Line.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Rotated(t *testing.T) {
	// round returns the nearest integer, rounding ties away from zero.
	// This is required because `math.Round` wasn't introduced until Go1.10
	round := func(x float64) float64 {
		t := math.Trunc(x)
		if math.Abs(x-t) >= 0.5 {
			return t + math.Copysign(1, x)
		}
		return t
	}
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		around Vec
		angle  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Line
	}{
		{
			name:   "Rotating around line center",
			fields: fields{A: V(1, 1), B: V(3, 3)},
			args:   args{around: V(2, 2), angle: math.Pi},
			want:   L(V(3, 3), V(1, 1)),
		},
		{
			name:   "Rotating around x-y origin",
			fields: fields{A: V(1, 1), B: V(3, 3)},
			args:   args{around: V(0, 0), angle: math.Pi},
			want:   L(V(-1, -1), V(-3, -3)),
		},
		{
			name:   "Rotating around line end",
			fields: fields{A: V(1, 1), B: V(3, 3)},
			args:   args{around: V(1, 1), angle: math.Pi},
			want:   L(V(1, 1), V(-1, -1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			// Have to round the results, due to floating-point in accuracies.  Results are correct to approximately
			// 10 decimal places.
			got := l.Rotated(tt.args.around, tt.args.angle)
			if round(got.A.X) != tt.want.A.X ||
				round(got.B.X) != tt.want.B.X ||
				round(got.A.Y) != tt.want.A.Y ||
				round(got.B.Y) != tt.want.B.Y {
				t.Errorf("Line.Rotated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Scaled(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		scale float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Line
	}{
		{
			name:   "Scaling by 1",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{scale: 1},
			want:   L(V(0, 0), V(10, 10)),
		},
		{
			name:   "Scaling by >1",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{scale: 2},
			want:   L(V(-5, -5), V(15, 15)),
		},
		{
			name:   "Scaling by <1",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{scale: 0.5},
			want:   L(V(2.5, 2.5), V(7.5, 7.5)),
		},
		{
			name:   "Scaling by -1",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{scale: -1},
			want:   L(V(10, 10), V(0, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.Scaled(tt.args.scale); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.Scaled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_ScaledXY(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	type args struct {
		around Vec
		scale  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Line
	}{
		{
			name:   "Scaling by 1 around origin",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{around: ZV, scale: 1},
			want:   L(V(0, 0), V(10, 10)),
		},
		{
			name:   "Scaling by >1 around origin",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{around: ZV, scale: 2},
			want:   L(V(0, 0), V(20, 20)),
		},
		{
			name:   "Scaling by <1 around origin",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{around: ZV, scale: 0.5},
			want:   L(V(0, 0), V(5, 5)),
		},
		{
			name:   "Scaling by -1 around origin",
			fields: fields{A: V(0, 0), B: V(10, 10)},
			args:   args{around: ZV, scale: -1},
			want:   L(V(0, 0), V(-10, -10)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.ScaledXY(tt.args.around, tt.args.scale); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line.ScaledXY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_String(t *testing.T) {
	type fields struct {
		A Vec
		B Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Getting string",
			fields: fields{A: V(0, 0), B: V(1, 1)},
			want:   "Line(Vec(0, 0), Vec(1, 1))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				A: tt.fields.A,
				B: tt.fields.B,
			}
			if got := l.String(); got != tt.want {
				t.Errorf("Line.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
