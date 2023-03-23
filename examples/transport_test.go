/**
 * refer to: git@github.com:go-resty/resty.git
 */
package examples

import (
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
)

func TestTransportSet(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	session := requests.R()
	// tsp:= otelhttp.NewTransport(http.DefaultTransport)
	tsp := http.DefaultTransport.(*http.Transport).Clone()
	tsp.MaxIdleConnsPerHost = 1
	tsp.MaxIdleConns = 1
	tsp.MaxConnsPerHost = 1
	session.SetTransport(tsp)

	resp, err := session.Get(ts.URL + "/sleep/11")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatal(resp.Text())
	}
}
