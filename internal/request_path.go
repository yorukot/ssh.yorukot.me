package internal

import (
	"strings"

	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
)

func resolveRequestPath(input string, blogPosts []content.BlogPost) string {
	requestPath := pathutil.NormalizePath(input)
	if requestPath == "/" || requestPath == "/blog" || strings.HasPrefix(requestPath, "/blog/") {
		return requestPath
	}

	slug := strings.Trim(requestPath, "/")
	if strings.Contains(slug, "/") {
		return requestPath
	}

	if post, err := content.FindPost(blogPosts, "/blog/"+slug); err == nil {
		return post.Path
	}

	return requestPath
}
