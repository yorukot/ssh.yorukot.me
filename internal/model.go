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

	headerContent string
	headerHeight  int
	footerContent string
	footerHeight  int

	renderedContent string
	renderedWidth   int
	renderedBase    int
	renderedHeight  int
	renderedBg      string
	renderedPage    string

	blogLineStarts  []int
	blogLineHeights []int

	scrollbarContent     string
	scrollbarHeight      int
	scrollbarTotal       int
	scrollbarVisible     int
	scrollbarYOffset     int
	scrollbarThumbHeight int
	scrollbarTop         int
	scrollbarBg          string
}

func (m *Model) contentWidth(hasScrollbar bool) int {
	availableWidth := m.innerWidth - styles.InnerBoxPaddingSide*2
	if hasScrollbar {
		availableWidth -= constants.ScrollbarGap + constants.ScrollbarWidth
	}

	return max(availableWidth, constants.MinContentWidth)
}

func (m *Model) contentHeight() int {
	return max(1, m.innerHeight-styles.InnerBoxPaddingTop*2-m.headerHeight-m.footerHeight)
}

func (m *Model) refreshChrome() {
	m.headerContent = m.header.Render()
	m.headerHeight = lipgloss.Height(m.headerContent)
	m.footerContent = m.footer.Render()
	m.footerHeight = lipgloss.Height(m.footerContent)
	m.scrollbarContent = ""
}
