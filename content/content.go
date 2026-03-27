package content

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
)

type Content struct {
	HeaderTitle string
}

func GetContent() Content {
	return Content{
		HeaderTitle: "Yorukot",
	}
}

func MarkdownContent(path string) (string, error) {
	baseDir := filepath.Join("content", "markdown")
	cleanPath := strings.Trim(pathutil.NormalizePath(path), "/")

	searchDir := baseDir
	if cleanPath != "" {
		searchDir = filepath.Join(baseDir, filepath.FromSlash(cleanPath))
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return "", err
	}

	var mdFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		mdFiles = append(mdFiles, entry.Name())
	}

	if len(mdFiles) == 0 {
		return "", errors.New("no markdown file found")
	}

	sort.Strings(mdFiles)

	content, err := os.ReadFile(filepath.Join(searchDir, mdFiles[0]))
	if err != nil {
		return "", err
	}

	return string(content), nil
}
