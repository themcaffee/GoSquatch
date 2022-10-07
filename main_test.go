package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func cleanup(dist string) {
	os.RemoveAll(dist)
}

func TestCheckError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected check to panic, got %v", r)
		}
	}()

	check(fmt.Errorf("test error"))
}

func TestInitApp(t *testing.T) {
	srcTest := "src_test"
	distTest := "dist_test"
	defer cleanup(distTest)
	app, err := InitApp(srcTest, distTest)
	if err != nil {
		t.Errorf("expected InitApp to return no error, got %v", err)
	}
	if app.SrcDir != srcTest {
		t.Errorf("expected SrcDir to be %v, got %v", srcTest, app.SrcDir)
	}
	if app.DistDir != distTest {
		t.Errorf("expected DistDir to be %v, got %v", distTest, app.DistDir)
	}
	if len(app.SiteTemplate) == 0 {
		t.Errorf("expected SiteTemplate to be populated")
	}
	if len(app.Pages) == 0 {
		t.Errorf("expected Pages to be populated")
	}
	if len(app.Layouts) == 0 {
		t.Errorf("expected Layouts to be populated")
	}
	// check that non-page files were moved
	_, err = os.Stat(filepath.Join(distTest, "static", "main.css"))
	if err != nil {
		t.Errorf("expected main.css to exist, got %v", err)
	}
	// check that non-page md files were not moved
	_, err = os.Stat(filepath.Join(distTest, "README.html"))
	if err == nil {
		t.Errorf("expected README.html to not exist, got %v", err)
	}
	_, err = os.Stat(filepath.Join(distTest, "README.md"))
	if err == nil {
		t.Errorf("expected README.md to not exist, got %v", err)
	}
}

func TestInitAppEmptyTemplate(t *testing.T) {
	srcTest := "src_test_empty_template"
	distTest := "dist_test"
	defer cleanup(distTest)
	_, err := InitApp(srcTest, distTest)
	if err == nil {
		t.Errorf("expected InitApp to return an error, got %v", err)
	}
}

func TestRenderPage(t *testing.T) {
	srcTest := "src_test"
	distTest := "dist_test"
	defer cleanup(distTest)
	app, _ := InitApp(srcTest, distTest)
	for _, page := range app.Pages {
		err := app.renderPage(page)
		if err != nil {
			t.Errorf("expected renderPage to return no error, got %v", err)
		}
	}
	// check that the files were created
	_, err := os.Stat(filepath.Join(distTest, "index.html"))
	if err != nil {
		t.Errorf("expected index.html to exist, got %v", err)
	}
	_, err = os.Stat(filepath.Join(distTest, "pages", "example.html"))
	if err != nil {
		t.Errorf("expected example.html to exist, got %v", err)
	}
}

func TestGetPage(t *testing.T) {
	srcTest := "src_test"
	distTest := "dist_test"
	defer cleanup(distTest)
	app, _ := InitApp(srcTest, distTest)
	page, err := app.getPage(filepath.Join(srcTest, "pages", "example.md"))
	if err != nil {
		t.Errorf("expected getPage to return no error, got %v", err)
	}
	if page.Title != "Example page title" {
		t.Errorf("expected page title to be 'Example page title', got %v", page.Title)
	}
	if page.Layout != "pages" {
		t.Errorf("expected page layout to be 'pages', got %v", page.Layout)
	}
}
