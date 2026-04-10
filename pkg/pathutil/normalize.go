package pathutil

import "strings"

func NormalizePath(path string) string {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" || trimmedPath == "/" {
		return "/"
	}

	normalized := "/" + strings.Trim(trimmedPath, "/")
	if normalized == "" {
		return "/"
	}

	return normalized
}
