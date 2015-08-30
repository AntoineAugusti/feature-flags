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

func (interactor *FeatureService) UpdateFeature(featureKey string, newFeature FeatureFlag) (feature FeatureFlag, err error) {
	_ = interactor.DB.Update(func(tx *bolt.Tx) error {

		if feature, err = getFeature(tx, featureKey); err != nil {
			return err
		}

		feature.Enabled = newFeature.Enabled

		if len(newFeature.Users) > 0 {
			feature.Users = newFeature.Users
		}

		if len(newFeature.Groups) > 0 {
			feature.Groups = newFeature.Groups
		}

		if newFeature.Percentage > 0 {
			feature.Percentage = newFeature.Percentage
		}

		if err = putFeature(tx, feature); err != nil {
			return err
		}

		return nil
	})

	return
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
