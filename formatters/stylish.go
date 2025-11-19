package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

// Рекурсивный разбор дифов
func formatOutputStylish(diffs []parsers.Diff, level int) string {
	var out string

	spacing_braces := strings.Repeat(" ", (level * 4))
	spacing := strings.Repeat(" ", (level+1)*4-2)
	out += "{\n"
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})
	var newKey string
	var oldKey string
	for _, diff := range diffs {
		// Преобразование значения в строку
		if diff.New == nil {
			newKey = "null"
		} else {
			newKey = fmt.Sprintf("%v", diff.New)
		}
		if diff.Old == nil {
			oldKey = "null"
		} else {
			oldKey = fmt.Sprintf("%v", diff.Old)
		}
		switch diff.State {
		case parsers.Equal:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s  %s: %s\n", spacing, diff.Key, newKey)
			} else {
				out += fmt.Sprintf("%s  %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		case parsers.Updated:
			switch diff.FromTo {
			case parsers.ValueToValue:
				out += fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, oldKey)
				out += fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, newKey)
			case parsers.MapToValue:
				out += fmt.Sprintf("%s- %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
				out += fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, newKey)
			case parsers.ValueToMap:
				out += fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, oldKey)
				out += fmt.Sprintf("%s+ %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			default:
				out += fmt.Sprintf("%s  %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		case parsers.Added:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, newKey)
			} else {
				out += fmt.Sprintf("%s+ %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		case parsers.Removed:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, oldKey)
			} else {
				out += fmt.Sprintf("%s- %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		}
	}
	out += spacing_braces + "}\n"
	return out
}
