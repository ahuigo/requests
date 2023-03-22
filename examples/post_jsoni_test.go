package examples

import (
	ejson "encoding/json"
	"testing"

	"github.com/ahuigo/requests"
)

// Post Json: application/json
func TestPostJsonInterface(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	anyTypeData := map[string]interface{}{
		"name": "Alex",
	}
	json := requests.Jsoni(anyTypeData)
	resp, err := requests.Post(ts.URL+"/post", json)
	if err != nil {
		t.Error(err)
	}

	// parse data
	var data = struct {
		Body string
	}{}
	resp.Json(&data)

	// is expected results
	jsonData, _ := ejson.Marshal(json) // if data.Data!= "{\"name\":\"Alex\"}"{
	if data.Body != string(jsonData) {
		t.Fatalf("expected: %s,\ninvalid body:%s", string(jsonData), resp.Text())
	}
}
