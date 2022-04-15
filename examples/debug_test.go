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

	// requests.R().SetDebug()		debug requests and response(no body)
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

func TestDebugRequestAndResponse(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	session := requests.R().SetDebug()
	resp, err := session.Post(ts.URL+"/post",
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
	//debug curl requests
	curl := resp.GetDumpCurl()
	if !strings.Contains(curl, "Cookie: count=1") {
		t.Fatal("bad curl:", curl)
	}
	//debug response
	dumpResponse := resp.GetDumpResponse()
	if !strings.Contains(curl, `'{"name":"ahuigo"}'`) {
		t.Fatal("bad dump response:", dumpResponse)
	}
}
