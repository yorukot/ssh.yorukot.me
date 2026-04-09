package stainmd

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
)

func TestRenderBasicMarkdown(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"# Hello",
		"",
		"This is a [link](https://example.com) and `code`.",
		"",
		"> quoted text",
		"",
		"- first item",
		"- second item",
		"",
		"```go",
		"fmt.Println(\"hi\")",
		"```",
	}, "\n")

	out, err := renderer.Render(input, 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)

	checks := []string{
		"Hello",
		"This is a ",
		"code",
		"│ quoted text",
		"• first item",
		"• second item",
		"fmt.Println(\"hi\")",
	}

	for _, want := range checks {
		if !strings.Contains(plain, want) {
			t.Fatalf("rendered output missing %q\noutput:\n%s", want, out)
		}
	}

	if !strings.Contains(out, xansi.SetHyperlink("https://example.com")) {
		t.Fatalf("expected OSC8 hyperlink start sequence, got:\n%s", out)
	}

	if !strings.Contains(out, xansi.ResetHyperlink()) {
		t.Fatalf("expected OSC8 hyperlink reset sequence, got:\n%s", out)
	}

	if strings.Contains(plain, "https://example.com)") || strings.Contains(plain, " link https://example.com") {
		t.Fatalf("rendered output should not include visible markdown link destination\noutput:\n%s", out)
	}
}

func TestRenderOrderedList(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("1. first\n2. second", 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)

	if !strings.Contains(plain, "1. first") {
		t.Fatalf("expected ordered list first item, got:\n%s", out)
	}

	if !strings.Contains(plain, "2. second") {
		t.Fatalf("expected ordered list second item, got:\n%s", out)
	}
}

func TestRenderWrapsParagraphs(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("This is a fairly long line that should wrap when width is small.", 20)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if !strings.Contains(out, "\n") {
		t.Fatalf("expected wrapped output to contain a newline, got:\n%s", out)
	}
}

func TestRenderTightListInlineMarkdown(t *testing.T) {
	renderer := New()
	renderer.Content.InlineCode = renderer.Content.InlineCode.Bold(true)

	out, err := renderer.Render("- item with `code`", 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if strings.Contains(out, "`code`") {
		t.Fatalf("expected inline markdown in tight list item to be rendered, got:\n%s", out)
	}

	if !strings.Contains(out, "code") {
		t.Fatalf("expected rendered inline code text, got:\n%s", out)
	}
}

func TestRenderFencedCodeBlockHighlightsSyntax(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("```go\nfmt.Println(\"hi\")\n```", 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)
	if !strings.Contains(plain, "fmt.Println(\"hi\")") {
		t.Fatalf("expected code block content, got:\n%s", out)
	}

	if countANSICodes(out) < 4 {
		t.Fatalf("expected highlighted output to contain multiple ANSI sequences, got:\n%s", out)
	}
}

func TestRenderTable(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"| Name | Status |",
		"| --- | --- |",
		"| Links | Ready |",
		"| Tables | New |",
	}, "\n")

	out, err := renderer.Render(input, 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)
	checks := []string{
		"+",
		"| Name",
		"Status |",
		"Links",
		"Tables",
	}

	for _, want := range checks {
		if !strings.Contains(plain, want) {
			t.Fatalf("expected table output to contain %q, got:\n%s", want, out)
		}
	}
}

func countANSICodes(value string) int {
	return len(ansiPattern.FindAllString(value, -1))
}
