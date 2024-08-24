package urlshortener

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type pathToUrl struct {
	Path string
	Url  string
}

// Return handle the data and return as Handler
func MapHandle(pathsToUrl map[string]string, fallback http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrl[r.URL.Path]

		if ok {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}

	})
}

// Handle yaml files
func Yamlhandler(yml []byte, fallback http.Handler) (http.Handler, error) {

	parsedYaml, err := parseYaml(yml)

	if err != nil {
		return nil, err
	}

	mappedYaml := buildMap(parsedYaml)
	yamlHandler := MapHandle(mappedYaml, fallback)

	return yamlHandler, nil
}

// Map builder
func buildMap(data []pathToUrl) map[string]string {

	mapVar := make(map[string]string)

	for _, data := range data {
		mapVar[data.Path] = data.Url
	}

	return mapVar

}

// yaml parser
func parseYaml(yml []byte) ([]pathToUrl, error) {
	var pathToUrls []pathToUrl
	err := yaml.Unmarshal(yml, &pathToUrls)

	if err != nil {
		return nil, err
	}
	return pathToUrls, nil
}

// Handle JSON files
func JSONHandler(jsonData []byte, fallback http.Handler) (http.Handler, error) {

	parsedJson, err := parseJSON(jsonData)

	if err != nil {
		return nil, err
	}

	mappedJSON := buildMap(*parsedJson)
	jsonHandler := MapHandle(mappedJSON, fallback)

	return jsonHandler, nil
}

// JSON parser
func parseJSON(jsonData []byte) (*[]pathToUrl, error) {
	var pathsToUrl = []pathToUrl{}
	err := json.Unmarshal(jsonData, &pathsToUrl)

	if err != nil {
		return nil, err
	}

	return &pathsToUrl, nil

}

func DBHandler(fallback http.Handler) (http.Handler, error) {
	db, err := dbConnect()

	if err != nil {
		return nil, err
	}

	var pathsToUrl = []pathToUrl{}
	query := "SELECT * FROM CI_URLS ORDER BY PATH ASC"

	err = db.Select(&pathsToUrl, query)

	if err != nil {
		return nil, err
	}

	mappedDB := buildMap(pathsToUrl)

	dbHandler := MapHandle(mappedDB, fallback)

	return dbHandler, nil
}

func dbConnect() (*sqlx.DB, error) {
	dbDetails := "user=postgres dbname=testruns password=testpost sslmode=disable host=localhost"

	db, err := sqlx.Connect("postgres", dbDetails)

	if err != nil {
		return nil, err
	}

	return db, nil

}
