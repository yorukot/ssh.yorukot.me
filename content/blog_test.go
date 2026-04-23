package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBlogPostsFromDirSkipsDraftPostsAndSortsPublishedPosts(t *testing.T) {
	root := t.TempDir()

	writeBlogPost(t, filepath.Join(root, "visible-old", "index.md"), `---
title: Visible Old
publish_date: 2024-01-02
post_slug: visible-old
---

Visible old content.
`)
	writeBlogPost(t, filepath.Join(root, "draft-hidden", "index.md"), `---
title: Draft Hidden
publish_date: 2025-05-06
post_slug: draft-hidden
drafted: true
draft: true
---

Draft content.
`)
	writeBlogPost(t, filepath.Join(root, "visible-new", "index.md"), `---
title: Visible New
publish_date: 2025-07-08
post_slug: visible-new
---

Visible new content.
`)

	posts, err := blogPostsFromDir(root)
	if err != nil {
		t.Fatalf("blogPostsFromDir returned error: %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 visible posts, got %d", len(posts))
	}
	if posts[0].Path != "/blog/visible-new" {
		t.Fatalf("posts[0].Path = %q, want %q", posts[0].Path, "/blog/visible-new")
	}
	if posts[1].Path != "/blog/visible-old" {
		t.Fatalf("posts[1].Path = %q, want %q", posts[1].Path, "/blog/visible-old")
	}

	if _, err := FindPost(posts, "/blog/draft-hidden"); err == nil {
		t.Fatal("expected drafted post to be absent from FindPost results")
	}
}

func TestBlogPostsFromDirAllowsDraftPostToShareSlugWithPublishedPost(t *testing.T) {
	root := t.TempDir()

	writeBlogPost(t, filepath.Join(root, "visible-post", "index.md"), `---
title: Visible
publish_date: 2024-01-02
post_slug: shared-slug
---

Visible content.
`)
	writeBlogPost(t, filepath.Join(root, "draft-post", "index.md"), `---
title: Draft
publish_date: 2025-01-02
post_slug: shared-slug
draft: true
---

Draft content.
`)

	posts, err := blogPostsFromDir(root)
	if err != nil {
		t.Fatalf("blogPostsFromDir returned error: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 visible post, got %d", len(posts))
	}
	if posts[0].Path != "/blog/shared-slug" {
		t.Fatalf("posts[0].Path = %q, want %q", posts[0].Path, "/blog/shared-slug")
	}
}

func writeBlogPost(t *testing.T, path, body string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}
