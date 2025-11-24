package parsers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Статус сравнения значений ключа
const (
	Equal = iota
	Updated
	Removed
	Added
	Root
)

// Структура для хранения результата сравнения конфигов
// DiffChild поле для вложенных структур, если один из ключей оказался такой структурой
type Diff struct {
	State     int    `json:"state"`
	Key       string `json:"key"`
	OldIsMap  bool   `json:"oldismap"`
	NewIsMap  bool   `json:"newismap"`
	Old       any    `json:"old,omitempty"`
	New       any    `json:"new,omitempty"`
	DiffChild []Diff `json:"diff_child,omitempty"`
}

// Парсинг конфигов, преобразование их в виде map[string]any
// Затем получение списка изменений []Diffs
func ParseConfigs(pathBefore string, pathAfter string) (Diff, error) {
	diff := Diff{
		State:    Root,
		Key:      "",
		NewIsMap: true,
	}
	cfgBefore, err := parseConfig(pathBefore)
	if err != nil {
		return diff, err
	}
	cfgAfter, err := parseConfig(pathAfter)
	if err != nil {
		return diff, err
	}
	mapBefore, ok := cfgBefore.(map[string]any)
	if !ok {
		err := fmt.Errorf("something wrong with config %s", pathBefore)
		return diff, err
	}
	mapAfter, ok := cfgAfter.(map[string]any)
	if !ok {
		err := fmt.Errorf("something wrong with config %s", pathAfter)
		return diff, err
	}
	diff.DiffChild = diffFromMaps(mapBefore, mapAfter)
	return diff, nil
}

// Преобразование конфига
func parseConfig(path string) (any, error) {
	var cnf any
	var err error

	switch {
	case strings.HasSuffix(path, ".json"):
		err = ParseJsonConfig(path, &cnf)
	case strings.HasSuffix(path, ".yml"):
		err = ParseYamlConfig(path, &cnf)
	default:
		err = fmt.Errorf("unknown file format %s", path)
	}

	if err != nil {
		return nil, err
	}
	return cnf, nil
}

// Преобразование JSON
func ParseJsonConfig(path string, cnf *any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file: %v", err)
		}
	}()
	dec := json.NewDecoder(f)
	dec.UseNumber()
	if err := dec.Decode(&cnf); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

// Преобразование Yaml
func ParseYamlConfig(path string, cnf *any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file: %v", err)
		}
	}()
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
	// в diffs попадают результаты сравнения
	var diffs []Diff
	// Список ключей в двух мапах
	uniqueKeys := make(map[string]bool)
	for key := range a {
		uniqueKeys[key] = true
	}
	for key := range b {
		uniqueKeys[key] = true
	}
	// Сравнение по полученному списку ключей
	for key := range uniqueKeys {
		diffs = append(diffs, compareValues(a, b, key))
	}
	return diffs
}

// Сравнение двух значений в мапах
func compareValues(old map[string]any, new map[string]any, key string) Diff {
	// Получение потенциальных map из ключей
	newMap, newIsMap := new[key].(map[string]any)
	oldMap, oldIsMap := old[key].(map[string]any)
	// Получение значений из ключей
	newValue, newExists := new[key]
	oldValue, oldExists := old[key]

	state := Equal
	diffChild := []Diff{}

	// Получение статусов изменения значений ключа
	switch {
	case newExists && oldExists:
		if !reflect.DeepEqual(newValue, oldValue) {
			state = Updated
		}
	case newExists:
		state = Added
	default:
		state = Removed
	}

	// Обработка случаев, когда одно из значений является вложенной структурой
	switch {
	case newIsMap && oldIsMap:
		diffChild = diffFromMaps(oldMap, newMap)
	case newIsMap:
		diffChild = diffFromMaps(newMap, newMap)
	case oldIsMap:
		diffChild = diffFromMaps(oldMap, oldMap)
	}

	// это попадет в массив diffs функции diffFromMaps
	return Diff{
		State:     state,
		NewIsMap:  newIsMap,
		OldIsMap:  oldIsMap,
		Key:       key,
		Old:       oldValue,
		New:       newValue,
		DiffChild: diffChild,
	}
}
