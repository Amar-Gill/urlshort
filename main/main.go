package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Amar-Gill/urlshort"
)

func main() {
	yamlFile := flag.String("yaml", "default.yaml", "a yaml file in the form:\n  - path: ...\n    url: ...\n")
	jsonFile := flag.String("json", "default.json", "a json file in the form [ { \"path\": ..., \"url\": ... } ]")
	flag.Parse()

	yamlData, err := os.ReadFile(*yamlFile)
	if err != nil {
		panic(err)
	}

	jsonData, err := os.ReadFile(*jsonFile)
	if err != nil {
		panic(err)
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yamlHandler, err := urlshort.YAMLHandler(yamlData, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonData, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
