package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Content struct {
	HeaderTitle   string
	HeaderTagline string
	FooterQuote   string
	FooterLinks   []FooterLink
}

type FooterLink struct {
	Label string
	URL   string
}

func GetContent() Content {
	return Content{
		HeaderTitle:   "Yorukot",
		HeaderTagline: "Open-source developer",
		FooterQuote:   "Get busy living, or get busy dying",
		FooterLinks: []FooterLink{
			{Label: "Email", URL: "mailto:hi@yorukot.me"},
			{Label: "GitHub", URL: "https://github.com/yorukot"},
			{Label: "Telegram", URL: "https://t.me/yorukot"},
			{Label: "Discord", URL: "https://dc.yorukot.me"},
			{Label: "Ko-fi", URL: "https://donate.yorukot.me"},
		},
	}
}

func HomePage() (string, error) {
	body, err := os.ReadFile(filepath.Join("content", "markdown", "intro.md"))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

// TODO: This sucks we need to change this
func BlogIndex(posts []BlogPost) string {
	if len(posts) == 0 {
		return "# Blog\n\nNo posts yet."
	}

	var builder strings.Builder
	builder.WriteString("# Blog\n\n")
	builder.WriteString("Posts you can open from this SSH site.\n\n")

	for _, post := range posts {
		builder.WriteString("## [")
		builder.WriteString(post.Title)
		builder.WriteString("](")
		builder.WriteString(post.Path)
		builder.WriteString(")\n\n")

		if post.PublishDate != "" {
			builder.WriteString("Published: ")
			builder.WriteString(post.PublishDate)
			builder.WriteString("\n\n")
		}

		if post.Description != "" {
			builder.WriteString(post.Description)
			builder.WriteString("\n\n")
		}

		builder.WriteString("Path: `")
		builder.WriteString(post.Path)
		builder.WriteString("`\n\n")
	}

	return strings.TrimSpace(builder.String())
}

func NotFoundPage(path string) string {
	return fmt.Sprintf("# Not Found\n\nNo page found for `%s`.\n\nTry `/` or `/blog`.", path)
}

func ErrorPage(path string, err error) string {
	return fmt.Sprintf("# Error\n\nFailed to load `%s`.\n\n%s", path, err)
}
