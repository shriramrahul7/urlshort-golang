package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	urlshort "gophercises.com/url-short"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	yamlFile := flag.String("yaml", "./../paths.yaml", "specify the location of the yaml file with the paths and urls.")
	jsonFile := flag.String("json", "paths.json", "provide a json file path which would contain the paths and the URLs")
	modeOfUse := flag.String("mode", "yaml", "Specify the mode of usage. ('json' or 'yaml') ")
	flag.Parse()

	var filePath string
	mode := strings.ToLower(*modeOfUse)

	if mode == "yaml" {
		filePath = *yamlFile
	} else {
		filePath = *jsonFile
	}
	//reading the filePath from ioutil package and storing the bytes in fileBytes
	// FIXME: try to implement buffer reading so that you won't have to read all bytes in to the memory.
	fileBytes, err := ioutil.ReadFile(filePath)
	check(err)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/youtube":        "https://youtube.com",
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/urlshort":       "https://github.com/gophercises/urlshort",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	jsonHandler, _ := urlshort.JSONHandler([]byte(fileBytes), mapHandler)
	// check(err1)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `

	yamlHandler, _ := urlshort.YAMLHandler([]byte(fileBytes), mapHandler)
	// check(err2)

	var handler http.HandlerFunc
	if mode == "yaml" {
		handler = yamlHandler
	} else {
		handler = jsonHandler
	}
	fmt.Println("Starting the server on :8000")
	http.ListenAndServe(":8000", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the url shortner")
}
