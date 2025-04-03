package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractValuesFromString(t *testing.T) {
	testdata := []struct {
		in  string
		out string
	}{
		{"abc", "abc"},
		{"sbv\n\t    ", "sbv"},
		{"zero\x00foobar", "zero\nfoobar"},
		{" zero\x00foobar \t\n", " zero\nfoobar"},
		{"zero\x00zero\x00foobar", "zero\nzero\nfoobar"},
	}
	for _, tt := range testdata {
		t.Run(tt.in, func(t *testing.T) {
			s := ExtractValueFromString(tt.in)
			require.Equal(t, tt.out, s)
		})
	}
}

func TestIsCorrectEnvVarName(t *testing.T) {
	testdata := []struct {
		in  string
		out bool
	}{
		{"abc", true},
		{"ABC123_", true},
		{"aA12_", true},
		{"_SDFD", true},
		{"999", false},
		{"9EEE", false},
		{"", false},
		{"A#", false},
		{"AeeeeКириллица", false},
	}
	for _, tt := range testdata {
		t.Run(tt.in, func(t *testing.T) {
			r := IsCorrectEnvVarName(tt.in)
			require.Equal(t, tt.out, r)
		})
	}
}
