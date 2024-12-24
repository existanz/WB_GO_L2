package sort

import (
	"reflect"
	"testing"
)

func TestCompareDefault(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"a", "b", -1},
		{"b", "a", 1},
		{"a", "a", 0},
		{"", "a", -1},
		{"a", "", 1},
		{"akula", "akella", 1},
	}

	for _, test := range tests {
		if res := compareDefault(test.a, test.b); res != test.want {
			t.Errorf("CompareDefault(%q, %q) = %d, want %d", test.a, test.b, res, test.want)
		}
	}
}

func TestCompareNumeric(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"a", "b", -1},
		{"b", "12", 1},
		{"13", "b", -1},
		{"011", "2", 1},
		{"", "a", -1},
		{"a", "", 1},
		{"akula", "akella", 1},
		{"1", "2", -1},
		{"2", "1", 1},
		{"2", "2", 0},
	}

	for _, test := range tests {
		if res := compareNumeric(test.a, test.b); res != test.want {
			t.Errorf("CompareNumeric(%q, %q) = %d, want %d", test.a, test.b, res, test.want)
		}
	}
}

func TestCompareMonth(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"a", "b", -1},
		{"b", "january", -1},
		{"january", "mar", -1},
		{"january", "january", 0},
		{"", "january", -1},
		{"oct", "dec", -1},
		{"december", "october", 1},
		{"december", "december", 0},
		{"jun", "june", 0},
		{"jul", "jun", 1},
	}

	for _, test := range tests {
		if res := compareMonth(test.a, test.b); res != test.want {
			t.Errorf("CompareMonth(%q, %q) = %d, want %d", test.a, test.b, res, test.want)
		}
	}
}

func TestCompareNumWithSuffix(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"10", "10", 0},
		{"10", "10k", -1},
		{"10M", "10K", 1},
		{"0110", "2k", -1},
		{"", "a", -1},
		{"a", "", 1},
		{"10u", "10g", 1},
		{"1", "2", -1},
		{"2", "1", 1},
		{"2", "2", 0},
	}

	for _, test := range tests {
		if res := compareNumWithSuffix(test.a, test.b); res != test.want {
			t.Errorf("CompareNumWithSuffix(%q, %q) = %d, want %d", test.a, test.b, res, test.want)
		}
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h", "a", "b"}, []string{"a", "b", "c", "d", "e", "f", "g", "h"}},
		{[]string{"a", "a", "a", "a", "a", "a", "a", "a", "b", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c", "d", "e", "e", "e", "e", "z", "z"}, []string{"a", "b", "c", "d", "e", "z"}},
		{[]string{}, []string{}},
	}

	for _, test := range tests {
		result := Unique(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Unique(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestTrimSpaces(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a   ", "b    ", "c   ", "d   ", "e", "f", "g", "h", "a", "b"}, []string{"a", "b", "c", "d", "e", "f", "g", "h", "a", "b"}},
		{[]string{"a", "a", "a", "a", "a", "a", "a ", "a", "b  ", "b"}, []string{"a", "a", "a", "a", "a", "a", "a", "a", "b", "b"}},
		{[]string{" ", " ", " ", " ", " ", " ", "e", "e", "z", "z"}, []string{"e", "e", "z", "z"}},
		{[]string{""}, []string{}},
	}

	for _, test := range tests {
		result := TrimSpaces(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("TrimSpaces(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestComparatorCompare(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"a b a", "b b c", -1},
		{"b 7 b", "1 2 12", 1},
		{"7 r 13", "q w b", -1},
		{"011", "2", -1},
		{"", "a", -1},
		{"a", "", 1},
		{"akula", "akella", 1},
		{"1", "2 r 12", 1},
		{"2", "1", 1},
	}

	for _, test := range tests {
		c := &comparator{
			collumnID: 2,
			compare:   compareDefault,
		}
		if res := c.Compare(test.a, test.b); res != test.want {
			t.Errorf("Comparator Compare(%q, %q) = %d, want %d", test.a, test.b, res, test.want)
		}
	}
}
