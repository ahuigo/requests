package examples

import (
	ejson "encoding/json"
	"strings"
	"testing"

	"github.com/ahuigo/requests"
)

// Post QueryString and content-type: none
// curl -X POST "https://www.httpbin.org/post?name=ahuigo"
func TestPostParams(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	resp, err := requests.Post(
		ts.URL+"/post",
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
		t.Fatal("invalid response body:", resp.Text())
	}
}

// Post Form UrlEncoded data: application/x-www-form-urlencoded
// curl -H 'Content-Type: application/x-www-form-urlencoded' https://www.httpbin.org/post -d 'name=ahuigo'
func TestPostFormUrlEncode(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()
	resp, err := requests.Post(
		ts.URL+"/post",
		requests.Datas{
			"name": "ahuigo",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Body string
	}{}
	resp.Json(&data)
	if data.Body != "name=ahuigo" {
		t.Fatal("invalid response body:", resp.Text())
	}
}

// Test POST:  multipart/form-data; boundary=....
// curl https://www.httpbin.org/post -F 'name=ahuigo'
func TestPostFormData(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()
	resp, err := requests.Post(
		ts.URL+"/post",
		requests.FormData{
			"name": "ahuigo",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Body string
	}{}
	resp.Json(&data)
	if !strings.Contains(data.Body, "form-data; name=\"name\"\r\n\r\nahuigo\r\n") {
		t.Error("invalid response body:", resp.Text())
		t.Error("invalid response body:", data.Body)
	}
}

// Post Json: application/json
// curl -H "Content-Type: application/json" https://www.httpbin.org/post -d '{"name":"Alex"}'
func TestPostJson(t *testing.T) {
	println("Test POST: post json data")
	// You can also use `json := requests.Jsoni(anyTypeData)`
	json := requests.Json{
		"name": "Alex",
	}
	resp, err := requests.Post("https://www.httpbin.org/post", json)
	if err != nil {
		t.Fatal(err)
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
// curl -H "Content-Type: text/plain" https://www.httpbin.org/post -d 'raw data: Hi, Jack!'
func TestRawBytes(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	println("Test POST: post bytes data")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post(ts.URL+"/post", []byte(rawText))
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Body string
	}{}
	resp.Json(&data)
	if data.Body != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Raw String: text/plain
// curl -H "Content-Type: text/plain" https://www.httpbin.org/post -d 'raw data: Hi, Jack!'
func TestRawString1(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post(ts.URL+"/post", rawText,
		requests.Header{"Content-Type": "text/plain"},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data interface{}
	resp.Json(&data)
	if data.(map[string]interface{})["body"].(string) != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Raw String2: text/plain
// curl -H "Content-Type: text/plain" https://www.httpbin.org/post -d 'raw data: Hi, Jack!'
func TestRawString2(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post(
		ts.URL+"/post",
		requests.ContentTypePlain,
		rawText,
	)
	if err != nil {
		t.Fatal(err)
	}
	var data interface{}
	resp.Json(&data)
	if data.(map[string]interface{})["body"].(string) != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

// TestPostEncodedString: application/x-www-form-urlencoded
// curl -H 'Content-Type: application/x-www-form-urlencoded' http://0:4500/post -d 'name=Alex&age=29'
func TestPostEncodedString(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	resp, err := requests.Post(ts.URL+"/post", "name=Alex&age=29")
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Body string
	}{}
	resp.Json(&data)
	if data.Body != "name=Alex\u0026age=29" {
		t.Error("invalid response body:", resp.Text())
	}
}
