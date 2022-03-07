package examples

import (
	"strings"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestAuth(t *testing.T) {
	ts := createEchoServer()
	defer ts.Close()
	// test authentication usernae,password
	resp, err := requests.Get(
		ts.URL+"/echo",
		requests.Auth{"httpwatch", "foo"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Text(), "Authorization: Basic ") {
		t.Fatal("bad auth body:\n" + resp.Text())
	}
	// this save file test PASS
	// resp.SaveFile("auth.jpeg")
}
