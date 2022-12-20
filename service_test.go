package main

import (
	"net/http"
	"net/http/httptest"
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

	assert.Nil(t, add("s11", "s12", "s13", "http://example.org"))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/s11/s12/s13", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.org", w.Body.String())
}

func TestAssetRoute(t *testing.T) {

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/index", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data-server=\"https://korap.ids-mannheim.de/\"")
	assert.Contains(t, w.Body.String(), "<title>External Provider</title>")
}
