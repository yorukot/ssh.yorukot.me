package mkrender

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yorukot/ssh.yorukot.me/pkg/stainmd"
)

const defaultImageManifestPath = "content/blog_image_manifest.json"

type Renderer struct {
	imageManifest map[string]string
}

type imageManifestFile struct {
	Site   string            `json:"site"`
	Images map[string]string `json:"images"`
}

var rawTagLine = regexp.MustCompile(`^</?[A-Za-z][^>]*>$`)
var rawMarkerLine = regexp.MustCompile(`^\{\{[^{}]+\}\}$`)

func New() Renderer {
	return Renderer{
		imageManifest: loadImageManifest(defaultImageManifestPath),
	}
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
			return resolveMarkdownImagePath(sourcePath, destination, r.imageManifest)
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

func resolveMarkdownImagePath(sourcePath, destination string, manifest map[string]string) string {
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
	resolved := filepath.Clean(filepath.Join(filepath.Dir(sourcePath), pathPart))
	if manifestURL := manifest[filepath.ToSlash(resolved)]; manifestURL != "" {
		return manifestURL + suffix
	}
	return resolved + suffix
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

func loadImageManifest(path string) map[string]string {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var manifest imageManifestFile
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil
	}
	if len(manifest.Images) == 0 {
		return nil
	}

	images := make(map[string]string, len(manifest.Images))
	for path, url := range manifest.Images {
		path = strings.TrimSpace(filepath.ToSlash(path))
		url = strings.TrimSpace(url)
		if path == "" || url == "" {
			continue
		}
		images[path] = url
	}
	return images
}
