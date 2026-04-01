package header

import (
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"strings"
)

type Model struct {
	title   string
	tagline string
	path    string
	width   int
	bg      string
}

func New(w int, bg, path string) Model {
	data := content.GetContent()
	
	return Model{
		title:   data.HeaderTitle,
		tagline: data.HeaderTagline,
		path:    path,
		// We need to * 2 for the both side
		width:   w - styles.InnerBoxPaddingSide * 2, 
		bg:      bg,
	}
}

func (h Model) Render() string {
	route := h.path
	if strings.TrimSpace(route) == "" {
		route = "/"
	}

	title := styles.HeaderTitle(h.bg).Render(h.title)
	tagline := styles.HeaderTagline(h.bg).Render(h.tagline)
	meta := styles.HeaderMeta(h.bg).Render(route)
	rule := styles.HeaderRule(h.bg).Render(strings.Repeat("─", max(h.width, 1)))

	stack := lipgloss.JoinVertical(lipgloss.Left, title, tagline, "", meta, rule)
	return lipgloss.NewStyle().Width(h.width).Align(lipgloss.Left).Render(stack)
}
