package utils

import (
	"path"
	"strings"
)

// Join will join path together which will make sure not leading or trailing "/"
func Join(in ...string) string {
	// Add "/" to list specific prefix.
	cp := path.Join(in...)
	// Trim "/" to prevent object start or end with "/"
	cp = strings.Trim(cp, "/")

	return cp
}

// Relative will calculate the relative path for full by prefix.
func Relative(full, prefix string) string {
	// Remove the prefix from full.
	cp := strings.TrimPrefix(strings.Trim(full, "/"), strings.Trim(prefix, "/"))
	// Trim "/" to prevent object start or end with "/"
	cp = strings.Trim(cp, "/")

	return cp
}
