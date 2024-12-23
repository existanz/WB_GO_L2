package main

import (
	"encoding/json"
	"testing"
)

func TestAnagrams(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{
			[]string{"листОк", "тяпка", "слиток", "Пятка", "пятак", "листок", "", "пятак", "столик"},
			`{"листок":["листок","слиток","столик"],"тяпка":["пятак","пятка","тяпка"]}`,
		},
		{
			[]string{"РВОТА", "отвар", "втора", "отрок", "автор", "товар", "рвота", "рокОт", "Автор", "ЗАтор", "каБан", "банка"},
			`{"кабан":["банка","кабан"],"отрок":["отрок","рокот"],"рвота":["автор","втора","отвар","рвота","товар"]}`,
		},
		{
			[]string{},
			`{}`,
		},
		{
			[]string{"даздраперма"},
			`{}`,
		},
	}

	for _, test := range tests {
		result := Anagrams(test.input)
		json, err := json.Marshal(result)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if string(json) != test.expected {
			t.Errorf("Anagrams(%q) = %q; want %q", test.input, string(json), test.expected)
		}
	}
}
