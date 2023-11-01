package main

import "testing"

func Test_extract(t *testing.T) {
	var table = []struct {
		input       string
		expectedOut string
		err         bool
	}{
		{
			input:       `a\`,
			expectedOut: ``,
			err:         true,
		},
		{
			input:       `a4bc2d5e`,
			expectedOut: `aaaabccddddde`,
			err:         false,
		},
		{
			input:       `45`,
			expectedOut: ``,
			err:         true,
		},
		{
			input:       `qwe\45`,
			expectedOut: `qwe44444`,
			err:         false,
		},
	}

	for _, test := range table {
		out, err := string_unpacking(test.input)
		if out != test.expectedOut || test.err && err == nil || !test.err && err != nil {
			t.Error("Wrong result")
		}
	}
}
