package db

import (
	"log"

	"github.com/boltdb/bolt"
)

func GetBucketName() string {
	return "features"
}

func GenerateDefaultBucket(name string, db *bolt.DB) {
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
