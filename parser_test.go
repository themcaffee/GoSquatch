package main

import (
	"strings"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func TestParserHeading(t *testing.T) {
	app := App{
		SiteTemplate:  "site.html",
		SrcDir:        "src",
		DistDir:       "dist",
		Layouts:       make(map[string]string),
		Pages:         make([]Page, 0),
		IgnoreFolders: make(map[string]bool),
		IgnoreFiles:   make(map[string]bool),
		ThemeConfig:   ThemeConfig{},
	}
	md := []byte("# Heading One\n## Heading Two\n### Heading Three\n#### Heading Four\n##### Heading Five\n###### Heading Six")
	// render the markdown file
	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: app.renderHook,
	}
	renderer := html.NewRenderer(opts)
	output := string(markdown.ToHTML(md, nil, renderer))
	output = strings.ReplaceAll(output, "\n", "")
	output = strings.ReplaceAll(output, "\t", "")
	output = strings.ReplaceAll(output, " ", "")
	if output != "<h1>HeadingOne</h1><h2>HeadingTwo</h2><h3>HeadingThree</h3><h4>HeadingFour</h4><h5>HeadingFive</h5><h6>HeadingSix</h6>" {
		t.Errorf("Expected: %s, got: %s", "<h1>Heading One</h1><h2>Heading Two</h2><h3>Heading Three</h3><h4>Heading Four</h4><h5>Heading Five</h5><h6>Heading Six</h6>", output)
	}
}
