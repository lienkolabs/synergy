package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	bytes := markdown.Render(doc, renderer)
	if len(bytes) > 0 {
		return string(bytes)
	}
	return ""
}

func PrettyDuration(d time.Duration) string {
	if d.Seconds() < 60 {
		if d.Seconds() < 2 {
			return "1 second"
		}
		return fmt.Sprintf("%.0f seconds", d.Seconds())
	}
	if d.Minutes() < 60 {
		if d.Minutes() < 2 {
			return "1 minute"
		}
		return fmt.Sprintf("%.0f minutes", d.Minutes())
	}
	if d.Hours() < 24 {
		if d.Hours() < 2 {
			return "1 hour"
		}
		return fmt.Sprintf("%.0f hours", d.Hours())
	}
	if d.Hours()/24.0 < 2 {
		return "1 day"
	}
	return fmt.Sprintf("%.0f days", d.Hours()/24.0)
}

func FileType(fileName string) string {
	pos := strings.LastIndex(fileName, ".")
	if pos < len(fileName) && pos > 0 {
		return fileName[pos+1:]
	}
	return ""
}
