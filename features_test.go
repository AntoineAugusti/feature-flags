package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestGetFeaturesEmpty(t *testing.T) {
	db := getTestDB()
	defer closeDB(db)
	features, err := getService(db).GetFeatures()

	// No features
	assert.Equal(t, len(features), 0)
	assert.Nil(t, err)
}

func TestAddFeature(t *testing.T) {
	db := getTestDB()
	defer closeDB(db)
	features, _ := getService(db).GetFeatures()

	// No features
	assert.Equal(t, len(features), 0)

	// Create a new feature
	err := getService(db).AddFeature(getDummyFeature())
	assert.Nil(t, err)

	// We can get the feature
	features, _ = getService(db).GetFeatures()
	assert.Equal(t, len(features), 1)
	assert.Equal(t, features[0].Key, "foo")

	// I cannot add a feature with the same key
	err = getService(db).AddFeature(getDummyFeature())
	assert.Equal(t, err.Error(), "Feature already exists")
}

func TestGetFeature(t *testing.T) {
	db := getTestDB()
	defer closeDB(db)

	// Create a new feature
	_ = getService(db).AddFeature(getDummyFeature())

	// Get an existing feature
	f, err := getService(db).GetFeature("foo")
	assert.Nil(t, err)
	assert.Equal(t, f.Key, "foo")

	// Try to find an unexisting feature
	f, err = getService(db).GetFeature("bar")
	assert.Equal(t, err.Error(), "Unable to find feature")
	assert.Equal(t, len(f.Key), 0)
}

func TestUpdateFeature(t *testing.T) {
	db := getTestDB()
	defer closeDB(db)

	// Create a new feature
	_ = getService(db).AddFeature(getDummyFeature())

	newFeature := getDummyFeature()
	newFeature.Enabled = true
	newFeature.Users = []uint32{1, 2}
	newFeature.Groups = []string{"c", "d"}
	newFeature.Percentage = uint32(22)

	// Update the feature
	f, err := getService(db).UpdateFeature(newFeature.Key, newFeature)
	assert.Nil(t, err)
	assert.Equal(t, f.Enabled, true)
	assert.Equal(t, f.Users, []uint32{1, 2})
	assert.Equal(t, f.Groups, []string{"c", "d"})
	assert.Equal(t, f.Percentage, uint32(22))

	// Update an unexisting feature
	_, err = getService(db).UpdateFeature("bar", newFeature)
	assert.NotNil(t, err)
}

func TestRemoveFeature(t *testing.T) {
	db := getTestDB()
	defer closeDB(db)

	// Create a new feature
	err := getService(db).AddFeature(getDummyFeature())
	features, _ := getService(db).GetFeatures()
	assert.Equal(t, len(features), 1)

	// Delete the feature
	err = getService(db).RemoveFeature("foo")
	features, _ = getService(db).GetFeatures()
	assert.Nil(t, err)
	assert.Equal(t, len(features), 0)
}

func TestFeatureExists(t *testing.T) {
	db := getTestDB()
	defer closeDB(db)

	// Create a new feature
	_ = getService(db).AddFeature(getDummyFeature())
	assert.Equal(t, getService(db).FeatureExists("foo"), true)

	// Delete the feature
	_ = getService(db).RemoveFeature("foo")
	assert.Equal(t, getService(db).FeatureExists("foo"), false)
}

func getService(db *bolt.DB) *FeatureService {
	return &FeatureService{db}
}

func getTestDB() *bolt.DB {
	db, err := bolt.Open(getDBPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	generateDefaultBucket(getBucketName(), db)

	return db
}

func getDBPath() string {
	return "/tmp/test.db"
}

func getDummyFeature() FeatureFlag {
	return FeatureFlag{
		Key:        "foo",
		Enabled:    false,
		Users:      []uint32{22},
		Groups:     []string{},
		Percentage: 42,
	}
}

func closeDB(db *bolt.DB) {
	if err := os.Remove(getDBPath()); err != nil {
		panic(err)
	}
	db.Close()
}
