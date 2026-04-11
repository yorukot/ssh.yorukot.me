package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yorukot/ssh.yorukot.me/pkg/stainmd"
)

const defaultMarkdownPath = "pkg/stainmd/demo/special-recruit.md"

func main() {
	inputPath := flag.String("input", defaultMarkdownPath, "markdown file to render")
	outputPath := flag.String("output", "", "optional file path to write rendered output")
	width := flag.Int("width", 72, "render width")
	flag.Parse()

	renderer := stainmd.New()

	markdown, err := loadMarkdown(*inputPath)
	if err != nil {
		log.Fatal(err)
	}

	out, err := renderer.Render(markdown, *width)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rendering %s\n\n", *inputPath)
	fmt.Print(out)

	if *outputPath != "" {
		if err := os.WriteFile(*outputPath, []byte(out), 0o644); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n\nSaved rendered output to %s\n", *outputPath)
	}
}

func loadMarkdown(path string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
