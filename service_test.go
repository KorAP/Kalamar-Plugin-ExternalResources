package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappingRoute(t *testing.T) {

	dir := t.TempDir()

	InitDB(dir)
	defer closeDB()
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/s10/s10/s10", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "No entry found", w.Body.String())

	assert.Nil(t, add(db, "s11", "s12", "s13", "sueddeutsche", "http://example.org"))

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

func TestIndexer(t *testing.T) {

	dir := t.TempDir()

	dbx := initDB(dir)
	defer dbx.Close()

	// Test index plain
	file, err := os.Open("testdata/sz_mapping_example1.csv")
	assert.Nil(t, err)
	defer file.Close()
	indexDB(file, dbx)

	// Test index gzip
	file, err = os.Open("testdata/sz_mapping_example2.csv.gz")
	assert.Nil(t, err)
	defer file.Close()
	var gzipr io.Reader
	gzipr, err = gzip.NewReader(file)
	assert.Nil(t, err)
	indexDB(gzipr, dbx)

	txn := dbx.NewTransaction(true)
	defer txn.Discard()

	item, err := txn.Get([]byte("U92/JAN/00001"))
	assert.Nil(t, err)
	err = item.Value(func(val []byte) error {
		assert.Equal(t, string(val), "S&uuml;ddeutsche Zeitung,https://archiv.szarchiv.de/Portal/restricted/Start.act?articleId=A800000")
		return nil
	})
	assert.Nil(t, err)

	item, err = txn.Get([]byte("U92/JAN/00003"))
	assert.Nil(t, err)
	err = item.Value(func(val []byte) error {
		assert.Equal(t, string(val), "S&uuml;ddeutsche Zeitung,https://archiv.szarchiv.de/Portal/restricted/Start.act?articleId=A800010")
		return nil
	})
	assert.Nil(t, err)

	item, err = txn.Get([]byte("U92/FEB/00003"))
	assert.Nil(t, err)
	err = item.Value(func(val []byte) error {
		assert.Equal(t, string(val), "S&uuml;ddeutsche Zeitung,https://archiv.szarchiv.de/Portal/restricted/Start.act?articleId=A806912")
		return nil
	})
	assert.Nil(t, err)

}
