package main

import (
	"dev05/internal/config"
	"dev05/internal/files"
	"dev05/internal/grep"
	"fmt"
	"os"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines, err := files.ReadLines(cfg.Filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfg.Count {
		fmt.Println(grep.CountGrep(cfg, lines))
		return
	}

	out := grep.Grep(cfg, lines)

	for _, v := range out {
		fmt.Fprintln(os.Stdout, v)
	}

}
