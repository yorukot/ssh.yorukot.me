package header

import (
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

type Header struct {
	title string
	width int
	height int
}

func New(w, h int) Header {
	data := content.GetContent()

	return Header{
		title: data.HeaderTitle,
		width: w - 2,
		height: h - 2,
	}
}

func (h Header) Render() string {
	return styles.HeaderBox(h.width).Align(lipgloss.Left).Render(h.title)
}
