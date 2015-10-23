package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnabled(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    true,
		Users:      []uint32{},
		Groups:     []string{},
		Percentage: 20,
	}

	assert.True(t, f.IsEnabled())
	assert.False(t, f.IsPartiallyEnabled())

	// Disable the feature
	f.Enabled = false

	assert.False(t, f.IsEnabled())
	assert.True(t, f.IsPartiallyEnabled())
}

func TestValidate(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    false,
		Users:      []uint32{},
		Groups:     []string{},
		Percentage: 101,
	}

	err := f.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "Percentage must be between 0 and 100", err.Error())

	f.Percentage = 50
	f.Key = "ab"
	err = f.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "Feature key must be between 3 and 50 characters", err.Error())

	f.Key = "a&6"
	err = f.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "Feature key must only contain digits, lowercase letters and underscores", err.Error())

	f.Key = "foo"
	assert.Nil(t, f.Validate())
}

func TestPartiallyEnabled(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    false,
		Users:      []uint32{},
		Groups:     []string{},
		Percentage: 20,
	}

	assert.True(t, f.IsPartiallyEnabled())

	f.Percentage = 0
	f.Groups = []string{"a"}
	assert.True(t, f.IsPartiallyEnabled())

	f.Groups = []string{}
	f.Users = []uint32{22}
	assert.True(t, f.IsPartiallyEnabled())

	f.Percentage = 100
	assert.False(t, f.IsPartiallyEnabled())
	assert.True(t, f.IsEnabled())
}

func TestGroupHasAccess(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    false,
		Users:      []uint32{42},
		Groups:     []string{"bar"},
		Percentage: 20,
	}
	// Make sure the feature is not enabled
	assert.False(t, f.IsEnabled())

	assert.True(t, f.GroupHasAccess("bar"))
	assert.False(t, f.GroupHasAccess("baz"))

	f.Groups = []string{"bar", "baz"}
	assert.True(t, f.GroupHasAccess("baz"))

	f.Enabled = true
	assert.True(t, f.GroupHasAccess("klm"))

	f.Groups = []string{}
	f.Percentage = 100
	f.Enabled = false
	assert.True(t, f.GroupHasAccess("test"))
}

func TestUserHasAccess(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    false,
		Users:      []uint32{42},
		Groups:     []string{},
		Percentage: 20,
	}
	// Make sure the feature is not enabled
	assert.False(t, f.IsEnabled())

	assert.True(t, f.UserHasAccess(42))
	assert.False(t, f.UserHasAccess(1337))

	f.Users = []uint32{42, 1337}
	assert.True(t, f.UserHasAccess(1337))

	f.Enabled = true
	assert.True(t, f.UserHasAccess(222))

	f.Users = []uint32{}
	f.Percentage = 100
	f.Enabled = false
	assert.True(t, f.UserHasAccess(222))
}
