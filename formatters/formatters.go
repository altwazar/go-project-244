package formatters

import "code/parsers"

func FormatOut(diff parsers.Diff, format string) string {
	var out string
	switch format {
	case "stylish":
		out = formatOutputStylish(diff)
	case "plain":
		out = formatOutputPlain(diff)
	case "json":
		out = formatOutputJSON(diff)
	}
	return out
}

func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}
