package main

import "testing"

func TestEnabled(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    true,
		Users:      []uint32{},
		Groups:     []string{},
		Percentage: 20,
	}
	if !f.IsEnabled() {
		t.Fatalf("Feature should be enabled")
	}

	if f.IsPartiallyEnabled() {
		t.Fatalf("Enabled feature should not be partially enabled")
	}

	// Disable the feature
	f.Enabled = false
	if f.IsEnabled() {
		t.Fatalf("Feature should be disabled")
	}
}

func TestPartiallyEnabled(t *testing.T) {
	f := FeatureFlag{
		Key:        "foo",
		Enabled:    false,
		Users:      []uint32{},
		Groups:     []string{},
		Percentage: 20,
	}

	if !f.IsPartiallyEnabled() {
		t.Fatalf("Feature should be partially enabled because of the percentage")
	}
	f.Percentage = 0

	f.Groups = []string{"a"}
	if !f.IsPartiallyEnabled() {
		t.Fatalf("Feature should be partially enabled because of the groups")
	}
	f.Groups = []string{}

	f.Users = []uint32{22}
	if !f.IsPartiallyEnabled() {
		t.Fatalf("Feature should be partially enabled because of the users")
	}
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
	if f.IsEnabled() {
		t.Fatalf("Feature should not be enabled")
	}

	if !f.GroupHasAccess("bar") {
		t.Fatalf("Group bar should have access")
	}
	if f.GroupHasAccess("baz") {
		t.Fatalf("Group baz should not have access")
	}

	f.Groups = []string{"bar", "baz"}
	if !f.GroupHasAccess("baz") {
		t.Fatalf("Group baz should have access")
	}

	f.Enabled = true
	if !f.GroupHasAccess("klm") {
		t.Fatalf("Group klm should have access")
	}

	f.Groups = []string{}
	f.Percentage = 100
	f.Enabled = false
	if !f.GroupHasAccess("test") {
		t.Fatalf("Group test should have access")
	}
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
	if f.IsEnabled() {
		t.Fatalf("Feature should not be enabled")
	}

	if !f.UserHasAccess(42) {
		t.Fatalf("User 42 should have access")
	}
	if f.UserHasAccess(1337) {
		t.Fatalf("User 1337 should not have access")
	}

	f.Users = []uint32{42, 1337}
	if !f.UserHasAccess(1337) {
		t.Fatalf("User 1337 should have access")
	}

	f.Enabled = true
	if !f.UserHasAccess(222) {
		t.Fatalf("User 222 should have access")
	}

	f.Users = []uint32{}
	f.Percentage = 100
	f.Enabled = false
	if !f.UserHasAccess(222) {
		t.Fatalf("User 222 should have access")
	}
}
