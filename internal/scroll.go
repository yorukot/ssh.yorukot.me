package internal

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
)

func (m *Model) contentViewportHeight() int {
	headerHeight := lipgloss.Height(header.New(m.innerWidth, m.innerHeight).Render())
	footerHeight := lipgloss.Height(m.renderHelpFooter(max(m.innerWidth-5, 20)))
	return max(m.innerHeight-headerHeight-footerHeight-5, 3)
}

func (m *Model) renderScrollableContent(content string, width int, headerContent string) string {
	contentStyle := lipgloss.NewStyle().Width(width)
	wrapped := contentStyle.Render(content)
	lines := strings.Split(wrapped, "\n")
	footerHeight := lipgloss.Height(m.renderHelpFooter(width))
	viewportHeight := max(m.innerHeight-lipgloss.Height(headerContent)-footerHeight-5, 3)
	m.scrollOffset = clampScrollOffset(m.scrollOffset, len(lines), viewportHeight)

	end := min(m.scrollOffset+viewportHeight, len(lines))
	visible := lines[m.scrollOffset:end]
	for len(visible) < viewportHeight {
		visible = append(visible, "")
	}

	contentPane := lipgloss.NewStyle().Width(width).Height(viewportHeight).Render(strings.Join(visible, "\n"))
	scrollbar := m.renderScrollbar(len(lines), viewportHeight)
	return lipgloss.JoinHorizontal(lipgloss.Top, contentPane, " ", scrollbar)
}

func clampScrollOffset(offset, totalLines, viewportHeight int) int {
	maxOffset := max(totalLines-viewportHeight, 0)
	return min(max(offset, 0), maxOffset)
}

func (m Model) renderScrollbar(totalLines, viewportHeight int) string {
	trackColor := "240"
	thumbColor := "252"
	if m.bg == "light" {
		trackColor = "252"
		thumbColor = "240"
	}

	track := lipgloss.NewStyle().Foreground(lipgloss.Color(trackColor))
	thumb := lipgloss.NewStyle().Foreground(lipgloss.Color(thumbColor)).Bold(true)
	if totalLines <= viewportHeight || viewportHeight <= 0 {
		return lipgloss.NewStyle().Height(max(viewportHeight, 1)).Render(track.Render("│"))
	}

	thumbHeight := max(viewportHeight*viewportHeight/max(totalLines, 1), 1)
	maxOffset := max(totalLines-viewportHeight, 1)
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
