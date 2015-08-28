package main

import (
	"hash/crc32"
	"strconv"
)

type FeatureFlag struct {
	Key        string   `json:"key"`
	Enabled    bool     `json:"enabled"`
	Users      []uint32 `json:"users"`
	Groups     []string `json:"groups"`
	Percentage uint32   `json:"percentage"`
}

func (f *FeatureFlag) isEnabled() bool {
	return f.Enabled
}

func (f *FeatureFlag) isPartiallyEnabled() bool {
	return !f.Enabled && (f.hasUsers() || f.hasGroups() || f.hasPercentage())
}

func (f *FeatureFlag) groupHasAccess(group string) bool {
	return f.isEnabled() || (f.isPartiallyEnabled() && f.groupInGroups(group))
}

func (f *FeatureFlag) UserHasAccess(user uint32) bool {
	// A user has access:
	// - if the feature is enabled
	// - if the feature is partially enabled and he has been given access explicity
	// - if the feature is partially enabled and he is in the allowed percentage
	return f.isEnabled() || (f.isPartiallyEnabled() && (f.userInUsers(user) || f.userIsAllowedByPercentage(user)))
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
	return crc32.ChecksumIEEE(uint32ToBytes(user))%100 < f.Percentage
}

func (f *FeatureFlag) userInUsers(user uint32) bool {
	return intInSlice(user, f.Users)
}

func (f *FeatureFlag) groupInGroups(group string) bool {
	return stringInSlice(group, f.Groups)
}

func uint32ToBytes(u uint32) []byte {
	return []byte(strconv.FormatUint(uint64(u), 10))
}

func intInSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
