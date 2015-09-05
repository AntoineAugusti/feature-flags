package repos

import (
	"encoding/json"
	"fmt"

	db "github.com/antoineaugusti/feature-flags/db"
	m "github.com/antoineaugusti/feature-flags/models"
	"github.com/boltdb/bolt"
)

func PutFeature(tx *bolt.Tx, feature m.FeatureFlag) error {
	features := tx.Bucket([]byte(db.GetBucketName()))

	bytes, err := json.Marshal(feature)
	if err != nil {
		return err
	}

	err = features.Put([]byte(feature.Key), bytes)
	if err != nil {
		return err
	}

	return nil
}

func GetFeatures(tx *bolt.Tx) (m.FeatureFlags, error) {
	featuresBucket := tx.Bucket([]byte(db.GetBucketName()))
	cursor := featuresBucket.Cursor()

	features := make(m.FeatureFlags, 0)

	for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
		feature := m.FeatureFlag{}

		err := json.Unmarshal(value, &feature)
		if err != nil {
			return nil, err
		}
		features = append(features, feature)
	}

	return features, nil
}

func FeatureExists(tx *bolt.Tx, featureKey string) bool {
	features := tx.Bucket([]byte(db.GetBucketName()))
	bytes := features.Get([]byte(featureKey))
	return bytes != nil
}

func GetFeature(tx *bolt.Tx, featureKey string) (m.FeatureFlag, error) {
	features := tx.Bucket([]byte(db.GetBucketName()))

	bytes := features.Get([]byte(featureKey))
	if bytes == nil {
		return m.FeatureFlag{}, fmt.Errorf("Unable to find feature")
	}

	feature := m.FeatureFlag{}

	err := json.Unmarshal(bytes, &feature)
	if err != nil {
		return m.FeatureFlag{}, err
	}

	return feature, nil
}

func RemoveFeature(tx *bolt.Tx, featureKey string) error {
	features := tx.Bucket([]byte(db.GetBucketName()))
	return features.Delete([]byte(featureKey))
}
