package styles

import "charm.land/lipgloss/v2"

func QuitText() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#949494"))
}