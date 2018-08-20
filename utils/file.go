package utils

import (
	"os"
	"path/filepath"
)

// CreateFile will create a file recursively.
func CreateFile(p string) (f *os.File, err error) {
	p, err = Expand(p)
	if err != nil {
		return
	}

	err = os.MkdirAll(filepath.Dir(p), os.ModeDir|0777)
	if err != nil {
		return
	}

	return os.Create(p)
}
