package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mattn/go-jsonpointer"
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

func add(corpusID, docID, textID string, provider string, url string) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(corpusID+"/"+docID+"/"+textID), []byte(provider+","+url))
		return err
	})
	return err
}

func initDB(dir string) {
	if db != nil {
		return
	}
	var err error
	db, err = badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		log.Fatal(err)
	}
}

func closeDB() {
	db.Close()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	korapServer := os.Getenv("KORAP_SERVER")
	if korapServer == "" {
		korapServer = "https://korap.ids-mannheim.de"
	}

	var pluginManifest map[string]any
	json.Unmarshal([]byte(`{
		"name" : "External Provider",
		"desc" : "Buy content from an external provider",
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

	externalProvider := os.Getenv("KORAP_EXTERNAL_PROVIDER")
	if externalProvider == "" {
		externalProvider = "https://korap.ids-mannheim.de/plugin/external/"
	}
	jsonpointer.Set(pluginManifest, "/embed/0/onClick/template", externalProvider)

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
			"korapServer": korapServer,
		})
	})

	// Return provider information
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

	initDB("db")
	defer closeDB()

	// Index csv file
	if len(os.Args) > 1 {

		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(file)

		txn := db.NewTransaction(true)

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
		err = txn.Commit()

		if err != nil {
			log.Fatal("Unable to commit")
		}

		return
	}
	r := setupRouter()

	port := os.Getenv("KORAP_EXTERNAL_PROVIDER_PORT")
	if port == "" {
		port = "5722"
	}

	log.Fatal(http.ListenAndServe(":"+port, r))
}
