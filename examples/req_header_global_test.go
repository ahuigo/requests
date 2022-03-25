package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// Set session headers
func TestSendGlobalHeader2(t *testing.T) {
	session := requests.R()

	headerK := "User-Agent"
	headerV := "Custom-Test-Go-User-Agent"
	req, err := session.SetGlobalHeader(headerK, headerV).BuildRequest("post", "http://baidu.com/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
	if req.Header.Get(headerK) != headerV {
		t.Fatalf("Expected header %s is %s", headerK, headerV)
	}
}
