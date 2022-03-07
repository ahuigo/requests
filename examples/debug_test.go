package examples

import (
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
)

func TestGetDebug(t *testing.T) {
	session := requests.NewSession().SetDebug(true)
	resp, err := session.Post("https://httpbin.org/post",
		requests.Json{
			"name": "ahuigo",
		},
		&http.Cookie{
			Name:  "count",
			Value: "1",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatalf("bad response")
	}
}
