package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

func main() {
	address := flag.String("a", ":8080", "address to listen")
	boltLocation := flag.String("d", "bolt.db", "location of the database file")
	flag.Parse()

	// Open the DB connection
	db, err := bolt.Open(*boltLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// Close the DB connection on exit
	defer db.Close()

	// Generate the default bucket
	generateDefaultBucket(getBucketName(), db)

	api := APIHandler{FeatureService{db}}

	// Create and listen for the HTTP server
	router := NewRouter(&api)
	log.Fatal(http.ListenAndServe(*address, router))
}

func getBucketName() string {
	return "features"
}

func generateDefaultBucket(name string, db *bolt.DB) {
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
