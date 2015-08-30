package main

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

func putFeature(tx *bolt.Tx, feature FeatureFlag) error {
	features := tx.Bucket([]byte(getBucketName()))

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

func getFeatures(tx *bolt.Tx) (FeatureFlags, error) {
	featuresBucket := tx.Bucket([]byte(getBucketName()))
	cursor := featuresBucket.Cursor()

	features := make(FeatureFlags, 0)

	for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
		feature := FeatureFlag{}

		err := json.Unmarshal(value, &feature)
		if err != nil {
			return nil, err
		}
		features = append(features, feature)
	}

	return features, nil
}

func featureExists(tx *bolt.Tx, featureKey string) bool {
	features := tx.Bucket([]byte(getBucketName()))
	bytes := features.Get([]byte(featureKey))
	return bytes != nil
}

func getFeature(tx *bolt.Tx, featureKey string) (FeatureFlag, error) {
	features := tx.Bucket([]byte(getBucketName()))

	bytes := features.Get([]byte(featureKey))
	if bytes == nil {
		return FeatureFlag{}, fmt.Errorf("Unable to find feature")
	}

	feature := FeatureFlag{}

	err := json.Unmarshal(bytes, &feature)
	if err != nil {
		return FeatureFlag{}, err
	}

	return feature, nil
}

func removeFeature(tx *bolt.Tx, featureKey string) error {
	features := tx.Bucket([]byte(getBucketName()))
	return features.Delete([]byte(featureKey))
}
