package main

import (
	"log"
	"net/http"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
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
		}

		err = item.Value(func(v []byte) error {
			c.String(http.StatusOK, string(v))
			return nil
		})

		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}

		return nil
	})

	if err != nil {
		c.String(http.StatusNotFound, err.Error())
	}

	c.String(http.StatusExpectationFailed, err.Error())
}

func add(corpusID, docID, textID string, url string) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(corpusID+"/"+docID+"/"+textID), []byte(url))
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
	r.GET("/:corpus_id/:doc_id/:text_id", CheckSaleUrl)
	return r
}

func main() {
	initDB("db")
	r := setupRouter()
	r.Run(":8080")
	defer closeDB()
}
