package examples

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ahuigo/requests"
)

func TestGetDebug(t *testing.T) {
	ts := createHttpbinServer(0)
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
	ts := createHttpbinServer(0)
	defer ts.Close()

	session := requests.R().SetDebugBody()
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
	if !strings.Contains(curl, "Cookie: count=1") || !strings.Contains(curl, "curl -X POST") {
		t.Fatal("bad curl:", curl)
	}
	//debug response
	dumpResponse := resp.GetDumpResponse()
	if !strings.Contains(dumpResponse, `"body":"{\"name\":\"ahuigo\"}"`) || !strings.Contains(dumpResponse, "Content-Type: application/json") {
		t.Fatal("bad dump response:", dumpResponse)
	}
}
