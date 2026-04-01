package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
)

func HeaderBox(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Padding(constants.HeaderBoxPaddingTop, constants.HeaderBoxPaddingSide).
		BorderStyle(lipgloss.RoundedBorder())
}

func InnerBox(w, h int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Padding(constants.InnerBoxPaddingTop, constants.InnerBoxPaddingSide)
}
