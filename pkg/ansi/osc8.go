package ansi

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
)

func OSC8Link(url, label string) string {
	if strings.TrimSpace(url) == "" {
		return label
	}

	return xansi.SetHyperlink(url) + label + xansi.ResetHyperlink()
}
