package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	db "github.com/antoineaugusti/golang-feature-flags/db"
	h "github.com/antoineaugusti/golang-feature-flags/http"
	s "github.com/antoineaugusti/golang-feature-flags/services"
	"github.com/boltdb/bolt"
)

func main() {
	address := flag.String("a", ":8080", "address to listen")
	boltLocation := flag.String("d", "bolt.db", "location of the database file")
	flag.Parse()

	// Open the DB connection
	database, err := bolt.Open(*boltLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// Close the DB connection on exit
	defer database.Close()

	// Generate the default bucket
	db.GenerateDefaultBucket(db.GetBucketName(), database)

	api := h.APIHandler{s.FeatureService{database}}

	// Create and listen for the HTTP server
	router := h.NewRouter(&api)
	log.Fatal(http.ListenAndServe(*address, router))
}
