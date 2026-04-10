package pathutil

import "strings"

func ParentPath(path string) string {
	normalized := NormalizePath(path)
	if normalized == "/" {
		return "/"
	}

	idx := strings.LastIndex(normalized, "/")
	if idx <= 0 {
		return "/"
	}

	return normalized[:idx]
}
