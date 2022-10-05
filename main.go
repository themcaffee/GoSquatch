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
	LayoutsDir   map[string]string
}

type Page struct {
	Title  string
	Body   string
	Layout string
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

func getPage(fp string) (Page, error) {
	page := Page{}
	// read the markdown file
	md, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("Could not read file: ", fp)
		return page, err
	}

	// render the markdown file
	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: renderHook,
	}
	renderer := html.NewRenderer(opts)
	page.Body = string(markdown.ToHTML(md, nil, renderer))

	// Get metadata
	lines := strings.Split(string(md), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[_metadata_:title]:- \"") {
			title := strings.TrimPrefix(line, "[_metadata_:title]:- \"")
			page.Title = strings.TrimSuffix(title, "\"")
		}
		if strings.HasPrefix(line, "[_metadata_:layout]:- \"") {
			layout := strings.TrimPrefix(line, "[_metadata_:layout]:- \"")
			page.Layout = strings.TrimSuffix(layout, "\"")
		}
	}
	if page.Title == "" {
		return page, fmt.Errorf("no title found in %v", fp)
	}
	if page.Layout == "" {
		return page, fmt.Errorf("no layout found in %v", fp)
	}
	return page, nil
}

func (app App) renderPage(fp string) (err error) {
	page, err := getPage(fp)
	if err != nil {
		fmt.Println("Could not get title for file: ", fp)
		return err
	}

	innerLayoutByte, err := os.ReadFile(app.LayoutsDir["layout_"+page.Layout])
	if err != nil {
		fmt.Printf("error reading template file at %v: %v\n", app.LayoutsDir["layout_"+page.Layout], err)
		return err
	}
	innerLayout := string(innerLayoutByte)

	// Parse inner template
	t := template.New("page")
	t, err = t.Parse(innerLayout)
	if err != nil {
		fmt.Printf("error parsing template file at %v: %v\n", app.LayoutsDir[page.Layout], err)
		return err
	}
	t = template.Must(t, err)
	var inner bytes.Buffer
	err = t.Execute(&inner, page)
	if err != nil {
		fmt.Println("error executing template: ", err)
		return err
	}
	page.Body = inner.String()

	// Parse outer template
	t = template.New("Render")
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
	relpath, err := filepath.Rel(app.SrcDir, fp)
	if err != nil {
		fmt.Println("Could not get relative path: ", err)
		return err
	}
	newFilePath := relpath[:len(relpath)-3] + ".html"
	newFilePath = filepath.Join(app.DistDir, newFilePath)
	if err := os.MkdirAll(filepath.Dir(newFilePath), 0755); err != nil {
		fmt.Println("Could not create directory: ", err)
		return err
	}
	err = os.WriteFile(newFilePath, processed.Bytes(), 0644)
	if err != nil {
		fmt.Println("Could not write file: ", err)
		return err
	}
	return err
}

func (app App) getAllPages() ([]string, error) {
	// get all markdown files
	var mdFiles []string
	pageDir := filepath.Join(app.SrcDir)
	filepath.Walk(pageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			mdFiles = append(mdFiles, path)
		}
		return nil
	})
	return mdFiles, nil
}

func getAllLayouts(srcDir string) (map[string]string, error) {
	// get all layouts
	layouts := make(map[string]string)
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			name := filepath.Base(path)
			name = name[:len(name)-5]
			layouts[name] = path
		}
		return nil
	})
	return layouts, nil
}

func InitApp(srcDir string, distDir string) (App, error) {
	app := App{SrcDir: srcDir, DistDir: distDir}
	// Remove existing dist directory
	os.RemoveAll(distDir)
	// Create new dist directory
	os.Mkdir(distDir, 0755)
	var err error
	app.LayoutsDir, err = getAllLayouts(srcDir)
	if err != nil {
		return app, err
	}
	// cache the page layout
	pageTemplateByte, err := os.ReadFile(app.LayoutsDir["layout"])
	if err != nil {
		fmt.Printf("error reading template file at %v: %v\n", app.LayoutsDir["layout"], err)
		return app, err
	}
	app.PageTemplate = string(pageTemplateByte)
	return app, nil
}

func main() {
	fmt.Println("Starting build...")
	// Get input variables from Github Actions
	srcDir := os.Getenv("INPUT_SRCDIR")
	if len(srcDir) == 0 {
		srcDir = "src"
	}
	distDir := os.Getenv("INPUT_DISTDIR")
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
