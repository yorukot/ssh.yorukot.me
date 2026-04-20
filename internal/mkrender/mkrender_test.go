package mkrender

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
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
		if got := resolveMarkdownImagePath(sourcePath, input, nil); got != want {
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
		if got := resolveMarkdownImagePath(sourcePath, input, nil); got != input {
			t.Fatalf("resolveMarkdownImagePath(%q) = %q, want unchanged", input, got)
		}
	}
}

func TestResolveMarkdownImagePathUsesManifestForLocalImages(t *testing.T) {
	sourcePath := "content/markdown/blog/before-you-build-a-tui-or-cli-app/index.md"
	manifest := map[string]string{
		"content/markdown/blog/before-you-build-a-tui-or-cli-app/superfile-hackernews.png": "https://yorukot.me/_astro/superfile-hackernews.DfsmIV1R_1MOWGh.webp",
	}

	got := resolveMarkdownImagePath(sourcePath, "./superfile-hackernews.png", manifest)
	want := "https://yorukot.me/_astro/superfile-hackernews.DfsmIV1R_1MOWGh.webp"
	if got != want {
		t.Fatalf("resolveMarkdownImagePath() = %q, want %q", got, want)
	}
}

func TestResolveMarkdownImagePathPreservesSuffixAfterManifestLookup(t *testing.T) {
	sourcePath := "content/markdown/blog/how-to-custom-gnome-terminal/index.md"
	manifest := map[string]string{
		"content/markdown/blog/how-to-custom-gnome-terminal/intro.webp": "https://yorukot.me/_astro/intro.hash.webp",
	}

	got := resolveMarkdownImagePath(sourcePath, "./intro.webp#caption", manifest)
	want := "https://yorukot.me/_astro/intro.hash.webp#caption"
	if got != want {
		t.Fatalf("resolveMarkdownImagePath() = %q, want %q", got, want)
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

	want := xansi.SetHyperlink("src/content/blog/how-to-custom-gnome-terminal/intro.webp")
	if !strings.Contains(out, want) {
		t.Fatalf("expected rendered image to link to resolved path %q\noutput:\n%s", want, out)
	}
	if !strings.Contains(out, "Image: ") || !strings.Contains(out, "intro image") {
		t.Fatalf("expected rendered image to include image hint\noutput:\n%s", out)
	}
}

func TestRenderWithSourceUsesManifestImageURL(t *testing.T) {
	renderer := Renderer{
		imageManifest: map[string]string{
			"content/markdown/blog/how-to-custom-gnome-terminal/intro.webp": "https://yorukot.me/_astro/intro.BiMzQ8pt_Z1eIygR.webp",
		},
	}

	out, err := renderer.RenderWithSource(
		"![intro image](./intro.webp)",
		"content/markdown/blog/how-to-custom-gnome-terminal/index.md",
		160,
		"dark",
	)
	if err != nil {
		t.Fatalf("RenderWithSource returned error: %v", err)
	}

	want := "https://yorukot.me/_astro/intro.BiMzQ8pt_Z1eIygR.webp"
	if !strings.Contains(out, xansi.SetHyperlink(want)) {
		t.Fatalf("expected rendered image to link to manifest URL %q\noutput:\n%s", want, out)
	}
	if strings.Contains(out, " "+want) {
		t.Fatalf("expected rendered image to hide manifest URL behind OSC8 label\noutput:\n%s", out)
	}
}

func TestLoadImageManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	body := `{
  "site": "https://yorukot.me",
  "images": {
    "content/markdown/blog/post/image.png": "https://yorukot.me/_astro/image.hash.webp"
  }
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	manifest := loadImageManifest(path)
	got := manifest["content/markdown/blog/post/image.png"]
	want := "https://yorukot.me/_astro/image.hash.webp"
	if got != want {
		t.Fatalf("loadImageManifest()[image] = %q, want %q", got, want)
	}
}
