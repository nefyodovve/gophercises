package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort/urlshort"
	"net/http"
	"os"
)

var yamlPath string
var jsonPath string

func main() {
	handleFlags()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	if len(yamlPath) != 0 {
		yaml, err := os.ReadFile(yamlPath)
		if err != nil {
			panic(err)
		}
		handler, err = urlshort.YAMLHandler(yaml, handler)
		if err != nil {
			panic(err)
		}
	}
	if len(jsonPath) != 0 {
		json, err := os.ReadFile(jsonPath)
		if err != nil {
			panic(err)
		}
		handler, err = urlshort.JSONHandler(json, handler)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func handleFlags() {
	flag.StringVar(&yamlPath, "yaml", "", "yaml file with list of shortened urls")
	flag.StringVar(&jsonPath, "json", "", "json file with list of shortened urls")
	flag.Parse()
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
