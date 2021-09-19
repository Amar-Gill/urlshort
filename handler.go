package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(byteSlice []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(byteSlice)

	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYaml)

	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(byteSlice []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(byteSlice)

	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedJSON)

	return MapHandler(pathMap, fallback), nil
}

type pathUrl struct {
	Path string
	Url  string
}

func parseYaml(byteSlice []byte) ([]pathUrl, error) {
	var data []pathUrl

	err := yaml.Unmarshal(byteSlice, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseJSON(byteSlice []byte) ([]pathUrl, error) {
	var data []pathUrl

	err := json.Unmarshal(byteSlice, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func buildMap(structSlice []pathUrl) map[string]string {
	m := make(map[string]string)

	for _, v := range structSlice {
		m[v.Path] = v.Url
	}

	return m
}
