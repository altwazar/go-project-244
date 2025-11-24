package formatters

import (
	"code/parsers"
	"fmt"
	"sort"
	"strings"
)

// Рекурсивный разбор изменений в плоском формате
func formatOutputPlain(rootDiff parsers.Diff) string {
	out := formatDiffChildPlain(rootDiff.DiffChild, "")
	return strings.TrimSuffix(out, "\n")
}

func formatDiffChildPlain(diffs []parsers.Diff, path string) string {
	builder := &plainBuilder{
		path:  path,
		diffs: diffs,
	}
	return builder.build()
}

type plainBuilder struct {
	path  string
	diffs []parsers.Diff
}

func (b *plainBuilder) build() string {
	b.sortDiffs()

	var out string
	for _, diff := range b.diffs {
		out += b.buildDiffLine(diff)
	}
	return out
}

func (b *plainBuilder) sortDiffs() {
	sort.Slice(b.diffs, func(i, j int) bool {
		return b.diffs[i].Key < b.diffs[j].Key
	})
}

func (b *plainBuilder) buildDiffLine(diff parsers.Diff) string {
	switch diff.State {
	case parsers.Updated:
		return b.buildUpdatedLine(diff)
	case parsers.Added:
		return b.buildAddedLine(diff)
	case parsers.Removed:
		return b.buildRemovedLine(diff)
	default:
		return ""
	}
}

func (b *plainBuilder) buildUpdatedLine(diff parsers.Diff) string {
	pathKey := b.buildPathKey(diff.Key)

	switch {
	case !diff.OldIsMap && !diff.NewIsMap:
		return fmt.Sprintf("Property '%s' was updated. From %s to %s\n",
			pathKey, b.formatValue(diff.Old), b.formatValue(diff.New))
	case diff.OldIsMap && !diff.NewIsMap:
		return fmt.Sprintf("Property '%s' was updated. From [complex value] to %s\n",
			pathKey, b.formatValue(diff.New))
	case !diff.OldIsMap && diff.NewIsMap:
		return fmt.Sprintf("Property '%s' was updated. From %s to [complex value]\n",
			pathKey, b.formatValue(diff.Old))
	default:
		return formatDiffChildPlain(diff.DiffChild, pathKey)
	}
}

func (b *plainBuilder) buildAddedLine(diff parsers.Diff) string {
	pathKey := b.buildPathKey(diff.Key)

	if !diff.OldIsMap && !diff.NewIsMap {
		return fmt.Sprintf("Property '%s' was added with value: %s\n",
			pathKey, b.formatValue(diff.New))
	}
	return fmt.Sprintf("Property '%s' was added with value: [complex value]\n", pathKey)
}

func (b *plainBuilder) buildRemovedLine(diff parsers.Diff) string {
	pathKey := b.buildPathKey(diff.Key)
	return fmt.Sprintf("Property '%s' was removed\n", pathKey)
}

func (b *plainBuilder) buildPathKey(key string) string {
	if b.path != "" {
		return b.path + "." + key
	}
	return key
}

func (b *plainBuilder) formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	if isString(value) {
		return fmt.Sprintf("'%v'", value)
	}
	return fmt.Sprintf("%v", value)
}
