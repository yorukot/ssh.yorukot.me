package internal

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/footer"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
	"github.com/yorukot/ssh.yorukot.me/internal/mkrender"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

type Model struct {
	width       int
	height      int
	innerWidth  int
	innerHeight int
	keys        keymap.Bindings

	path         string
	selectedBlog int

	bg string

	colorProfile string

	blogs []content.BlogPost

	ready    bool
	main     viewport.Model
	header   header.Model
	footer   footer.Model
	markdown mkrender.Renderer

	renderedContent string
	renderedWidth   int
	renderedBg      string
	renderedPage    string

	blogLineStarts  []int
	blogLineHeights []int
}

func (m *Model) contentWidth() int {
	availableWidth := m.innerWidth - styles.InnerBoxPaddingSide*2
	if m.hasScrollbar() {
		availableWidth -= constants.ScrollbarGap + constants.ScrollbarWidth
	}

	return max(availableWidth, constants.MinContentWidth)
}

func (m *Model) contentHeight() int {
	headerHeight := lipgloss.Height(m.header.Render())
	footerHeight := lipgloss.Height(m.footer.Render())
	return max(1, m.innerHeight-styles.InnerBoxPaddingTop*2-headerHeight-footerHeight)
}

func (m *Model) hasScrollbar() bool {
	total := max(m.main.TotalLineCount(), 1)
	visible := max(m.main.VisibleLineCount(), 1)
	return total > visible
}
