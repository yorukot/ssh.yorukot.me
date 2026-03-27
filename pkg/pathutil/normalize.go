package pathutil

import "strings"

func NormalizePath(path string) string {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" || trimmedPath == "/" {
		return "/"
	}

	return "/" + strings.TrimLeft(trimmedPath, "/")
}
