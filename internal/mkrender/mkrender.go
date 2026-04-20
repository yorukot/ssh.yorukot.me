package mkrender

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yorukot/ssh.yorukot.me/pkg/stainmd"
)

type Renderer struct{}

var rawTagLine = regexp.MustCompile(`^</?[A-Za-z][^>]*>$`)
var rawMarkerLine = regexp.MustCompile(`^\{\{[^{}]+\}\}$`)

func New() Renderer {
	return Renderer{}
}

func (r Renderer) Render(markdown string, width int, bg string) (string, error) {
	return r.RenderWithSource(markdown, "", width, bg)
}

func (r Renderer) RenderWithSource(markdown, sourcePath string, width int, bg string) (string, error) {
	if width < 1 {
		width = 1
	}

	renderer := stainmd.MochaStyles()
	if bg == "light" {
		renderer = stainmd.LatteStyles()
	}
	if sourcePath != "" {
		renderer.ImagePathResolver = func(destination string) string {
			return resolveMarkdownImagePath(sourcePath, destination)
		}
	}

	return renderer.Render(stripRawTagLines(markdown), width)
}

func stripRawTagLines(markdown string) string {
	lines := strings.Split(markdown, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if rawTagLine.MatchString(trimmed) || rawMarkerLine.MatchString(trimmed) {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, "\n")
}

func resolveMarkdownImagePath(sourcePath, destination string) string {
	destination = strings.TrimSpace(destination)
	if destination == "" || strings.HasPrefix(destination, "/") {
		return destination
	}

	if parsed, err := url.Parse(destination); err == nil && parsed.Scheme != "" {
		return destination
	}

	if !strings.HasPrefix(destination, "./") && !strings.HasPrefix(destination, "../") {
		return destination
	}

	pathPart, suffix := splitPathSuffix(destination)
	return filepath.Clean(filepath.Join(filepath.Dir(sourcePath), pathPart)) + suffix
}

func splitPathSuffix(value string) (string, string) {
	cut := len(value)
	for _, marker := range []string{"?", "#"} {
		if index := strings.Index(value, marker); index >= 0 && index < cut {
			cut = index
		}
	}
	return value[:cut], value[cut:]
}
