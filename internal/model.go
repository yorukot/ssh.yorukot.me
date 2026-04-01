package internal

import (
	"charm.land/bubbles/v2/viewport"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
)

type Model struct {
	width        int
	height       int
	innerWidth   int
	innerHeight  int
	scrollOffset int
	keys         keymap.Bindings

	path string

	bg string

	colorProfile string

	blogs []content.BlogPost

	main   viewport.Model
	header header.Model
}
