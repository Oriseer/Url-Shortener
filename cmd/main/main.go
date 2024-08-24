package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	urlshortener "github.com/Oriseer/go_url_short/internal/urlShortener"
)

func main() {
	mux := defaultMux()
	yamlFile := flag.String("yf", "urlpaths.yaml", "Yaml File containing the urlpaths")
	jsonFile := flag.String("jf", "urlpaths.json", "JSON File containing the urlpaths")

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshortener.MapHandle(pathsToUrls, mux)

	yml, err := readYaml(*yamlFile)

	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := readJSON(*jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	yamlHandler, err := urlshortener.Yamlhandler(yml, mapHandler)

	if err != nil {
		log.Fatal(err)
		return
	}

	jsonHandler, err := urlshortener.JSONHandler(jsonData, yamlHandler)

	if err != nil {
		log.Fatal(err)
		return
	}

	dbHandler, err := urlshortener.DBHandler(jsonHandler)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Server starting.....")
	err = http.ListenAndServe("localhost:8080", dbHandler)

	if err != nil {
		log.Fatal(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello world")
}

func readYaml(y string) ([]byte, error) {
	ymlFile, err := os.ReadFile(y)

	if err != nil {
		return nil, err
	}

	return []byte(ymlFile), nil

}

func readJSON(j string) ([]byte, error) {
	jsonFile, err := os.ReadFile(j)

	if err != nil {
		return nil, err
	}

	return jsonFile, nil
}
