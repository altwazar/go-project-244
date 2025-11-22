package code

import (
	"code/formatters"
	"code/parsers"
)

// Сравнение конфигов, получаю diffs из конфигов и отдаю в форматер.
func GenDiff(pathBefore string, pathAfter string, format string) (string, error) {
	var out string
	diff, err := parsers.ParseConfigs(pathBefore, pathAfter)
	if err != nil {
		return "", err
	}
	out = formatters.FormatOut(diff, format)
	return out, nil
}
