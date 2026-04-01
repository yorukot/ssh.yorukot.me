package styles

import (
	"charm.land/lipgloss/v2"
)

const (
	InnerBoxPaddingTop int = 1
	InnerBoxPaddingSide int = 1
)

func InnerBox(w, h int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Padding(InnerBoxPaddingTop, InnerBoxPaddingSide)
}
