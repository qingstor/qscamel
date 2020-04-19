package utils

import (
	"strings"
)

// Join will join path together which will make sure not leading or trailing "/"
func Join(in ...string) string {
	x := make([]string, 0)
	for k, v := range in {
		if k == 0 {
			v = strings.TrimPrefix(v, "/")
		}
		// Trim all trailing "/"
		v = strings.TrimRight(v, "/")

		// Ignore empty string after trim.
		if v == "" {
			continue
		}

		x = append(x, v)
	}
	if len(x) == 0 {
		return ""
	}

	// Add "/" to list specific prefix.
	cp := strings.Join(x, "/")
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
