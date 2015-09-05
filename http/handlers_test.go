package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	db "github.com/antoineaugusti/golang-feature-flags/db"
	s "github.com/antoineaugusti/golang-feature-flags/services"
	"github.com/boltdb/bolt"
	"github.com/jmoiron/jsonq"
	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	reader   io.Reader
	base     string
	database *bolt.DB
)

func onStart() {
	database = getTestDB()
	server = httptest.NewServer(NewRouter(&APIHandler{s.FeatureService{database}}))
	base = fmt.Sprintf("%s/features", server.URL)
}

func onFinish() {
	database.Close()
	if err := os.Remove(getDBPath()); err != nil {
		panic(err)
	}
}

func TestAddFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	reader = strings.NewReader(getDummyFeaturePayload())
	request, err := http.NewRequest("POST", base, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// 201 response
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Structure is ok
	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		false,
		[]int{},
		[]string{"dev", "admin"},
		0,
	)

	// Add a feature with the same key
	reader = strings.NewReader(getDummyFeaturePayload())
	request, err = http.NewRequest("POST", base, reader)
	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	assertResponseWithStatusAndMessage(t, res, http.StatusBadRequest, "invalid_feature", "Feature already exists")
}

func TestGetFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	// Add the default dummy feature
	reader = strings.NewReader(getDummyFeaturePayload())
	postRequest, _ := http.NewRequest("POST", base, reader)
	if _, err := http.DefaultClient.Do(postRequest); err != nil {
		panic(err)
	}

	// Try to get the default dummy feature
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", base, "homepage_v2"), nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// 200 response
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Structure is ok
	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		false,
		[]int{},
		[]string{"dev", "admin"},
		0,
	)

	// Try to show an unexisting feature
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/%s", base, "notfound"), nil)
	res, _ = http.DefaultClient.Do(request)

	assert404Response(t, res)
}

func TestDeleteFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	// Add the default dummy feature
	reader = strings.NewReader(getDummyFeaturePayload())
	postRequest, _ := http.NewRequest("POST", base, reader)
	if _, err := http.DefaultClient.Do(postRequest); err != nil {
		panic(err)
	}

	// Try to delete the default dummy feature
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", base, "homepage_v2"), nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// 200 response
	assertResponseWithStatusAndMessage(t, res, http.StatusOK, "feature_deleted", "The feature was successfully deleted")

	// Try to delete an unexisting feature
	request, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/%s", base, "notfound"), nil)
	res, _ = http.DefaultClient.Do(request)

	assert404Response(t, res)
}

func assert404Response(t *testing.T, res *http.Response) {
	assertResponseWithStatusAndMessage(t, res, http.StatusNotFound, "feature_not_found", "The feature was not found")
}

func assertResponseWithStatusAndMessage(t *testing.T, res *http.Response, code int, status, message string) {
	assert.Equal(t, res.StatusCode, code)

	jq := extractJSON(res)
	assert.Equal(t, status, getJSONString(jq, "status"))
	assert.Equal(t, message, getJSONString(jq, "message"))
}

func getDummyFeaturePayload() string {
	return `{
      "key":"homepage_v2",
      "enabled":false,
      "users":[],
      "groups":[
         "dev",
         "admin"
      ],
      "percentage":0
    }`
}

func assertJSONMatchesStructure(t *testing.T, res *http.Response, key string, enabled bool, users []int, groups []string, percentage int) {
	jq := extractJSON(res)

	assert.Equal(t, key, getJSONString(jq, "key"))
	assert.False(t, enabled, getJSONBool(jq, "enabled"))
	assert.Equal(t, users, getJSONArrayOfInts(jq, "users"))
	assert.Equal(t, groups, getJSONArrayOfStrings(jq, "groups"))
	assert.Equal(t, percentage, getJSONInt(jq, "percentage"))
}

func getJSONString(jq *jsonq.JsonQuery, key string) string {
	v, _ := jq.String(key)
	return v
}

func getJSONInt(jq *jsonq.JsonQuery, key string) int {
	v, _ := jq.Int(key)
	return v
}

func getJSONBool(jq *jsonq.JsonQuery, key string) bool {
	v, _ := jq.Bool(key)
	return v
}

func getJSONArrayOfStrings(jq *jsonq.JsonQuery, key string) []string {
	v, _ := jq.ArrayOfStrings(key)
	return v
}

func getJSONArrayOfInts(jq *jsonq.JsonQuery, key string) []int {
	v, _ := jq.ArrayOfInts(key)
	return v
}

func extractJSON(res *http.Response) *jsonq.JsonQuery {
	data := map[string]interface{}{}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	return jsonq.NewQuery(data)
}

func getTestDB() *bolt.DB {
	boltDB, err := bolt.Open(getDBPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	db.GenerateDefaultBucket(db.GetBucketName(), boltDB)

	return boltDB
}

func getDBPath() string {
	return "/tmp/test.db"
}
