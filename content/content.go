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
	PGPKeyID      string
}

type FooterLink struct {
	Label   string
	Content string
	URL     string
}

func GetContent() Content {
	return Content{
		HeaderTitle:   "Yorukot",
		HeaderTagline: "Open-source developer",
		FooterQuote:   "Get busy living, or get busy dying",
		FooterLinks: []FooterLink{
			{Label: "Email", Content: "hi@yorukot.me", URL: "mailto:hi@yorukot.me"},
			{Label: "GitHub", Content: "github.com/yorukot", URL: "https://yorukot.me/github"},
			{Label: "Telegram", Content: "t.me/yorukot", URL: "https://yorukot.me/telegram"},
			{Label: "Discord", Content: "yoru.kot", URL: "https://yorukot.me/discord"},
			{Label: "Ko-fi", Content: "ko-fi.com/yorukot", URL: "https://yorukot.me/sponsor"},
			{Label: "OpenPGP", Content: "F0188B9BF901C94E", URL: "https://yorukot.me/gpg"},
			{Label: "Special", Content: "ssh.yorukot.me repo", URL: "https://github.com/yorukot/ssh.yorukot.me"},
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

func NotFoundPage(path string) string {
	return fmt.Sprintf("# Not Found\n\nNo page found for `%s`.\n\nTry `/` or `/blog`.", path)
}

func ErrorPage(path string, err error) string {
	return fmt.Sprintf("# Error\n\nFailed to load `%s`.\n\n%s", path, err)
}
