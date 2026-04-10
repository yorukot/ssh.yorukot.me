package internal

import (
	"strings"

	"github.com/yorukot/ssh.yorukot.me/content"
)

func (m *Model) isBlogPost() bool {
	return strings.HasPrefix(m.path, "/blog/")
}

func (m *Model) pageContent() string {
	switch {
	case m.path == "/":
		home, err := content.HomePage()
		if err != nil {
			return content.ErrorPage(m.path, err)
		}
		return home
	case m.isBlogPost():
		post, err := content.FindPost(m.blogs, m.path)
		if err != nil {
			return content.NotFoundPage(m.path)
		}
		return post.Content
	default:
		return content.NotFoundPage(m.path)
	}
}
