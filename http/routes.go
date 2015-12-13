package http

import "net/http"

type Route struct {
	// A human readable name for the route
	Name string
	// The HTTP method
	Method string
	// The URL
	Pattern string
	// The handler for this endpoint
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func getRoutes(api APIHandler) Routes {
	return Routes{
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
		// curl -H "Content-Type: application/json" -X POST -d '{"groups":"foo"}' -X GET http://localhost:8080/features/feature_test/access
		Route{
			"FeatureAccess",
			"POST",
			"/features/{featureKey}/access",
			api.FeatureAccess,
		},
		// curl -H "Content-Type: application/json" -X PATCH -d '{"percentage": 42}' http://localhost:8080/features/blah
		Route{
			"FeatureEdit",
			"PATCH",
			"/features/{featureKey}",
			api.FeatureEdit,
		},
	}
}
