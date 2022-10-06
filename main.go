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
	SiteTemplate string
	SrcDir       string
	DistDir      string
	Layouts      map[string]string
	Pages        []Page
}

type Page struct {
	Title    string
	Body     string
	Layout   string
	Filepath string
}

type InvalidPageError struct {
	s string
}

func (e InvalidPageError) Error() string {
	return e.s
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
	page := Page{Filepath: fp}
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

	// If the page metadata cannot be found, return an error to skip the page
	// This is useful for markdown that are not pages
	if page.Title == "" {
		return page, InvalidPageError{s: fmt.Sprintf("no title found in %v", fp)}
	}
	if page.Layout == "" {
		return page, InvalidPageError{s: fmt.Sprintf("no layout found in %v", fp)}
	}
	return page, nil
}

func (app App) renderPage(page Page) (err error) {
	innerLayout, ok := app.Layouts[page.Layout]
	if !ok {
		// Skip the page if the layout is not found
		fmt.Printf("Could not find layout %v for page %v", page.Layout, page.Filepath)
		return nil
	}

	// Parse inner template
	t := template.New("page")
	t, err = t.Parse(innerLayout)
	if err != nil {
		fmt.Printf("error parsing template file at %v: %v\n", app.Layouts[page.Layout], err)
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
	t, err = t.Parse(app.SiteTemplate)
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
	relpath, err := filepath.Rel(app.SrcDir, page.Filepath)
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

func (app *App) parseSrcDirectory() error {
	app.Layouts = make(map[string]string)
	app.Pages = make([]Page, 0)
	err := filepath.Walk(app.SrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == app.DistDir {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		base := filepath.Base(path)
		// parse the layouts
		if base == "layout.html" {
			layoutByte, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Could not read file: ", path)
				return err
			}
			app.SiteTemplate = string(layoutByte)
		} else if ext == ".html" && strings.HasPrefix(base, "layout_") {
			name := filepath.Base(path)
			name = strings.TrimSuffix(name, ".html")
			name = strings.TrimPrefix(name, "layout_")
			layoutByte, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("error reading template file at %v: %v\n", path, err)
				return err
			}
			app.Layouts[name] = string(layoutByte)
		} else if ext == ".md" {
			page, err := getPage(path)
			// Skip pages we can't read because they could be README, LICENSE, drafts, etc.
			if _, ok := err.(InvalidPageError); ok {
				return nil
			} else if err != nil {
				fmt.Println("Could not read file: ", path)
				return err
			}
			app.Pages = append(app.Pages, page)
		} else {
			// Copy any other file to the dist directory
			relpath, err := filepath.Rel(app.SrcDir, path)
			if err != nil {
				fmt.Println("Could not get relative path: ", err)
				return err
			}
			newFilePath := filepath.Join(app.DistDir, relpath)
			if err := os.MkdirAll(filepath.Dir(newFilePath), 0755); err != nil {
				fmt.Println("Could not create directory: ", err)
				return err
			}
			source, err := os.Open(path)
			if err != nil {
				fmt.Println("Could not open source file: ", err)
				return err
			}
			defer source.Close()
			destination, err := os.Create(newFilePath)
			if err != nil {
				fmt.Println("Could not create destination file: ", err)
				return err
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err != nil {
				fmt.Println("Could not copy file: ", err)
				return err
			}
		}
		return nil
	})
	return err
}

func InitApp(srcDir string, distDir string) (App, error) {
	app := App{SrcDir: srcDir, DistDir: distDir}
	// Remove existing dist directory
	os.RemoveAll(distDir)
	// Create new dist directory
	os.Mkdir(distDir, 0755)
	err := app.parseSrcDirectory()
	if err != nil {
		return app, err
	}
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
	for _, page := range app.Pages {
		err = app.renderPage(page)
		check(err)
	}
	fmt.Println("Build complete!")
}
