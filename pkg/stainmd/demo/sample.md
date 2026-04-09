# stainmd Demo

This sample shows the Markdown features that `stainmd` currently renders in the terminal.

## Inline Formatting

You can render plain text, *emphasis*, **strong text**, inline code like `go test ./...`, and label-only [OSC 8 links](https://github.com/yorukot/ssh.yorukot.me).

Autolinks also work: <https://catppuccin.com/palette>.

---

## Quotes

> Markdown in a terminal feels good when spacing, wrapping, and color are restrained.
>
> A good renderer should make structure obvious without shouting.

## Lists

- Unordered lists render with bullets
- Inline formatting still works inside list items, including `inline code`
- Images render as a text label plus path

1. Ordered lists keep their numbering
2. Nested formatting like **bold text** still works
3. Links inside lists work too: [demo link](https://example.com/demo)

## Images

![Architecture sketch](./architecture.png)

## Tables

| Feature | Status | Notes |
| --- | --- | --- |
| Links | Ready | Uses OSC 8 label-only hyperlinks |
| Code blocks | Ready | Fenced blocks use syntax highlighting |
| Tables | New | Pipe tables now render in the terminal |

## Code Blocks

```go
package main

import "fmt"

func main() {
	message := "hello from stainmd"
	fmt.Println(message)
}
```

```json
{
  "theme": "catppuccin-mocha",
  "osc8": true,
  "syntaxHighlight": true
}
```

### Final Notes

The sample intentionally includes a mix of short and long lines so wrapping stays visible in the demo output.
