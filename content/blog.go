package content

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
)

type BlogPost struct {
	Title       string   `yaml:"title"`
	Authors     []string `yaml:"authors"`
	Tags        []string `yaml:"tags"`
	Categories  []string `yaml:"categories"`
	PublishDate string   `yaml:"date"`
	UpdatedDate string   `yaml:"updated_date"`
	Description string   `yaml:"description"`
	Path        string
}

func MarkdownContent(path string) (string, error) {
	normalizedPath := pathutil.NormalizePath(path)

	baseDir := filepath.Join("content", "markdown")
	cleanPath := strings.Trim(normalizedPath, "/")

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

func BlogPosts() ([]BlogPost, error) {
	baseDir := filepath.Join("content", "markdown", "blog")
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	posts := make([]BlogPost, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		post, err := loadBlogPost(filepath.Join(baseDir, entry.Name()), entry.Name())
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		leftDate, leftOK := parseBlogDate(posts[i].PublishDate)
		rightDate, rightOK := parseBlogDate(posts[j].PublishDate)

		switch {
		case leftOK && rightOK && !leftDate.Equal(rightDate):
			return leftDate.After(rightDate)
		case leftOK != rightOK:
			return leftOK
		default:
			return posts[i].Path < posts[j].Path
		}
	})

	return posts, nil
}

func loadBlogPost(dir, slug string) (BlogPost, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return BlogPost{}, err
	}

	var mdFiles []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		mdFiles = append(mdFiles, entry.Name())
	}

	if len(mdFiles) == 0 {
		return BlogPost{}, errors.New("no markdown file found")
	}

	sort.Strings(mdFiles)
	body, err := os.ReadFile(filepath.Join(dir, mdFiles[0]))
	if err != nil {
		return BlogPost{}, err
	}

	post := BlogPost{
		Title: slug,
		Path:  "/blog/" + slug,
	}

	_, metadata, err := parseBlogFrontMatter(body)
	if err != nil {
		return BlogPost{}, err
	}

	if strings.TrimSpace(metadata.Title) == "" {
		return BlogPost{}, errors.New("missing required front matter field: title")
	}

	post.Title = metadata.Title
	post.Description = metadata.Description
	post.PublishDate = metadata.PublishDate

	if post.Description == "" {
		post.Description = "A new post in the blog archive."
	}

	return post, nil
}

func parseBlogFrontMatter(body []byte) ([]byte, BlogPost, error) {
	var metadata BlogPost
	rest, err := frontmatter.Parse(strings.NewReader(string(body)), &metadata)
	if err != nil {
		return nil, BlogPost{}, err
	}
	return rest, metadata, nil
}

func parseBlogDate(value string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
