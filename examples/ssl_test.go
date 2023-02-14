package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

func TestSkipSsl(t *testing.T) {
	// 1. create tls test server
	ts := createHttpbinServer(true)
	defer ts.Close()

	session := requests.R()

	// 2. fake CA certificate
	// session.SetCaCert("conf/rootCA.crt")

	// 3. skip ssl
	session = session.SkipSsl(true)

	// 4. send get request
	resp, err := session.Get(ts.URL + "/get?a=1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatal(resp.Text())
	}
}
