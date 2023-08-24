package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"path/filepath"
	"regexp"

	"github.com/gomarkdown/markdown/ast"
)

type SquatchConfig struct {
	DistDir       string      `json:"dist"`
	IgnoreFolders string      `json:"ignore_folders"`
	IgnoreFiles   string      `json:"ignore_files"`
	ThemeConfig   ThemeConfig `json:"theme"`
}

type ThemeConfig struct {
	BlockTemplate    string         `json:"block_template"`
	List             List           `json:"list"`
	ListItem         ListItem       `json:"list_item"`
	Paragraph        string         `json:"paragraph"`
	Math             string         `json:"math"`
	MathBlock        string         `json:"math_block"`
	Heading          Heading        `json:"heading"`
	HorizontalRule   string         `json:"horizontal_rule"`
	Emph             string         `json:"emph"`
	Strong           string         `json:"strong"`
	Del              string         `json:"del"`
	Link             Link           `json:"link"`
	CrossReference   CrossReference `json:"cross_reference"`
	Citation         Citation       `json:"citation"`
	Image            string         `json:"image"`
	Text             string         `json:"text"`
	HTMLBlock        string         `json:"html_block"`
	CodeBlock        CodeBlock      `json:"code_block"`
	Softbreak        string         `json:"softbreak"`
	Hardbreak        string         `json:"hardbreak"`
	NonBlockingSpace string         `json:"non_blocking_space"`
	Code             string         `json:"code"`
	HTMLSpan         string         `json:"html_span"`
	Table            string         `json:"table"`
	TableCell        TableCell      `json:"table_cell"`
	TableHeader      string         `json:"table_header"`
	TableBody        string         `json:"table_body"`
	TableRow         string         `json:"table_row"`
	TableFooter      string         `json:"table_footer"`
	Caption          string         `json:"caption"`
	CaptionFigure    string         `json:"caption_figure"`
	Callout          Callout        `json:"callout"`
	Index            Index          `json:"index"`
	Subscript        string         `json:"subscript"`
	Superscript      string         `json:"superscript"`
	Footnotes        string         `json:"footnotes"`
}

type List struct {
}

type ListItem struct {
}

type Heading struct {
	Level Level `json:"level"`
}

type Level struct {
	One   string `json:"1"`
	Two   string `json:"2"`
	Three string `json:"3"`
	Four  string `json:"4"`
	Five  string `json:"5"`
	Six   string `json:"6"`
}

type Link struct {
}

type CrossReference struct {
}

type Citation struct {
}

type CodeBlock struct {
}

type TableCell struct {
}

type Callout struct {
}

type Index struct {
}

func getSquatchConfig(fp string) (SquatchConfig, error) {
	var configStruct SquatchConfig
	// read the config file
	config, err := os.ReadFile(fp)
	if err != nil {
		configStruct = SquatchConfig{DistDir: "dist", IgnoreFolders: "", IgnoreFiles: ""}
		return configStruct, nil
	}
	err = json.Unmarshal(config, &configStruct)
	if err != nil {
		fmt.Println("Could not parse config file: ", fp)
		return configStruct, err
	}
	return configStruct, nil
}

func (app App) renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// TODO: implement this
	if _, ok := node.(*ast.List); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.ListItem); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Paragraph); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Math); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.MathBlock); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Heading); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.HorizontalRule); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Emph); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Strong); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Del); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Link); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.CrossReference); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Citation); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Image); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Text); ok {
		wikiLinkRegexCompiled := regexp.MustCompile("/[[(.+?(|.+?)?)]]([W])/g")
		wikiLink := wikiLinkRegexCompiled.ReplaceAllStringFunc(string(node.(*ast.Text).Literal), func(match string) string {
			matches := wikiLinkRegexCompiled.FindStringSubmatch(match)
			p1, p2, p3 := matches[1], matches[2], matches[3]
			// p1 = Whole match, p2 = Link text, p3 = ?Whitespace
			var slug string
			if strings.Contains(p1, "|") {
				slug = p1[:strings.Index(p1, "|")]
			} else {
				slug = p1
			}
			var filePathArr []string
			filepath.Walk("/src", func(path string, info os.FileInfo, err error) error {
				// TODO: Load the source directory from parameters or config
				if err != nil {
					return err
				}
				if !info.IsDir() && info.Name() == slug+".md" {
					filePathArr = append(filePathArr, path[:len(path)-3])
				}
				return nil
			})
			var href string
			if len(filePathArr) == 1 {
				href = filePathArr[0]
			} else {
				href = slug // If href = slug, then it's a broken link
			}
			var linkText string
			if p2 != "" {
				linkText = p2[1:]
			} else {
				linkText = p1
			}

			linkMarkdown := fmt.Sprintf(`[%s](%s)%s`, linkText, href, p3)
			return linkMarkdown
		})

		node.(*ast.Text).Literal = []byte(wikiLink)
		// false so the default renderer handles the added link
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.HTMLBlock); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.CodeBlock); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Softbreak); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Hardbreak); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.NonBlockingSpace); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Code); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.HTMLSpan); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.Table); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.TableCell); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.TableHeader); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.TableBody); ok {
		return ast.GoToNext, false
	} else if _, ok := node.(*ast.TableRow); ok {
		return ast.GoToNext, false
	} else {
		return ast.GoToNext, false
	}
}
