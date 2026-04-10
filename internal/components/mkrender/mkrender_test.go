package mkrender

import (
	"strings"
	"testing"
)

func TestRenderUsesStainmd(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"# Hello",
		"",
		"{{notice}}",
		"Read this",
		"{{noticed}}",
		"",
		"Use ~~old~~ `new`.",
	}, "\n")

	out, err := renderer.Render(input, 40, "dark")
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	checks := []string{"Hello", "Read this", "old", "new"}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("rendered output missing %q\noutput:\n%s", want, out)
		}
	}

	if strings.Contains(out, "{{notice}}") || strings.Contains(out, "{{noticed}}") {
		t.Fatalf("expected rendered output to hide notice markers\noutput:\n%s", out)
	}
}
