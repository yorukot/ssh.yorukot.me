package internal

import (
	"strings"
)

func (m *Model) isBlogPost() bool {
	return strings.HasPrefix(m.path, "/blog/")
}
