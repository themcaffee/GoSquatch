package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gomarkdown/markdown/ast"
)

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

func getThemeConfig(fp string) (ThemeConfig, error) {
	var configStruct ThemeConfig
	// read the config file
	config, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("Could not read file: ", fp)
		return configStruct, err
	}
	err = json.Unmarshal(config, &configStruct)
	if err != nil {
		fmt.Println("Could not parse config file: ", fp)
		return configStruct, err
	}
	return configStruct, nil
}

func (app App) renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
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
		level := strconv.Itoa(node.(*ast.Heading).Level)
		if !entering {
			w.Write([]byte("</h" + level + ">\n"))
			return ast.GoToNext, true
		}
		switch level {
		case "1":
			node := fmt.Sprintf("<h1 class=\"%s\">", app.ThemeConfig.Heading.Level.One)
			w.Write([]byte(node))
		case "2":
			node := fmt.Sprintf("<h2 class=\"%s\">", app.ThemeConfig.Heading.Level.Two)
			w.Write([]byte(node))
		case "3":
			node := fmt.Sprintf("<h3 class=\"%s\">", app.ThemeConfig.Heading.Level.Three)
			w.Write([]byte(node))
		case "4":
			node := fmt.Sprintf("<h4 class=\"%s\">", app.ThemeConfig.Heading.Level.Four)
			w.Write([]byte(node))
		case "5":
			node := fmt.Sprintf("<h5 class=\"%s\">", app.ThemeConfig.Heading.Level.Five)
			w.Write([]byte(node))
		case "6":
			node := fmt.Sprintf("<h6 class=\"%s\">", app.ThemeConfig.Heading.Level.Six)
			w.Write([]byte(node))
		}
		return ast.GoToNext, true
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
		if !entering {
			w.Write([]byte("</figure>"))
		}
		src := string(node.(*ast.Image).Destination)
		c := node.(*ast.Image).GetChildren()[0]
		alt := string(c.AsLeaf().Literal)
		if entering && alt != "" {
			node := fmt.Sprintf("<figure class=\"%s\"><img src=\"%s\" alt=\"%s\">", app.ThemeConfig.Image, src, alt)
			w.Write([]byte(node))
		} else {
			node := fmt.Sprintf("<figure class=\"%s\"><img src=\"%s\">", app.ThemeConfig.Image, src)
			w.Write([]byte(node))
		}
		return ast.GoToNext, true
	} else if _, ok := node.(*ast.Text); ok {
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
