package internal

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

func (m *Model) contentViewportHeight() int {
	headerHeight := lipgloss.Height(header.New(m.innerWidth, m.bg, m.path).Render())
	helpHeight := lipgloss.Height(m.renderHelpFooter(max(m.innerWidth-constants.ContentWidthInset, constants.MinContentWidth)))
	return max(m.innerHeight-headerHeight-helpHeight-constants.LayoutVerticalSpacing, constants.MinViewportHeight)
}

func (m *Model) renderScrollableContent(width int, headerContent, helpView string) string {
	lines := m.wrappedLines
	footerHeight := lipgloss.Height(helpView)
	viewportHeight := max(m.innerHeight-lipgloss.Height(headerContent)-footerHeight-constants.LayoutVerticalSpacing, constants.MinViewportHeight)
	m.scrollOffset = clampScrollOffset(m.scrollOffset, len(lines), viewportHeight)

	end := min(m.scrollOffset+viewportHeight, len(lines))
	visible := lines[m.scrollOffset:end]
	for len(visible) < viewportHeight {
		visible = append(visible, "")
	}

	contentPane := lipgloss.NewStyle().Width(width).Height(viewportHeight).Render(strings.Join(visible, "\n"))
	scrollbar := m.renderScrollbar(len(lines), viewportHeight)
	if scrollbar == "" {
		return contentPane
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, contentPane, " ", scrollbar)
}

func clampScrollOffset(offset, totalLines, viewportHeight int) int {
	maxOffset := max(totalLines-viewportHeight, constants.MinScrollOffset)
	return min(max(offset, constants.MinScrollOffset), maxOffset)
}

func (m Model) renderScrollbar(totalLines, viewportHeight int) string {
	track := styles.ScrollbarTrack(m.bg)
	thumb := styles.ScrollbarThumb(m.bg)
	if totalLines <= viewportHeight || viewportHeight <= constants.MinScrollOffset {
		return ""
	}

	thumbHeight := max(viewportHeight*viewportHeight/max(totalLines, constants.MinScrollbarHeight), constants.MinScrollbarHeight)
	maxOffset := max(totalLines-viewportHeight, constants.MinScrollbarOffset)
	scrollOffset := clampScrollOffset(m.scrollOffset, totalLines, viewportHeight)
	thumbTop := (viewportHeight - thumbHeight) * scrollOffset / maxOffset

	lines := make([]string, viewportHeight)
	for i := range lines {
		lines[i] = track.Render("│")
		if i >= thumbTop && i < thumbTop+thumbHeight {
			lines[i] = thumb.Render("█")
		}
	}

	return strings.Join(lines, "\n")
}
