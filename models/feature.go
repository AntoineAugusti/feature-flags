package models

import (
	"fmt"
	"hash/crc32"
	"regexp"

	helpers "github.com/antoineaugusti/feature-flags/helpers"
)

type FeatureFlag struct {
	Key        string   `json:"key"`
	Enabled    bool     `json:"enabled"`
	Users      []uint32 `json:"users"`
	Groups     []string `json:"groups"`
	Percentage uint32   `json:"percentage"`
}

type FeatureFlags []FeatureFlag

func (f *FeatureFlag) Validate() error {
	// Validate percentage
	if f.Percentage < 0 || f.Percentage > 100 {
		return fmt.Errorf("Percentage must be between 0 and 100")
	}

	// Validate key
	if len(f.Key) < 3 || len(f.Key) > 50 {
		return fmt.Errorf("Feature key must be between 3 and 50 characters")
	}

	if !regexp.MustCompile(`^[a-z0-9_]*$`).MatchString(f.Key) {
		return fmt.Errorf("Feature key must only contain digits, lowercase letters and underscores")
	}
	return nil
}

func (f *FeatureFlag) IsEnabled() bool {
	return f.Enabled
}

func (f *FeatureFlag) IsPartiallyEnabled() bool {
	return !f.IsEnabled() && (f.hasUsers() || f.hasGroups() || f.hasPercentage())
}

func (f *FeatureFlag) GroupHasAccess(group string) bool {
	return f.IsEnabled() || f.Percentage == 100 || (f.IsPartiallyEnabled() && f.groupInGroups(group))
}

func (f *FeatureFlag) UserHasAccess(user uint32) bool {
	// A user has access:
	// - if the feature is enabled
	// - if the feature is partially enabled and he has been given access explicity
	// - if the feature is partially enabled and he is in the allowed percentage
	return f.IsEnabled() || (f.IsPartiallyEnabled() && (f.userInUsers(user) || f.userIsAllowedByPercentage(user)))
}

func (f *FeatureFlag) hasUsers() bool {
	return len(f.Users) > 0
}

func (f *FeatureFlag) hasGroups() bool {
	return len(f.Groups) > 0
}

func (f *FeatureFlag) hasPercentage() bool {
	return f.Percentage > 0
}

func (f *FeatureFlag) userIsAllowedByPercentage(user uint32) bool {
	return crc32.ChecksumIEEE(helpers.Uint32ToBytes(user))%100 < f.Percentage
}

func (f *FeatureFlag) userInUsers(user uint32) bool {
	return helpers.IntInSlice(user, f.Users)
}

func (f *FeatureFlag) groupInGroups(group string) bool {
	return helpers.StringInSlice(group, f.Groups)
}
