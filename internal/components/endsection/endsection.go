package endsection

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
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
		sections = append(sections, renderLinks(m.content.FooterLinks, m.width, m.bg)...)
	}

	stack := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return lipgloss.NewStyle().Width(max(m.width, 1)).Align(lipgloss.Left).Render(stack)
}

func renderLinks(links []content.FooterLink, width int, bg string) []string {
	lines := make([]string, 0, len(links))

	for _, link := range links {
		if link.Label == "" || link.URL == "" {
			continue
		}

		value := link.Content
		if strings.TrimSpace(value) == "" {
			value = link.URL
		}

		lines = append(lines, renderLine(link.Label, value, link.URL, width, bg))
	}

	return lines
}

func renderLine(label, value, url string, width int, bg string) string {
	labelText := fmt.Sprintf("%-10s", label)
	labelPart := styles.EndSectionLabel(bg).Render(labelText)
	separatorPart := styles.EndSectionSeparator(bg).Render(" : ")
	prefix := labelPart + separatorPart
	prefixWidth := lipgloss.Width(prefix)
	availableValueWidth := max(1, width-prefixWidth)

	linkStyle := styles.EndSectionLink(bg)
	if strings.TrimSpace(url) != "" {
		linkStyle = linkStyle.Hyperlink(url)
	}
	visibleValue := linkStyle.Render(value)

	wrappedValue := lipgloss.Wrap(visibleValue, availableValueWidth, "")
	valueLines := strings.Split(wrappedValue, "\n")
	for i := 1; i < len(valueLines); i++ {
		valueLines[i] = strings.Repeat(" ", prefixWidth) + valueLines[i]
	}

	if len(valueLines) == 0 {
		return prefix
	}
	valueLines[0] = prefix + valueLines[0]
	return strings.Join(valueLines, "\n")
}
