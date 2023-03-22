package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// Get example: fetch json response
func TestGetJson(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	// resp, err := requests.Get("https://httpbin.org/json")
	resp, err := requests.Get(ts.URL + "/get")
	if err == nil {
		var json map[string]interface{}
		err = resp.Json(&json)
		t.Logf("response json:%#v\n", json)
	}
	if err != nil {
		t.Fatal(err)
	}
}

// Get with params
func TestGetParams(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	params := requests.Params{"name": "ahuigo", "page": "1"}
	resp, err := requests.Get(ts.URL+"/get", params)

	if err != nil {
		t.Fatal(err)
	}
	if err == nil {
		type HbResponse struct {
			Args map[string]string `json:"args"`
		}
		json := &HbResponse{}
		if err := resp.Json(&json); err != nil {
			t.Fatalf("bad json:%s", resp.Text())
		}
		if json.Args["name"] != "ahuigo" {
			t.Fatal("Invalid response: " + resp.Text())
		}
	}
}

// Support array args like: ids=id1&ids=id2&ids=id3
func TestGetParamArray(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	params := requests.Params{"name": "ahuigo"}
	paramsArray := requests.ParamsArray{
		"ids": []string{"id1", "id2"},
	}
	resp, err := requests.Get(ts.URL+"/get", params, paramsArray)

	if err != nil {
		t.Fatal(err)
	}
	if err == nil {
		type HbResponse struct {
			Args map[string]string `json:"args"`
		}
		json := &HbResponse{}
		if err := resp.Json(&json); err != nil {
			t.Fatalf("bad json:%s", resp.Text())
		}
		if json.Args["name"] != "ahuigo" || json.Args["ids"] != "id1,id2" {
			t.Fatal("Invalid response: " + resp.Text())
		}
	}
}

func TestGetWithHeader(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	params := requests.Params{"name": "ahuigo"}
	// var header requests.Header
	header := requests.Header(nil)
	resp, err := requests.Get(ts.URL+"/get", params, header)

	if err != nil {
		t.Fatal(err)
	}
	if err == nil {
		type HbResponse struct {
			Args map[string]string `json:"args"`
		}
		json := &HbResponse{}
		if err := resp.Json(&json); err != nil {
			t.Fatalf("bad json:%s", resp.Text())
		}
		t.Log(resp.Text())
	}
}
