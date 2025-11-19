package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Equal = iota
	Updated
	Removed
	Added
)
const (
	ValueToValue = iota
	ValueToMap
	MapToValue
	MapToMap
)

type Diff struct {
	State     int
	FromTo    int
	Key       string
	Old       any
	New       any
	DiffChild []Diff
}

func ParseJsonConfig(path string, cnf *any) error {
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

func ParseYamlConfig(path string, cnf *any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cnf); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

// Сравнение двух карт и получение из них структуры дифа
// На текущий момент эта функция для сравнения ключей вызывает compareValues
// А если значения мапы, то compareValues вызывает снова diffFromMaps
// Надо разобрать рекурсивный вызов
func diffFromMaps(a map[string]any, b map[string]any) []Diff {
	var diffs []Diff
	uniqueKeys := make(map[string]bool)
	for key := range a {
		uniqueKeys[key] = true
	}
	for key := range b {
		uniqueKeys[key] = true
	}
	for key := range uniqueKeys {
		diffs = append(diffs, compareValues(a, b, key))
	}
	return diffs
}

// Сравнение двух значений в мапах
func compareValues(old map[string]any, new map[string]any, key string) Diff {
	newMap, newIsMap := new[key].(map[string]any)
	oldMap, oldIsMap := old[key].(map[string]any)
	newKey, newExists := new[key]
	oldKey, oldExists := old[key]
	state := Equal
	fromto := ValueToValue
	diffChild := []Diff{}
	if newExists && oldExists {
		if !reflect.DeepEqual(newKey, oldKey) {
			state = Updated
		}
	} else if newExists {
		state = Added
	} else {
		state = Removed
	}
	if newIsMap && oldIsMap {
		diffChild = diffFromMaps(oldMap, newMap)
		fromto = MapToMap
	} else if newIsMap {
		diffChild = diffFromMaps(newMap, newMap)
		fromto = ValueToMap
	} else if oldIsMap {
		diffChild = diffFromMaps(oldMap, oldMap)
		fromto = MapToValue
	}
	return Diff{
		State:     state,
		FromTo:    fromto,
		Key:       key,
		Old:       oldKey,
		New:       newKey,
		DiffChild: diffChild,
	}
}
func ParseConfigs(pathBefore string, pathAfter string) ([]Diff, error) {
	var diffs []Diff
	cfgBefore, err := parseConfig(pathBefore)
	if err != nil {
		return diffs, err
	}
	cfgAfter, err := parseConfig(pathAfter)
	if err != nil {
		return diffs, err
	}
	mapBefore, ok := cfgBefore.(map[string]any)
	if !ok {
		err := fmt.Errorf("something wrong with config %s", pathBefore)
		return diffs, err
	}
	mapAfter, ok := cfgAfter.(map[string]any)
	if !ok {
		err := fmt.Errorf("something wrong with config %s", pathAfter)
		return diffs, err
	}
	diffs = diffFromMaps(mapBefore, mapAfter)
	return diffs, nil
}
func parseConfig(path string) (any, error) {
	var cnf any
	if strings.HasSuffix(path, ".json") {
		err := ParseJsonConfig(path, &cnf)
		if err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(path, ".yml") {
		err := ParseYamlConfig(path, &cnf)
		if err != nil {
			return nil, err
		}
	} else {
		err := fmt.Errorf("unknown file format %s", path)
		return nil, err
	}
	return cnf, nil
}
