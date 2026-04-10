package styles

import "charm.land/lipgloss/v2"

func ScrollbarTrack(bg string) lipgloss.Style {
	color := "240"
	if bg == "light" {
		color = "252"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func ScrollbarThumb(bg string) lipgloss.Style {
	color := "252"
	if bg == "light" {
		color = "240"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}

func FooterText(bg string) lipgloss.Style {
	color := "250"
	if bg == "light" {
		color = "241"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func FooterLink(bg string) lipgloss.Style {
	color := "255"
	if bg == "light" {
		color = "238"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}

func FooterSeparator(bg string) lipgloss.Style {
	color := "244"
	if bg == "light" {
		color = "247"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func EndSectionQuote(bg string) lipgloss.Style {
	color := "224"
	if bg == "light" {
		color = "174"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Italic(true)
}

func EndSectionLabel(bg string) lipgloss.Style {
	color := "183"
	if bg == "light" {
		color = "97"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}

func EndSectionLink(bg string) lipgloss.Style {
	color := "117"
	if bg == "light" {
		color = "32"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}

func EndSectionSeparator(bg string) lipgloss.Style {
	color := "240"
	if bg == "light" {
		color = "251"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func HeaderTitle(bg string) lipgloss.Style {
	color := "255"
	if bg == "light" {
		color = "232"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}

func HeaderTagline(bg string) lipgloss.Style {
	color := "245"
	if bg == "light" {
		color = "242"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func HeaderMeta(bg string) lipgloss.Style {
	color := "250"
	if bg == "light" {
		color = "240"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func HeaderRule(bg string) lipgloss.Style {
	color := "238"
	if bg == "light" {
		color = "250"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func BlogIndexTitle(bg string) lipgloss.Style {
	color := "255"
	if bg == "light" {
		color = "232"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}

func BlogIndexIntro(bg string) lipgloss.Style {
	color := "247"
	if bg == "light" {
		color = "242"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func BlogIndexCard(bg string) lipgloss.Style {
	border := "238"
	if bg == "light" {
		border = "252"
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderForeground(lipgloss.Color(border)).
		PaddingLeft(1)
}

func BlogIndexCardActive(bg string) lipgloss.Style {
	border := "117"
	if bg == "light" {
		border = "31"
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeft(true).
		BorderForeground(lipgloss.Color(border)).
		PaddingLeft(1)
}

func BlogIndexMeta(bg string) lipgloss.Style {
	color := "244"
	if bg == "light" {
		color = "240"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func BlogIndexDescription(bg string) lipgloss.Style {
	color := "250"
	if bg == "light" {
		color = "243"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func BlogIndexPath(bg string) lipgloss.Style {
	color := "246"
	if bg == "light" {
		color = "241"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func BlogIndexPathActive(bg string) lipgloss.Style {
	color := "252"
	if bg == "light" {
		color = "236"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}
