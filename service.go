package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

var db *badger.DB

func CheckSaleUrl(c *gin.Context) {

	corpusID := c.Param("corpus_id")
	docID := c.Param("doc_id")
	textID := c.Param("text_id")

	c.Header("Access-Control-Allow-Origin", "null")
	c.Header("Access-Control-Allow-Credentials", "null")
	c.Header("Vary", "Origin")

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

	//
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "null")
		c.Header("Access-Control-Allow-Credentials", "null")
		c.Header("Vary", "Origin")
		c.HTML(http.StatusOK, "main.html", gin.H{
			"korapServer": "https://korap.ids-mannheim.de/instance/test",
		})
	})

	r.HEAD("/:corpus_id/:doc_id/:text_id", CheckSaleUrl)
	r.GET("/:corpus_id/:doc_id/:text_id", CheckSaleUrl)
	r.Static("/assets", "./assets")
	return r
}

func main() {
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
	log.Fatal(http.ListenAndServe(":5722", r))
}
