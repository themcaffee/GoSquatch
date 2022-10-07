package main

import (
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
		ThemeConfig: ThemeConfig{
			Heading: Heading{
				Level: Level{
					One:   "heading-one",
					Two:   "heading-two",
					Three: "heading-three",
					Four:  "heading-four",
					Five:  "heading-five",
					Six:   "heading-six",
				},
			},
		},
	}
	md := []byte("# Heading One\n## Heading Two\n### Heading Three\n#### Heading Four\n##### Heading Five\n###### Heading Six")
	// render the markdown file
	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: app.renderHook,
	}
	renderer := html.NewRenderer(opts)
	output := string(markdown.ToHTML(md, nil, renderer))
	if output != "<h1 class=\"heading-one\">Heading One</h1>\n<h2 class=\"heading-two\">Heading Two</h2>\n<h3 class=\"heading-three\">Heading Three</h3>\n<h4 class=\"heading-four\">Heading Four</h4>\n<h5 class=\"heading-five\">Heading Five</h5>\n<h6 class=\"heading-six\">Heading Six</h6>\n" {
		t.Errorf("Expected <h1 class=\"heading-one\">Heading One</h1>\\n<h2 class=\"heading-two\">Heading Two</h2>\\n<h3 class=\"heading-three\">Heading Three</h3>\\n<h4 class=\"heading-four\">Heading Four</h4>\\n<h5 class=\"heading-five\">Heading Five</h5>\\n<h6 class=\"heading-six\">Heading Six</h6>\\n, got %s", output)
	}
}
