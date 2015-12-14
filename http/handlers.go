package http

import (
	"encoding/json"
	"net/http"

	m "github.com/antoineaugusti/feature-flags/models"
	services "github.com/antoineaugusti/feature-flags/services"
	"github.com/gorilla/mux"
)

// Handles incoming requests
type APIHandler struct {
	FeatureService services.FeatureService
}

// A simple structure to respond with error messages
type APIMessage struct {
	// The HTTP status code
	code int
	// A status message
	Status string `json:"status"`
	// A human readable message
	Message string `json:"message"`
}

// Describes the request when checking the access to a feature
type AccessRequest struct {
	Groups []string `json:"groups"`
	User   uint32   `json:"user"`
}

func (handler APIHandler) FeatureIndex(w http.ResponseWriter, r *http.Request) {
	features, err := handler.FeatureService.GetFeatures()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(features); err != nil {
		panic(err)
	}
}

func (handler APIHandler) FeatureShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the feature exists
	if !handler.featureExists(vars["featureKey"]) {
		writeNotFound(w)
		return
	}

	// Fetch the feature
	feature, err := handler.FeatureService.GetFeature(vars["featureKey"])
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(feature); err != nil {
		panic(err)
	}
}

func (handler APIHandler) FeaturesAccess(w http.ResponseWriter, r *http.Request) {
	var ar AccessRequest

	// Get all features in the bucket
	features, err := handler.FeatureService.GetFeatures()
	if err != nil {
		panic(err)
	}

	// Decode the access request
	err = json.NewDecoder(r.Body).Decode(&ar)
	if err != nil {
		writeUnprocessableEntity(err, w)
		return
	}

	// Keep only accessible features
	accessibleFeatures := make(m.FeatureFlags, 0)
	for _, feature := range features {
		if hasAccessToFeature(feature, ar) {
			accessibleFeatures = append(accessibleFeatures, feature)
		}
	}

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accessibleFeatures); err != nil {
		panic(err)
	}
}

func (handler APIHandler) FeatureAccess(w http.ResponseWriter, r *http.Request) {
	var ar AccessRequest
	vars := mux.Vars(r)

	// Check if the feature exists
	if !handler.featureExists(vars["featureKey"]) {
		writeNotFound(w)
		return
	}

	// Fetch the feature
	feature, err := handler.FeatureService.GetFeature(vars["featureKey"])
	if err != nil {
		panic(err)
	}

	// Decode the access request
	err = json.NewDecoder(r.Body).Decode(&ar)
	if err != nil {
		writeUnprocessableEntity(err, w)
		return
	}

	if hasAccessToFeature(feature, ar) {
		writeMessage(http.StatusOK, "has_access", "The user has access to the feature", w)
	} else {
		writeMessage(http.StatusOK, "not_access", "The user does not have access to the feature", w)
	}
}

func (handler APIHandler) FeatureRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the feature exists
	if !handler.featureExists(vars["featureKey"]) {
		writeNotFound(w)
		return
	}

	// Delete it
	err := handler.FeatureService.RemoveFeature(vars["featureKey"])
	if err != nil {
		panic(err)
	}

	writeMessage(http.StatusOK, "feature_deleted", "The feature was successfully deleted", w)
}

func (handler APIHandler) FeatureCreate(w http.ResponseWriter, r *http.Request) {
	var feature m.FeatureFlag

	if err := json.NewDecoder(r.Body).Decode(&feature); err != nil {
		writeUnprocessableEntity(err, w)
		return
	}

	if err := feature.Validate(); err != nil {
		writeMessage(400, "invalid_feature", err.Error(), w)
		return
	}

	err := handler.FeatureService.AddFeature(feature)
	if err != nil && err.Error() == "Feature already exists" {
		writeMessage(400, "invalid_feature", err.Error(), w)
		return
	}

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(feature); err != nil {
		panic(err)
	}
}

func (handler APIHandler) FeatureEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the feature exists
	if !handler.featureExists(vars["featureKey"]) {
		writeNotFound(w)
		return
	}

	// Fetch the feature
	newFeature, err := handler.FeatureService.GetFeature(vars["featureKey"])
	if err != nil {
		panic(err)
	}

	// Update the overwritten fields of the feature
	if err = json.NewDecoder(r.Body).Decode(&newFeature); err != nil {
		writeUnprocessableEntity(err, w)
		return
	}

	// Validate given values
	if err := newFeature.Validate(); err != nil {
		writeMessage(400, "invalid_feature", err.Error(), w)
		return
	}

	newFeature, err = handler.FeatureService.UpdateFeature(vars["featureKey"], newFeature)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(newFeature); err != nil {
		panic(err)
	}
}

func (handler APIHandler) featureExists(featureKey string) bool {
	return handler.FeatureService.FeatureExists(featureKey)
}

func getJsonHeader() string {
	return "application/json"
}

func writeNotFound(w http.ResponseWriter) {
	writeMessage(http.StatusNotFound, "feature_not_found", "The feature was not found", w)
}

func writeUnprocessableEntity(err error, w http.ResponseWriter) {
	writeMessage(422, "invalid_json", "Cannot decode the given JSON payload", w)
}

func writeMessage(code int, status string, message string, w http.ResponseWriter) {
	apiMessage := APIMessage{code, status, message}
	bytes, _ := json.Marshal(apiMessage)

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(apiMessage.code)
	w.Write(bytes)
}

func hasAccessToFeature(feature m.FeatureFlag, ar AccessRequest) bool {
	// Handle trivial case
	if feature.IsEnabled() {
		return true
	}

	// Access thanks to a group?
	if len(ar.Groups) > 0 {
		for _, group := range ar.Groups {
			if feature.GroupHasAccess(group) {
				return true
			}
		}
	}

	// Access thanks to the user?
	if ar.User > 0 {
		return feature.UserHasAccess(ar.User)
	}

	return false
}
