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
	// Формирование отступов
	spacing_braces := strings.Repeat(" ", (level * 4))
	spacing := strings.Repeat(" ", (level+1)*4-2)

	out += "{\n"
	// Сортировка по ключам
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})
	// Значения ключей в нужном формате
	var newValue string
	var oldValue string
	for _, diff := range diffs {
		if diff.New == nil {
			newValue = "null"
		} else {
			newValue = fmt.Sprintf("%v", diff.New)
		}
		if diff.Old == nil {
			oldValue = "null"
		} else {
			oldValue = fmt.Sprintf("%v", diff.Old)
		}
		// Формирование тела вывода, для вложенных структур рекурсивно вызывается эта же функция
		switch diff.State {
		case parsers.Equal:
			if !diff.OldIsMap && !diff.NewIsMap {
				out += fmt.Sprintf("%s  %s: %s\n", spacing, diff.Key, newValue)
			} else {
				out += fmt.Sprintf("%s  %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		case parsers.Updated:
			switch {
			case !diff.OldIsMap && !diff.NewIsMap:
				out += fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, oldValue)
				out += fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, newValue)
			case diff.OldIsMap && !diff.NewIsMap:
				out += fmt.Sprintf("%s- %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
				out += fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, newValue)
			case !diff.OldIsMap && diff.NewIsMap:
				out += fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, oldValue)
				out += fmt.Sprintf("%s+ %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			default:
				out += fmt.Sprintf("%s  %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		case parsers.Added:
			if !diff.OldIsMap && !diff.NewIsMap {
				out += fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, newValue)
			} else {
				out += fmt.Sprintf("%s+ %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		case parsers.Removed:
			if !diff.OldIsMap && !diff.NewIsMap {
				out += fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, oldValue)
			} else {
				out += fmt.Sprintf("%s- %s: %s", spacing, diff.Key, formatOutputStylish(diff.DiffChild, level+1))
			}
		}
	}
	out += spacing_braces + "}\n"
	return out
}
