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

func init() {
	database = getTestDB()
	server = httptest.NewServer(NewRouter(&APIHandler{s.FeatureService{database}}))
	base = fmt.Sprintf("%s/features", server.URL)
}

func TestAddFeatureFlag(t *testing.T) {
	defer closeDB()

	payload := `{
      "key":"homepage_v2",
      "enabled":false,
      "users":[],
      "groups":[
         "dev",
         "admin"
      ],
      "percentage":0
    }`

	reader = strings.NewReader(payload)
	request, err := http.NewRequest("POST", base, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// 201 response
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	// Structure is ok
	assertJSONMatchesStructure(
		t, res,
		"homepage_v2",
		false,
		[]int{},
		[]string{"dev", "admin"},
		0,
	)
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

func closeDB() {
	if err := os.Remove(getDBPath()); err != nil {
		panic(err)
	}
	database.Close()
}
