package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

func TestRetryCondition(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	// retry 3 times
	maxRetries := 3
	r := requests.R().SetDebug().
		SetRetryCount(maxRetries).
		SetRetryCondition(func(resp *requests.Response, err error) bool {
			if err != nil {
				return true
			}
			var json map[string]interface{}
			resp.Json(&json)
			return json["headers"] != "a"
		})

	resp, err := r.Post(ts.URL+"/post", []byte("alex"))
	if err != nil {
		t.Fatal(err)
	}

	if resp.Attempt != maxRetries {
		t.Fatalf("Attemp %d not equal to %d", resp.Attempt, maxRetries)
	}

	var json map[string]string
	resp.Json(&json)
	if json["body"] != "alex" {
		t.Fatalf("Bad response body:%s", resp.Text())
	}
	if json["method"] != "POST" {
		t.Fatalf("Bad request method:%s", resp.Text())
	}
}

func TestRetryConditionFalse(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	// retry 3 times
	r := requests.R().
		SetRetryCount(3).
		SetRetryCondition(func(resp *requests.Response, err error) bool {
			return false
		})

	resp, err := r.Get(ts.URL + "/get")
	if err != nil {
		t.Fatal(err)
	}

	if resp.Attempt != 0 {
		t.Fatalf("Attemp %d not equal to %d", resp.Attempt, 0)
	}

	var json map[string]interface{}
	resp.Json(&json)
	t.Logf("response json:%#v\n", json["headers"])
}
