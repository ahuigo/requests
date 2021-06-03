package examples

import (
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestProxy(t *testing.T) {
	println("5. Get: with proxy")
	req := requests.Requests()
	// req.Proxy("http://192.168.1.190:8888")
	req.Get("https://www.httpbin.org/cookies/set?freeform=1234")
	req.Get("https://www.httpbin.org")
}