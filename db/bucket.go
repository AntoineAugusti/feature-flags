package db

import (
	"log"

	"github.com/boltdb/bolt"
)

// Get the name of the bucket
func GetBucketName() string {
	return "features"
}

// Generate the default bucket if it does not exist yet
func GenerateDefaultBucket(name string, db *bolt.DB) {
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
