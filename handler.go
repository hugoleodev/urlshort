package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if pathToUrl, ok := pathsToUrls[r.URL.String()]; ok {
			http.Redirect(w, r, pathToUrl, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type RedirectMap struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(data []byte) ([]RedirectMap, error) {
	var redirectMaps []RedirectMap

	err := yaml.Unmarshal(data, &redirectMaps)

	if err != nil {
		return nil, err
	}

	return redirectMaps, nil
}

func buildMap(parsedYAML []RedirectMap) map[string]string {

	var pathMap = make(map[string]string)

	for _, rm := range parsedYAML {
		pathMap[rm.Path] = rm.Url
	}

	return pathMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yml)

	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}
