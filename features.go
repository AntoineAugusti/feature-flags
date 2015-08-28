package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type FeatureService struct {
	DB *bolt.DB
}

func (interactor *FeatureService) AddFeature(newFeature FeatureFlag) error {
	return interactor.DB.Update(func(tx *bolt.Tx) error {

		feature, err := getFeature(tx, newFeature.Key)
		if err != nil && err.Error() != "Unable to find feature" {
			return err
		}

		if len(feature.Key) > 0 {
			return fmt.Errorf("Feature already exists")
		}

		err = putFeature(tx, newFeature)
		if err != nil {
			return err
		}

		return nil
	})
}

func (interactor *FeatureService) GetFeatures() (features []FeatureFlag, err error) {
	_ = interactor.DB.View(func(tx *bolt.Tx) error {

		features, err = getFeatures(tx)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (interactor *FeatureService) GetFeature(featureKey string) (feature FeatureFlag, err error) {
	_ = interactor.DB.View(func(tx *bolt.Tx) error {

		feature, err = getFeature(tx, featureKey)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (interactor *FeatureService) UpdateFeature(key string, enabled bool, users []uint32, groups []string, percentage uint32) error {
	return interactor.DB.Update(func(tx *bolt.Tx) error {

		feature, err := getFeature(tx, key)
		if err != nil {
			return err
		}

		feature.Enabled = enabled

		if len(users) > 0 {
			feature.Users = users
		}

		if len(groups) > 0 {
			feature.Groups = groups
		}

		if percentage > 0 {
			feature.Percentage = percentage
		}

		err = putFeature(tx, feature)
		if err != nil {
			return err
		}

		return nil
	})
}

func (interactor *FeatureService) RemoveFeature(featureKey string) error {
	return interactor.DB.Update(func(tx *bolt.Tx) error {
		return removeFeature(tx, featureKey)
	})
}

func (interactor *FeatureService) FeatureExists(featureKey string) (exists bool) {
	_ = interactor.DB.View(func(tx *bolt.Tx) error {
		exists = featureExists(tx, featureKey)
		return nil
	})

	return
}
