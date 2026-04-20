package internal

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/internal/components/blogindex"
	"github.com/yorukot/ssh.yorukot.me/internal/components/endsection"
	"github.com/yorukot/ssh.yorukot.me/internal/components/footer"
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
	m.footer = footer.New(m.innerWidth, m.bg, m.keys)
	m.syncViewport()
}

func (m *Model) syncViewport() {
	contentHeight := m.contentHeight()
	baseWidth := m.contentWidth(false)
	pageKey := m.path
	if m.path == "/blog" {
		pageKey = fmt.Sprintf("%s:%d", m.path, m.selectedBlog)
	}

	needsRender := m.renderedContent == "" ||
		m.renderedBase != baseWidth ||
		m.renderedHeight != contentHeight ||
		m.renderedBg != m.bg ||
		m.renderedPage != pageKey

	contentWidth := m.renderedWidth
	if needsRender {
		renderedContent, metrics := m.renderContent(baseWidth)
		contentWidth = baseWidth
		if lipgloss.Height(renderedContent) > contentHeight {
			contentWidth = m.contentWidth(true)
			renderedContent, metrics = m.renderContent(contentWidth)
		}

		m.blogLineStarts = metrics.LineStarts
		m.blogLineHeights = metrics.LineHeights
		m.renderedContent = renderedContent
		m.renderedWidth = contentWidth
		m.renderedBase = baseWidth
		m.renderedHeight = contentHeight
		m.renderedBg = m.bg
		m.renderedPage = pageKey
	}

	if !m.ready {
		m.main = viewport.New(viewport.WithWidth(contentWidth), viewport.WithHeight(contentHeight))
		m.main.MouseWheelEnabled = false
		m.main.FillHeight = true
		m.main.KeyMap.Up = m.keys.Up
		m.main.KeyMap.Down = m.keys.Down
		m.main.KeyMap.Left.Unbind()
		m.main.KeyMap.Right.Unbind()
		m.ready = true
	} else {
		m.main.SetWidth(contentWidth)
		m.main.SetHeight(contentHeight)
	}

	if m.main.GetContent() != m.renderedContent {
		m.main.SetContent(m.renderedContent)
	}

	m.syncBlogSelectionViewport()
}

func (m *Model) renderContent(width int) (string, blogindex.Metrics) {
	appendEndSection := func(content string) string {
		if m.path == "/" {
			return content
		}

		endContent := endsection.New(width, m.bg).Render()
		return lipgloss.JoinVertical(lipgloss.Left, content, "", endContent)
	}

	if m.path == "/blog" {
		content, metrics := blogindex.Render(m.blogs, width, m.selectedBlog, m.bg)
		return appendEndSection(content), metrics
	}

	page := m.pageMarkdown()
	content, err := m.markdown.RenderWithSource(page.Content, page.SourcePath, width, m.bg)
	if err != nil {
		content = lipgloss.Wrap(page.Content, width, "")
	}

	return appendEndSection(content), blogindex.Metrics{}
}

func (m *Model) syncBlogSelectionViewport() {
	if m.path != "/blog" || len(m.blogLineStarts) == 0 || m.selectedBlog >= len(m.blogLineStarts) {
		return
	}

	start := m.blogLineStarts[m.selectedBlog]
	height := m.blogLineHeights[m.selectedBlog]
	if height <= 0 {
		height = 1
	}
	end := start + height - 1

	top := m.main.YOffset()
	visibleHeight := max(m.main.Height(), 1)
	bottom := top + visibleHeight - 1

	switch {
	case start < top:
		m.main.SetYOffset(start)
	case end > bottom:
		m.main.SetYOffset(end - visibleHeight + 1)
	}
}

func (m *Model) scrollbarView() string {

	h := m.main.Height()
	if h <= 0 {
		return ""
	}

	total := max(m.main.TotalLineCount(), 1)
	visible := max(m.main.VisibleLineCount(), 1)
	if total <= visible {
		return ""
	}

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
