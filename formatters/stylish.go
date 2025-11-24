package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

const (
	indentSize = 4
	prefixSize = 2
)

func formatOutputStylish(rootDiff parsers.Diff) string {
	out := formatDiffChildStylish(rootDiff.DiffChild, 0)
	return strings.TrimSuffix(out, "\n")
}

// Рекурсивный разбор дифов
func formatDiffChildStylish(diffs []parsers.Diff, level int) string {
	builder := &stylishBuilder{
		level: level,
		diffs: diffs,
	}
	return builder.build()
}

type stylishBuilder struct {
	level int
	diffs []parsers.Diff
}

func (b *stylishBuilder) build() string {
	spacingBraces := strings.Repeat(" ", (b.level * indentSize))
	spacing := strings.Repeat(" ", (b.level+1)*indentSize-prefixSize)

	var out string
	out += "{\n"

	b.sortDiffs()
	out += b.buildDiffLines(spacing)
	out += spacingBraces + "}\n"

	return out
}

func (b *stylishBuilder) sortDiffs() {
	sort.Slice(b.diffs, func(i, j int) bool {
		return b.diffs[i].Key < b.diffs[j].Key
	})
}

func (b *stylishBuilder) buildDiffLines(spacing string) string {
	var out string
	for _, diff := range b.diffs {
		out += b.buildDiffLine(diff, spacing)
	}
	return out
}

func (b *stylishBuilder) buildDiffLine(diff parsers.Diff, spacing string) string {
	switch diff.State {
	case parsers.Equal:
		return b.buildEqualLine(diff, spacing)
	case parsers.Updated:
		return b.buildUpdatedLine(diff, spacing)
	case parsers.Added:
		return b.buildAddedLine(diff, spacing)
	case parsers.Removed:
		return b.buildRemovedLine(diff, spacing)
	default:
		return ""
	}
}

func (b *stylishBuilder) buildEqualLine(diff parsers.Diff, spacing string) string {
	if !diff.OldIsMap && !diff.NewIsMap {
		return fmt.Sprintf("%s  %s: %s\n", spacing, diff.Key, b.formatValue(diff.New))
	}
	return fmt.Sprintf("%s  %s: %s", spacing, diff.Key, formatDiffChildStylish(diff.DiffChild, b.level+1))
}

func (b *stylishBuilder) buildUpdatedLine(diff parsers.Diff, spacing string) string {
	switch {
	case !diff.OldIsMap && !diff.NewIsMap:
		return b.buildSimpleUpdated(diff, spacing)
	case diff.OldIsMap && !diff.NewIsMap:
		return b.buildMapToSimpleUpdated(diff, spacing)
	case !diff.OldIsMap && diff.NewIsMap:
		return b.buildSimpleToMapUpdated(diff, spacing)
	default:
		return b.buildNestedUpdated(diff, spacing)
	}
}

func (b *stylishBuilder) buildSimpleUpdated(diff parsers.Diff, spacing string) string {
	return fmt.Sprintf("%s- %s: %s\n%s+ %s: %s\n",
		spacing, diff.Key, b.formatValue(diff.Old),
		spacing, diff.Key, b.formatValue(diff.New))
}

func (b *stylishBuilder) buildMapToSimpleUpdated(diff parsers.Diff, spacing string) string {
	return fmt.Sprintf("%s- %s: %s%s+ %s: %s\n",
		spacing, diff.Key, formatDiffChildStylish(diff.DiffChild, b.level+1),
		spacing, diff.Key, b.formatValue(diff.New))
}

func (b *stylishBuilder) buildSimpleToMapUpdated(diff parsers.Diff, spacing string) string {
	return fmt.Sprintf("%s- %s: %s\n%s+ %s: %s",
		spacing, diff.Key, b.formatValue(diff.Old),
		spacing, diff.Key, formatDiffChildStylish(diff.DiffChild, b.level+1))
}

func (b *stylishBuilder) buildNestedUpdated(diff parsers.Diff, spacing string) string {
	return fmt.Sprintf("%s  %s: %s", spacing, diff.Key, formatDiffChildStylish(diff.DiffChild, b.level+1))
}

func (b *stylishBuilder) buildAddedLine(diff parsers.Diff, spacing string) string {
	if !diff.OldIsMap && !diff.NewIsMap {
		return fmt.Sprintf("%s+ %s: %s\n", spacing, diff.Key, b.formatValue(diff.New))
	}
	return fmt.Sprintf("%s+ %s: %s", spacing, diff.Key, formatDiffChildStylish(diff.DiffChild, b.level+1))
}

func (b *stylishBuilder) buildRemovedLine(diff parsers.Diff, spacing string) string {
	if !diff.OldIsMap && !diff.NewIsMap {
		return fmt.Sprintf("%s- %s: %s\n", spacing, diff.Key, b.formatValue(diff.Old))
	}
	return fmt.Sprintf("%s- %s: %s", spacing, diff.Key, formatDiffChildStylish(diff.DiffChild, b.level+1))
}

func (b *stylishBuilder) formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", value)
}
