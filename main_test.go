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
	if len(app.PageTemplate) == 0 {
		t.Errorf("expected PageTemplate to be populated")
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

func TestGetAllPages(t *testing.T) {
	srcTest := "src_test"
	distTest := "dist_test"
	defer cleanup(distTest)
	app, _ := InitApp(srcTest, distTest)
	pages, err := app.getAllPages()
	if err != nil {
		t.Errorf("expected getAllPages to return no error, got %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %v", len(pages))
	}
}

func TestRenderPage(t *testing.T) {
	srcTest := "src_test"
	distTest := "dist_test"
	defer cleanup(distTest)
	app, _ := InitApp(srcTest, distTest)
	pages, _ := app.getAllPages()
	for _, page := range pages {
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
	_, err = os.Stat(filepath.Join(distTest, "example.html"))
	if err != nil {
		t.Errorf("expected example.html to exist, got %v", err)
	}
}
