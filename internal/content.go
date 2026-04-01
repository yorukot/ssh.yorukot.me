package internal

import (
	"errors"
	"strings"

	"github.com/yorukot/ssh.yorukot.me/content"
)

func (m *Model) renderContent() (string, error) {
	if m.path == "/blog" {
	}
	if strings.HasPrefix(m.path, "/blog") {
		post, err := content.FindPost(m.blogs, m.path)
		if err != nil {
			return "", err
		}
		return post.Content, nil
	}
	
	return "", errors.New("content not found")
}