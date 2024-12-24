package main

import (
	"dev03/internal/config"
	"dev03/internal/files"
	"dev03/internal/sort"
	"fmt"
	"os"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

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

	if cfg.OnlyCheck {
		if sort.IsSorted(lines, cfg) {
			fmt.Printf("File %s is sorted\n", cfg.Filename)
		} else {
			fmt.Printf("File %s is not sorted\n", cfg.Filename)
		}
	}

	if !sort.IsSorted(lines, cfg) {
		sort.Sort(lines, cfg)
		err = files.WriteLines(cfg.Filename+"_sorted", lines)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
