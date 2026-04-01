package internal

import (
	"charm.land/bubbles/v2/help"
	contentpkg "github.com/yorukot/ssh.yorukot.me/content"
)

type Model struct {
	width        int
	height       int
	innerWidth   int
	innerHeight  int
	scrollOffset int
	contentWidth int
	help         help.Model
	keys         keyMap
	profile      string
	bg           string
	path         string
	rawMarkdown  string
	renderedBody string
	wrappedLines []string

	footerQuoteIndex    int
	footerQuoteDeleting bool
	footerQuotePause    int
	footerCursorVisible bool
	blogPosts           []contentpkg.BlogPost
	blogIndex           int
	blogLineOffsets     []int

	// cached relate
	cachedPath  string
	cachedBg    string
	cachedWidth int
}
