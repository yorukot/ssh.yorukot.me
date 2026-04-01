package header

import (
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"strings"
)

type Header struct {
	title   string
	tagline string
	path    string
	width   int
	bg      string
}

func New(w int, bg, path string) Header {
	data := content.GetContent()

	return Header{
		title:   data.HeaderTitle,
		tagline: data.HeaderTagline,
		path:    path,
		width:   w - constants.HeaderFrameInset,
		bg:      bg,
	}
}

func (h Header) Render() string {
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
