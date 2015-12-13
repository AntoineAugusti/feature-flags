package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Bind routes to handlers and create a router
func NewRouter(api APIHandler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range getRoutes(api) {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
