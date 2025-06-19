package def

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Enterprise string

func (e Enterprise) Validate(rulesDir string) error {
	if e == "" {
		return fmt.Errorf("enterprise value cannot be empty")
	}

	// Check if the string has a .yml or .yaml extension
	if !strings.HasSuffix(string(e), ".yml") && !strings.HasSuffix(string(e), ".yaml") {
		return fmt.Errorf("enterprise value must be a path to a YAML file (with .yml or .yaml extension): %s", e)
	}

	// Check if the file exists
	filePath := fmt.Sprintf("%s/%s", rulesDir, string(e))
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of enterprise file '%s': %w", filePath, err)
	}

	_, err = os.Stat(absFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("enterprise file does not exist: %s", filePath)
		}
		return fmt.Errorf("error checking if enterprise file exists: %w", err)
	}

	return nil
}
