package internal

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

func (m *Model) windowsSizeChange(msg tea.WindowSizeMsg) {
	m.height = msg.Height
	m.width = msg.Width
	m.innerHeight = m.height
	m.innerWidth = min(m.width, constants.MaxContentWidth)

	// Model refresh
	m.header = header.New(m.innerWidth, m.bg, m.path)
	m.syncViewport()
}

func (m *Model) syncViewport() {
	// TODO: we need to update this since we are not going to only use scroll bar in the blog content
	if !m.isBlogPost() {
		return
	}

	post, err := content.FindPost(m.blogs, m.path)
	if err != nil {
		return
	}

	contentWidth := m.contentWidth()
	contentHeight := m.contentHeight()
	renderedContent, err := m.markdown.Render(post.Content, contentWidth, m.bg)
	if err != nil {
		renderedContent = lipgloss.Wrap(post.Content, contentWidth, "")
	}

	if !m.ready {
		m.main = viewport.New(viewport.WithWidth(contentWidth), viewport.WithHeight(contentHeight))
		m.main.MouseWheelEnabled = true
		m.main.FillHeight = true
		m.main.SoftWrap = false
		m.main.KeyMap.Up = m.keys.Up
		m.main.KeyMap.Down = m.keys.Down
		m.ready = true
	} else {
		m.main.SetWidth(contentWidth)
		m.main.SetHeight(contentHeight)
	}

	if m.main.GetContent() != renderedContent {
		m.main.SetContent(renderedContent)
	}
}

func (m *Model) scrollbarView() string {
	h := m.main.Height()
	if h <= 0 {
		return ""
	}

	total := max(m.main.TotalLineCount(), 1)
	visible := max(m.main.VisibleLineCount(), 1)
	thumbHeight := max(1, visible*h/total)
	maxTop := max(0, h-thumbHeight)
	top := int(m.main.ScrollPercent() * float64(maxTop))

	lines := make([]string, 0, h)
	for i := range h {
		if i >= top && i < top+thumbHeight {
			lines = append(lines, styles.ScrollbarThumb(m.bg).Render("█"))
			continue
		}
		lines = append(lines, styles.ScrollbarTrack(m.bg).Render("│"))
	}

	return strings.Join(lines, "\n")
}
