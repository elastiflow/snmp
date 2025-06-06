package reader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Definition interface {
	Type() string
	Validate() error
}

func ReadAndValidate[T Definition](dirPath string) (map[string]T, error) {
	filePaths, err := walkDirectory(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	definitions := make(map[string]T)

	for _, filePath := range filePaths {
		definitionsInFile, parseErr := read[T](filePath)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse definitions: %w", err)
		}

		err = validateUniqueness(definitions, definitionsInFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate definitions in file %s: %w", filePath, err)
		}

		err = validateDefinitions(definitionsInFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate definitions in file %s: %w", filePath, err)
		}

		mergeDefinitions(definitions, definitionsInFile)
	}

	return definitions, nil
}

func read[T Definition](filePath string) (map[string]T, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	definitions := make(map[string]T)
	err = yaml.Unmarshal(fileBytes, &definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}

	return definitions, nil
}

func validateUniqueness[T Definition](definitionsA, definitionsB map[string]T) error {
	for id := range definitionsB {
		if definition, ok := definitionsA[id]; ok {
			return fmt.Errorf("found duplicate %s with id %s", definition.Type(), id)
		}
	}
	return nil
}

func validateDefinitions[T Definition](definitions map[string]T) error {
	for id, definition := range definitions {
		err := definition.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate %s with id %s: %w", definition.Type(), id, err)
		}
	}
	return nil
}

func mergeDefinitions[T Definition](definitionsA, definitionsB map[string]T) {
	for id, definition := range definitionsB {
		definitionsA[id] = definition
	}
}
