package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Cols  []int
	Delim string
	Sep   bool
}

func LoadConfig() (Config, error) {
	cfg := Config{}

	f := flag.String("f", "0", "fields - выбрать поля (колонки), через запятую")
	d := flag.String("d", "\t", "delimiter - использовать другой разделитель")
	s := flag.Bool("s", false, "separated - только строки с разделителем")

	flag.Parse()
	cols, err := parseStringToInts(*f)

	if err != nil {
		return cfg, fmt.Errorf("invalid fields: %w", err)
	}
	cfg.Cols = cols

	cfg.Delim = *d
	cfg.Sep = *s

	return cfg, nil

}

func parseStringToInts(s string) ([]int, error) {
	var result []int
	for _, v := range strings.Split(s, ",") {
		i, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return nil, fmt.Errorf("invalid integer in list: %w", err)
		}
		result = append(result, i)
	}
	return result, nil
}
