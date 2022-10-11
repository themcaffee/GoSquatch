package main

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	ping(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "pong" {
		t.Fatalf("expected pong, got %s", string(data))
	}
}

func TestGetLivePageIndex(t *testing.T) {
	srcDir := "src_test"
	defer cleanup("dist")
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	app, err := InitApp(srcDir)
	if err != nil {
		t.Fatal(err)
	}
	Build(srcDir)
	app.getLivePage(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if len(string(data)) == 0 {
		t.Fatalf("expected index, got %s", string(data))
	}
}

func TestGetLivePage(t *testing.T) {
	srcDir := "src_test"
	defer cleanup("dist")
	req := httptest.NewRequest("GET", "/pages/example", nil)
	w := httptest.NewRecorder()
	app, err := InitApp(srcDir)
	if err != nil {
		t.Fatal(err)
	}
	Build(srcDir)
	app.getLivePage(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if len(string(data)) == 0 {
		t.Fatalf("expected about, got %s", string(data))
	}
}
