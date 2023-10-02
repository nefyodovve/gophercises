package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pathsToUrls != nil {
			for key, val := range pathsToUrls {
				if key == r.RequestURI {
					http.Redirect(w, r, val, http.StatusSeeOther)
				}
			}
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	if len(yml) == 0 {
		return MapHandler(nil, fallback), nil
	}
	data := []struct {
		Path string `yaml:"path"`
		Url  string `yaml:"url"`
	}{}
	err := yaml.Unmarshal(yml, &data)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, entry := range data {
		m[entry.Path] = entry.Url
	}
	return MapHandler(m, fallback), nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	if len(jsonData) == 0 {
		return MapHandler(nil, fallback), nil
	}
	data := []struct {
		Path string `json:"path"`
		Url  string `json:"url"`
	}{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, entry := range data {
		m[entry.Path] = entry.Url
	}
	return MapHandler(m, fallback), nil
}
