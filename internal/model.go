package internal

import "charm.land/bubbles/v2/help"

type Model struct {
	width        int
	height       int
	innerWidth   int
	innerHeight  int
	scrollOffset int
	help         help.Model
	keys         keyMap
	profile      string
	bg           string
	path         string
}
