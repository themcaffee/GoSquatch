package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

type App struct {
	SiteTemplate  string
	SrcDir        string
	DistDir       string
	Layouts       map[string]string
	Pages         []Page
	IgnoreFolders map[string]bool
	IgnoreFiles   map[string]bool
	ThemeConfig   ThemeConfig
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

func (app App) getPage(fp string) (Page, error) {
	page := Page{Filepath: fp}
	// read the markdown file
	md, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("Could not read file: ", fp)
		return page, err
	}

	lines := strings.Split(string(md), "\n")
	contentStart := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				fm := lines[1:i]
				for _, l := range fm {
					l = strings.TrimSpace(l)
					if strings.HasPrefix(l, "title:") {
						page.Title = strings.TrimSpace(strings.TrimPrefix(l, "title:"))
					} else if strings.HasPrefix(l, "layout:") {
						page.Layout = strings.TrimSpace(strings.TrimPrefix(l, "layout:"))
					}
				}
				contentStart = i + 1
				break
			}
		}
	}

	for _, line := range lines[contentStart:] {
		if strings.HasPrefix(line, "[_metadata_:title]:- \"") {
			title := strings.TrimPrefix(line, "[_metadata_:title]:- \"")
			page.Title = strings.TrimSuffix(title, "\"")
		}
		if strings.HasPrefix(line, "[_metadata_:layout]:- \"") {
			layout := strings.TrimPrefix(line, "[_metadata_:layout]:- \"")
			page.Layout = strings.TrimSuffix(layout, "\"")
		}
	}

	// render the markdown file (without frontmatter)
	bodyLines := lines[contentStart:]
	mdBody := strings.Join(bodyLines, "\n")
	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: app.renderHook,
	}
	renderer := html.NewRenderer(opts)
	page.Body = string(markdown.ToHTML([]byte(mdBody), nil, renderer))

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

		// Ignore directories and files
		if info.IsDir() {
			if _, ok := app.IgnoreFolders[info.Name()]; ok {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if _, ok := app.IgnoreFiles[info.Name()]; ok {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// parse the layouts
		ext := filepath.Ext(path)
		base := filepath.Base(path)
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
			page, err := app.getPage(path)
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

func InitApp(srcDir string) (App, error) {
	app := App{SrcDir: srcDir}
	// Parse the theme config
	squatchConfig, err := getSquatchConfig(filepath.Join(app.SrcDir, ".squatch"))
	if err != nil {
		return app, err
	}
	app.DistDir = squatchConfig.DistDir
	// load the list of folders to ignore
	app.IgnoreFolders = map[string]bool{app.DistDir: true}
	for _, folder := range squatchConfig.IgnoreFolders {
		if folder == "" {
			continue
		}
		app.IgnoreFolders[folder] = true
	}
	app.IgnoreFiles = map[string]bool{}
	// load the list of files to ignore
	for _, file := range squatchConfig.IgnoreFiles {
		if file == "" {
			continue
		}
		app.IgnoreFiles[file] = true
	}
	// Remove existing dist directory
	os.RemoveAll(app.DistDir)
	// Create new dist directory
	os.Mkdir(app.DistDir, 0755)
	err = app.parseSrcDirectory()
	if err != nil {
		return app, err
	}
	return app, nil
}

func (app App) printDistFolder() {
	filepath.Walk(app.DistDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
}

func Build(srcDir string) {
	fmt.Println("Starting build...")
	// Get input variables from Github Actions
	srcDirEnv := os.Getenv("INPUT_SRCDIR")
	if len(srcDirEnv) != 0 {
		srcDir = srcDirEnv
	} else if len(srcDir) == 0 {
		srcDir = "src"
	}

	// Initialize the app
	app, err := InitApp(srcDir)
	check(err)

	// Convert all pages
	for _, page := range app.Pages {
		err = app.renderPage(page)
		check(err)
	}
	fmt.Println("Build complete! Dist folder:")
	app.printDistFolder()
}

func main() {
	var srcDir string
	flag.StringVar(&srcDir, "src-dir", "src", "Source directory")
	var port string
	flag.StringVar(&port, "port", "8080", "Port to run the live server on")
	liveServerPtr := flag.Bool("live-server", false, "Run a live server")
	flag.Parse()
	if *liveServerPtr {
		LiveServer(srcDir, port)
	} else {
		Build(srcDir)
	}
}
