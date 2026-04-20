package stainmd

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
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

	afterLink := xansi.ResetHyperlink() + renderer.Content.Paragraph.Render(" and ")
	if !strings.Contains(out, afterLink) {
		t.Fatalf("expected paragraph style to continue after hyperlink reset\noutput:\n%s", out)
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

func TestRenderRestoresParagraphStyleAfterInlineStyles(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"before **bold** after",
		"before *soft* after",
		"before `code` after",
		"before [link](https://example.com) after",
		"before **bold [link](https://example.com) tail** after",
	}, "\n\n")

	out, err := renderer.Render(input, 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	checks := []string{
		renderer.Content.Strong.Render("bold") + renderer.Content.Paragraph.Render(" after"),
		xansi.ResetHyperlink() + renderer.Content.Paragraph.Render(" after"),
		xansi.ResetHyperlink() + renderer.Content.Strong.Render(" tail") + renderer.Content.Paragraph.Render(" after"),
	}

	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("expected inline style to restore following text style for %q\noutput:\n%s", want, out)
		}
	}

	if strings.Count(out, renderer.Content.Paragraph.Render(" after")) < 4 {
		t.Fatalf("expected paragraph style to be restored after inline segments\noutput:\n%s", out)
	}
}

func TestRenderRestoresParagraphStyleAfterCJKStrongText(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("**你以為的特選**：啊不就成績不好的人選擇逃避讀書的管道嗎，都是一群壞孩子。", 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	wantStrong := renderer.Content.Strong.Render("你以為的特選")
	if !strings.Contains(out, wantStrong) {
		t.Fatalf("expected CJK strong text to use strong style\noutput:\n%s", out)
	}

	wantParagraph := renderer.Content.Paragraph.Render("：啊不就成績不好的人選擇逃避讀書的管道嗎，都是一群壞孩子。")
	if !strings.Contains(out, wantParagraph) {
		t.Fatalf("expected paragraph style to continue after CJK strong text\noutput:\n%s", out)
	}
}

func TestRenderRestoresParagraphStyleAcrossSoftLineBreak(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("第一行\n第二行", 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	want := renderer.Content.Paragraph.Render("第一行") + "\n" + renderer.Content.Paragraph.Render("第二行")
	want = renderer.Content.Paragraph.Render("第一行") + " " + renderer.Content.Paragraph.Render("第二行")
	if !strings.Contains(out, want) {
		t.Fatalf("expected paragraph style to continue across soft line break\noutput:\n%s", out)
	}

	plain := stripANSI(out)
	if !strings.Contains(plain, "第一行 第二行") {
		t.Fatalf("expected visible soft line break content\noutput:\n%s", out)
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

func TestRenderInlineCodeUsesSinglePadding(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("before `code` after", 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	want := renderer.Content.Paragraph.Render("before ") + renderer.Content.InlineCode.Render("code") + renderer.Content.Paragraph.Render(" after")
	if !strings.Contains(out, want) {
		t.Fatalf("expected inline code to render with a single layer of padding\noutput:\n%s", out)
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

func TestRenderFencedCodeBlockWrapsLongLines(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("```go\nconst veryLongIdentifier = \"abcdefghijklmnopqrstuvwxyz0123456789\"\n```", 24)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)
	if strings.Contains(plain, "veryLongIdentifier = \"abcdefghijklmnopqrstuvwxyz0123456789\"") {
		t.Fatalf("expected long fenced code line to wrap, got:\n%s", out)
	}
	if strings.Count(plain, "\n") < 2 {
		t.Fatalf("expected wrapped fenced code block to span multiple lines, got:\n%s", out)
	}
}

func TestRenderIndentedCodeBlockWrapsLongLines(t *testing.T) {
	renderer := New()

	out, err := renderer.Render("    some_super_long_code_token_without_spaces_but_with_underscores_to_force_wrapping", 22)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)
	if strings.Contains(plain, "some_super_long_code_token_without_spaces_but_with_underscores_to_force_wrapping") {
		t.Fatalf("expected indented code block line to wrap, got:\n%s", out)
	}
	if strings.Count(plain, "\n") < 1 {
		t.Fatalf("expected wrapped indented code block to contain a newline, got:\n%s", out)
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
		"╭",
		"│",
		"├",
		"╯",
		"Name",
		"Status",
		"Links",
		"Tables",
	}

	for _, want := range checks {
		if !strings.Contains(plain, want) {
			t.Fatalf("expected table output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderTableLeavesRoomForTerminalWrap(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"| Name | Status |",
		"| --- | --- |",
		"| Links | Ready |",
	}, "\n")

	out, err := renderer.Render(input, 24)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	for _, line := range strings.Split(stripANSI(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if width := lipgloss.Width(line); width > 20 {
			t.Fatalf("expected table line width to leave viewport headroom, got %d for %q\noutput:\n%s", width, line, out)
		}
	}
}

func TestRenderTableWithCustomBorder(t *testing.T) {
	renderer := New()
	renderer.Content.Table.Border = lipgloss.ThickBorder()

	input := strings.Join([]string{
		"| Name | Status |",
		"| --- | --- |",
		"| Links | Ready |",
	}, "\n")

	out, err := renderer.Render(input, 80)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)
	checks := []string{
		"┏",
		"┃",
		"Name",
		"┣",
		"┗",
	}

	for _, want := range checks {
		if !strings.Contains(plain, want) {
			t.Fatalf("expected custom table output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderTableWrapsLongCells(t *testing.T) {
	renderer := New()

	input := strings.Join([]string{
		"| Name | Description |",
		"| --- | --- |",
		"| Links | This table cell should wrap instead of getting clipped by the viewport |",
	}, "\n")

	out, err := renderer.Render(input, 32)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	plain := stripANSI(out)
	checks := []string{
		"This table",
		"cell",
		"should",
		"wrap",
		"instead",
		"viewport",
	}

	for _, want := range checks {
		if !strings.Contains(plain, want) {
			t.Fatalf("expected wrapped table output to contain %q, got:\n%s", want, out)
		}
	}

	if strings.Contains(plain, "This table cell should wrap instead of getting clipped by the viewport |") {
		t.Fatalf("expected long table cell to wrap, got:\n%s", out)
	}
	if strings.Count(plain, "\n") < 6 {
		t.Fatalf("expected wrapped table row to span multiple lines, got:\n%s", out)
	}
}

func countANSICodes(value string) int {
	return len(ansiPattern.FindAllString(value, -1))
}
