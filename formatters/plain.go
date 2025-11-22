package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

// Рекурсивный разбор изменений в плоском формате
func formatOutputPlain(diffs []parsers.Diff, path string) string {
	var out string
	// Сортировка по ключам
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})
	var newKey string
	var oldKey string
	var pathKey string
	for _, diff := range diffs {
		// Значения ключей в нужном формате
		if diff.New == nil {
			newKey = "null"
		} else if isString(diff.New) {
			newKey = fmt.Sprintf("'%v'", diff.New)
		} else {
			newKey = fmt.Sprintf("%v", diff.New)
		}
		if diff.Old == nil {
			oldKey = "null"
		} else if isString(diff.Old) {
			oldKey = fmt.Sprintf("'%v'", diff.Old)
		} else {
			oldKey = fmt.Sprintf("%v", diff.Old)
		}
		// Путь ключа для plain вывода
		if path != "" {
			pathKey = path + "." + diff.Key
		} else {
			pathKey = diff.Key
		}
		// Формирование тела вывода, для вложенных структур рекурсивно вызывается эта же функция
		switch diff.State {
		case parsers.Updated:
			switch {
			case !diff.OldIsMap && !diff.NewIsMap:
				out += fmt.Sprintf("Property '%s' was updated. From %s to %s\n", pathKey, oldKey, newKey)
			case diff.OldIsMap && !diff.NewIsMap:
				out += fmt.Sprintf("Property '%s' was updated. From [complex value] to %s\n", pathKey, newKey)
			case !diff.OldIsMap && diff.NewIsMap:
				out += fmt.Sprintf("Property '%s' was updated. From %s to [complex value]\n", pathKey, oldKey)
			default:
				out += formatOutputPlain(diff.DiffChild, pathKey)
			}
		case parsers.Added:
			if !diff.OldIsMap && !diff.NewIsMap {
				out += fmt.Sprintf("Property '%s' was added with value: %s\n", pathKey, newKey)
			} else {
				out += fmt.Sprintf("Property '%s' was added with value: [complex value]\n", pathKey)
			}
		case parsers.Removed:
			out += fmt.Sprintf("Property '%s' was removed\n", pathKey)
		}
	}
	if path == "" {
		out = strings.TrimPrefix(out, "\n")
	}
	return out
}
