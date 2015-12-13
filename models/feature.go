package models

import (
	"fmt"
	"hash/crc32"
	"regexp"

	helpers "github.com/antoineaugusti/feature-flags/helpers"
)

// Represents a feature flag
type FeatureFlag struct {
	// The key of a feature flag
	Key string `json:"key"`
	// Tell if a feature flag is enabled. If set to false,
	// the feature flag can still be partially enabled thanks to
	// the Users, Groups and Percentage properties
	Enabled bool `json:"enabled"`
	// Gives access to a feature to specific user IDs
	Users []uint32 `json:"users"`
	// Gives access to a feature to specific groups
	Groups []string `json:"groups"`
	// Gives access to a feature to a percentage of users
	Percentage uint32 `json:"percentage"`
}

type FeatureFlags []FeatureFlag

// Self validate the properties of a feature flag
func (f FeatureFlag) Validate() error {
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

// Check if a feature flag is enabled
func (f FeatureFlag) IsEnabled() bool {
	return f.Enabled || f.Percentage == 100
}

// Check if a feature flag is partially enabled
func (f FeatureFlag) IsPartiallyEnabled() bool {
	return !f.IsEnabled() && (f.hasUsers() || f.hasGroups() || f.hasPercentage())
}

// Check if a group has access to a feature
func (f FeatureFlag) GroupHasAccess(group string) bool {
	return f.IsEnabled() || (f.IsPartiallyEnabled() && f.groupInGroups(group))
}

// Check if a user has access to a feature
func (f FeatureFlag) UserHasAccess(user uint32) bool {
	// A user has access:
	// - if the feature is enabled
	// - if the feature is partially enabled and he has been given access explicity
	// - if the feature is partially enabled and he is in the allowed percentage
	return f.IsEnabled() || (f.IsPartiallyEnabled() && (f.userInUsers(user) || f.userIsAllowedByPercentage(user)))
}

// Tell if specific users have access to the feature
func (f FeatureFlag) hasUsers() bool {
	return len(f.Users) > 0
}

// Tell if specific groups have access to the feature
func (f FeatureFlag) hasGroups() bool {
	return len(f.Groups) > 0
}

// Tell if a specific percentage of users has access to the feature
func (f FeatureFlag) hasPercentage() bool {
	return f.Percentage > 0
}

// Check if a user has access to the feature thanks to the percentage value
func (f FeatureFlag) userIsAllowedByPercentage(user uint32) bool {
	return crc32.ChecksumIEEE(helpers.Uint32ToBytes(user))%100 < f.Percentage
}

// Check if a user is in the list of allowed users
func (f FeatureFlag) userInUsers(user uint32) bool {
	return helpers.IntInSlice(user, f.Users)
}

// Check if a group is in the list of allowed groups
func (f FeatureFlag) groupInGroups(group string) bool {
	return helpers.StringInSlice(group, f.Groups)
}
