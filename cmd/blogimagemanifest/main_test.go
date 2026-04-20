package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildManifestMapsLocalImagesToAstroURLs(t *testing.T) {
	root := t.TempDir()
	blogRoot := filepath.Join(root, "content", "markdown", "blog")
	distRoot := filepath.Join(root, "dist")

	postDir := filepath.Join(blogRoot, "before-you-build-a-tui-or-cli-app")
	writeFile(t, filepath.Join(postDir, "index.md"), `---
title: Test
post_slug: before-you-build-a-tui-or-cli-app
---

![TUI example](./tui-example.webp)
![External](https://i.imgur.com/xyN5XY2.webp)
![CLI example](./docker-example.png)
`)

	writeFile(t, filepath.Join(distRoot, "blog", "before-you-build-a-tui-or-cli-app", "index.html"), `
<article>
  <img src="/_astro/tui-example.AAA.webp" alt="TUI example">
  <img src="https://i.imgur.com/xyN5XY2.webp" alt="External">
  <img src="/_astro/docker-example.BBB.webp" alt="CLI example">
</article>
`)

	manifest, err := buildManifest(blogRoot, distRoot, "https://yorukot.me")
	if err != nil {
		t.Fatalf("buildManifest returned error: %v", err)
	}

	checks := map[string]string{
		filepath.ToSlash(filepath.Join(postDir, "tui-example.webp")):   "https://yorukot.me/_astro/tui-example.AAA.webp",
		filepath.ToSlash(filepath.Join(postDir, "docker-example.png")): "https://yorukot.me/_astro/docker-example.BBB.webp",
	}

	if len(manifest.Images) != len(checks) {
		t.Fatalf("expected %d image mappings, got %d: %#v", len(checks), len(manifest.Images), manifest.Images)
	}
	for path, want := range checks {
		if got := manifest.Images[path]; got != want {
			t.Fatalf("manifest.Images[%q] = %q, want %q", path, got, want)
		}
	}
}

func TestBuildManifestFailsWhenAstroImagesDoNotMatchLocalImages(t *testing.T) {
	root := t.TempDir()
	blogRoot := filepath.Join(root, "content", "markdown", "blog")
	distRoot := filepath.Join(root, "dist")

	postDir := filepath.Join(blogRoot, "post")
	writeFile(t, filepath.Join(postDir, "index.md"), `---
title: Test
post_slug: post
---

![First](./first.png)
![Second](./second.png)
`)
	writeFile(t, filepath.Join(distRoot, "blog", "post", "index.html"), `<img src="/_astro/first.hash.webp">`)

	if _, err := buildManifest(blogRoot, distRoot, "https://yorukot.me"); err == nil {
		t.Fatal("expected buildManifest to fail when _astro image count does not match local image count")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}
