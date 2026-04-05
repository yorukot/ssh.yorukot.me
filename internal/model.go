package internal

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/components/mkrender"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

type Model struct {
	width       int
	height      int
	innerWidth  int
	innerHeight int
	keys        keymap.Bindings

	path string

	bg string

	colorProfile string

	blogs []content.BlogPost

	ready    bool
	main     viewport.Model
	header   header.Model
	markdown mkrender.Renderer
}

func (m *Model) contentWidth() int {
	availableWidth := m.innerWidth - styles.InnerBoxPaddingSide*2
	if m.isBlogPost() {
		availableWidth -= constants.ScrollbarGap + constants.ScrollbarWidth
	}

	return max(availableWidth, constants.MinContentWidth)
}

func (m *Model) contentHeight() int {
	headerHeight := lipgloss.Height(m.header.Render())
	return max(1, m.innerHeight-styles.InnerBoxPaddingTop*2-headerHeight)
}
