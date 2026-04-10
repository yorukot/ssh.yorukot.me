package footer

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

type Model struct {
	keys  keymap.Bindings
	width int
	bg    string
}

func New(w int, bg string, keys keymap.Bindings) Model {
	return Model{
		keys:  keys,
		width: w - styles.InnerBoxPaddingSide*2,
		bg:    bg,
	}
}

func (f Model) Render() string {
	rule := styles.FooterSeparator(f.bg).Render(strings.Repeat("─", max(f.width, 1)))
	help := renderHelp(f.keys.ShortHelp(), f.bg)
	stack := lipgloss.JoinVertical(lipgloss.Left, rule, help)

	return lipgloss.NewStyle().Width(max(f.width, 1)).Align(lipgloss.Left).Render(stack)
}

func renderHelp(bindings []key.Binding, bg string) string {
	parts := make([]string, 0, len(bindings))
	separator := styles.FooterSeparator(bg).Render(" • ")

	for _, binding := range bindings {
		help := binding.Help()
		if help.Key == "" && help.Desc == "" {
			continue
		}

		parts = append(parts, lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.FooterLink(bg).Render(help.Key),
			styles.FooterText(bg).Render(" "+help.Desc),
		))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, joinWithSeparator(parts, separator)...)
}

func joinWithSeparator(items []string, separator string) []string {
	if len(items) == 0 {
		return nil
	}

	parts := make([]string, 0, len(items)*2-1)
	for i, item := range items {
		if i > 0 {
			parts = append(parts, separator)
		}
		parts = append(parts, item)
	}

	return parts
}
