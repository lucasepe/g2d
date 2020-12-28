package cmd

import "testing"

func TestLastPathSegment(t *testing.T) {
	test := []struct {
		uri  string
		want string
	}{
		{"http://google.com/foo/bar?a=1&b=2", "bar"},
		{"/some/path/to/remove/file.g2d", "file"},
		{"/Documents/g2D Scripts/star.g2d", "star"},
		{"http://google.com/foo/rosette.g2d?a=1&b=2", "rosette"},
	}

	for _, tt := range test {
		t.Run(tt.uri, func(t *testing.T) {
			got, err := lastPathSegment(tt.uri)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}
