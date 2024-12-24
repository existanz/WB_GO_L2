package cut

import (
	"dev06/internal/config"
	"strings"
)

func Cut(s string, cfg config.Config) (string, error) {
	if cfg.Sep && strings.Contains(s, cfg.Delim) || len(cfg.Cols) == 0 {
		return "", nil
	}
	strs := strings.Split(s, cfg.Delim)
	result := make([]string, 0, len(cfg.Cols))
	for _, v := range cfg.Cols {
		if v < 0 || v >= len(strs) {
			return "", nil
		}
		result = append(result, strs[v])
	}
	return strings.Join(result, cfg.Delim), nil
}
