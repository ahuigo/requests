package examples

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ahuigo/requests"
)

func TestGetDebug(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	// requests.R().SetDebug()		debug requests and response
	// requests.R().SetDebugBody()	debug requests and response(with body)
	session := requests.R().SetDebugBody()
	var resp *requests.Response
	var err error
	output := requests.IoCaptureOutput(func() {
		resp, err = session.Post(ts.URL+"/post",
			requests.Json{
				"name": "ahuigo",
			},
			&http.Cookie{
				Name:  "count",
				Value: "1",
			},
		)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, `{"name":"ahuigo"}`) {
		t.Fatalf("can not find debug body")
	}
	if resp.Text() == "" {
		t.Fatalf("bad response")
	}
}
