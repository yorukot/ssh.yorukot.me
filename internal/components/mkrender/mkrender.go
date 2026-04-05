package mkrender

import "github.com/charmbracelet/glamour"

type Renderer struct{}

func New() Renderer {
	return Renderer{}
}

func (Renderer) Render(markdown string, width int, bg string) (string, error) {
	if width < 1 {
		width = 1
	}

	style := glamour.DarkStyle
	if bg == "light" {
		style = glamour.LightStyle
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(width),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		return "", err
	}

	return renderer.Render(markdown)
}
