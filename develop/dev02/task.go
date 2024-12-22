package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func DecodeString(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	if unicode.IsDigit(rune(s[0])) {
		return "", fmt.Errorf("incorrect string")
	}

	var prev rune
	var escaped bool
	var b strings.Builder
	for _, r := range s {
		if escaped {
			prev = r
			escaped = !escaped
			continue
		}
		if n, err := strconv.Atoi(string(r)); err == nil {
			if prev == rune(0) {
				return "", fmt.Errorf("incorrect string")
			}
			ss := strings.Repeat(string(prev), n)
			b.WriteString(ss)
			prev = rune(0)
			continue
		}
		if prev != rune(0) {
			b.WriteRune(prev)
		}
		escaped = (r == '\\' && prev != '\\')
		prev = r
	}
	if prev != rune(0) {
		b.WriteRune(prev)
	}
	return b.String(), nil
}
