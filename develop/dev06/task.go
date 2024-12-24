package main

import (
	"bufio"
	"dev06/internal/config"
	"dev06/internal/cut"
	"fmt"
	"io"
	"os"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Write your rows, when finish write EOF (control+D)")
	Run(cfg)
}

func Run(cfg config.Config) {
	reader := bufio.NewReader(os.Stdin)
	out := make([]string, 0)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		cuted, err := cut.Cut(line, cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cuted == "" {
			continue
		}
		out = append(out, cuted)
	}
	for _, v := range out {
		fmt.Fprintln(os.Stdout, v)
	}
}
