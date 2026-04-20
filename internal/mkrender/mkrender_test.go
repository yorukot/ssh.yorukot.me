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

func TestResolveMarkdownImagePathResolvesRelativeImagePaths(t *testing.T) {
	sourcePath := "src/content/blog/before-you-build-a-tui-or-cli-app/index.md"
	tests := map[string]string{
		"./tui-example.webp":    "src/content/blog/before-you-build-a-tui-or-cli-app/tui-example.webp",
		"./docker-example.png":  "src/content/blog/before-you-build-a-tui-or-cli-app/docker-example.png",
		"../shared/example.png": "src/content/blog/shared/example.png",
	}

	for input, want := range tests {
		if got := resolveMarkdownImagePath(sourcePath, input); got != want {
			t.Fatalf("resolveMarkdownImagePath(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestResolveMarkdownImagePathLeavesExternalAndAbsoluteImagePaths(t *testing.T) {
	sourcePath := "src/content/blog/how-to-custom-gnome-terminal/index.md"
	tests := []string{
		"https://i.imgur.com/xyN5XY2.webp",
		"/blog-assets/how-to-custom-gnome-terminal/intro.webp",
		"mailto:hi@yorukot.me",
	}

	for _, input := range tests {
		if got := resolveMarkdownImagePath(sourcePath, input); got != input {
			t.Fatalf("resolveMarkdownImagePath(%q) = %q, want unchanged", input, got)
		}
	}
}

func TestRenderWithSourceUsesResolvedImagePaths(t *testing.T) {
	renderer := New()

	out, err := renderer.RenderWithSource(
		"![intro image](./intro.webp)",
		"src/content/blog/how-to-custom-gnome-terminal/index.md",
		160,
		"dark",
	)
	if err != nil {
		t.Fatalf("RenderWithSource returned error: %v", err)
	}

	want := "src/content/blog/how-to-custom-gnome-terminal/intro.webp"
	if !strings.Contains(out, want) {
		t.Fatalf("expected rendered image to use resolved path %q\noutput:\n%s", want, out)
	}
}
