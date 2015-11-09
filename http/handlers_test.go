package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	db "github.com/antoineaugusti/feature-flags/db"
	m "github.com/antoineaugusti/feature-flags/models"
	s "github.com/antoineaugusti/feature-flags/services"
	"github.com/boltdb/bolt"
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
	server = httptest.NewServer(NewRouter(&APIHandler{FeatureService: s.FeatureService{DB: database}}))
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
		[]int{2},
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
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		false,
		[]int{2},
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
	res, _ := http.DefaultClient.Do(request)

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

func TestAccessFeatureFlag(t *testing.T) {
	onStart()
	defer onFinish()

	// Add the default dummy feature
	createDummyFeatureFlag()

	// Access thanks to the user ID
	reader = strings.NewReader(`{"user":2}`)
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s/access", base, "homepage_v2"), reader)
	res, _ := http.DefaultClient.Do(request)

	assertAccessToTheFeature(t, res)

	// No access because of the user ID
	reader = strings.NewReader(`{"user":3}`)
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/%s/access", base, "homepage_v2"), reader)
	res, _ = http.DefaultClient.Do(request)

	assertNoAccessToTheFeature(t, res)

	// Access thanks to the group
	reader = strings.NewReader(`{"user":3, "groups":["dev", "foo"]}`)
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/%s/access", base, "homepage_v2"), reader)
	res, _ = http.DefaultClient.Do(request)

	assertAccessToTheFeature(t, res)
}

func TestListFeatureFlags(t *testing.T) {
	var features m.FeatureFlags
	onStart()
	defer onFinish()

	request, _ := http.NewRequest("GET", base, nil)
	res, _ := http.DefaultClient.Do(request)

	json.NewDecoder(res.Body).Decode(&features)
	assert.Equal(t, 0, len(features))

	// Add the default dummy feature
	createDummyFeatureFlag()

	request, _ = http.NewRequest("GET", base, nil)
	res, _ = http.DefaultClient.Do(request)

	json.NewDecoder(res.Body).Decode(&features)
	assert.Equal(t, 1, len(features))
	assert.Equal(t, "homepage_v2", features[0].Key)
}

func assertAccessToTheFeature(t *testing.T, res *http.Response) {
	assertResponseWithStatusAndMessage(t, res, http.StatusOK, "has_access", "The user has access to the feature")
}

func assertNoAccessToTheFeature(t *testing.T, res *http.Response) {
	assertResponseWithStatusAndMessage(t, res, http.StatusOK, "not_access", "The user does not have access to the feature")
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
	var apiMessage APIMessage
	assert.Equal(t, res.StatusCode, code)

	json.NewDecoder(res.Body).Decode(&apiMessage)
	assert.Equal(t, status, apiMessage.Status)
	assert.Equal(t, message, apiMessage.Message)
}

func getDummyFeaturePayload() string {
	return `{
      "key":"homepage_v2",
      "enabled":false,
      "users":[2],
      "groups":[
         "dev",
         "admin"
      ],
      "percentage":0
    }`
}

func assertJSONMatchesStructure(t *testing.T, res *http.Response, key string, enabled bool, users []int, groups []string, percentage int) {
	var feature m.FeatureFlag
	json.NewDecoder(res.Body).Decode(&feature)

	assert.Equal(t, key, feature.Key)
	assert.Equal(t, enabled, feature.Enabled)
	assert.Equal(t, intsToUints32(users), feature.Users)
	assert.Equal(t, groups, feature.Groups)
	assert.Equal(t, uint32(percentage), feature.Percentage)
}

func intsToUints32(numbers []int) []uint32 {
	res := make([]uint32, 0)
	for _, nb := range numbers {
		res = append(res, uint32(nb))
	}
	return res
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
