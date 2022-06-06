package urlshort

import (
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
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

func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls := parseYAML(yamlData)
	pathsToUrls := buildMap(urls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYAML(yamlData []byte) []pathUrl {
	var pathUrls []pathUrl

	err := yaml.Unmarshal(yamlData, &pathUrls)

	if err != nil {
		panic(err)
	}

	return pathUrls
}

func buildMap(urls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)

	for _, path := range urls {
		pathsToUrls[path.Path] = path.URL
	}

	return pathsToUrls
}

type pathUrl struct {
	Path string `yaml:"path,omitempty"`
	URL  string `yaml:"url,omitempty"`
}
