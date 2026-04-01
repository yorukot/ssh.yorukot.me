package internal

import (
	"strings"

	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
)

func (m *Model) navigateTo(path string) {
	nextPath := pathutil.NormalizePath(path)
	if nextPath == m.path {
		return
	}

	m.path = nextPath
	m.scrollOffset = constants.MinScrollOffset
	if nextPath == "/blog" {
		m.scrollOffset = max(m.blogIndex-1, constants.MinScrollOffset)
	}
	m.rawMarkdown = ""
	m.renderedBody = ""
	m.wrappedLines = nil
	m.cachedPath = ""
	m.cachedWidth = 0
	m.blogLineOffsets = nil
}

func (m *Model) navigateBack() {
	if m.isBlogIndex() {
		m.navigateTo("/")
		return
	}

	if strings.HasPrefix(m.path, "/blog/") {
		m.navigateTo("/blog")
		return
	}

	m.navigateTo("/")
}
