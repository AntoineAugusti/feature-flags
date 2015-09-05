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

	res := createDummyFeatureFlag()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		false,
		[]int{},
		[]string{"dev", "admin"},
		0,
	)

	// Add a feature with the same key
	res = createDummyFeatureFlag()

	assertResponseWithStatusAndMessage(t, res, http.StatusBadRequest, "invalid_feature", "Feature already exists")

	// Add with an invalid JSON payload
	reader = strings.NewReader("{foo:bar}")
	request, _ := http.NewRequest("POST", base, reader)
	res, _ = http.DefaultClient.Do(request)

	assert422Response(t, res)
}

func TestGetFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	// Add the default dummy feature
	createDummyFeatureFlag()

	// Get the default dummy feature
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", base, "homepage_v2"), nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		false,
		[]int{},
		[]string{"dev", "admin"},
		0,
	)

	// Show an unexisting feature
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/%s", base, "notfound"), nil)
	res, _ = http.DefaultClient.Do(request)

	assert404Response(t, res)
}

func TestDeleteFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	// Add the default dummy feature
	createDummyFeatureFlag()

	// Delete the default dummy feature
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", base, "homepage_v2"), nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	assertResponseWithStatusAndMessage(t, res, http.StatusOK, "feature_deleted", "The feature was successfully deleted")

	// Delete an unexisting feature
	request, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/%s", base, "notfound"), nil)
	res, _ = http.DefaultClient.Do(request)

	assert404Response(t, res)
}

func TestEditFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	// Add the default dummy feature
	createDummyFeatureFlag()

	// Edit the default dummy feature
	payload := `{
      "enabled":true,
      "users":[1,2],
      "groups":[
         "a",
         "b"
      ],
      "percentage":42
    }`

	reader = strings.NewReader(payload)
	request, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", base, "homepage_v2"), reader)
	res, _ := http.DefaultClient.Do(request)

	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		true,
		[]int{1, 2},
		[]string{"a", "b"},
		42,
	)

	// Edit an unexisting feature
	request, _ = http.NewRequest("PATCH", fmt.Sprintf("%s/%s", base, "notfound"), reader)
	res, _ = http.DefaultClient.Do(request)

	assert404Response(t, res)

	// Edit with an invalid JSON payload
	reader = strings.NewReader("{foo:bar}")
	request, _ = http.NewRequest("PATCH", fmt.Sprintf("%s/%s", base, "homepage_v2"), reader)
	res, _ = http.DefaultClient.Do(request)

	assert422Response(t, res)

	// Edit with an invalid percentage
	reader = strings.NewReader(`{"percentage":101}`)
	request, _ = http.NewRequest("PATCH", fmt.Sprintf("%s/%s", base, "homepage_v2"), reader)
	res, _ = http.DefaultClient.Do(request)

	assertResponseWithStatusAndMessage(t, res, http.StatusBadRequest, "invalid_feature", "Percentage must be between 0 and 100")
}

func createDummyFeatureFlag() *http.Response {
	reader = strings.NewReader(getDummyFeaturePayload())
	postRequest, _ := http.NewRequest("POST", base, reader)
	res, err := http.DefaultClient.Do(postRequest)
	if err != nil {
		panic(err)
	}

	return res
}

func assert422Response(t *testing.T, res *http.Response) {
	assertResponseWithStatusAndMessage(t, res, 422, "invalid_json", "Cannot decode the given JSON payload")
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
	assert.Equal(t, enabled, getJSONBool(jq, "enabled"))
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
