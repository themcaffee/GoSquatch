package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func filewatch(srcDir string) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go watchLoop(watcher, srcDir)

	// Add a path.
	err = watcher.Add(srcDir)
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{}) // Block forever
}

func watchLoop(w *fsnotify.Watcher, srcDir string) {
	var (
		// Wait 100ms for new events; each new event resets the timer.
		waitFor = 100 * time.Millisecond

		// Keep track of the timers, as path → timer.
		mu     sync.Mutex
		timers = make(map[string]*time.Timer)

		// Callback we run.
		buildEvent = func(e fsnotify.Event) {
			Build(srcDir)

			// Don't need to remove the timer if you don't have a lot of files.
			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Printf("error: %v", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Printf("event: %v", e)

			// Get timer.
			mu.Lock()
			t, ok := timers[e.Name]
			mu.Unlock()

			// No timer yet, so create one.
			if !ok {
				t = time.AfterFunc(math.MaxInt64, func() { buildEvent(e) })
				t.Stop()

				mu.Lock()
				timers[e.Name] = t
				mu.Unlock()
			}

			// Reset the timer for this path, so it will start from 100ms again.
			t.Reset(waitFor)
		}
	}
}

func checkLive(w http.ResponseWriter, err error) {
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func (app App) getLivePage(w http.ResponseWriter, r *http.Request) {
	// get the path from the url
	cleanPath := filepath.Clean(r.URL.Path)
	fp := filepath.Join(app.DistDir, cleanPath+".html")

	// return the index page if the path is empty
	if cleanPath == "/" {
		fp = filepath.Join(app.DistDir, "index.html")
	}

	// 404 if the requested page doesn't exist
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 404 if the requested page is a directory
	if info, err := os.Stat(fp); err == nil && info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// read the markdown file
	fileData, err := os.ReadFile(fp)
	checkLive(w, err)

	w.Header().Set("Content-Type", "text/html")
	w.Write(fileData)
}

func ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "pong")
}

func LiveServer(srcDir string, port string) {
	app, err := InitApp(srcDir)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	go filewatch(app.SrcDir)
	Build(app.SrcDir)

	// serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// serve pages
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.getLivePage)
	mux.HandleFunc("/ping", ping)
	err = http.ListenAndServe(":"+port, mux)

	// handle server closing
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
