package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// temporary, to provide a mock server for the HTML5 client while the actual
// storage is still under development here.
const (
	exampleId = "0bc88096498be52d7909bcec815d37b3"
)

var fileRoute = regexp.MustCompile("^/files/([^/]+)$")

func serveHttp() error {
	http.HandleFunc("/", route)

	addr := ":1080"
	log.Printf("serving clients at %s", addr)

	return http.ListenAndServe(addr, nil)
}

func route(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s", r.Method, r.URL.RequestURI())

	w.Header().Set("Server", "tusd")

	if r.Method == "POST" && r.URL.Path == "/files" {
		createFile(w, r)
	} else if match := fileRoute.FindStringSubmatch(r.URL.Path); match != nil {
		switch r.Method {
		case "HEAD":
			w.Header().Set("X-Resume", "bytes=0-99")
		case "GET":
			reply(w, http.StatusNotImplemented, "File download")
		case "PUT":
			reply(w, http.StatusOK, "chunk created")
		default:
			reply(w, http.StatusMethodNotAllowed, "Invalid http method")
		}
	} else {
		reply(w, http.StatusNotFound, "No matching route")
	}
}

func reply(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, "%d - %s: %s\n", code, http.StatusText(code), message)
}

func createFile(w http.ResponseWriter, r *http.Request) {
	contentRange, err := parseContentRange(r.Header.Get("Content-Range"))
	if err != nil {
		reply(w, http.StatusBadRequest, err.Error())
		return
	}

	if contentRange.Size == -1 {
		reply(w, http.StatusBadRequest, "Content-Range must indicate total file size.")
		return
	}

	if contentRange.End != -1 {
		reply(w, http.StatusNotImplemented, "File data in initial request.")
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Location", "/files/"+exampleId)
	w.WriteHeader(http.StatusCreated)
}