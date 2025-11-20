package formatters

import "code/parsers"

func FormatOut(diffs []parsers.Diff, format string) string {
	var out string
	switch format {
	case "stylish":
		out = formatOutputStylish(diffs, 0)
	case "plain":
		out = formatOutputPlain(diffs, "")
	case "json":
		out = formatOutputJson(diffs, 0, false)
	}
	return out
}
