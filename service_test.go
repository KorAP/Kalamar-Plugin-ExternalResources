package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappingRoute(t *testing.T) {

	dir := t.TempDir()

	initDB(dir)
	defer closeDB()
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/s10/s10/s10", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "No entry found", w.Body.String())

	assert.Nil(t, add("s11", "s12", "s13", "sueddeutsche", "http://example.org"))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/s11/s12/s13", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "null")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "null")
	assert.Equal(t, w.Header().Get("Vary"), "Origin")
	assert.Equal(t, "sueddeutsche,http://example.org", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/s11/s12/s14", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "No entry found", w.Body.String())
}

func TestWidgetRoute(t *testing.T) {

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "null")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "null")
	assert.Equal(t, w.Header().Get("Vary"), "Origin")
	assert.Contains(t, w.Body.String(), "data-server=\"https://korap.ids-mannheim.de\"")
	assert.Contains(t, w.Body.String(), "<title>External Resources</title>")

	os.Setenv("KORAP_SERVER", "https://korap.ids-mannheim.de/instance/test")

	router = setupRouter()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data-server=\"https://korap.ids-mannheim.de/instance/test\"")
	assert.Contains(t, w.Body.String(), "<title>External Resources</title>")
}

func TestManifestRoute(t *testing.T) {

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/plugin.json", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, w.Body.String(), "permissions")
	assert.Contains(t, w.Body.String(), "/plugin/external")

	os.Setenv("KORAP_EXTERNAL_RESOURCES", "https://korap.ids-mannheim.de/plugin/fun")

	router = setupRouter()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/plugin.json", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, w.Body.String(), "permissions")
	assert.Contains(t, w.Body.String(), "/plugin/fun")

}