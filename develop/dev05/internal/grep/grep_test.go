package grep

import (
	"dev05/internal/config"
	"reflect"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		input    []string
		pattrn   string
		ci       bool
		fixed    bool
		invert   bool
		expected []string
	}{
		{
			[]string{"a   ", "b    ", "c   ", "d   ", "A", "f", "g", "h", "a", "b"},
			"a",
			false,
			false,
			false,
			[]string{"a   ", "a"},
		},
		{
			[]string{"a   ", "b    ", "c   ", "d   ", "A", "f", "g", "h", "a", "b"},
			"a",
			true,
			true,
			false,
			[]string{"A", "a"},
		},
		{
			[]string{"a", "B tree", "a", "a", "a", "a", "a ", "a", "b  ", "b"},
			"B",
			true,
			false,
			false,
			[]string{"B tree", "b  ", "b"},
		},
		{
			[]string{" ", " ", " ", " ", " ", " ", "e", "e", "z", "z"},
			"z",
			false,
			false,
			true,
			[]string{" ", " ", " ", " ", " ", " ", "e", "e"},
		},
		{
			[]string{""},
			"_",
			false,
			false,
			false,
			[]string{},
		},
	}

	for _, test := range tests {
		cfg := config.Config{Pattern: test.pattrn, CI: test.ci, Fixed: test.fixed, Invert: test.invert}
		result := Grep(cfg, test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Grep(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}
