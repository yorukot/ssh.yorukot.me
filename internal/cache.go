package internal

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	contentpkg "github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/footer"
	markdownpkg "github.com/yorukot/ssh.yorukot.me/pkg/markdown"
)

func (m *Model) cached(width int) {
	if width <= 0 {
		m.wrappedLines = []string{""}
		return
	}

	if m.cachedPath != m.path || m.rawMarkdown == "" {
		if m.isBlogIndex() {
			m.rawMarkdown = ""
			m.renderedBody = ""
			m.cachedPath = m.path
			m.wrappedLines = nil
		} else {
			markdownContent, err := contentpkg.MarkdownContent(m.path)
			if err != nil {
				markdownContent = fmt.Sprintf("\n\nfailed to load markdown for %s\n\n%s", m.path, err)
			}

			m.rawMarkdown = markdownContent
			m.cachedPath = m.path
			m.renderedBody = ""
			m.wrappedLines = nil
		}
	}

	if !m.isBlogIndex() && (m.cachedWidth != width || m.cachedBg != m.bg || m.renderedBody == "") {
		m.renderedBody = markdownpkg.New(width, m.bg).Render(m.rawMarkdown)
		m.cachedWidth = width
		m.cachedBg = m.bg
		m.wrappedLines = nil
	}

	if m.wrappedLines != nil {
		return
	}

	body := strings.TrimSpace(m.renderedBody)
	if m.isBlogIndex() {
		body = strings.TrimSpace(m.renderBlogIndex(width))
		m.syncBlogScrollOffset()
	}
	footerView := footer.New(width, m.bg, m.footerQuoteText(), m.footerCursorVisible).Render()
	sections := []string{body}
	if strings.TrimSpace(footerView) != "" {
		sections = append(sections, footerView)
	}

	wrapped := lipgloss.NewStyle().Width(width).Render(strings.Join(sections, "\n\n"))
	m.wrappedLines = strings.Split(wrapped, "\n")
}
