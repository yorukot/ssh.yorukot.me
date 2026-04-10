package blogindex

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
)

func HandleKey(msg tea.KeyMsg, keys keymap.Bindings, posts []content.BlogPost, selected int) (int, string, bool) {
	if len(posts) == 0 {
		return selected, "", false
	}

	switch {
	case key.Matches(msg, keys.Up):
		return min(max(selected-1, 0), len(posts)-1), "", true
	case key.Matches(msg, keys.Down):
		return min(max(selected+1, 0), len(posts)-1), "", true
	case key.Matches(msg, keys.Enter):
		return selected, posts[selected].Path, true
	default:
		return selected, "", false
	}
}
