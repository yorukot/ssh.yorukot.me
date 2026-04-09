package stainmd

import "charm.land/lipgloss/v2"

type HeaderStyle struct {
	HeadingOne   lipgloss.Style
	HeadingTwo   lipgloss.Style
	HeadingThree lipgloss.Style
	HeadingFour  lipgloss.Style
	HeadingFive  lipgloss.Style
	HeadingSix   lipgloss.Style
}

func (r Renderer) headingStyle(level int) lipgloss.Style {
	switch level {
	case 1:
		return r.Header.HeadingOne
	case 2:
		return r.Header.HeadingTwo
	case 3:
		return r.Header.HeadingThree
	case 4:
		return r.Header.HeadingFour
	case 5:
		return r.Header.HeadingFive
	default:
		return r.Header.HeadingSix
	}
}
