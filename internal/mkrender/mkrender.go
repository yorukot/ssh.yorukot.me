package mkrender

import (
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

func (Renderer) Render(markdown string, width int, bg string) (string, error) {
	if width < 1 {
		width = 1
	}

	renderer := stainmd.MochaStyles()
	if bg == "light" {
		renderer = stainmd.LatteStyles()
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
