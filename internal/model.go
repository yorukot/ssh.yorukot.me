package internal

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/footer"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
	"github.com/yorukot/ssh.yorukot.me/internal/mkrender"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"time"
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
	renderedBase    int
	renderedHeight  int
	renderedBg      string
	renderedPage    string

	blogLineStarts  []int
	blogLineHeights []int

	lastWheelButton tea.MouseButton
	wheelBurstLines int
	lastWheelAt     time.Time
}

func (m *Model) contentWidth(hasScrollbar bool) int {
	availableWidth := m.innerWidth - styles.InnerBoxPaddingSide*2
	if hasScrollbar {
		availableWidth -= constants.ScrollbarGap + constants.ScrollbarWidth
	}

	return max(availableWidth, constants.MinContentWidth)
}

func (m *Model) contentHeight() int {
	headerHeight := lipgloss.Height(m.header.Render())
	footerHeight := lipgloss.Height(m.footer.Render())
	return max(1, m.innerHeight-styles.InnerBoxPaddingTop*2-headerHeight-footerHeight)
}
