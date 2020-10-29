package main

import "testing"

func TestParseVersion(t *testing.T) {
	cases := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "valid",
			input: []byte("12.11.18\n"),
			want:  "12.11.18",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseVersion(tc.input)
			if err != nil {
				t.Error(err)
			}
			if got != tc.want {
				t.Errorf("got: %v - want: %v", got, tc.want)
			}
		})
	}
}
