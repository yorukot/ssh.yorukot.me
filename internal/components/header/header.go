package header

import (
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

type Header struct {
	Title string
	Width int
}

func New(w int) Header {
	data := content.GetContent()

	return Header{
		Title: data.HeaderTitle,
		Width: w,
	}
}

func (h Header) Render() string {
	return styles.HeaderBox(h.Width).Align(lipgloss.Left).Render(h.Title)
}
