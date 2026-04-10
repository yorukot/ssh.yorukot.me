package mkrender

import (
	"github.com/yorukot/ssh.yorukot.me/pkg/stainmd"
)

type Renderer struct{}

func New() Renderer {
	return Renderer{}
}

func (Renderer) Render(markdown string, width int, bg string) (string, error) {
	if width < 1 {
		width = 1
	}

	renderer := stainmd.MochaStyles()
	if bg == "light" {
		renderer = stainmd.LatteStyles()
	}

	return renderer.Render(markdown, width)
}
