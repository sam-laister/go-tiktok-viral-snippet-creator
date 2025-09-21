package helper

import (
	"fmt"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	return err == nil && fileInfo.IsDir()
}

func GetFilesInDirectory(path string) ([]string, error) {
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, fmt.Sprintf("%s/%s", path, entry.Name()))
	}

	return files, nil
}

func CreateDirectoryIfNotExists(path string) error {
	if Exists(path) {
		return nil
	}

	return os.MkdirAll(path, 0755)
}
