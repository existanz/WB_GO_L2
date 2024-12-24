package sort

import (
	"dev03/internal/config"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func compareDefault(a, b string) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func compareNumeric(a, b string) int {
	aInt, errA := strconv.Atoi(a)
	bInt, errB := strconv.Atoi(b)
	if errA != nil && errB != nil {
		return compareDefault(a, b)
	}
	if errA != nil {
		return 1
	}
	if errB != nil {
		return -1
	}
	if aInt > bInt {
		return 1
	}
	if aInt < bInt {
		return -1
	}
	return 0
}

var months = map[string]time.Month{
	"jan":       time.January,
	"january":   time.January,
	"feb":       time.February,
	"february":  time.February,
	"mar":       time.March,
	"march":     time.March,
	"apr":       time.April,
	"april":     time.April,
	"may":       time.May,
	"jun":       time.June,
	"june":      time.June,
	"jul":       time.July,
	"july":      time.July,
	"aug":       time.August,
	"august":    time.August,
	"sep":       time.September,
	"september": time.September,
	"oct":       time.October,
	"october":   time.October,
	"nov":       time.November,
	"november":  time.November,
	"dec":       time.December,
	"december":  time.December,
}

func compareMonth(a, b string) int {
	aMonth := months[strings.ToLower(a)]
	bMonth := months[strings.ToLower(b)]
	if aMonth == 0 && bMonth == 0 {
		return compareDefault(a, b)
	}
	if aMonth > bMonth {
		return 1
	}
	if aMonth < bMonth {
		return -1
	}
	return 0
}

func compareNumWithSuffix(a, b string) int {
	a = strings.ToUpper(a)
	b = strings.ToUpper(b)
	numA, suffixA, errA := parseNumWithSuffix(a)
	numB, suffixB, errB := parseNumWithSuffix(b)

	if errA != nil && errB != nil {
		return compareDefault(a, b)
	}
	if errA != nil {
		return 1
	}
	if errB != nil {
		return -1
	}

	numA = applySuffix(numA, suffixA)
	numB = applySuffix(numB, suffixB)

	if numA < numB {
		return -1
	} else if numA > numB {
		return 1
	} else {
		return 0
	}
}

func parseNumWithSuffix(s string) (float64, string, error) {
	re := regexp.MustCompile(`^(\d+(\.\d+)?)([KMGT]?)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 4 {
		return 0, "", fmt.Errorf("invalid format: %s", s)
	}
	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid number: %s", matches[1])
	}
	return num, matches[3], nil
}

func applySuffix(num float64, suffix string) float64 {
	switch suffix {
	case "K":
		return num * 1024
	case "M":
		return num * 1024 * 1024
	case "G":
		return num * 1024 * 1024 * 1024
	default:
		return num
	}
}

func Unique(values []string) []string {
	unique := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := unique[value]; !ok {
			unique[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

func TrimSpaces(values []string) []string {
	result := make([]string, 0, len(values))
	for _, val := range values {
		val = strings.TrimSpace(val)
		if len(val) == 0 {
			continue
		}
		result = append(result, val)
	}
	return result
}

type comparator struct {
	collumnID int
	compare   func(a, b string) int
}

func (c *comparator) Compare(a, b string) int {
	aWords := strings.Fields(a)
	bWords := strings.Fields(b)
	if len(aWords) <= c.collumnID && len(bWords) <= c.collumnID {
		return c.compare(a, b)
	}
	if len(aWords) <= c.collumnID {
		return 1
	}
	if len(bWords) <= c.collumnID {
		return -1
	}
	aWord := aWords[c.collumnID]
	bWord := bWords[c.collumnID]
	return c.compare(aWord, bWord)
}

func Sort(lines []string, cfg config.Config) {
	c := comparator{
		collumnID: cfg.ColumnID,
		compare:   compareDefault,
	}
	switch cfg.SortBy {

	case config.NumSort:
		c.compare = compareNumeric
	case config.MonthSort:
		c.compare = compareMonth
	case config.HumanSort:
		c.compare = compareNumWithSuffix
	default:
		c.compare = compareDefault
	}
	slices.SortFunc(lines, c.Compare)

	if cfg.Desc {
		slices.Reverse(lines)
	}
}

func IsSorted(lines []string, cfg config.Config) bool {
	c := comparator{
		collumnID: cfg.ColumnID,
		compare:   compareDefault,
	}
	switch cfg.SortBy {
	case config.NumSort:
		c.compare = compareNumeric
	case config.MonthSort:
		c.compare = compareMonth
	case config.HumanSort:
		c.compare = compareNumWithSuffix
	default:
		c.compare = compareDefault
	}
	sign := 1
	if cfg.Desc {
		sign = -1
	}
	return slices.IsSortedFunc(lines, func(a, b string) int {
		return sign * c.Compare(a, b)
	})
}
