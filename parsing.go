package code

import (
	"code/formatters"
	"code/parsers"
)

type pair struct {
	First  any
	Second any
}

func CompareConfigs(pathBefore string, pathAfter string, format string) (string, error) {
	var out string
	diffs, err := parsers.ParseConfigs(pathBefore, pathAfter)
	if err != nil {
		return "", err
	}
	out = formatters.FormatOut(diffs, format)
	return out, nil
}
