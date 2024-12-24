package config

import (
	"flag"
	"fmt"
)

type SortType int

const (
	DefaultSort SortType = iota
	NumSort
	MonthSort
	HumanSort
)

type Config struct {
	Filename    string
	SortBy      SortType
	ColumnID    int
	Desc        bool
	Unique      bool
	IgnoreBlank bool
	OnlyCheck   bool
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	k := flag.Int("k", 0, "указание колонки для сортировки")
	n := flag.Bool("n", false, "сортировать по числовому значению")
	r := flag.Bool("r", false, "сортировать в обратном порядке")
	u := flag.Bool("u", false, "не выводить повторяющиеся строки")
	M := flag.Bool("M", false, "сортировать по названию месяца")
	b := flag.Bool("b", false, "игнорировать хвостовые пробелы")
	c := flag.Bool("c", false, "проверять отсортированы ли данные")
	h := flag.Bool("h", false, "сортировать по числовому значению с учётом суффиксов")
	flag.Parse()

	cfg.ColumnID = *k

	if *n {
		cfg.SortBy = NumSort
	}
	if *M {
		cfg.SortBy = MonthSort
	}
	if *h {
		cfg.SortBy = HumanSort
	}

	cfg.Desc = *r
	cfg.Unique = *u
	cfg.IgnoreBlank = *b
	cfg.OnlyCheck = *c

	if flag.NArg() < 1 {
		return cfg, fmt.Errorf("Please provide a filename")
	}

	cfg.Filename = flag.Arg(0)
	return cfg, nil
}
