package parsers

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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
