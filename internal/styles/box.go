package styles

import (
	"charm.land/lipgloss/v2"
)

func HeaderBox(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder())
}

func FullScreenBox(w, h int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Padding(1, 1)
}