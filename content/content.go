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
			{Label: "Email", Content: "hi@yorukot.me", URL: "https://yorukot.me/mail"},
			{Label: "GitHub", URL: "https://yorukot.me/gh"},
			{Label: "Telegram", URL: "https://yorukot.me/tg"},
			{Label: "Discord", URL: "https://yorukot.me/dc"},
			{Label: "Ko-fi", URL: "https://yorukot.me/donate"},
			{Label: "OpenPGP", Content: "F0188B9BF901C94E", URL: "https://yorukot.me/gpg"},
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
