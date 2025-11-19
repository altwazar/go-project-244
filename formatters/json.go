package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

// Рекурсивный разбор дифов
func formatOutputJson(diffs []parsers.Diff, level int) string {
	var out string

	spacing_braces := strings.Repeat(" ", (level * 4))
	spacing := strings.Repeat(" ", (level+1)*4-2)
	out += "[\n"

	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})
	var newValue string
	var oldValue string
	for _, diff := range diffs {
		// Преобразование значения в строку
		key := fmt.Sprintf("\"%v\"", diff.Key)
		if diff.New == nil {
			newValue = "null"
		} else if isString(diff.New) {
			newValue = fmt.Sprintf("\"%v\"", diff.New)
		} else {
			newValue = fmt.Sprintf("%v", diff.New)
		}
		if diff.Old == nil {
			oldValue = "null"
		} else if isString(diff.Old) {
			oldValue = fmt.Sprintf("\"%v\"", diff.Old)
		} else {
			oldValue = fmt.Sprintf("%v", diff.Old)
		}
		out += spacing + "{\n"
		out += fmt.Sprintf("%s\"state\": %v,\n", spacing, diff.State)
		out += fmt.Sprintf("%s\"key\": %v,\n", spacing, key)
		switch diff.State {
		case parsers.Equal:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s\"Value\": %v,\n", spacing, oldValue)
			} else {
				out += fmt.Sprintf("%s\"Value\": %s,\n", spacing, formatOutputJson(diff.DiffChild, level+1))
			}
		case parsers.Updated:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s\"oldValue\": %v,\n", spacing, oldValue)
				out += fmt.Sprintf("%s\"newValue\": %v,\n", spacing, newValue)
			} else if diff.FromTo == parsers.MapToValue {
				out += fmt.Sprintf("%s\"oldValue\": %s,\n", spacing, formatOutputJson(diff.DiffChild, level+1))
				out += fmt.Sprintf("%s\"newValue\": %v,\n", spacing, newValue)
			} else if diff.FromTo == parsers.ValueToMap {
				out += fmt.Sprintf("%s\"oldValue\": %v,\n", spacing, oldValue)
				out += fmt.Sprintf("%s\"newValue\": %s,\n", spacing, formatOutputJson(diff.DiffChild, level+1))
			} else {
				out += fmt.Sprintf("%s\"diff\": %s,\n", spacing, formatOutputJson(diff.DiffChild, level+1))
			}
		case parsers.Added:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s\"newValue\": %v,\n", spacing, newValue)
			} else {
				out += fmt.Sprintf("%s\"newValue\": %s,\n", spacing, formatOutputJson(diff.DiffChild, level+1))
			}
		case parsers.Removed:
			if diff.FromTo == parsers.ValueToValue {
				out += fmt.Sprintf("%s\"oldValue\": %v,\n", spacing, oldValue)
			} else {
				out += fmt.Sprintf("%s\"oldValue\": %s,\n", spacing, formatOutputJson(diff.DiffChild, level+1))
			}
		}
		out = strings.TrimSuffix(out, ",\n") + "\n"
		out += spacing + "},\n"
	}

	out = strings.TrimSuffix(out, ",\n") + "\n"
	if level == 0 {
		out += spacing_braces + "]\n"
	} else {
		out += spacing_braces + "]"
	}
	return out
}
