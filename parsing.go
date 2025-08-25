package code

import (
	"encoding/json"
	"fmt"
	"os"
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
	onlyBefore, onlyAfter, both := diffKeys(mapBefore, mapAfter)
	all := make([]string, 0, len(onlyAfter)+len(onlyBefore)+len(both))
	all = append(all, onlyAfter...)
	all = append(all, onlyBefore...)
	all = append(all, both...)
	// Отсортированные ключи из двух карт
	slices.Sort(all)
	for _, key := range both {
		if reflect.DeepEqual(mapBefore[key], mapAfter[key]) {
			outMap[key] = fmt.Sprintf("\t  %s: %v\n", key, mapBefore[key])
		} else {
			outMap[key] = fmt.Sprintf("\t- %s: %v\n\t+ %s: %v\n", key, mapBefore[key], key, mapAfter[key])
		}
	}
	for _, key := range onlyBefore {
		outMap[key] = fmt.Sprintf("\t- %s: %v\n", key, mapBefore[key])
	}
	for _, key := range onlyAfter {
		outMap[key] = fmt.Sprintf("\t+ %s: %v\n", key, mapAfter[key])
	}
	out += "{\n"
	for _, key := range all {
		out += outMap[key]
	}
	out += "}\n"
	return out, nil
}

func diffKeys(a map[string]any, b map[string]any) (keysAOnly, keysBOnly, keysBoth []string) {
	for k := range a {
		if _, ok := b[k]; ok {
			keysBoth = append(keysBoth, k)
		} else {
			keysAOnly = append(keysAOnly, k)
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			keysBOnly = append(keysBOnly, k)
		}
	}
	return
}

func parseConfig(path string) (any, error) {
	var cnf any
	if strings.HasSuffix(path, ".json") {
		err := parseJsonConfig(path, &cnf)
		if err != nil {
			return nil, err
		}

	} else {
		err := fmt.Errorf("unknown file format %s", path)
		return nil, err
	}
	return cnf, nil
}

func parseJsonConfig(path string, cnf *any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.UseNumber()
	if err := dec.Decode(&cnf); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
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
