package utils

import (
	"os"
	"path/filepath"
)

func ListFileRecursive(root string) ([]string, error) {
	var fileNames []string
	contents, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, content := range contents {
		if content.IsDir() {
			subContents, err := ListFileRecursive(filepath.Join(root, content.Name()))
			if err != nil {
				return nil, err
			}

			fileNames = append(fileNames, subContents...)
		} else {
			fileNames = append(fileNames, filepath.Join(root, content.Name()))
		}
	}

	return fileNames, nil
}
