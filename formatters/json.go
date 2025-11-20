package formatters

import (
	"code/parsers"
	"encoding/json"
	"log"
)

// Преобразование []diff в JSON, затем в строку
func formatOutputJson(diffs []parsers.Diff, level int) string {
	jsonData, err := json.MarshalIndent(diffs, "", "  ")
	if err != nil {
		log.Printf("warning: failed to marshal diff to json: %v", err)
	}
	return string(jsonData)
}

// Рекурсивный разбор дифов
// func formatOutputJson(diffs []parsers.Diff, level int) string {
// 	var out string

// 	// spacingBraces := strings.Repeat(" ", (level * 4))
// 	spacingKeys := strings.Repeat(" ", (level)*4+4)
// 	spacingCurvedBraces := strings.Repeat(" ", (level)*4+2)
// 	spacingBraces := strings.Repeat(" ", (level)*4)
// 	out += "[\n"

// 	sort.Slice(diffs, func(i, j int) bool {
// 		return diffs[i].Key < diffs[j].Key
// 	})
// 	var newValue string
// 	var oldValue string
// 	var state string
// 	for _, diff := range diffs {
// 		// Преобразование значения в строку в нужном виде
// 		key := fmt.Sprintf("\"%v\"", diff.Key)
// 		if diff.New == nil {
// 			newValue = "null"
// 		} else if isString(diff.New) {
// 			newValue = fmt.Sprintf("\"%v\"", diff.New)
// 		} else {
// 			newValue = fmt.Sprintf("%v", diff.New)
// 		}
// 		if diff.Old == nil {
// 			oldValue = "null"
// 		} else if isString(diff.Old) {
// 			oldValue = fmt.Sprintf("\"%v\"", diff.Old)
// 		} else {
// 			oldValue = fmt.Sprintf("%v", diff.Old)
// 		}
// 		switch diff.State {
// 		case parsers.Updated:
// 			state = "\"updated\""
// 		case parsers.Added:
// 			state = "\"added\""
// 		case parsers.Removed:
// 			state = "\"removed\""
// 		case parsers.Equal:
// 			state = "\"equal\""
// 		}
// 		out += spacingCurvedBraces + "{\n"
// 		out += fmt.Sprintf("%s\"state\": %v,\n", spacingKeys, state)
// 		out += fmt.Sprintf("%s\"key\": %v,\n", spacingKeys, key)
// 		switch diff.State {
// 		case parsers.Equal:
// 			if diff.FromTo == parsers.ValueToValue {
// 				out += fmt.Sprintf("%s\"value\": %v,\n", spacingKeys, oldValue)
// 			} else {
// 				out += fmt.Sprintf("%s\"value\": %s,\n", spacingKeys, formatOutputJson(diff.DiffChild, level+1))
// 			}
// 		case parsers.Updated:
// 			switch diff.FromTo {
// 			case parsers.ValueToValue:
// 				out += fmt.Sprintf("%s\"oldValue\": %v,\n", spacingKeys, oldValue)
// 				out += fmt.Sprintf("%s\"newValue\": %v,\n", spacingKeys, newValue)
// 			case parsers.MapToValue:
// 				out += fmt.Sprintf("%s\"oldValue\": %s,\n", spacingKeys, formatOutputJson(diff.DiffChild, level+1))
// 				out += fmt.Sprintf("%s\"newValue\": %v,\n", spacingKeys, newValue)
// 			case parsers.ValueToMap:
// 				out += fmt.Sprintf("%s\"oldValue\": %v,\n", spacingKeys, oldValue)
// 				out += fmt.Sprintf("%s\"newValue\": %s,\n", spacingKeys, formatOutputJson(diff.DiffChild, level+1))
// 			default:
// 				out += fmt.Sprintf("%s\"diff\": %s,\n", spacingKeys, formatOutputJson(diff.DiffChild, level+1))
// 			}
// 		case parsers.Added:
// 			if diff.FromTo == parsers.ValueToValue {
// 				out += fmt.Sprintf("%s\"newValue\": %v,\n", spacingKeys, newValue)
// 			} else {
// 				out += fmt.Sprintf("%s\"newValue\": %s,\n", spacingKeys, formatOutputJson(diff.DiffChild, level+1))
// 			}
// 		case parsers.Removed:
// 			if diff.FromTo == parsers.ValueToValue {
// 				out += fmt.Sprintf("%s\"oldValue\": %v,\n", spacingKeys, oldValue)
// 			} else {
// 				out += fmt.Sprintf("%s\"oldValue\": %s,\n", spacingKeys, formatOutputJson(diff.DiffChild, level+1))
// 			}
// 		}
// 		out = strings.TrimSuffix(out, ",\n") + "\n"
// 		out += spacingCurvedBraces + "},\n"
// 	}

// 	out = strings.TrimSuffix(out, ",\n") + "\n"
// 	if level == 0 {
// 		out += "]\n"
// 	} else {
// 		out += spacingBraces + "]"
// 	}
// 	return out
// }
