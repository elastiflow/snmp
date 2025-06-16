package reader

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

func walkDirectory(dirPath string) ([]string, error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of directory '%s': %w", dirPath, err)
	}

	dir := &directory{}
	err = filepath.WalkDir(absDirPath, dir.walk)
	if err != nil {
		return nil, fmt.Errorf("error while walking directory '%s' - %w", dirPath, err)
	}
	return dir.filePaths, nil
}

type directory struct {
	filePaths []string
}

func (d *directory) walk(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !entry.IsDir() && isYamlFile(path) {
		d.filePaths = append(d.filePaths, path)
	}
	return nil
}

func isYamlFile(filePath string) bool {
	yamlExtension := regexp.MustCompile(`\.ya?ml$`)
	return yamlExtension.MatchString(filePath)
}
