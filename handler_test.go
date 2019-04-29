package urlshort_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hugoleodev/urlshort"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func checkRedirect(t *testing.T, rr *httptest.ResponseRecorder, url string) {
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPermanentRedirect)
	}

	if location := rr.Header().Get("Location"); location != url {
		t.Errorf("handler returned wrong header location: got %v want %v", location, url)
	}
}

func checkHandlerPathMap(t *testing.T, path string, url string, handler http.Handler) {

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	checkRedirect(t, rr, url)
}

func TestMapHandler(t *testing.T) {

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	handler := urlshort.MapHandler(pathsToUrls, defaultMux())

	for path, url := range pathsToUrls {
		checkHandlerPathMap(t, path, url, handler)
	}

	req, err := http.NewRequest("GET", "/nonexist", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestYAMLHandler(t *testing.T) {

	yaml := `
- path: /google
  url: https://www.google.com
- path: /facebook
  url: https://www.facebook.com
- path: /twitter
  url: https://www.twitter.com
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), defaultMux())

	if err != nil {
		t.Fatal(err)
	}

	pathsToUrls := map[string]string{
		"/google":   "https://www.google.com",
		"/facebook": "https://www.facebook.com",
		"/twitter":  "https://www.twitter.com",
	}

	for path, url := range pathsToUrls {
		checkHandlerPathMap(t, path, url, yamlHandler)
	}
}
