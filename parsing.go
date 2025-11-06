package code

import (
	"code/parsers"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type pair struct {
	First  any
	Second any
}

func CompareConfigs(pathBefore string, pathAfter string, format string) (string, error) {
	var out string
	cfgBefore, err := parseConfig(pathBefore)
	if err != nil {
		return "", err
	}
	cfgAfter, err := parseConfig(pathAfter)
	if err != nil {
		return "", err
	}
	mapBefore, ok := cfgBefore.(map[string]any)
	if !ok {
		err := fmt.Errorf("something wrong with config %s", pathBefore)
		return "", err
	}
	mapAfter, ok := cfgAfter.(map[string]any)
	if !ok {
		err := fmt.Errorf("something wrong with config %s", pathAfter)
		return "", err
	}
	out = formatOutput(mapBefore, mapAfter, 1)
	return out, nil
}

func formatOutput(mapBefore map[string]any, mapAfter map[string]any, identlvl int) string {
	// Карта для форматированного вывода по ключам
	var out string
	outMap := make(map[string]string)
	keysRemoved, keysAdded, keysEqual, keysUpdated := diffKeys(mapBefore, mapAfter)
	spacing := strings.Repeat(" ", identlvl*4-2)
	for key := range keysEqual {
		if mapVal, isMap := mapBefore[key].(map[string]any); isMap {
			outMap[key] = fmt.Sprintf("%s  %s: %s\n",
				spacing,
				key,
				formatOutput(mapVal, mapVal, identlvl+1))
		} else {
			outMap[key] = fmt.Sprintf("%s  %s: %v\n", spacing, key, mapBefore[key])
		}
	}
	for key := range keysUpdated {
		if beforeMap, beforeIsMap := mapBefore[key].(map[string]any); beforeIsMap {
			if afterMap, afterIsMap := mapAfter[key].(map[string]any); afterIsMap {
				// Оба значения - map[string]any
				outMap[key] = fmt.Sprintf("%s  %s: %s\n",
					spacing,
					key,
					formatOutput(beforeMap, afterMap, identlvl+1))
			} else {
				// Только before является map
				outMap[key] = fmt.Sprintf("%s- %s: %s\n%s+  %s: %v\n",
					spacing,
					key,
					formatOutput(beforeMap, beforeMap, identlvl+1),
					spacing,
					key,
					mapAfter[key])
			}
		} else if afterMap, afterIsMap := mapAfter[key].(map[string]any); afterIsMap {
			// Только after является map
			// outMap[key] = formatOutput(map[string]any{"value": mapBefore[key]}, afterMap, identlvl+1)
			outMap[key] = fmt.Sprintf("%s- %s: %v\n%s+  %s: %s\n",
				spacing,
				key,
				mapBefore[key],
				spacing,
				key,
				formatOutput(afterMap, afterMap, identlvl+1))
		} else {
			outMap[key] = fmt.Sprintf(
				"%s- %s: %v\n%s+ %s: %v\n",
				spacing,
				key,
				mapBefore[key],
				spacing,
				key,
				mapAfter[key],
			)
		}

	}
	for key := range keysRemoved {
		if mapVal, isMap := mapBefore[key].(map[string]any); isMap {
			outMap[key] = fmt.Sprintf("%s- %s: %s\n",
				spacing,
				key,
				formatOutput(mapVal, mapVal, identlvl+1))
		} else {
			outMap[key] = fmt.Sprintf("%s- %s: %v\n", spacing, key, mapBefore[key])
		}
	}
	for key := range keysAdded {
		if mapVal, isMap := mapAfter[key].(map[string]any); isMap {
			outMap[key] = fmt.Sprintf("%s+ %s: %s\n",
				spacing,
				key,
				formatOutput(mapVal, mapVal, identlvl+1))
		} else {
			outMap[key] = fmt.Sprintf("%s+ %s: %v\n", spacing, key, mapAfter[key])
		}
	}

	all := make([]string, 0, len(keysAdded)+len(keysRemoved)+len(keysEqual)+len(keysUpdated))
	for key := range keysAdded {
		all = append(all, key)
	}
	for key := range keysRemoved {
		all = append(all, key)
	}
	for key := range keysEqual {
		all = append(all, key)
	}
	for key := range keysUpdated {
		all = append(all, key)
	}
	// Отсортированные ключи из двух карт
	slices.Sort(all)
	out += "{\n"
	for _, key := range all {
		out += outMap[key]
	}
	out += strings.Repeat(" ", 4*(identlvl-1)) + "}"
	return out
}

func diffKeys(a map[string]any, b map[string]any) (
	map[string]any,
	map[string]any,
	map[string]any,
	map[string]pair,
) {
	keysRemoved := make(map[string]any)
	keysAdded := make(map[string]any)
	keysEqual := make(map[string]any)
	keysUpdated := make(map[string]pair)
	for k := range a {
		if _, ok := b[k]; ok {
			if reflect.DeepEqual(a[k], b[k]) {
				keysEqual[k] = a[k]
			} else {
				keysUpdated[k] = pair{
					First:  a[k],
					Second: b[k],
				}
			}
		} else {
			keysRemoved[k] = a[k]
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			keysAdded[k] = b[k]
		}
	}
	return keysRemoved, keysAdded, keysEqual, keysUpdated
}

func parseConfig(path string) (any, error) {
	var cnf any
	if strings.HasSuffix(path, ".json") {
		err := parsers.ParseJsonConfig(path, &cnf)
		if err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(path, ".yml") {
		err := parsers.ParseYamlConfig(path, &cnf)
		if err != nil {
			return nil, err
		}
	} else {
		err := fmt.Errorf("unknown file format %s", path)
		return nil, err
	}
	return cnf, nil
}

// Вывод содержимого конфига
func walkConfig(prefix string, v any) {
	switch x := v.(type) {
	case map[string]any:
		for k, vv := range x {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			walkConfig(key, vv)
		}
	case []any:
		for i, vv := range x {
			var key string
			if prefix == "" {
				key = fmt.Sprintf("[%d]", i)
			} else {
				key = fmt.Sprintf("%s[%d]", prefix, i)
			}
			walkConfig(key, vv)
		}
	default:
		fmt.Printf("%s = %v (type=%T)\n", prefix, x, x)
	}
}

// Получение значения из конфига
func lookupInConfig(cfg any, key string) (any, bool) {
	// По ключу могут быть разные значения
	var res any
	if key == "" {
		return cfg, true
	}
	m, ok := cfg.(map[string]any)
	if !ok {
		return nil, false
	}
	var exists bool
	res, exists = m[key]
	if !exists {
		return nil, false
	}
	return res, true
}
