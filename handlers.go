package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type APIHandler struct {
	FeatureService FeatureService
}

type APIMessage struct {
	code    int
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (handler *APIHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!\n")
}

func (handler *APIHandler) FeatureIndex(w http.ResponseWriter, r *http.Request) {
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

func (handler *APIHandler) FeatureShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the feature exists
	if !handler.FeatureExists(vars["featureKey"]) {
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

func (handler *APIHandler) FeatureRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the feature exists
	if !handler.FeatureExists(vars["featureKey"]) {
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

func (handler *APIHandler) FeatureCreate(w http.ResponseWriter, r *http.Request) {
	var feature FeatureFlag
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &feature); err != nil {
		w.Header().Set("Content-Type", getJsonHeader())
		w.WriteHeader(422) // Unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	err = handler.FeatureService.AddFeature(feature)
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

func (handler *APIHandler) FeatureExists(featureKey string) bool {
	return handler.FeatureService.FeatureExists(featureKey)
}

func getJsonHeader() string {
	return "application/json"
}

func writeNotFound(w http.ResponseWriter) {
	writeMessage(http.StatusNotFound, "feature_not_found", "The feature was not found", w)
}

func writeMessage(code int, status string, message string, w http.ResponseWriter) {
	apiMessage := APIMessage{code, status, message}
	bytes, _ := json.Marshal(apiMessage)

	w.Header().Set("Content-Type", getJsonHeader())
	w.WriteHeader(apiMessage.code)
	w.Write(bytes)
}
