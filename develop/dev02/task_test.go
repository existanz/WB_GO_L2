package main

import "testing"

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"алгосы3\\4\\5", "алгосыыы45", false},
		{"qwe\\45", "qwe44444", false},
		{"qwe0a5", "qwaaaaa", false},
		{"q🙃e\\\\5", "q🙃e\\\\\\\\\\", false},
		{"qwe45", "", true},
		{"a\\4b", "a4b", false},
		{"a\\bc", "abc", false},
		{"a\\", "a\\", false},
		{"\\3", "3", false},
		{"\\", "\\", false},
		{"abc\\d\\2", "abcd2", false},
		{"\\1", "1", false},
	}

	for _, test := range tests {
		result, err := DecodeString(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("Decode(%q) expected error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("Decode(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
