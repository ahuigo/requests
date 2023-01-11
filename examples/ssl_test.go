package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

func TestSkipSsl(t *testing.T) {
	session := requests.R()
	session = session.SkipSsl(true)
	resp, err:= session.Get("https://www.httpbin.org/get")
	if err!= nil {
		t.Fatal(err)
	}
	if resp.Text()==""{
		t.Fatal(resp.Text())
	}
}