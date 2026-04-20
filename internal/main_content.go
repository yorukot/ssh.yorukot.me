package internal

import (
	"strings"

	"github.com/yorukot/ssh.yorukot.me/content"
)

type pageMarkdown struct {
	Content    string
	SourcePath string
}

func (m *Model) isBlogPost() bool {
	return strings.HasPrefix(m.path, "/blog/")
}

func (m *Model) pageMarkdown() pageMarkdown {
	switch {
	case m.path == "/":
		home, err := content.HomePage()
		if err != nil {
			return pageMarkdown{Content: content.ErrorPage(m.path, err)}
		}
		return pageMarkdown{
			Content:    home,
			SourcePath: "content/markdown/intro.md",
		}
	case m.isBlogPost():
		post, err := content.FindPost(m.blogs, m.path)
		if err != nil {
			return pageMarkdown{Content: content.NotFoundPage(m.path)}
		}
		return pageMarkdown{
			Content:    post.Content,
			SourcePath: post.SourcePath,
		}
	default:
		return pageMarkdown{Content: content.NotFoundPage(m.path)}
	}
}
