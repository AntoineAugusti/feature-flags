package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func getRoutes(api *APIHandler) Routes {
	return Routes{
		Route{
			"Welcome",
			"GET",
			"/",
			api.Welcome,
		},
		Route{
			"FeatureIndex",
			"GET",
			"/features",
			api.FeatureIndex,
		},
		// curl -H "Content-Type: application/json" -X POST -d '{"key":"blah","enabled":true, "users":[22,42], "groups":["foo", "bar"], "percentage": null}' http://localhost:8080/features
		Route{
			"FeatureCreate",
			"POST",
			"/features",
			api.FeatureCreate,
		},
		// curl -X "DELETE" http://localhost:8080/features/blah
		Route{
			"FeatureRemove",
			"DELETE",
			"/features/{featureKey}",
			api.FeatureRemove,
		},
		Route{
			"FeatureShow",
			"GET",
			"/features/{featureKey}",
			api.FeatureShow,
		},
		// curl -H "Content-Type: application/json" -d '{"groups":"foo"}' -X GET http://localhost:8080/features/feature_test/access
		Route{
			"FeatureAccess",
			"GET",
			"/features/{featureKey}/access",
			api.FeatureAccess,
		},
	}
}
