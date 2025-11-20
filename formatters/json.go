package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

// Рекурсивное формирование json вывода
// noDiff - для обработки вывода вложенной структуры без сравнения
func formatOutputJson(diffs []parsers.Diff, level int, noDiff bool) string {
	var out string
	// Строка отступ
	spacing := strings.Repeat(" ", 2)
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})
	// Шаблон в начале на нулевом уровне
	if level == 0 {
		out += "{\n"
		out += spacing + "\"key\": \"\",\n"
		out += spacing + "\"type\": \"root\",\n"
		out += spacing + "\"children\": [\n"
	}
	// Перебор диффов. Начальный отступ ри нулевом уровне
	if level != 0 {
		out += diffRange(diffs, level, spacing, noDiff)
	} else {
		out += diffRange(diffs, level+2, spacing, noDiff)
	}
	// Удаляется запятая в конце перечислений
	out = strings.TrimSuffix(out, ",\n") + "\n"
	switch {
	// Обработка вывода нулевого уровня
	case level == 0 && !noDiff:
		out += spacing + "]\n}\n"
	case level == 0 && noDiff:
		out += "}\n"
	// закрытие списка после вывода самого списка диффов
	case level != 0 && !noDiff:
		out += strings.Repeat(spacing, level-1) + "]"
	}
	return out
}

// Обработка списка диффов
func diffRange(diffs []parsers.Diff, level int, spacing string, noDiff bool) string {
	var out string
	for _, diff := range diffs {
		// Кроме нулевого уровня отступы меняются в функциях вывода
		if noDiff {
			out += formatNoDiff(diff, level, spacing)
		} else {
			out += formatDiff(diff, level, spacing)
		}
	}
	return out
}

// Формирование вывода различий
func formatDiff(diff parsers.Diff, level int, spacing string) string {
	// У ключей и скобочек разный отступ
	spacingBraces := strings.Repeat(spacing, level)
	spacingKeys := strings.Repeat(spacing, level+1)
	var out string
	var newValue string
	var oldValue string
	var state string
	// Преобразование значений в нужную строку для вывода
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
	// Преобразование статуса сравнения в нужный вывод
	switch diff.State {
	case parsers.Updated:
		if diff.OldIsMap && diff.NewIsMap {
			state = "nested"
		} else {
			state = "changed"
		}
	case parsers.Added:
		state = "added"
	case parsers.Removed:
		state = "deleted"
	case parsers.Equal:
		state = "unchanged"
	}
	// Начальная скобочка и фиксированные поля key и type
	out += spacingBraces + "{\n"
	out += fmt.Sprintf("%s\"key\": \"%s\",\n", spacingKeys, diff.Key)
	out += fmt.Sprintf("%s\"type\": \"%s\",\n", spacingKeys, state)
	// Формирование вывода в зависимости от статуса
	// Если в качестве значения вложенная структура, то { на той же строке
	// и уровень смещения увеличивается на 1
	// Если Nested с разбором диффа, то на отдельный и уровень на 2
	switch state {
	case "unchanged":
		if !diff.OldIsMap {
			out += fmt.Sprintf("%s\"value1\": %v,\n", spacingKeys, oldValue)
		} else {
			out += fmt.Sprintf("%s\"value1\": {\n%s\n%s},\n",
				spacingKeys,
				formatOutputJson(diff.DiffChild, level+1, true),
				spacingKeys)
		}
	case "changed":
		switch {
		case !diff.OldIsMap && !diff.NewIsMap:
			if diff.Old != nil {
				out += fmt.Sprintf("%s\"value1\": %v,\n", spacingKeys, oldValue)
			}
			if diff.New != nil {
				out += fmt.Sprintf("%s\"value2\": %v,\n", spacingKeys, newValue)
			}
		case diff.OldIsMap && !diff.NewIsMap:
			out += fmt.Sprintf("%s\"value1\": {\n%s%s},\n", spacingKeys,
				formatOutputJson(diff.DiffChild, level+1, true), spacingKeys)
			if diff.New != nil {
				out += fmt.Sprintf("%s\"value2\": %v,\n", spacingKeys, newValue)
			}
		case !diff.OldIsMap && diff.NewIsMap:
			if diff.Old != nil {
				out += fmt.Sprintf("%s\"value1\": %v,\n", spacingKeys, oldValue)
			}
			out += fmt.Sprintf("%s\"value2\": {\n%s%s},\n", spacingKeys,
				formatOutputJson(diff.DiffChild, level+1, true), spacingKeys)
		}
	case "nested":
		out += fmt.Sprintf("%s\"children\": [\n%s,\n", spacingKeys,
			formatOutputJson(diff.DiffChild, level+2, false))
	case "added":
		if !diff.NewIsMap {
			out += fmt.Sprintf("%s\"value2\": %v,\n", spacingKeys, newValue)
		} else {
			out += fmt.Sprintf("%s\"value2\": {\n%s%s},\n", spacingKeys,
				formatOutputJson(diff.DiffChild, level+1, true), spacingKeys)
		}
	case "deleted":
		if !diff.OldIsMap {
			out += fmt.Sprintf("%s\"value1\": %v,\n", spacingKeys, oldValue)
		} else {
			out += fmt.Sprintf("%s\"value1\": {\n%s%s},\n",
				spacingKeys,
				formatOutputJson(diff.DiffChild, level+1, true),
				spacingKeys)
		}
	}
	out = strings.TrimSuffix(out, ",\n") + "\n"
	out += spacingBraces + "},\n"
	return out
}

func formatNoDiff(diff parsers.Diff, level int, spacing string) string {
	spacingKeys := strings.Repeat(spacing, level+1)
	var out string
	var newValue string
	if diff.New == nil {
		newValue = "null"
	} else if isString(diff.New) {
		newValue = fmt.Sprintf("\"%v\"", diff.New)
	} else {
		newValue = fmt.Sprintf("%v", diff.New)
	}
	if !diff.NewIsMap {
		out += fmt.Sprintf("%s\"%s\": %v,\n", spacingKeys, diff.Key, newValue)
	} else {
		out += fmt.Sprintf("%s\"%s\": {\n%s%s},\n",
			spacingKeys,
			diff.Key,
			formatOutputJson(diff.DiffChild, level+1, true),
			spacingKeys)
	}
	return out
}
