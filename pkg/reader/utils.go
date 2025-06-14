package reader

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

func walkDirectory(dirPath string) ([]string, error) {
	test, _ := filepath.Abs(dirPath)

	dir := &directory{}
	err := filepath.WalkDir(test, dir.walk)
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
