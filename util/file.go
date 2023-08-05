package util

import (
	"os"
	"path/filepath"
)

func OpenOrCreateFileAndNestedFolders(filePath string) *os.File {
	if err := os.MkdirAll(filepath.Dir(filePath), 0770); err != nil {
		panic(err)
	}
	if file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777); err != nil {
		panic(err)
	} else {
		return file
	}
}
