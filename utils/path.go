package utils

import (
	"strings"
)

// Join will join path together which will make sure not leading or trailing "/"
func Join(in ...string) string {
	if len(in) == 0 {
		return ""
	}
	if strings.HasSuffix(in[0], "/") {
		if in[0] == "/" {
			return strings.TrimSuffix(strings.Join(in[1:], "/"), "/")
		}
		res := ""
		for k, v := range in {
			if k == 0 {
				v = strings.TrimPrefix(v, "/")
			} else {
				v = strings.TrimPrefix(v, "/")
				v = strings.TrimSuffix(v, "/")

				// Ignore empty string after trim.
				if v == "" {
					continue
				}
			}
			res += v
		}
		return res
	}
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
	if strings.HasSuffix(prefix, "/") {
		return strings.TrimSuffix(strings.TrimPrefix(full, strings.TrimPrefix(strings.TrimSuffix(prefix, "/"), "/")), "/")
	}
	// Remove the prefix from full.
	cp := strings.TrimPrefix(strings.Trim(full, "/"), strings.Trim(prefix, "/"))
	// Trim "/" to prevent object start or end with "/"
	cp = strings.Trim(cp, "/")

	return cp
}

func GetRelativePathStrict(prefix, full string) string {
	if len(prefix) > 1 && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:]
	}

	rel := strings.TrimPrefix(full, prefix)

	return rel
}

func RebuildPath(prefix, rel string) string {
	if prefix == "/" && rel == "/" {
		return ""
	}

	if len(prefix) > 1 && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:]
	}

	return prefix + rel
}
