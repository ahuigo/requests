package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/requests"
	"github.com/davecgh/go-spew/spew"
)

func TestClose(t *testing.T) {
	req := requests.R()
	for i := 0; i < 10; i++ {
		_, err := req.Post(
			"http://localhost:1337/requests",
			requests.Datas{"SrcIp": "4312"})
		fmt.Printf("\r%d %v", i, err)
		req.Close()
	}

	spew.Dump(req)
	fmt.Println("10 times get test end.")
}

func TestTimeout(t *testing.T) {
	req := requests.R().SetTimeout(time.Second)
	req.Get("http://golang.org")
}
