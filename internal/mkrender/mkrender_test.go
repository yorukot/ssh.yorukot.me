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

func TestRenderIgnoresRawTagLines(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"# Hello",
		"",
		"<Alert type=\"info\">",
		"Read this",
		"</Alert>",
	}, "\n")

	out, err := renderer.Render(input, 40, "dark")
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if !strings.Contains(out, "Read this") {
		t.Fatalf("expected rendered output to keep tag body\noutput:\n%s", out)
	}

	if strings.Contains(out, "<Alert") || strings.Contains(out, "</Alert>") {
		t.Fatalf("expected rendered output to ignore raw tag lines\noutput:\n%s", out)
	}
	if strings.Contains(out, "&lt;Alert") {
		t.Fatalf("expected rendered output to avoid escaped raw tag lines\noutput:\n%s", out)
	}
}
