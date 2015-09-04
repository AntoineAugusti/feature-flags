package services

import (
	"fmt"

	m "github.com/antoineaugusti/golang-feature-flags/models"
	repos "github.com/antoineaugusti/golang-feature-flags/repos"
	"github.com/boltdb/bolt"
)

type FeatureService struct {
	DB *bolt.DB
}

func (interactor *FeatureService) AddFeature(newFeature m.FeatureFlag) error {
	return interactor.DB.Update(func(tx *bolt.Tx) error {

		feature, err := repos.GetFeature(tx, newFeature.Key)
		if err != nil && err.Error() != "Unable to find feature" {
			return err
		}

		if len(feature.Key) > 0 {
			return fmt.Errorf("Feature already exists")
		}

		err = repos.PutFeature(tx, newFeature)
		if err != nil {
			return err
		}

		return nil
	})
}

func (interactor *FeatureService) GetFeatures() (features m.FeatureFlags, err error) {
	_ = interactor.DB.View(func(tx *bolt.Tx) error {

		features, err = repos.GetFeatures(tx)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (interactor *FeatureService) GetFeature(featureKey string) (feature m.FeatureFlag, err error) {
	_ = interactor.DB.View(func(tx *bolt.Tx) error {

		feature, err = repos.GetFeature(tx, featureKey)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (interactor *FeatureService) UpdateFeature(featureKey string, newFeature m.FeatureFlag) (feature m.FeatureFlag, err error) {
	_ = interactor.DB.Update(func(tx *bolt.Tx) error {

		if feature, err = repos.GetFeature(tx, featureKey); err != nil {
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

		if err = repos.PutFeature(tx, feature); err != nil {
			return err
		}

		return nil
	})

	return
}

func (interactor *FeatureService) RemoveFeature(featureKey string) error {
	return interactor.DB.Update(func(tx *bolt.Tx) error {
		return repos.RemoveFeature(tx, featureKey)
	})
}

func (interactor *FeatureService) FeatureExists(featureKey string) (exists bool) {
	_ = interactor.DB.View(func(tx *bolt.Tx) error {
		exists = repos.FeatureExists(tx, featureKey)
		return nil
	})

	return
}
