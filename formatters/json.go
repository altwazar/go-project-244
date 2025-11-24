package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

// Рекурсивное формирование json вывода
func formatOutputJSON(rootDiff parsers.Diff) string {
	spacing := strings.Repeat(" ", 2)
	out := formatDiff(rootDiff, 0, spacing)
	out = strings.TrimSuffix(out, ",\n")
	return out
}

// Обработка списка диффов
func diffRange(diffs []parsers.Diff, level int, spacing string, noDiff bool) string {
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})

	var out string
	for _, diff := range diffs {
		if noDiff {
			out += formatNoDiff(diff, level, spacing)
		} else {
			out += formatDiff(diff, level, spacing)
		}
	}

	out = strings.TrimSuffix(out, ",\n") + "\n"
	if !noDiff {
		out += strings.Repeat(spacing, level-1) + "]"
	}
	return out
}

// Формирование вывода различий
func formatDiff(diff parsers.Diff, level int, spacing string) string {
	builder := &diffBuilder{
		diff:          diff,
		level:         level,
		spacing:       spacing,
		spacingBraces: strings.Repeat(spacing, level),
		spacingKeys:   strings.Repeat(spacing, level+1),
	}

	return builder.build()
}

type diffBuilder struct {
	diff          parsers.Diff
	level         int
	spacing       string
	spacingBraces string
	spacingKeys   string
}

func (b *diffBuilder) build() string {
	state := b.getState()

	var out string
	out += b.spacingBraces + "{\n"
	out += fmt.Sprintf("%s\"key\": \"%s\",\n", b.spacingKeys, b.diff.Key)
	out += fmt.Sprintf("%s\"type\": \"%s\",\n", b.spacingKeys, state)
	out += b.buildValueSection(state)
	out = strings.TrimSuffix(out, ",\n") + "\n"
	out += b.spacingBraces + "},\n"

	return out
}

func (b *diffBuilder) getState() string {
	switch b.diff.State {
	case parsers.Updated:
		if b.diff.OldIsMap && b.diff.NewIsMap {
			return "nested"
		}
		return "changed"
	case parsers.Added:
		return "added"
	case parsers.Removed:
		return "deleted"
	case parsers.Equal:
		return "unchanged"
	case parsers.Root:
		return "root"
	default:
		return "unknown"
	}
}

func (b *diffBuilder) buildValueSection(state string) string {
	switch state {
	case "unchanged":
		return b.buildUnchangedValue()
	case "changed":
		return b.buildChangedValue()
	case "nested", "root":
		return b.buildNestedValue()
	case "added":
		return b.buildAddedValue()
	case "deleted":
		return b.buildDeletedValue()
	default:
		return ""
	}
}

func (b *diffBuilder) buildUnchangedValue() string {
	if !b.diff.OldIsMap {
		return fmt.Sprintf("%s\"value1\": %v,\n", b.spacingKeys, b.formatValue(b.diff.Old))
	}
	return fmt.Sprintf("%s\"value1\": {\n%s\n%s},\n",
		b.spacingKeys,
		diffRange(b.diff.DiffChild, b.level+1, b.spacing, true),
		b.spacingKeys)
}

func (b *diffBuilder) buildChangedValue() string {
	var out string

	switch {
	case !b.diff.OldIsMap && !b.diff.NewIsMap:
		out += b.buildOldValueIfExists()
		out += b.buildNewValueIfExists()
	case b.diff.OldIsMap && !b.diff.NewIsMap:
		out += b.buildOldMapValue()
		out += b.buildNewValueIfExists()
	case !b.diff.OldIsMap && b.diff.NewIsMap:
		out += b.buildOldValueIfExists()
		out += b.buildNewMapValue()
	}

	return out
}

func (b *diffBuilder) buildNestedValue() string {
	return fmt.Sprintf("%s\"children\": [\n%s,\n", b.spacingKeys,
		diffRange(b.diff.DiffChild, b.level+2, b.spacing, false))
}

func (b *diffBuilder) buildAddedValue() string {
	if !b.diff.NewIsMap {
		return fmt.Sprintf("%s\"value2\": %v,\n", b.spacingKeys, b.formatValue(b.diff.New))
	}
	return fmt.Sprintf("%s\"value2\": {\n%s%s},\n", b.spacingKeys,
		diffRange(b.diff.DiffChild, b.level+1, b.spacing, true), b.spacingKeys)
}

func (b *diffBuilder) buildDeletedValue() string {
	if !b.diff.OldIsMap {
		return fmt.Sprintf("%s\"value1\": %v,\n", b.spacingKeys, b.formatValue(b.diff.Old))
	}
	return fmt.Sprintf("%s\"value1\": {\n%s%s},\n",
		b.spacingKeys,
		diffRange(b.diff.DiffChild, b.level+1, b.spacing, true),
		b.spacingKeys)
}

func (b *diffBuilder) buildOldValueIfExists() string {
	if b.diff.Old != nil {
		return fmt.Sprintf("%s\"value1\": %v,\n", b.spacingKeys, b.formatValue(b.diff.Old))
	}
	return ""
}

func (b *diffBuilder) buildNewValueIfExists() string {
	if b.diff.New != nil {
		return fmt.Sprintf("%s\"value2\": %v,\n", b.spacingKeys, b.formatValue(b.diff.New))
	}
	return ""
}

func (b *diffBuilder) buildOldMapValue() string {
	return fmt.Sprintf("%s\"value1\": {\n%s%s},\n", b.spacingKeys,
		diffRange(b.diff.DiffChild, b.level+1, b.spacing, true), b.spacingKeys)
}

func (b *diffBuilder) buildNewMapValue() string {
	return fmt.Sprintf("%s\"value2\": {\n%s%s},\n", b.spacingKeys,
		diffRange(b.diff.DiffChild, b.level+1, b.spacing, true), b.spacingKeys)
}

func (b *diffBuilder) formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	if isString(value) {
		return fmt.Sprintf("\"%v\"", value)
	}
	return fmt.Sprintf("%v", value)
}

func formatNoDiff(diff parsers.Diff, level int, spacing string) string {
	spacingKeys := strings.Repeat(spacing, level+1)
	var out string

	formattedValue := formatSimpleValue(diff.New)

	if !diff.NewIsMap {
		out += fmt.Sprintf("%s\"%s\": %v,\n", spacingKeys, diff.Key, formattedValue)
	} else {
		out += fmt.Sprintf("%s\"%s\": {\n%s%s},\n",
			spacingKeys,
			diff.Key,
			diffRange(diff.DiffChild, level+1, spacing, true),
			spacingKeys)
	}
	return out
}

func formatSimpleValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	if isString(value) {
		return fmt.Sprintf("\"%v\"", value)
	}
	return fmt.Sprintf("%v", value)
}
