package code

import (
	"code/parsers"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

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
	// Карта для форматированного вывода по ключам
	outMap := make(map[string]string)
	keysRemoved, keysAdded, keysEqual, keysUpdated := diffKeys(mapBefore, mapAfter)
	spacing := "  "
	for _, key := range keysEqual {
		outMap[key] = fmt.Sprintf("%s  %s: %v\n", spacing, key, mapBefore[key])
	}
	for _, key := range keysUpdated {
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
	for _, key := range keysRemoved {
		outMap[key] = fmt.Sprintf("%s- %s: %v\n", spacing, key, mapBefore[key])
	}
	for _, key := range keysAdded {
		outMap[key] = fmt.Sprintf("%s+ %s: %v\n", spacing, key, mapAfter[key])
	}

	all := make([]string, 0, len(keysAdded)+len(keysRemoved)+len(keysEqual)+len(keysUpdated))
	all = append(all, keysAdded...)
	all = append(all, keysRemoved...)
	all = append(all, keysEqual...)
	all = append(all, keysUpdated...)
	// Отсортированные ключи из двух карт
	slices.Sort(all)
	out += "{\n"
	for _, key := range all {
		out += outMap[key]
	}
	out += "}\n"
	return out, nil
}

func diffKeys(a map[string]any, b map[string]any) (keysRemoved, keysAdded, keysEqual, keysUpdated []string) {
	for k := range a {
		if _, ok := b[k]; ok {
			if reflect.DeepEqual(a[k], b[k]) {
				keysEqual = append(keysEqual, k)
			} else {
				keysUpdated = append(keysUpdated, k)
			}
		} else {
			keysRemoved = append(keysRemoved, k)
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			keysAdded = append(keysAdded, k)
		}
	}
	return
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
