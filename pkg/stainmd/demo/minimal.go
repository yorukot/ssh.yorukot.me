package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yorukot/ssh.yorukot.me/pkg/stainmd"
)

const defaultMarkdownPath = "pkg/stainmd/demo/sample.md"

func main() {
	renderer := stainmd.New()
	path := defaultMarkdownPath

	markdown, err := loadMarkdown(path)
	if err != nil {
		log.Fatal(err)
	}

	out, err := renderer.Render(markdown, 72)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rendering %s\n\n", path)
	fmt.Print(out)
}

func loadMarkdown(path string) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
