package endsection

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"github.com/yorukot/ssh.yorukot.me/pkg/stainmd"
)

type Model struct {
	content content.Content
	width   int
	bg      string
}

func New(w int, bg string) Model {
	return Model{
		content: content.GetContent(),
		width:   w,
		bg:      bg,
	}
}

func (m Model) Render() string {
	sections := []string{
		styles.EndSectionSeparator(m.bg).Render(strings.Repeat("─", max(m.width, 1))),
	}

	if m.content.FooterQuote != "" {
		sections = append(sections, styles.EndSectionQuote(m.bg).Render(m.content.FooterQuote))
	}

	if len(m.content.FooterLinks) > 0 {
		sections = append(sections, renderLinks(m.content.FooterLinks, m.bg)...)
	}

	stack := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return lipgloss.NewStyle().Width(max(m.width, 1)).Align(lipgloss.Left).Render(stack)
}

func renderLinks(links []content.FooterLink, bg string) []string {
	lines := make([]string, 0, len(links))

	for _, link := range links {
		if link.Label == "" || link.URL == "" {
			continue
		}

		value := link.Content
		if strings.TrimSpace(value) == "" {
			value = link.URL
		}

		lines = append(lines, renderLine(link.Label, value, link.URL, bg))
	}

	return lines
}

func renderLine(label, value, url, bg string) string {
	visibleValue := styles.EndSectionLink(bg).Render(value)
	if strings.TrimSpace(url) != "" {
		visibleValue = stainmd.OSC8Link(url, visibleValue)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.EndSectionLabel(bg).Render(fmt.Sprintf("%-10s", label)),
		styles.EndSectionSeparator(bg).Render(" : "),
		visibleValue,
	)
}
