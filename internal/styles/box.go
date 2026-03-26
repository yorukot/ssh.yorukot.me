package styles

import "charm.land/lipgloss/v2"

func BoxHeader(width int) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder())
	return style
}
