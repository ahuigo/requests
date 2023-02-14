package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
)

// Test Response
func TestResponse(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	// resp, _ := requests.Get("https://httpbin.org/get")
	resp, _ := requests.Get(ts.URL + "/get")
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Time:", resp.Time())
	fmt.Println("Size:", resp.Size())
	fmt.Println("Headers:")
	for key, value := range resp.Header() {
		fmt.Println(key, "=", value)
	}
	fmt.Println("Cookies:")
	for i, cookie := range resp.Cookies() {
		fmt.Printf("cookie%d: name:%s value:%s\n", i, cookie.Name, cookie.Value)
	}

}

// Test response headers
func TestResponseHeader(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	// resp, _ := requests.Get("https://httpbin.org/get")
	resp, _ := requests.Get(ts.URL + "/get")

	if resp.R.Header.Get("content-type") != "application/json" {
		t.Fatal("bad response header")
	}

	println("content-type:", resp.R.Header.Get("content-type"))
}

// Test response body
func TestResponseBody(t *testing.T) {
	ts := createHttpbinServer(false)
	defer ts.Close()

	// resp, _ := requests.Get("https://httpbin.org/get")
	resp, _ := requests.Get(ts.URL + "/get")
	println(resp.Body())
	println(resp.Text())
}
