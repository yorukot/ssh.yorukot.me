package internal

import (
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
)

func (m *Model) goToPath(path string) {
	m.path = pathutil.NormalizePath(path)
	m.header = header.New(m.innerWidth, m.bg, m.path)
	m.refreshChrome()
	m.syncViewport()
	m.main.GotoTop()
}
