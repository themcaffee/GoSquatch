package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type App struct {
	PageTemplate string
	SrcDir       string
	DistDir      string
}

type Page struct {
	Title string
	Body  string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if _, ok := node.(*ast.Heading); ok {
		level := strconv.Itoa(node.(*ast.Heading).Level)
		if entering && level == "1" {
			w.Write([]byte(`<h1 class="title is-1 has-text-centered">`))
		} else if entering {
			w.Write([]byte("<h" + level + ">"))
		} else {
			w.Write([]byte("</h" + level + ">"))
		}
		return ast.GoToNext, true
	} else if _, ok := node.(*ast.Image); ok {
		src := string(node.(*ast.Image).Destination)
		c := node.(*ast.Image).GetChildren()[0]
		alt := string(c.AsLeaf().Literal)
		if entering && alt != "" {
			w.Write([]byte(`<figure class="image is-5by3"><img src="` + src + `" alt="` + alt + `">`))
		} else if entering {
			w.Write([]byte(`<figure class="image is-5by3"><img src="` + src + `">`))
		} else {
			w.Write([]byte(`</figure>`))
		}
		return ast.SkipChildren, true
	} else {
		return ast.GoToNext, false
	}
}

func getTitle(md []byte) (string, error) {
	lines := strings.Split(string(md), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[_metadata_:title]:- \"") {
			title := strings.TrimPrefix(line, "[_metadata_:title]:- \"")
			return strings.TrimSuffix(title, "\""), nil
		}
	}
	return "", fmt.Errorf("could not find title")
}

func (app App) renderPage(fp string) (err error) {
	filename := filepath.Base(fp)
	filename = filename[:len(filename)-3]

	// read the markdown file
	md, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("Could not read file: ", fp)
		return err
	}

	// render the markdown file
	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: renderHook,
	}
	renderer := html.NewRenderer(opts)
	output := string(markdown.ToHTML(md, nil, renderer))

	title, err := getTitle(md)
	if err != nil {
		fmt.Println("Could not get title for file: ", fp)
		return err
	}

	// render the page with a template
	page := Page{Title: title, Body: output}
	t := template.New("Render")
	t, err = t.Parse(app.PageTemplate)
	if err != nil {
		fmt.Println("Could not parse template: ", err)
		return err
	}
	t = template.Must(t, err)
	var processed bytes.Buffer
	err = t.Execute(&processed, page)
	if err != nil {
		fmt.Println("Could not execute template: ", err)
		return err
	}

	// write the page to a file
	newFileName := filename + ".html"
	newFilePath := filepath.Join(app.DistDir, newFileName)
	err = os.WriteFile(newFilePath, processed.Bytes(), 0644)
	if err != nil {
		fmt.Println("Could not write file: ", err)
		return err
	}
	return err
}

func (app App) getAllPages() ([]string, error) {
	// get all markdown files
	pageDir := filepath.Join(app.SrcDir, "pages")
	files, err := os.ReadDir(pageDir)
	if err != nil {
		return nil, err
	}
	var mdFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			mdFiles = append(mdFiles, filepath.Join(app.SrcDir, "pages", file.Name()))
		}
	}
	return mdFiles, err
}

func InitApp(srcDir string, distDir string) (App, error) {
	app := App{SrcDir: srcDir, DistDir: distDir}
	// Remove existing dist directory
	os.RemoveAll(distDir)
	// Create new dist directory
	os.Mkdir(distDir, 0755)
	// cache the page layout
	templateLocation := filepath.Join(srcDir, "layout.html")
	pageTemplateByte, err := os.ReadFile(templateLocation)
	if err != nil {
		fmt.Printf("error reading template file at %v: %v\n", err, templateLocation)
		return app, err
	}
	app.PageTemplate = string(pageTemplateByte)
	return app, nil
}

func main() {
	fmt.Println("Starting build...")
	// Get input variables from Github Actions
	srcDir := os.Getenv("INPUT_SRC_DIR")
	if len(srcDir) == 0 {
		srcDir = "src"
	}
	distDir := os.Getenv("INPUT_DIST_DIR")
	if len(distDir) == 0 {
		distDir = "dist"
	}

	// Initialize the app
	app, err := InitApp(srcDir, distDir)
	check(err)

	// Convert all pages
	pages, err := app.getAllPages()
	check(err)
	for _, page := range pages {
		err = app.renderPage(page)
		check(err)
	}
	fmt.Println("Build complete!")
}
