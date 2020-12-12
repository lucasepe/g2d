package geometry

import (
	"math"
	"reflect"
	"testing"
)

// closeEnough will shift the decimal point by the accuracy required, truncates the results and compares them.
// Effectively this compares two floats to a given decimal point.
//  Example:
//  closeEnough(100.125342432, 100.125, 2) == true
//  closeEnough(math.Pi, 3.14, 2) == true
//  closeEnough(0.1234, 0.1245, 3) == false
func closeEnough(got, expected float64, decimalAccuracy int) bool {
	gotShifted := got * math.Pow10(decimalAccuracy)
	expectedShifted := expected * math.Pow10(decimalAccuracy)

	return math.Trunc(gotShifted) == math.Trunc(expectedShifted)
}

func TestC(t *testing.T) {
	type args struct {
		radius float64
		center Vec
	}
	tests := []struct {
		name string
		args args
		want Circle
	}{
		{
			name: "C(): positive radius",
			args: args{radius: 10, center: ZV},
			want: Circle{Radius: 10, Center: ZV},
		},
		{
			name: "C(): zero radius",
			args: args{radius: 0, center: ZV},
			want: Circle{Radius: 0, Center: ZV},
		},
		{
			name: "C(): negative radius",
			args: args{radius: -5, center: ZV},
			want: Circle{Radius: -5, Center: ZV},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := C(tt.args.center, tt.args.radius); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_String(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Circle.String(): positive radius",
			fields: fields{radius: 10, center: ZV},
			want:   "Circle(Vec(0, 0), 10.00)",
		},
		{
			name:   "Circle.String(): zero radius",
			fields: fields{radius: 0, center: ZV},
			want:   "Circle(Vec(0, 0), 0.00)",
		},
		{
			name:   "Circle.String(): negative radius",
			fields: fields{radius: -5, center: ZV},
			want:   "Circle(Vec(0, 0), -5.00)",
		},
		{
			name:   "Circle.String(): irrational radius",
			fields: fields{radius: math.Pi, center: ZV},
			want:   "Circle(Vec(0, 0), 3.14)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.String(); got != tt.want {
				t.Errorf("Circle.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Norm(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   Circle
	}{
		{
			name:   "Circle.Norm(): positive radius",
			fields: fields{radius: 10, center: ZV},
			want:   C(ZV, 10),
		},
		{
			name:   "Circle.Norm(): zero radius",
			fields: fields{radius: 0, center: ZV},
			want:   C(ZV, 0),
		},
		{
			name:   "Circle.Norm(): negative radius",
			fields: fields{radius: -5, center: ZV},
			want:   C(ZV, 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.Norm(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Norm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Area(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "Circle.Area(): positive radius",
			fields: fields{radius: 10, center: ZV},
			want:   100 * math.Pi,
		},
		{
			name:   "Circle.Area(): zero radius",
			fields: fields{radius: 0, center: ZV},
			want:   0,
		},
		{
			name:   "Circle.Area(): negative radius",
			fields: fields{radius: -5, center: ZV},
			want:   25 * math.Pi,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.Area(); got != tt.want {
				t.Errorf("Circle.Area() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Moved(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	type args struct {
		delta Vec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Circle
	}{
		{
			name:   "Circle.Moved(): positive movement",
			fields: fields{radius: 10, center: ZV},
			args:   args{delta: V(10, 20)},
			want:   C(V(10, 20), 10),
		},
		{
			name:   "Circle.Moved(): zero movement",
			fields: fields{radius: 10, center: ZV},
			args:   args{delta: ZV},
			want:   C(V(0, 0), 10),
		},
		{
			name:   "Circle.Moved(): negative movement",
			fields: fields{radius: 10, center: ZV},
			args:   args{delta: V(-5, -10)},
			want:   C(V(-5, -10), 10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.Moved(tt.args.delta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Moved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Resized(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	type args struct {
		radiusDelta float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Circle
	}{
		{
			name:   "Circle.Resized(): positive delta",
			fields: fields{radius: 10, center: ZV},
			args:   args{radiusDelta: 5},
			want:   C(V(0, 0), 15),
		},
		{
			name:   "Circle.Resized(): zero delta",
			fields: fields{radius: 10, center: ZV},
			args:   args{radiusDelta: 0},
			want:   C(V(0, 0), 10),
		},
		{
			name:   "Circle.Resized(): negative delta",
			fields: fields{radius: 10, center: ZV},
			args:   args{radiusDelta: -5},
			want:   C(V(0, 0), 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.Resized(tt.args.radiusDelta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Resized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Contains(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	type args struct {
		u Vec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Circle.Contains(): point on cicles' center",
			fields: fields{radius: 10, center: ZV},
			args:   args{u: ZV},
			want:   true,
		},
		{
			name:   "Circle.Contains(): point offcenter",
			fields: fields{radius: 10, center: V(5, 0)},
			args:   args{u: ZV},
			want:   true,
		},
		{
			name:   "Circle.Contains(): point on circumference",
			fields: fields{radius: 10, center: V(10, 0)},
			args:   args{u: ZV},
			want:   true,
		},
		{
			name:   "Circle.Contains(): point outside circle",
			fields: fields{radius: 10, center: V(15, 0)},
			args:   args{u: ZV},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.Contains(tt.args.u); got != tt.want {
				t.Errorf("Circle.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Union(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	type args struct {
		d Circle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Circle
	}{
		{
			name:   "Circle.Union(): overlapping circles",
			fields: fields{radius: 5, center: ZV},
			args:   args{d: C(ZV, 5)},
			want:   C(ZV, 5),
		},
		{
			name:   "Circle.Union(): separate circles",
			fields: fields{radius: 1, center: ZV},
			args:   args{d: C(V(0, 2), 1)},
			want:   C(V(0, 1), 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(tt.fields.center, tt.fields.radius)
			if got := c.Union(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Intersect(t *testing.T) {
	type fields struct {
		radius float64
		center Vec
	}
	type args struct {
		d Circle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Circle
	}{
		{
			name:   "Circle.Intersect(): intersecting circles",
			fields: fields{radius: 1, center: ZV},
			args:   args{d: C(V(1, 0), 1)},
			want:   C(V(0.5, 0), 1),
		},
		{
			name:   "Circle.Intersect(): non-intersecting circles",
			fields: fields{radius: 1, center: ZV},
			args:   args{d: C(V(3, 3), 1)},
			want:   C(V(1.5, 1.5), 0),
		},
		{
			name:   "Circle.Intersect(): first circle encompassing second",
			fields: fields{radius: 10, center: ZV},
			args:   args{d: C(V(3, 3), 1)},
			want:   C(ZV, 10),
		},
		{
			name:   "Circle.Intersect(): second circle encompassing first",
			fields: fields{radius: 1, center: V(-1, -4)},
			args:   args{d: C(ZV, 10)},
			want:   C(ZV, 10),
		},
		{
			name:   "Circle.Intersect(): matching circles",
			fields: fields{radius: 1, center: ZV},
			args:   args{d: C(ZV, 1)},
			want:   C(ZV, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C(
				tt.fields.center,
				tt.fields.radius,
			)
			if got := c.Intersect(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_IntersectPoints(t *testing.T) {
	type fields struct {
		Center Vec
		Radius float64
	}
	type args struct {
		l Line
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Vec
	}{
		{
			name:   "Line intersects circle at two points",
			fields: fields{Center: V(2, 2), Radius: 1},
			args:   args{L(V(0, 0), V(10, 10))},
			want:   []Vec{V(1.292, 1.292), V(2.707, 2.707)},
		},
		{
			name:   "Line intersects circle at one point",
			fields: fields{Center: V(-0.5, -0.5), Radius: 1},
			args:   args{L(V(0, 0), V(10, 10))},
			want:   []Vec{V(0.207, 0.207)},
		},
		{
			name:   "Line endpoint is circle center",
			fields: fields{Center: V(0, 0), Radius: 1},
			args:   args{L(V(0, 0), V(10, 10))},
			want:   []Vec{V(0.707, 0.707)},
		},
		{
			name:   "Both line endpoints within circle",
			fields: fields{Center: V(0, 0), Radius: 1},
			args:   args{L(V(0.2, 0.2), V(0.5, 0.5))},
			want:   []Vec{},
		},
		{
			name:   "Line does not intersect circle",
			fields: fields{Center: V(10, 0), Radius: 1},
			args:   args{L(V(0, 0), V(10, 10))},
			want:   []Vec{},
		},
		{
			name:   "Horizontal line intersects circle at two points",
			fields: fields{Center: V(5, 5), Radius: 1},
			args:   args{L(V(0, 5), V(10, 5))},
			want:   []Vec{V(4, 5), V(6, 5)},
		},
		{
			name:   "Vertical line intersects circle at two points",
			fields: fields{Center: V(5, 5), Radius: 1},
			args:   args{L(V(5, 0), V(5, 10))},
			want:   []Vec{V(5, 4), V(5, 6)},
		},
		{
			name:   "Left and down line intersects circle at two points",
			fields: fields{Center: V(5, 5), Radius: 1},
			args:   args{L(V(10, 10), V(0, 0))},
			want:   []Vec{V(5.707, 5.707), V(4.292, 4.292)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Circle{
				Center: tt.fields.Center,
				Radius: tt.fields.Radius,
			}
			got := c.IntersectionPoints(tt.args.l)
			for i, v := range got {
				if !closeEnough(v.X, tt.want[i].X, 2) || !closeEnough(v.Y, tt.want[i].Y, 2) {
					t.Errorf("Circle.IntersectPoints() = %v, want %v", v, tt.want[i])
				}
			}
		})
	}
}
