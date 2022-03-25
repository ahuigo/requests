package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// Send headers
func TestSendHeader(t *testing.T) {
	println("Test Get: send header")
	requests.Get(
		"http://www.zhanluejia.net.cn",
		requests.Header{"Referer": "http://www.jeapedu.com"},
	)
}
