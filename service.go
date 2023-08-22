package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v3"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mattn/go-jsonpointer"
	"golang.org/x/text/language"
)

var db *badger.DB

func CheckSaleUrl(c *gin.Context) {

	corpusID := c.Param("corpus_id")
	docID := c.Param("doc_id")
	textID := c.Param("text_id")

	err := db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(corpusID + "/" + docID + "/" + textID))

		if err != nil {
			c.String(http.StatusNotFound, "No entry found")
			return nil
		}

		err = item.Value(func(v []byte) error {
			c.String(http.StatusOK, string(v))
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.String(http.StatusNotFound, err.Error())
	}
}

func add(dbx *badger.DB, corpusID, docID, textID string, provider string, url string) error {
	err := dbx.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(corpusID+"/"+docID+"/"+textID), []byte(provider+","+url))
		return err
	})
	return err
}

func InitDB(dir string) {
	if db != nil {
		return
	}
	db = initDB(dir)
}

func initDB(dir string) *badger.DB {
	dbx, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		log.Fatal(err)
	}
	return dbx
}

func closeDB() {
	db.Close()
}

func IndexDB(ri io.Reader) error {
	return indexDB(ri, db)
}

// indexDB reads in a csv file and adds
// information to the database
func indexDB(ri io.Reader, dbx *badger.DB) error {

	r := csv.NewReader(ri)

	txn := dbx.NewTransaction(true)

	i := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if err := txn.Set([]byte(record[0]), []byte(record[1]+","+record[2])); err == badger.ErrTxnTooBig {
			log.Println("Commit", record[0], "after", i, "inserts")
			i = 0
			err = txn.Commit()
			if err != nil {
				log.Fatal("Unable to commit")
			}
			txn = db.NewTransaction(true)
			_ = txn.Set([]byte(record[0]), []byte(record[1]+","+record[2]))
		}
		i++
	}
	return txn.Commit()
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// apply i18n middleware
	r.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		RootPath:         "./i18n",
		AcceptLanguage:   []language.Tag{language.German, language.English},
		DefaultLanguage:  language.English,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	})))

	r.LoadHTMLGlob("templates/*")

	korapServer := os.Getenv("KORAP_SERVER")
	if korapServer == "" {
		korapServer = "https://korap.ids-mannheim.de"
	}

	var pluginManifest map[string]any
	json.Unmarshal([]byte(`{
		"name" : "External Resources",
		"desc" : "Retrieve content from an external provider",
		"embed" : [{
			"panel" : "match",
			"title" : "Full Text",
			"classes" : ["plugin", "cart"],
			"icon" : "\f07a",
			"onClick" : {
				"action" : "addWidget",
				"template":"",
				"permissions": [
					"scripts",
					"popups" 
				]
			}
		}]
	}`), &pluginManifest)

	externalResources := os.Getenv("KORAP_EXTERNAL_RESOURCES")
	if externalResources == "" {
		externalResources = "https://korap.ids-mannheim.de/plugin/external/"
	}
	jsonpointer.Set(pluginManifest, "/embed/0/onClick/template", externalResources)

	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", "null")
			h.Set("Access-Control-Allow-Credentials", "null")
			h.Set("Vary", "Origin")
		}
	}(),
	)

	// Return widget page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", gin.H{
			"korapServer":  korapServer,
			"title":        ginI18n.MustGetMessage(c, "fulltext"),
			"noAccess":     ginI18n.MustGetMessage(c, "noAccess"),
			"fromProvider": ginI18n.MustGetMessage(c, "fromProvider"),
		})
	})

	// Return resource information
	r.HEAD("/:corpus_id/:doc_id/:text_id", CheckSaleUrl)
	r.GET("/:corpus_id/:doc_id/:text_id", CheckSaleUrl)

	// Return plugin manifest
	r.GET("/plugin.json", func(c *gin.Context) {
		c.JSON(200, pluginManifest)
	})

	return r
}

func main() {
	if godotenv.Load() != nil {
		log.Println(".env file not loaded.")
	}

	InitDB("db")
	defer closeDB()

	// Index csv file
	if len(os.Args) > 1 {

		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		fileExt := filepath.Ext(os.Args[1])

		if fileExt == ".gz" || fileExt == ".csvz" {
			var gzipr io.Reader
			gzipr, err = gzip.NewReader(file)
			if err != nil {
				log.Fatal("Unable to open gzip file")
			} else {
				err = IndexDB(gzipr)
			}
		} else {
			err = IndexDB(file)
		}

		if err != nil {
			log.Fatal("Unable to commit")
		}
	}

	r := setupRouter()

	port := os.Getenv("KORAP_EXTERNAL_RESOURCES_PORT")
	if port == "" {
		port = "5722"
	}

	log.Println("Starting server on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
