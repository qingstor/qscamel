package utils

import (
	"os"
	"path"
)

// CreateFile will create a file recursively.
func CreateFile(p string) (f *os.File, err error) {
	p, err = Expand(p)
	if err != nil {
		return
	}

	err = os.MkdirAll(path.Dir(p), 0711)
	if err != nil {
		return
	}

	return os.Create(p)
}
