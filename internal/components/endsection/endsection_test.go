package endsection

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/yorukot/ssh.yorukot.me/content"
)

func TestRenderWrapsLongFooterLinkValues(t *testing.T) {
	m := Model{
		content: content.Content{
			FooterLinks: []content.FooterLink{{
				Label:   "GitHub",
				Content: "https://github.com/yorukot/this-is-a-very-long-link-value-that-should-wrap-cleanly",
				URL:     "https://github.com/yorukot/this-is-a-very-long-link-value-that-should-wrap-cleanly",
			}},
		},
		width: 28,
	}

	out := m.Render()
	plain := xansi.Strip(out)

	if strings.Count(plain, "\n") < 2 {
		t.Fatalf("expected wrapped footer output to span multiple lines, got:\n%s", out)
	}

	for _, line := range strings.Split(plain, "\n") {
		if xansi.StringWidth(line) > 28 {
			t.Fatalf("expected footer line width <= 28, got %d in line %q\nfull output:\n%s", xansi.StringWidth(line), line, out)
		}
	}
	if !strings.Contains(plain, "GitHub") {
		t.Fatalf("expected footer label to remain visible, got:\n%s", out)
	}
}
