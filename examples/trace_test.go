/**
 * refer to: git@github.com:go-resty/resty.git
 */
package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// test context: cancel multi
func TestTrace(t *testing.T) {
	ts := createHttpbinServer(0)
	defer ts.Close()

	params := requests.Params{"name": "ahuigo", "page": "1"}
	resp, err := requests.R().SetDebug().Get(ts.URL+"/get", params)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("connTime:%+v", resp.TraceInfo.ConnTime)
}
