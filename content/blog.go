package content

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
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
	Content     string
	RenderedContent string
}

// BlogPosts get all the blog posts and return it
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

		blogPath := filepath.Join(baseDir, entry.Name())
		post, err := loadBlogPost(blogPath, entry.Name())
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

// loadBlogPost load the blog base on the dir and return the BlogPost data
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

	_, post, err := parseBlogFrontMatter(body)
	if err != nil {
		return BlogPost{}, err
	}

	post.Content = string(body)
	post.Path = "/blog/" + slug
	
	return post, nil
}

func FindPost(blogPosts[]BlogPost, slug string) (BlogPost, error) {
	for _, post := range blogPosts {
		if post.Path == slug {
			return post, nil
		}
	}
	return BlogPost{}, errors.New("blog post not found")
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