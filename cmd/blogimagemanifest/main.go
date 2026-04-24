package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
)

type manifest struct {
	Site   string            `json:"site"`
	Images map[string]string `json:"images"`
}

type postMetadata struct {
	Slug  string `yaml:"post_slug"`
	Draft bool   `yaml:"draft"`
}

type localImage struct {
	Source string
}

var (
	markdownImagePattern = regexp.MustCompile(`!\[[^\]]*\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	htmlImagePattern     = regexp.MustCompile(`(?is)<img\b[^>]*\bsrc=["']([^"']+)["'][^>]*>`)
)

func main() {
	blogRoot := flag.String("blog-root", filepath.Join("content", "markdown", "blog"), "blog markdown root")
	distRoot := flag.String("dist", filepath.Join("yorukot.me", "dist"), "Astro dist root")
	siteURL := flag.String("site", "https://yorukot.me", "site domain for absolute image URLs")
	output := flag.String("output", filepath.Join("content", "blog_image_manifest.json"), "manifest output path")
	flag.Parse()

	manifest, err := buildManifest(*blogRoot, *distRoot, *siteURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate blog image manifest: %v\n", err)
		os.Exit(1)
	}

	body, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode blog image manifest: %v\n", err)
		os.Exit(1)
	}
	body = append(body, '\n')

	if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create manifest directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*output, body, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write blog image manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %d blog image mappings to %s\n", len(manifest.Images), *output)
}

func buildManifest(blogRoot, distRoot, siteURL string) (manifest, error) {
	siteURL = strings.TrimRight(siteURL, "/")
	if siteURL == "" {
		return manifest{}, errors.New("site URL is required")
	}

	markdownFiles, err := findMarkdownFiles(blogRoot)
	if err != nil {
		return manifest{}, err
	}

	result := manifest{
		Site:   siteURL,
		Images: make(map[string]string),
	}

	for _, markdownFile := range markdownFiles {
		post, err := postMetadataForFile(markdownFile)
		if err != nil {
			return manifest{}, err
		}
		if post.Draft {
			continue
		}

		localImages, err := localMarkdownImages(markdownFile)
		if err != nil {
			return manifest{}, err
		}
		localImages = uniqueLocalImages(localImages)
		if len(localImages) == 0 {
			continue
		}

		htmlPath, err := builtBlogHTMLPath(distRoot, post.Slug)
		if err != nil {
			return manifest{}, err
		}

		astroURLs, err := astroImageURLs(htmlPath)
		if err != nil {
			return manifest{}, err
		}
		if len(localImages) != len(astroURLs) {
			return manifest{}, fmt.Errorf("%s has %d local markdown images but %s has %d _astro webp images", markdownFile, len(localImages), htmlPath, len(astroURLs))
		}

		for i, image := range localImages {
			result.Images[filepath.ToSlash(image.Source)] = absoluteURL(siteURL, astroURLs[i])
		}
	}

	return result, nil
}

func findMarkdownFiles(blogRoot string) ([]string, error) {
	entries, err := os.ReadDir(blogRoot)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		postDir := filepath.Join(blogRoot, entry.Name())
		mdEntries, err := os.ReadDir(postDir)
		if err != nil {
			return nil, err
		}
		for _, mdEntry := range mdEntries {
			if mdEntry.IsDir() || filepath.Ext(mdEntry.Name()) != ".md" {
				continue
			}
			files = append(files, filepath.Join(postDir, mdEntry.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

func postMetadataForFile(markdownFile string) (postMetadata, error) {
	body, err := os.ReadFile(markdownFile)
	if err != nil {
		return postMetadata{}, err
	}

	var metadata postMetadata
	if _, err := frontmatter.Parse(strings.NewReader(string(body)), &metadata); err != nil {
		return postMetadata{}, err
	}

	metadata.Slug = strings.TrimSpace(metadata.Slug)
	if metadata.Slug == "" {
		metadata.Slug = filepath.Base(filepath.Dir(markdownFile))
	}
	return metadata, nil
}

func localMarkdownImages(markdownFile string) ([]localImage, error) {
	body, err := os.ReadFile(markdownFile)
	if err != nil {
		return nil, err
	}

	var images []localImage
	for _, match := range markdownImagePattern.FindAllStringSubmatch(string(body), -1) {
		destination := strings.TrimSpace(match[1])
		pathPart, _ := splitPathSuffix(destination)
		if !isLocalRelativePath(pathPart) {
			continue
		}

		source := filepath.Clean(filepath.Join(filepath.Dir(markdownFile), pathPart))
		images = append(images, localImage{Source: source})
	}

	return images, nil
}

func uniqueLocalImages(images []localImage) []localImage {
	seen := make(map[string]bool, len(images))
	unique := make([]localImage, 0, len(images))
	for _, image := range images {
		key := filepath.ToSlash(image.Source)
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, image)
	}
	return unique
}

func builtBlogHTMLPath(distRoot, slug string) (string, error) {
	candidates := []string{
		filepath.Join(distRoot, "blog", slug, "index.html"),
		filepath.Join(distRoot, "client", "blog", slug, "index.html"),
		filepath.Join(distRoot, "blog", slug+".html"),
		filepath.Join(distRoot, "client", "blog", slug+".html"),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	var found string
	err := filepath.WalkDir(distRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Base(path) != "index.html" {
			return nil
		}

		slashPath := filepath.ToSlash(path)
		if strings.HasSuffix(slashPath, "/blog/"+slug+"/index.html") {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if found != "" {
		return found, nil
	}

	return "", fmt.Errorf("built HTML for blog slug %q not found under %s", slug, distRoot)
}

func astroImageURLs(htmlPath string) ([]string, error) {
	body, err := os.ReadFile(htmlPath)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var urls []string
	for _, match := range htmlImagePattern.FindAllStringSubmatch(string(body), -1) {
		src := html.UnescapeString(strings.TrimSpace(match[1]))
		if !isAstroWebP(src) || seen[src] {
			continue
		}
		seen[src] = true
		urls = append(urls, src)
	}

	return urls, nil
}

func absoluteURL(siteURL, path string) string {
	if parsed, err := url.Parse(path); err == nil && parsed.Scheme != "" {
		return path
	}

	base := strings.TrimRight(siteURL, "/")
	if strings.HasPrefix(path, "/") {
		return base + path
	}
	return base + "/" + path
}

func isAstroWebP(src string) bool {
	parsed, err := url.Parse(src)
	if err != nil {
		return false
	}
	return strings.Contains(parsed.Path, "/_astro/") && strings.EqualFold(filepath.Ext(parsed.Path), ".webp")
}

func isLocalRelativePath(value string) bool {
	if value == "" || strings.HasPrefix(value, "/") {
		return false
	}

	if parsed, err := url.Parse(value); err == nil && parsed.Scheme != "" {
		return false
	}

	return strings.HasPrefix(value, "./") || strings.HasPrefix(value, "../")
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
