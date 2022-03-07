package examples

import (
	ejson "encoding/json"
	"testing"

	"github.com/ahuigo/requests"
)

// Post QueryString and content-type: none
func TestPostParams(t *testing.T) {
	println("Test POST: post params")
	resp, err := requests.Post(
		"https://www.httpbin.org/post",
		requests.Params{
			"name": "ahuigo",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Args struct {
			Name string
		}
	}{}
	_ = resp.Json(&data)
	if data.Args.Name != "ahuigo" {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Form UrlEncoded data: application/x-www-form-urlencoded
func TestPostFormUrlEncode(t *testing.T) {
	resp, err := requests.Post(
		"https://www.httpbin.org/post",
		requests.Datas{
			"name": "ahuigo",
		},
	)
	if err != nil {
		t.Error(err)
	}
	var data = struct {
		Form struct {
			Name string
		}
	}{}
	resp.Json(&data)
	if data.Form.Name != "ahuigo" {
		t.Error("invalid response body:", resp.Text())
	}
}

// Test POST:  multipart/form-data; boundary=....
func TestPostFormData(t *testing.T) {
	resp, err := requests.Post(
		"https://www.httpbin.org/post",
		requests.FormData{
			"name": "ahuigo",
		},
	)
	if err != nil {
		t.Error(err)
	}
	var data = struct {
		Form struct {
			Name string
		}
	}{}
	resp.Json(&data)
	if data.Form.Name != "ahuigo" {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Json: application/json
func TestPostJson(t *testing.T) {
	println("Test POST: post json data")
	// You can also use `json := requests.Jsoni(anyTypeData)`
	json := requests.Json{
		"name": "Alex",
	}
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

// Post Raw Bypes: text/plain
func TestRawBytes(t *testing.T) {
	println("Test POST: post bytes data")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post", []byte(rawText))
	if err != nil {
		t.Error(err)
	}
	var data = struct {
		Data string
	}{}
	resp.Json(&data)
	if data.Data != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Raw String: text/plain
func TestRawString1(t *testing.T) {
	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post", rawText,
		requests.Header{"Content-Type": "text/plain"},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data interface{}
	resp.Json(&data)
	if data.(map[string]interface{})["data"].(string) != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}
func TestRawString2(t *testing.T) {
	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post", rawText,
		requests.ContentTypePlain,
	)
	if err != nil {
		t.Fatal(err)
	}
	var data interface{}
	resp.Json(&data)
	if data.(map[string]interface{})["data"].(string) != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

func TestRawString3(t *testing.T) {
	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post",
		requests.ContentTypePlain,
		rawText,
	)
	if err != nil {
		t.Fatal(err)
	}
	var data interface{}
	resp.Json(&data)
	if data.(map[string]interface{})["data"].(string) != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

// TestPostEncodedString: application/x-www-form-urlencoded
func TestPostEncodedString(t *testing.T) {
	resp, err := requests.Post("https://www.httpbin.org/post", "name=Alex&age=29")
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Form struct {
			Name string
			Age  string
		}
	}{}
	resp.Json(&data)
	if data.Form.Age != "29" {
		t.Error("invalid response body:", resp.Text())
	}
}
