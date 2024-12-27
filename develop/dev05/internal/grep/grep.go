package grep

import (
	"dev05/internal/config"
	"fmt"
	"regexp"
)

func Grep(cfg config.Config, lines []string) []string {
	pattern := updatePattern(cfg.Pattern, cfg.CI, cfg.Fixed)

	ids := grepIds(lines, pattern, cfg.Invert)

	result := make([]string, 0, len(ids))

	for i := range ids {
		prev, cur, next := -1, ids[i], len(lines)
		if i > 0 {
			prev = ids[i-1] + cfg.After
		}
		if i < len(ids)-1 {
			next = ids[i+1]
		}

		low, high := max(min(prev+1, cur), cur-cfg.Before), min(next-1, cur+cfg.After)

		for i := low; i <= high; i++ {
			line := lines[i]
			if cfg.LineNum {
				line = fmt.Sprintf("%d: %s", i+1, line)
			}
			if cfg.Highlight && cur == i {
				line = fmt.Sprintf("\x1b[43m%s\x1b[0m", line)
			}
			result = append(result, line)
		}
	}
	return result
}

func CountGrep(cfg config.Config, lines []string) int {
	pattern := updatePattern(cfg.Pattern, cfg.CI, cfg.Fixed)
	return len(grepIds(lines, pattern, cfg.Invert))
}

func grepIds(lines []string, pattern string, invert bool) []int {
	result := make([]int, 0, len(lines))
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return result
	}
	for i, line := range lines {
		if invert != reg.MatchString(line) {
			result = append(result, i)
		}
	}
	return result
}

func updatePattern(pattern string, ci, fixed bool) string {
	if ci {
		pattern = "(?i)" + pattern
	}
	if fixed {
		pattern = "^" + pattern + "$"
	}
	return pattern
}
