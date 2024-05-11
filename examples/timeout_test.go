package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/ahuigo/requests"
	"github.com/davecgh/go-spew/spew"
)

func TestKeepaliveClose(t *testing.T) {
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

func TestBodyNotClose(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	// Do not close body
	session := requests.R().SetDoNotCloseBody(true)
	resp, err := session.Get(ts.URL + "/get")
	if err == nil {
		var json map[string]interface{}
		err = resp.Json(&json)
		t.Logf("response json:%#v\n", json)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp.R.Body.Close() // close body manually
}

func TestTimeout(t *testing.T) {
	req := requests.R().SetTimeout(time.Second)
	_, err := req.Get("http://golang.org")
	t.Log(err)
}
