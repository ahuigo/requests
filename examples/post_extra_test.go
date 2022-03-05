package examples

import (
	ejson "encoding/json"
	"testing"

	"github.com/ahuigo/requests"
)

// Post Json: application/json
func TestPostJsonInterface(t *testing.T) {
	anyTypeData := map[string]interface{}{
		"name": "Alex",
	}
	json := requests.Jsoni(anyTypeData)
	resp, err := requests.Post("https://www.httpbin.org/post", json)
	if err != nil {
		t.Error(err)
	}

	// parse data
	var data = struct {
		Data string
	}{}
	resp.Json(&data)

	// is expected results
	jsonData, _ := ejson.Marshal(json) // if data.Data!= "{\"name\":\"Alex\"}"{
	if data.Data != string(jsonData) {
		t.Error("invalid response body:", resp.Text())
	}
}
