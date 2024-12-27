package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Filename  string
	Pattern   string
	After     int
	Before    int
	Count     bool
	CI        bool
	Invert    bool
	Fixed     bool
	LineNum   bool
	Highlight bool
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	A := flag.Int("A", 0, "after печатать +N строк после совпадения")
	B := flag.Int("B", 0, "before печатать +N строк до совпадения")
	C := flag.Int("C", 0, "context (A+B) печатать ±N строк вокруг совпадения")
	c := flag.Bool("c", false, "count (количество строк)")
	i := flag.Bool("i", false, "ignore-case (игнорировать регистр)")
	v := flag.Bool("v", false, "invert (вместо совпадения, исключать)")
	F := flag.Bool("F", false, "fixed, точное совпадение со строкой, не паттерн")
	n := flag.Bool("n", false, "line num, печатать номер строки")

	flag.Parse()

	cfg.After = *A
	cfg.Before = *B
	if *C > 0 {
		cfg.After = *C
		cfg.Before = *C
	}

	cfg.Count = *c
	cfg.CI = *i
	cfg.Invert = *v
	cfg.Fixed = *F
	cfg.LineNum = *n

	cfg.Highlight = true // для выделения строки совпадения

	if flag.NArg() < 2 {
		return cfg, fmt.Errorf("Usagee go-grep [-A N] [-B N] [-C N] [-c] [-i] [-v] [-F] [-n] pattern filename")
	}

	cfg.Pattern = flag.Arg(0)
	cfg.Filename = flag.Arg(1)
	return cfg, nil
}
