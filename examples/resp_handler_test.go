package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// Test response headers
func TestResponseHeaderTodo(t *testing.T) {
	resp, _ := requests.Get("https://httpbin.org/get")
	println("content-type:", resp.R.Header.Get("content-type"))
	//println(resp.Text())
}
