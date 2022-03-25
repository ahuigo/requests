package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// Test Session with cookie
func TestSessionWithCookie(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	sess := requests.R().SetDebug()
	_, err := sess.Get(ts.URL + "/cookie/count")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := sess.Get(ts.URL + "/cookie/count")
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCookie("count") != "2" {
		t.Fatal("Failed to set cookie count")
	}
}
