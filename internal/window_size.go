package internal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
)

func (m *Model) windowsSizeChange(msg tea.WindowSizeMsg) {
	m.height = msg.Height
	m.width = msg.Width
	m.innerHeight = m.height
	m.innerWidth = min(m.width, constants.MaxContentWidth)
	
	// Model refresh
	m.header = header.New(m.innerWidth, m.bg, m.path)
}
