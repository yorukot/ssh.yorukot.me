package footer

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"github.com/yorukot/ssh.yorukot.me/pkg/ansi"
)

type Footer struct {
	links         []content.FooterLink
	width         int
	bg            string
	quote         string
	cursorVisible bool
}

func New(w int, bg, quote string, cursorVisible bool) Footer {
	data := content.GetContent()

	return Footer{
		links:         data.FooterLinks,
		width:         w,
		bg:            bg,
		quote:         quote,
		cursorVisible: cursorVisible,
	}
}

func (f Footer) Render() string {
	if len(f.links) == 0 {
		return ""
	}

	widthLimit := max(f.width, constants.FooterMinWrapWidth)

	separator := styles.FooterSeparator(f.bg).Render(" • ")
	rows := make([]string, 0, len(f.links))
	current := make([]string, 0, len(f.links))
	currentWidth := 0
	separatorWidth := lipgloss.Width(" • ")

	for _, link := range f.links {
		part := ansi.OSC8Link(link.URL, styles.FooterLink(f.bg).Render(link.Label))
		partWidth := lipgloss.Width(link.Label)
		nextWidth := currentWidth + partWidth
		if len(current) > 0 {
			nextWidth += separatorWidth
		}

		if len(current) > 0 && nextWidth > widthLimit {
			rows = append(rows, strings.Join(current, separator))
			current = []string{part}
			currentWidth = partWidth
			continue
		}

		current = append(current, part)
		currentWidth = nextWidth
	}

	if len(current) > 0 {
		rows = append(rows, strings.Join(current, separator))
	}

	quote := styles.FooterText(f.bg).Render(f.quote)
	if f.cursorVisible {
		quote += styles.FooterText(f.bg).Render("_")
	}

	year := time.Now().Year()
	
	meta := styles.FooterText(f.bg).Render(fmt.Sprintf("Code samples are under the MIT License. © Copyright 2023-%d Yorukot", year))

	content := lipgloss.JoinVertical(lipgloss.Center, quote, strings.Join(rows, "\n"), meta)
	return lipgloss.NewStyle().Width(f.width).Align(lipgloss.Center).Render(content)
}
